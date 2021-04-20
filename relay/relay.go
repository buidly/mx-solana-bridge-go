package relay

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/ElrondNetwork/elrond-eth-bridge/bridge"
	"github.com/ElrondNetwork/elrond-eth-bridge/bridge/elrond"
	"github.com/ElrondNetwork/elrond-eth-bridge/bridge/eth"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-go/config"
	"github.com/ElrondNetwork/elrond-go/core"
	factoryMarshalizer "github.com/ElrondNetwork/elrond-go/marshal/factory"
	"github.com/ElrondNetwork/elrond-go/p2p"
	"github.com/ElrondNetwork/elrond-go/p2p/libp2p"
)

const (
	ActionsTopicName = "actions/1"
	JoinedAction     = "joined"

	PrivateTopicName = "private/1"

	Timeout             = 30 * time.Second
	MinSignaturePercent = 67
)

type State int

const (
	Join                  State = 0
	GetPendingTransaction State = 1
	Propose               State = 2
	WaitForSignatures     State = 3
	Execute               State = 4
	WaitForProposal       State = 5
	WaitForExecute        State = 6
)

type Peers []core.PeerID

type Timer interface {
	sleep(d time.Duration)
	after(d time.Duration) <-chan time.Time
	nowUnix() int64
}

type defaultTimer struct{}

func (s *defaultTimer) sleep(d time.Duration) {
	time.Sleep(d)
}

func (s *defaultTimer) after(d time.Duration) <-chan time.Time {
	return time.After(d)
}

func (s *defaultTimer) nowUnix() int64 {
	return time.Now().Unix()
}

type NetMessenger interface {
	ID() core.PeerID
	Bootstrap() error
	Addresses() []string
	RegisterMessageProcessor(string, p2p.MessageProcessor) error
	HasTopic(name string) bool
	CreateTopic(name string, createChannelForTopic bool) error
	Broadcast(topic string, buff []byte)
	SendToConnectedPeer(topic string, buff []byte, peerID core.PeerID) error
	Close() error
}

type Relay struct {
	mu sync.Mutex

	peers     Peers
	messenger NetMessenger
	timer     Timer
	log       logger.Logger

	ethBridge    bridge.Bridge
	elrondBridge bridge.Bridge

	initialState       State
	pendingTransaction *bridge.DepositTransaction
}

func NewRelay(config *Config, log logger.Logger) (*Relay, error) {
	ethBridge, err := eth.NewClient(config.Eth)
	if err != nil {
		return nil, err
	}

	elrondBridge, err := elrond.NewClient(config.Elrond)
	if err != nil {
		return nil, err
	}

	messenger, err := buildNetMessenger(config.P2P)
	if err != nil {
		return nil, err
	}

	return &Relay{
		peers:     make(Peers, 0),
		messenger: messenger,
		timer:     &defaultTimer{},
		log:       log,

		ethBridge:    ethBridge,
		elrondBridge: elrondBridge,
	}, nil
}

func (r *Relay) Start(ctx context.Context) error {
	if err := r.init(); err != nil {
		return nil
	}

	ch := make(chan State, 1)
	ch <- r.initialState

	for {
		select {
		case state := <-ch:
			switch state {
			case Join:
				go r.join(ch)
			case GetPendingTransaction:
				go r.getPendingTransaction(ctx, ch)
			case Propose:
				go r.propose(ctx, ch)
			case WaitForProposal:
				go r.waitForProposal(ctx, ch)
			case WaitForSignatures:
				go r.waitForSignatures(ctx, ch)
			case Execute:
				go r.execute(ctx, ch)
			case WaitForExecute:
				go r.waitForExecute(ctx, ch)
			}
		case <-ctx.Done():
			return r.Stop()
		}
	}
}

func (r *Relay) Stop() error {
	return r.messenger.Close()
}

// State

