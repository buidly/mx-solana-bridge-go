package sol

import "errors"

var (
	errQuorumNotReached                    = errors.New("quorum not reached")
	errInsufficientSftBalance              = errors.New("insufficient SFT balance")
	errInsufficientBalance                 = errors.New("insufficient balance")
	errPublicKeyCast                       = errors.New("error casting public key to ECDSA")
	errNilBroadcaster                      = errors.New("nil broadcaster")
	errNilSignaturesHolder                 = errors.New("nil signatures holder")
	errNilGasHandler                       = errors.New("nil gas handler")
	errInvalidGasLimit                     = errors.New("invalid gas limit")
	errNilSolClient                        = errors.New("nil sol client")
	errDepositsAndBatchDepositsCountDiffer = errors.New("deposits and batch.DepositsCount differs")
	errBatchNotInitialized                 = errors.New("batch not initialized")
	errMintAccountDidNotMatch              = errors.New("Mint Account did not match")
	errNilFee                              = errors.New("nil fee")
)
