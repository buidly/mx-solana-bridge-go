// Code generated by https://github.com/gagliardetto/anchor-go. DO NOT EDIT.

package tokens_safe

import (
	"errors"
	ag_binary "github.com/gagliardetto/binary"
	ag_solanago "github.com/gagliardetto/solana-go"
	ag_format "github.com/gagliardetto/solana-go/text/format"
	ag_treeout "github.com/gagliardetto/treeout"
)

// RemoveTokenFromWhitelist is the `removeTokenFromWhitelist` instruction.
type RemoveTokenFromWhitelist struct {

	// [0] = [WRITE, SIGNER] adminAuthority
	//
	// [1] = [WRITE] safeSettings
	//
	// [2] = [WRITE] tokenMint
	//
	// [3] = [WRITE] whitelistedToken
	ag_solanago.AccountMetaSlice `bin:"-"`
}

// NewRemoveTokenFromWhitelistInstructionBuilder creates a new `RemoveTokenFromWhitelist` instruction builder.
func NewRemoveTokenFromWhitelistInstructionBuilder() *RemoveTokenFromWhitelist {
	nd := &RemoveTokenFromWhitelist{
		AccountMetaSlice: make(ag_solanago.AccountMetaSlice, 4),
	}
	return nd
}

// SetAdminAuthorityAccount sets the "adminAuthority" account.
func (inst *RemoveTokenFromWhitelist) SetAdminAuthorityAccount(adminAuthority ag_solanago.PublicKey) *RemoveTokenFromWhitelist {
	inst.AccountMetaSlice[0] = ag_solanago.Meta(adminAuthority).WRITE().SIGNER()
	return inst
}

// GetAdminAuthorityAccount gets the "adminAuthority" account.
func (inst *RemoveTokenFromWhitelist) GetAdminAuthorityAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(0)
}

// SetSafeSettingsAccount sets the "safeSettings" account.
func (inst *RemoveTokenFromWhitelist) SetSafeSettingsAccount(safeSettings ag_solanago.PublicKey) *RemoveTokenFromWhitelist {
	inst.AccountMetaSlice[1] = ag_solanago.Meta(safeSettings).WRITE()
	return inst
}

// GetSafeSettingsAccount gets the "safeSettings" account.
func (inst *RemoveTokenFromWhitelist) GetSafeSettingsAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(1)
}

// SetTokenMintAccount sets the "tokenMint" account.
func (inst *RemoveTokenFromWhitelist) SetTokenMintAccount(tokenMint ag_solanago.PublicKey) *RemoveTokenFromWhitelist {
	inst.AccountMetaSlice[2] = ag_solanago.Meta(tokenMint).WRITE()
	return inst
}

// GetTokenMintAccount gets the "tokenMint" account.
func (inst *RemoveTokenFromWhitelist) GetTokenMintAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(2)
}

// SetWhitelistedTokenAccount sets the "whitelistedToken" account.
func (inst *RemoveTokenFromWhitelist) SetWhitelistedTokenAccount(whitelistedToken ag_solanago.PublicKey) *RemoveTokenFromWhitelist {
	inst.AccountMetaSlice[3] = ag_solanago.Meta(whitelistedToken).WRITE()
	return inst
}

// GetWhitelistedTokenAccount gets the "whitelistedToken" account.
func (inst *RemoveTokenFromWhitelist) GetWhitelistedTokenAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(3)
}

func (inst RemoveTokenFromWhitelist) Build() *Instruction {
	return &Instruction{BaseVariant: ag_binary.BaseVariant{
		Impl:   inst,
		TypeID: Instruction_RemoveTokenFromWhitelist,
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst RemoveTokenFromWhitelist) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *RemoveTokenFromWhitelist) Validate() error {
	// Check whether all (required) accounts are set:
	{
		if inst.AccountMetaSlice[0] == nil {
			return errors.New("accounts.AdminAuthority is not set")
		}
		if inst.AccountMetaSlice[1] == nil {
			return errors.New("accounts.SafeSettings is not set")
		}
		if inst.AccountMetaSlice[2] == nil {
			return errors.New("accounts.TokenMint is not set")
		}
		if inst.AccountMetaSlice[3] == nil {
			return errors.New("accounts.WhitelistedToken is not set")
		}
	}
	return nil
}

func (inst *RemoveTokenFromWhitelist) EncodeToTree(parent ag_treeout.Branches) {
	parent.Child(ag_format.Program(ProgramName, ProgramID)).
		//
		ParentFunc(func(programBranch ag_treeout.Branches) {
			programBranch.Child(ag_format.Instruction("RemoveTokenFromWhitelist")).
				//
				ParentFunc(func(instructionBranch ag_treeout.Branches) {

					// Parameters of the instruction:
					instructionBranch.Child("Params[len=0]").ParentFunc(func(paramsBranch ag_treeout.Branches) {})

					// Accounts of the instruction:
					instructionBranch.Child("Accounts[len=4]").ParentFunc(func(accountsBranch ag_treeout.Branches) {
						accountsBranch.Child(ag_format.Meta("  adminAuthority", inst.AccountMetaSlice.Get(0)))
						accountsBranch.Child(ag_format.Meta("    safeSettings", inst.AccountMetaSlice.Get(1)))
						accountsBranch.Child(ag_format.Meta("       tokenMint", inst.AccountMetaSlice.Get(2)))
						accountsBranch.Child(ag_format.Meta("whitelistedToken", inst.AccountMetaSlice.Get(3)))
					})
				})
		})
}

func (obj RemoveTokenFromWhitelist) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	return nil
}
func (obj *RemoveTokenFromWhitelist) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	return nil
}

// NewRemoveTokenFromWhitelistInstruction declares a new RemoveTokenFromWhitelist instruction with the provided parameters and accounts.
func NewRemoveTokenFromWhitelistInstruction(
	// Accounts:
	adminAuthority ag_solanago.PublicKey,
	safeSettings ag_solanago.PublicKey,
	tokenMint ag_solanago.PublicKey,
	whitelistedToken ag_solanago.PublicKey) *RemoveTokenFromWhitelist {
	return NewRemoveTokenFromWhitelistInstructionBuilder().
		SetAdminAuthorityAccount(adminAuthority).
		SetSafeSettingsAccount(safeSettings).
		SetTokenMintAccount(tokenMint).
		SetWhitelistedTokenAccount(whitelistedToken)
}