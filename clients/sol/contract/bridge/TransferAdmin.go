// Code generated by https://github.com/gagliardetto/anchor-go. DO NOT EDIT.

package bridge

import (
	"errors"
	ag_binary "github.com/gagliardetto/binary"
	ag_solanago "github.com/gagliardetto/solana-go"
	ag_format "github.com/gagliardetto/solana-go/text/format"
	ag_treeout "github.com/gagliardetto/treeout"
)

// TransferAdmin is the `transferAdmin` instruction.
type TransferAdmin struct {
	AdminAuthority *ag_solanago.PublicKey

	// [0] = [WRITE, SIGNER] adminAuthority
	//
	// [1] = [WRITE] bridgeSettings
	ag_solanago.AccountMetaSlice `bin:"-"`
}

// NewTransferAdminInstructionBuilder creates a new `TransferAdmin` instruction builder.
func NewTransferAdminInstructionBuilder() *TransferAdmin {
	nd := &TransferAdmin{
		AccountMetaSlice: make(ag_solanago.AccountMetaSlice, 2),
	}
	return nd
}

// SetAdminAuthority sets the "adminAuthority" parameter.
func (inst *TransferAdmin) SetAdminAuthority(adminAuthority ag_solanago.PublicKey) *TransferAdmin {
	inst.AdminAuthority = &adminAuthority
	return inst
}

// SetAdminAuthorityAccount sets the "adminAuthority" account.
func (inst *TransferAdmin) SetAdminAuthorityAccount(adminAuthority ag_solanago.PublicKey) *TransferAdmin {
	inst.AccountMetaSlice[0] = ag_solanago.Meta(adminAuthority).WRITE().SIGNER()
	return inst
}

// GetAdminAuthorityAccount gets the "adminAuthority" account.
func (inst *TransferAdmin) GetAdminAuthorityAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(0)
}

// SetBridgeSettingsAccount sets the "bridgeSettings" account.
func (inst *TransferAdmin) SetBridgeSettingsAccount(bridgeSettings ag_solanago.PublicKey) *TransferAdmin {
	inst.AccountMetaSlice[1] = ag_solanago.Meta(bridgeSettings).WRITE()
	return inst
}

// GetBridgeSettingsAccount gets the "bridgeSettings" account.
func (inst *TransferAdmin) GetBridgeSettingsAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(1)
}

func (inst TransferAdmin) Build() *Instruction {
	return &Instruction{BaseVariant: ag_binary.BaseVariant{
		Impl:   inst,
		TypeID: Instruction_TransferAdmin,
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst TransferAdmin) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *TransferAdmin) Validate() error {
	// Check whether all (required) parameters are set:
	{
		if inst.AdminAuthority == nil {
			return errors.New("AdminAuthority parameter is not set")
		}
	}

	// Check whether all (required) accounts are set:
	{
		if inst.AccountMetaSlice[0] == nil {
			return errors.New("accounts.AdminAuthority is not set")
		}
		if inst.AccountMetaSlice[1] == nil {
			return errors.New("accounts.BridgeSettings is not set")
		}
	}
	return nil
}

func (inst *TransferAdmin) EncodeToTree(parent ag_treeout.Branches) {
	parent.Child(ag_format.Program(ProgramName, ProgramID)).
		//
		ParentFunc(func(programBranch ag_treeout.Branches) {
			programBranch.Child(ag_format.Instruction("TransferAdmin")).
				//
				ParentFunc(func(instructionBranch ag_treeout.Branches) {

					// Parameters of the instruction:
					instructionBranch.Child("Params[len=1]").ParentFunc(func(paramsBranch ag_treeout.Branches) {
						paramsBranch.Child(ag_format.Param("AdminAuthority", *inst.AdminAuthority))
					})

					// Accounts of the instruction:
					instructionBranch.Child("Accounts[len=2]").ParentFunc(func(accountsBranch ag_treeout.Branches) {
						accountsBranch.Child(ag_format.Meta("adminAuthority", inst.AccountMetaSlice.Get(0)))
						accountsBranch.Child(ag_format.Meta("bridgeSettings", inst.AccountMetaSlice.Get(1)))
					})
				})
		})
}

func (obj TransferAdmin) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	// Serialize `AdminAuthority` param:
	err = encoder.Encode(obj.AdminAuthority)
	if err != nil {
		return err
	}
	return nil
}
func (obj *TransferAdmin) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	// Deserialize `AdminAuthority`:
	err = decoder.Decode(&obj.AdminAuthority)
	if err != nil {
		return err
	}
	return nil
}

// NewTransferAdminInstruction declares a new TransferAdmin instruction with the provided parameters and accounts.
func NewTransferAdminInstruction(
	// Parameters:
	adminAuthority ag_solanago.PublicKey,
	// Accounts:
	adminAuthorityAccount ag_solanago.PublicKey,
	bridgeSettings ag_solanago.PublicKey) *TransferAdmin {
	return NewTransferAdminInstructionBuilder().
		SetAdminAuthority(adminAuthority).
		SetAdminAuthorityAccount(adminAuthorityAccount).
		SetBridgeSettingsAccount(bridgeSettings)
}
