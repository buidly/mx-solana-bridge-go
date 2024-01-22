package solmultiversx

import (
	"bytes"
	"sync"

	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-solana-bridge-go/core"
)

type signaturesHolder struct {
	mut            sync.RWMutex
	signedMessages map[string]*core.SignedMessage
	ethMessages    []*core.EthereumSignature
	log            logger.Logger
}

// NewSignatureHolder creates a new signatureHolder
func NewSignatureHolder() *signaturesHolder {
	return &signaturesHolder{
		signedMessages: make(map[string]*core.SignedMessage),
		ethMessages:    make([]*core.EthereumSignature, 0),
		log:            core.NewLoggerWithIdentifier(logger.GetOrCreate("signaturesHolder"), "signaturesHolder"),
	}
}

// ProcessNewMessage will store the new messages
func (sh *signaturesHolder) ProcessNewMessage(msg *core.SignedMessage, ethMsg []*core.EthereumSignature) {
	if msg == nil || ethMsg == nil {
		return
	}

	sh.mut.Lock()
	defer sh.mut.Unlock()

	sh.signedMessages[msg.UniqueID()] = msg
	sh.log.Info("New eth signatures processed", "ethMsg", ethMsg)
	sh.ethMessages = append(sh.ethMessages, ethMsg...)
}

// AllStoredSignatures will return the stored signatures
func (sh *signaturesHolder) AllStoredSignatures() []*core.SignedMessage {
	sh.mut.RLock()
	defer sh.mut.RUnlock()

	result := make([]*core.SignedMessage, 0, len(sh.signedMessages))
	for _, msg := range sh.signedMessages {
		result = append(result, msg)
	}

	return result
}

// Signatures will provide all gathered signatures for a given message hash
func (sh *signaturesHolder) Signatures(msgHash []byte) [][]byte {
	sh.mut.RLock()
	defer sh.mut.RUnlock()

	uniqueEthSigs := make(map[string]struct{})
	for _, ethMsg := range sh.ethMessages {
		if bytes.Equal(ethMsg.MessageHash, msgHash) {
			uniqueEthSigs[string(ethMsg.Signature)] = struct{}{}
		}
	}

	result := make([][]byte, 0, len(sh.signedMessages))
	for sig := range uniqueEthSigs {
		result = append(result, []byte(sig))
	}

	return result
}

// ClearStoredSignatures will clear any stored signatures
func (sh *signaturesHolder) ClearStoredSignatures() {
	sh.mut.Lock()
	defer sh.mut.Unlock()

	sh.signedMessages = make(map[string]*core.SignedMessage)
	sh.ethMessages = make([]*core.EthereumSignature, 0)
}

// IsInterfaceNil returns true if there is no value under the interface
func (sh *signaturesHolder) IsInterfaceNil() bool {
	return sh == nil
}
