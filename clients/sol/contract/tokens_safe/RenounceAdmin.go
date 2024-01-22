// Code generated by https://github.com/gagliardetto/anchor-go. DO NOT EDIT.

package tokens_safe

import (
	"errors"
	ag_binary "github.com/gagliardetto/binary"
	ag_solanago "github.com/gagliardetto/solana-go"
	ag_format "github.com/gagliardetto/solana-go/text/format"
	ag_treeout "github.com/gagliardetto/treeout"
)

// RenounceAdmin is the `renounceAdmin` instruction.
type RenounceAdmin struct {

	// [0] = [WRITE, SIGNER] adminAuthority
	//
	// [1] = [WRITE] safeSettings
	ag_solanago.AccountMetaSlice `bin:"-"`
}

// NewRenounceAdminInstructionBuilder creates a new `RenounceAdmin` instruction builder.
func NewRenounceAdminInstructionBuilder() *RenounceAdmin {
	nd := &RenounceAdmin{
		AccountMetaSlice: make(ag_solanago.AccountMetaSlice, 2),
	}
	return nd
}

// SetAdminAuthorityAccount sets the "adminAuthority" account.
func (inst *RenounceAdmin) SetAdminAuthorityAccount(adminAuthority ag_solanago.PublicKey) *RenounceAdmin {
	inst.AccountMetaSlice[0] = ag_solanago.Meta(adminAuthority).WRITE().SIGNER()
	return inst
}

// GetAdminAuthorityAccount gets the "adminAuthority" account.
func (inst *RenounceAdmin) GetAdminAuthorityAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(0)
}

// SetSafeSettingsAccount sets the "safeSettings" account.
func (inst *RenounceAdmin) SetSafeSettingsAccount(safeSettings ag_solanago.PublicKey) *RenounceAdmin {
	inst.AccountMetaSlice[1] = ag_solanago.Meta(safeSettings).WRITE()
	return inst
}

// GetSafeSettingsAccount gets the "safeSettings" account.
func (inst *RenounceAdmin) GetSafeSettingsAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(1)
}

func (inst RenounceAdmin) Build() *Instruction {
	return &Instruction{BaseVariant: ag_binary.BaseVariant{
		Impl:   inst,
		TypeID: Instruction_RenounceAdmin,
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst RenounceAdmin) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *RenounceAdmin) Validate() error {
	// Check whether all (required) accounts are set:
	{
		if inst.AccountMetaSlice[0] == nil {
			return errors.New("accounts.AdminAuthority is not set")
		}
		if inst.AccountMetaSlice[1] == nil {
			return errors.New("accounts.SafeSettings is not set")
		}
	}
	return nil
}

func (inst *RenounceAdmin) EncodeToTree(parent ag_treeout.Branches) {
	parent.Child(ag_format.Program(ProgramName, ProgramID)).
		//
		ParentFunc(func(programBranch ag_treeout.Branches) {
			programBranch.Child(ag_format.Instruction("RenounceAdmin")).
				//
				ParentFunc(func(instructionBranch ag_treeout.Branches) {

					// Parameters of the instruction:
					instructionBranch.Child("Params[len=0]").ParentFunc(func(paramsBranch ag_treeout.Branches) {})

					// Accounts of the instruction:
					instructionBranch.Child("Accounts[len=2]").ParentFunc(func(accountsBranch ag_treeout.Branches) {
						accountsBranch.Child(ag_format.Meta("adminAuthority", inst.AccountMetaSlice.Get(0)))
						accountsBranch.Child(ag_format.Meta("  safeSettings", inst.AccountMetaSlice.Get(1)))
					})
				})
		})
}

func (obj RenounceAdmin) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	return nil
}
func (obj *RenounceAdmin) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	return nil
}

// NewRenounceAdminInstruction declares a new RenounceAdmin instruction with the provided parameters and accounts.
func NewRenounceAdminInstruction(
	// Accounts:
	adminAuthority ag_solanago.PublicKey,
	safeSettings ag_solanago.PublicKey) *RenounceAdmin {
	return NewRenounceAdminInstructionBuilder().
		SetAdminAuthorityAccount(adminAuthority).
		SetSafeSettingsAccount(safeSettings)
}