func (r *Relay) join(ch chan State) {
	rand.Seed(time.Now().UnixNano())
	v := rand.Intn(5)
	r.timer.sleep(time.Duration(v) * time.Second)
	r.messenger.Broadcast(ActionsTopicName, []byte(JoinedAction))
	ch <- GetPendingTransaction
}

func (r *Relay) getPendingTransaction(ctx context.Context, ch chan State) {
	r.pendingTransaction = r.ethBridge.GetPendingDepositTransaction(ctx)

	if r.pendingTransaction == nil {
		r.timer.sleep(Timeout / 10)
		ch <- GetPendingTransaction
	} else {
		ch <- Propose
	}
}

func (r *Relay) propose(ctx context.Context, ch chan State) {
	if r.amITheLeader() {
		r.elrondBridge.Propose(ctx, r.pendingTransaction)
		ch <- WaitForSignatures
	} else {
		ch <- WaitForProposal
	}
}

func (r *Relay) waitForProposal(ctx context.Context, ch chan State) {
	select {
	case <-r.timer.after(Timeout):
		if r.elrondBridge.WasProposed(ctx, r.pendingTransaction) {
			r.elrondBridge.Sign(ctx, r.pendingTransaction)
			ch <- WaitForSignatures
		} else {
			ch <- Propose
		}
	case <-ctx.Done():
		if err := r.Stop(); err != nil {
			r.log.Error(err.Error())
		}
	}
}

func (r *Relay) waitForSignatures(ctx context.Context, ch chan State) {
	select {
	case <-r.timer.after(Timeout):
		count := r.elrondBridge.SignersCount(ctx, r.pendingTransaction)
		minCountRequired := math.Ceil(float64(len(r.peers)) * MinSignaturePercent / 100)

		if count >= uint(minCountRequired) && count > 0 {
			ch <- Execute
		} else {
			ch <- WaitForSignatures
		}
	case <-ctx.Done():
		if err := r.Stop(); err != nil {
			r.log.Error(err.Error())
		}
	}
}

func (r *Relay) execute(ctx context.Context, ch chan State) {
	if r.amITheLeader() {
		hash, err := r.elrondBridge.Execute(ctx, r.pendingTransaction)

		if err != nil {
			r.log.Error(err.Error())
		}

		r.log.Info(fmt.Sprintf("Bridge transaction executed with hash %q", hash))
	}

	ch <- WaitForExecute
}

func (r *Relay) waitForExecute(ctx context.Context, ch chan State) {
	select {
	case <-r.timer.after(Timeout):
		if r.elrondBridge.WasExecuted(ctx, r.pendingTransaction) {
			ch <- GetPendingTransaction
		} else {
			ch <- Execute
		}
	case <-ctx.Done():
		if err := r.Stop(); err != nil {
			r.log.Error(err.Error())
		}
	}
}

// MessageProcessor

func (r *Relay) ProcessReceivedMessage(message p2p.MessageP2P, _ core.PeerID) error {
	r.log.Info(fmt.Sprintf("Got message on topic %q\n", message.Topic()))

	switch message.Topic() {
	case ActionsTopicName:
		r.log.Info(fmt.Sprintf("Action: %q\n", string(message.Data())))
		switch string(message.Data()) {
		case JoinedAction:
			r.addPeer(message.Peer())
			if err := r.broadcastTopology(message.Peer()); err != nil {
				r.log.Error(err.Error())
			}
		}
	case PrivateTopicName:
		if err := r.setTopology(message.Data()); err != nil {
			r.log.Error(err.Error())
		}
	}

	return nil
}

func (r *Relay) IsInterfaceNil() bool {
	return r == nil
}

func (r *Relay) broadcastTopology(toPeer core.PeerID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.peers) == 1 && r.peers[0] == r.messenger.ID() {
		return nil
	}

	var data bytes.Buffer
	enc := gob.NewEncoder(&data)
	if err := enc.Encode(r.peers); err != nil {
		return err
	}

	if err := r.messenger.SendToConnectedPeer(PrivateTopicName, data.Bytes(), toPeer); err != nil {
		return err
	}

	return nil
}

