// Code generated by https://github.com/gagliardetto/anchor-go. DO NOT EDIT.

package tokens_safe

import (
	"errors"
	ag_binary "github.com/gagliardetto/binary"
	ag_solanago "github.com/gagliardetto/solana-go"
	ag_format "github.com/gagliardetto/solana-go/text/format"
	ag_treeout "github.com/gagliardetto/treeout"
)

// InitializeSettings is the `initializeSettings` instruction.
type InitializeSettings struct {
	AdminAuthority  *ag_solanago.PublicKey
	BridgeAuthority *ag_solanago.PublicKey

	// [0] = [WRITE, SIGNER] adminAuthority
	//
	// [1] = [WRITE] safeSettings
	//
	// [2] = [WRITE] firstBatch
	//
	// [3] = [] systemProgram
	ag_solanago.AccountMetaSlice `bin:"-"`
}

// NewInitializeSettingsInstructionBuilder creates a new `InitializeSettings` instruction builder.
func NewInitializeSettingsInstructionBuilder() *InitializeSettings {
	nd := &InitializeSettings{
		AccountMetaSlice: make(ag_solanago.AccountMetaSlice, 4),
	}
	return nd
}

// SetAdminAuthority sets the "adminAuthority" parameter.
func (inst *InitializeSettings) SetAdminAuthority(adminAuthority ag_solanago.PublicKey) *InitializeSettings {
	inst.AdminAuthority = &adminAuthority
	return inst
}

// SetBridgeAuthority sets the "bridgeAuthority" parameter.
func (inst *InitializeSettings) SetBridgeAuthority(bridgeAuthority ag_solanago.PublicKey) *InitializeSettings {
	inst.BridgeAuthority = &bridgeAuthority
	return inst
}

// SetAdminAuthorityAccount sets the "adminAuthority" account.
func (inst *InitializeSettings) SetAdminAuthorityAccount(adminAuthority ag_solanago.PublicKey) *InitializeSettings {
	inst.AccountMetaSlice[0] = ag_solanago.Meta(adminAuthority).WRITE().SIGNER()
	return inst
}

// GetAdminAuthorityAccount gets the "adminAuthority" account.
func (inst *InitializeSettings) GetAdminAuthorityAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(0)
}

// SetSafeSettingsAccount sets the "safeSettings" account.
func (inst *InitializeSettings) SetSafeSettingsAccount(safeSettings ag_solanago.PublicKey) *InitializeSettings {
	inst.AccountMetaSlice[1] = ag_solanago.Meta(safeSettings).WRITE()
	return inst
}

// GetSafeSettingsAccount gets the "safeSettings" account.
func (inst *InitializeSettings) GetSafeSettingsAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(1)
}

// SetFirstBatchAccount sets the "firstBatch" account.
func (inst *InitializeSettings) SetFirstBatchAccount(firstBatch ag_solanago.PublicKey) *InitializeSettings {
	inst.AccountMetaSlice[2] = ag_solanago.Meta(firstBatch).WRITE()
	return inst
}

// GetFirstBatchAccount gets the "firstBatch" account.
func (inst *InitializeSettings) GetFirstBatchAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(2)
}

// SetSystemProgramAccount sets the "systemProgram" account.
func (inst *InitializeSettings) SetSystemProgramAccount(systemProgram ag_solanago.PublicKey) *InitializeSettings {
	inst.AccountMetaSlice[3] = ag_solanago.Meta(systemProgram)
	return inst
}

// GetSystemProgramAccount gets the "systemProgram" account.
func (inst *InitializeSettings) GetSystemProgramAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(3)
}

func (inst InitializeSettings) Build() *Instruction {
	return &Instruction{BaseVariant: ag_binary.BaseVariant{
		Impl:   inst,
		TypeID: Instruction_InitializeSettings,
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst InitializeSettings) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *InitializeSettings) Validate() error {
	// Check whether all (required) parameters are set:
	{
		if inst.AdminAuthority == nil {
			return errors.New("AdminAuthority parameter is not set")
		}
		if inst.BridgeAuthority == nil {
			return errors.New("BridgeAuthority parameter is not set")
		}
	}

	// Check whether all (required) accounts are set:
	{
		if inst.AccountMetaSlice[0] == nil {
			return errors.New("accounts.AdminAuthority is not set")
		}
		if inst.AccountMetaSlice[1] == nil {
			return errors.New("accounts.SafeSettings is not set")
		}
		if inst.AccountMetaSlice[2] == nil {
			return errors.New("accounts.FirstBatch is not set")
		}
		if inst.AccountMetaSlice[3] == nil {
			return errors.New("accounts.SystemProgram is not set")
		}
	}
	return nil
}

func (inst *InitializeSettings) EncodeToTree(parent ag_treeout.Branches) {
	parent.Child(ag_format.Program(ProgramName, ProgramID)).
		//
		ParentFunc(func(programBranch ag_treeout.Branches) {
			programBranch.Child(ag_format.Instruction("InitializeSettings")).
				//
				ParentFunc(func(instructionBranch ag_treeout.Branches) {

					// Parameters of the instruction:
					instructionBranch.Child("Params[len=2]").ParentFunc(func(paramsBranch ag_treeout.Branches) {
						paramsBranch.Child(ag_format.Param(" AdminAuthority", *inst.AdminAuthority))
						paramsBranch.Child(ag_format.Param("BridgeAuthority", *inst.BridgeAuthority))
					})

					// Accounts of the instruction:
					instructionBranch.Child("Accounts[len=4]").ParentFunc(func(accountsBranch ag_treeout.Branches) {
						accountsBranch.Child(ag_format.Meta("adminAuthority", inst.AccountMetaSlice.Get(0)))
						accountsBranch.Child(ag_format.Meta("  safeSettings", inst.AccountMetaSlice.Get(1)))
						accountsBranch.Child(ag_format.Meta("    firstBatch", inst.AccountMetaSlice.Get(2)))
						accountsBranch.Child(ag_format.Meta(" systemProgram", inst.AccountMetaSlice.Get(3)))
					})
				})
		})
}

func (obj InitializeSettings) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	// Serialize `AdminAuthority` param:
	err = encoder.Encode(obj.AdminAuthority)
	if err != nil {
		return err
	}
	// Serialize `BridgeAuthority` param:
	err = encoder.Encode(obj.BridgeAuthority)
	if err != nil {
		return err
	}
	return nil
}
func (obj *InitializeSettings) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	// Deserialize `AdminAuthority`:
	err = decoder.Decode(&obj.AdminAuthority)
	if err != nil {
		return err
	}
	// Deserialize `BridgeAuthority`:
	err = decoder.Decode(&obj.BridgeAuthority)
	if err != nil {
		return err
	}
	return nil
}

// NewInitializeSettingsInstruction declares a new InitializeSettings instruction with the provided parameters and accounts.
func NewInitializeSettingsInstruction(
	// Parameters:
	adminAuthority ag_solanago.PublicKey,
	bridgeAuthority ag_solanago.PublicKey,
	// Accounts:
	adminAuthorityAccount ag_solanago.PublicKey,
	safeSettings ag_solanago.PublicKey,
	firstBatch ag_solanago.PublicKey,
	systemProgram ag_solanago.PublicKey) *InitializeSettings {
	return NewInitializeSettingsInstructionBuilder().
		SetAdminAuthority(adminAuthority).
		SetBridgeAuthority(bridgeAuthority).
		SetAdminAuthorityAccount(adminAuthorityAccount).
		SetSafeSettingsAccount(safeSettings).
		SetFirstBatchAccount(firstBatch).
		SetSystemProgramAccount(systemProgram)
}