package chain

import (
	"fmt"
	"strings"
)

const (
	solanaToMultiversXName              = "SolanaToMultiversX"
	multiversXToSolanaName              = "MultiversXToSolana"
	baseLogIdTemplate                   = "%sMultiversX-Base"
	multiversXClientLogIdTemplate       = "%sMultiversX-MultiversXClient"
	multiversXDataGetterLogIdTemplate   = "%sMultiversX-MultiversXDataGetter"
	solanaClientLogIdTemplate           = "%sMultiversX-%sClient"
	multiversXRoleProviderLogIdTemplate = "%sMultiversX-MultiversXRoleProvider"
	solanaRoleProviderLogIdTemplate     = "%sMultiversX-%sRoleProvider"
	broadcasterLogIdTemplate            = "%sMultiversX-Broadcaster"
)

// Chain defines all the chain supported
type Chain string

const (
	// MultiversX is the string representation of the MultiversX chain
	MultiversX Chain = "msx"

	// Solana is the string representation of the Solana chain
	Solana Chain = "Solana"
)

// ToLower returns the lowercase string of chain
func (c Chain) ToLower() string {
	return strings.ToLower(string(c))
}

// SolanaToMultiversXName returns the string using chain value and solanaToMultiversXName
func (c Chain) SolanaToMultiversXName() string {
	return solanaToMultiversXName
}

// MultiversXToSolanaName returns the string using chain value and multiversXToSolanaName
func (c Chain) MultiversXToSolanaName() string {
	return multiversXToSolanaName
}

// BaseLogId returns the string using chain value and baseLogIdTemplate
func (c Chain) BaseLogId() string {
	return fmt.Sprintf(baseLogIdTemplate, c)
}

// MultiversXClientLogId returns the string using chain value and multiversXClientLogIdTemplate
func (c Chain) MultiversXClientLogId() string {
	return fmt.Sprintf(multiversXClientLogIdTemplate, c)
}

// MultiversXDataGetterLogId returns the string using chain value and multiversXDataGetterLogIdTemplate
func (c Chain) MultiversXDataGetterLogId() string {
	return fmt.Sprintf(multiversXDataGetterLogIdTemplate, c)
}

// SolanaClientLogId returns the string using chain value and solanaClientLogIdTemplate
func (c Chain) SolanaClientLogId() string {
	return fmt.Sprintf(solanaClientLogIdTemplate, c, c)
}

// MultiversXRoleProviderLogId returns the string using chain value and multiversXRoleProviderLogIdTemplate
func (c Chain) MultiversXRoleProviderLogId() string {
	return fmt.Sprintf(multiversXRoleProviderLogIdTemplate, c)
}

// SolanaRoleProviderLogId returns the string using chain value and solanaRoleProviderLogIdTemplate
func (c Chain) SolanaRoleProviderLogId() string {
	return fmt.Sprintf(solanaRoleProviderLogIdTemplate, c, c)
}

// BroadcasterLogId returns the string using chain value and broadcasterLogIdTemplate
func (c Chain) BroadcasterLogId() string {
	return fmt.Sprintf(broadcasterLogIdTemplate, c)
}