func (r *Relay) setTopology(data []byte) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// TODO: ignore if peers are already set
	if len(r.peers) > 1 {
		// ignore this call if we already have peers
		// TODO: find a better way here
		return nil
	}

	dec := gob.NewDecoder(bytes.NewReader(data))
	var topology Peers
	if err := dec.Decode(&topology); err != nil {
		return err
	}
	r.peers = topology

	return nil
}

func (r *Relay) addPeer(peerID core.PeerID) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// TODO: account for peers that rejoin
	if len(r.peers) == 0 || r.peers[len(r.peers)-1] < peerID {
		r.peers = append(r.peers, peerID)
		return
	}

	// TODO: can optimize via binary search
	for index, peer := range r.peers {
		if peer > peerID {
			r.peers = append(r.peers, "")
			copy(r.peers[index+1:], r.peers[index:])
			r.peers[index] = peerID
			break
		}
	}
}

// Helpers

func (r *Relay) init() error {
	if err := r.messenger.Bootstrap(); err != nil {
		return err
	}

	r.timer.sleep(10 * time.Second)
	r.log.Info(fmt.Sprint(r.messenger.Addresses()))

	if err := r.registerTopicProcessors(); err != nil {
		return nil
	}

	return nil
}

func (r *Relay) registerTopicProcessors() error {
	topics := []string{ActionsTopicName, PrivateTopicName}
	for _, topic := range topics {
		if !r.messenger.HasTopic(topic) {
			if err := r.messenger.CreateTopic(topic, true); err != nil {
				return err
			}
		}

		r.log.Info(fmt.Sprintf("Registered on topic %q", topic))
		if err := r.messenger.RegisterMessageProcessor(topic, r); err != nil {
			return err
		}
	}

	return nil
}

func (r *Relay) amITheLeader() bool {
	if len(r.peers) == 0 {
		return false
	} else {
		numberOfPeers := int64(len(r.peers))
		index := (r.timer.nowUnix() / int64(Timeout.Seconds())) % numberOfPeers

		return r.peers[index] == r.messenger.ID()
	}
}

func buildNetMessenger(cfg ConfigP2P) (NetMessenger, error) {
	internalMarshalizer, err := factoryMarshalizer.NewMarshalizer("gogo protobuf")
	if err != nil {
		panic(err)
	}

	nodeConfig := config.NodeConfig{
		Port:                       cfg.Port,
		Seed:                       cfg.Seed,
		MaximumExpectedPeerCount:   0,
		ThresholdMinConnectedPeers: 0,
	}
	peerDiscoveryConfig := config.KadDhtPeerDiscoveryConfig{
		Enabled:                          true,
		RefreshIntervalInSec:             5,
		ProtocolID:                       "/erd/relay/1.0.0",
		InitialPeerList:                  cfg.InitialPeerList,
		BucketSize:                       0,
		RoutingTableRefreshIntervalInSec: 300,
	}

	p2pConfig := config.P2PConfig{
		Node:                nodeConfig,
		KadDhtPeerDiscovery: peerDiscoveryConfig,
		Sharding: config.ShardingConfig{
			TargetPeerCount:         0,
			MaxIntraShardValidators: 0,
			MaxCrossShardValidators: 0,
			MaxIntraShardObservers:  0,
			MaxCrossShardObservers:  0,
			Type:                    "NilListSharder",
		},
	}

	args := libp2p.ArgsNetworkMessenger{
		Marshalizer:   internalMarshalizer,
		ListenAddress: libp2p.ListenAddrWithIp4AndTcp,
		P2pConfig:     p2pConfig,
		SyncTimer:     &libp2p.LocalSyncTimer{},
	}

	messenger, err := libp2p.NewNetworkMessenger(args)
	if err != nil {
		panic(err)
	}

	return messenger, nil
}
