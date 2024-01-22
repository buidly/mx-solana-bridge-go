package sol

import (
	"github.com/gagliardetto/solana-go"
)

type ProgramAddresses struct {
	BridgeProgramAddress     solana.PublicKey
	TokensSafeProgramAddress solana.PublicKey
}
