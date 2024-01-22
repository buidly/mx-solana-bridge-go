// Code generated by https://github.com/gagliardetto/anchor-go. DO NOT EDIT.

package tokens_safe

import (
	"errors"
	ag_binary "github.com/gagliardetto/binary"
	ag_solanago "github.com/gagliardetto/solana-go"
	ag_format "github.com/gagliardetto/solana-go/text/format"
	ag_treeout "github.com/gagliardetto/treeout"
)

// InitSupply is the `initSupply` instruction.
type InitSupply struct {
	Amount *uint64

	// [0] = [WRITE, SIGNER] adminAuthority
	//
	// [1] = [WRITE] safeSettings
	//
	// [2] = [WRITE] tokenMint
	//
	// [3] = [WRITE] whitelistedToken
	//
	// [4] = [WRITE] sendFrom
	//
	// [5] = [WRITE] sendTo
	//
	// [6] = [] tokenProgram
	//
	// [7] = [] associatedTokenProgram
	//
	// [8] = [] systemProgram
	ag_solanago.AccountMetaSlice `bin:"-"`
}

// NewInitSupplyInstructionBuilder creates a new `InitSupply` instruction builder.
func NewInitSupplyInstructionBuilder() *InitSupply {
	nd := &InitSupply{
		AccountMetaSlice: make(ag_solanago.AccountMetaSlice, 9),
	}
	return nd
}

// SetAmount sets the "amount" parameter.
func (inst *InitSupply) SetAmount(amount uint64) *InitSupply {
	inst.Amount = &amount
	return inst
}

// SetAdminAuthorityAccount sets the "adminAuthority" account.
func (inst *InitSupply) SetAdminAuthorityAccount(adminAuthority ag_solanago.PublicKey) *InitSupply {
	inst.AccountMetaSlice[0] = ag_solanago.Meta(adminAuthority).WRITE().SIGNER()
	return inst
}

// GetAdminAuthorityAccount gets the "adminAuthority" account.
func (inst *InitSupply) GetAdminAuthorityAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(0)
}

// SetSafeSettingsAccount sets the "safeSettings" account.
func (inst *InitSupply) SetSafeSettingsAccount(safeSettings ag_solanago.PublicKey) *InitSupply {
	inst.AccountMetaSlice[1] = ag_solanago.Meta(safeSettings).WRITE()
	return inst
}

// GetSafeSettingsAccount gets the "safeSettings" account.
func (inst *InitSupply) GetSafeSettingsAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(1)
}

// SetTokenMintAccount sets the "tokenMint" account.
func (inst *InitSupply) SetTokenMintAccount(tokenMint ag_solanago.PublicKey) *InitSupply {
	inst.AccountMetaSlice[2] = ag_solanago.Meta(tokenMint).WRITE()
	return inst
}

// GetTokenMintAccount gets the "tokenMint" account.
func (inst *InitSupply) GetTokenMintAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(2)
}

// SetWhitelistedTokenAccount sets the "whitelistedToken" account.
func (inst *InitSupply) SetWhitelistedTokenAccount(whitelistedToken ag_solanago.PublicKey) *InitSupply {
	inst.AccountMetaSlice[3] = ag_solanago.Meta(whitelistedToken).WRITE()
	return inst
}

// GetWhitelistedTokenAccount gets the "whitelistedToken" account.
func (inst *InitSupply) GetWhitelistedTokenAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(3)
}

// SetSendFromAccount sets the "sendFrom" account.
func (inst *InitSupply) SetSendFromAccount(sendFrom ag_solanago.PublicKey) *InitSupply {
	inst.AccountMetaSlice[4] = ag_solanago.Meta(sendFrom).WRITE()
	return inst
}

// GetSendFromAccount gets the "sendFrom" account.
func (inst *InitSupply) GetSendFromAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(4)
}

// SetSendToAccount sets the "sendTo" account.
func (inst *InitSupply) SetSendToAccount(sendTo ag_solanago.PublicKey) *InitSupply {
	inst.AccountMetaSlice[5] = ag_solanago.Meta(sendTo).WRITE()
	return inst
}

// GetSendToAccount gets the "sendTo" account.
func (inst *InitSupply) GetSendToAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(5)
}

// SetTokenProgramAccount sets the "tokenProgram" account.
func (inst *InitSupply) SetTokenProgramAccount(tokenProgram ag_solanago.PublicKey) *InitSupply {
	inst.AccountMetaSlice[6] = ag_solanago.Meta(tokenProgram)
	return inst
}

// GetTokenProgramAccount gets the "tokenProgram" account.
func (inst *InitSupply) GetTokenProgramAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(6)
}

// SetAssociatedTokenProgramAccount sets the "associatedTokenProgram" account.
func (inst *InitSupply) SetAssociatedTokenProgramAccount(associatedTokenProgram ag_solanago.PublicKey) *InitSupply {
	inst.AccountMetaSlice[7] = ag_solanago.Meta(associatedTokenProgram)
	return inst
}

// GetAssociatedTokenProgramAccount gets the "associatedTokenProgram" account.
func (inst *InitSupply) GetAssociatedTokenProgramAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(7)
}

// SetSystemProgramAccount sets the "systemProgram" account.
func (inst *InitSupply) SetSystemProgramAccount(systemProgram ag_solanago.PublicKey) *InitSupply {
	inst.AccountMetaSlice[8] = ag_solanago.Meta(systemProgram)
	return inst
}

// GetSystemProgramAccount gets the "systemProgram" account.
func (inst *InitSupply) GetSystemProgramAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(8)
}

func (inst InitSupply) Build() *Instruction {
	return &Instruction{BaseVariant: ag_binary.BaseVariant{
		Impl:   inst,
		TypeID: Instruction_InitSupply,
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst InitSupply) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *InitSupply) Validate() error {
	// Check whether all (required) parameters are set:
	{
		if inst.Amount == nil {
			return errors.New("Amount parameter is not set")
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
			return errors.New("accounts.TokenMint is not set")
		}
		if inst.AccountMetaSlice[3] == nil {
			return errors.New("accounts.WhitelistedToken is not set")
		}
		if inst.AccountMetaSlice[4] == nil {
			return errors.New("accounts.SendFrom is not set")
		}
		if inst.AccountMetaSlice[5] == nil {
			return errors.New("accounts.SendTo is not set")
		}
		if inst.AccountMetaSlice[6] == nil {
			return errors.New("accounts.TokenProgram is not set")
		}
		if inst.AccountMetaSlice[7] == nil {
			return errors.New("accounts.AssociatedTokenProgram is not set")
		}
		if inst.AccountMetaSlice[8] == nil {
			return errors.New("accounts.SystemProgram is not set")
		}
	}
	return nil
}

func (inst *InitSupply) EncodeToTree(parent ag_treeout.Branches) {
	parent.Child(ag_format.Program(ProgramName, ProgramID)).
		//
		ParentFunc(func(programBranch ag_treeout.Branches) {
			programBranch.Child(ag_format.Instruction("InitSupply")).
				//
				ParentFunc(func(instructionBranch ag_treeout.Branches) {

					// Parameters of the instruction:
					instructionBranch.Child("Params[len=1]").ParentFunc(func(paramsBranch ag_treeout.Branches) {
						paramsBranch.Child(ag_format.Param("Amount", *inst.Amount))
					})

					// Accounts of the instruction:
					instructionBranch.Child("Accounts[len=9]").ParentFunc(func(accountsBranch ag_treeout.Branches) {
						accountsBranch.Child(ag_format.Meta("        adminAuthority", inst.AccountMetaSlice.Get(0)))
						accountsBranch.Child(ag_format.Meta("          safeSettings", inst.AccountMetaSlice.Get(1)))
						accountsBranch.Child(ag_format.Meta("             tokenMint", inst.AccountMetaSlice.Get(2)))
						accountsBranch.Child(ag_format.Meta("      whitelistedToken", inst.AccountMetaSlice.Get(3)))
						accountsBranch.Child(ag_format.Meta("              sendFrom", inst.AccountMetaSlice.Get(4)))
						accountsBranch.Child(ag_format.Meta("                sendTo", inst.AccountMetaSlice.Get(5)))
						accountsBranch.Child(ag_format.Meta("          tokenProgram", inst.AccountMetaSlice.Get(6)))
						accountsBranch.Child(ag_format.Meta("associatedTokenProgram", inst.AccountMetaSlice.Get(7)))
						accountsBranch.Child(ag_format.Meta("         systemProgram", inst.AccountMetaSlice.Get(8)))
					})
				})
		})
}

func (obj InitSupply) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	// Serialize `Amount` param:
	err = encoder.Encode(obj.Amount)
	if err != nil {
		return err
	}
	return nil
}
func (obj *InitSupply) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	// Deserialize `Amount`:
	err = decoder.Decode(&obj.Amount)
	if err != nil {
		return err
	}
	return nil
}

// NewInitSupplyInstruction declares a new InitSupply instruction with the provided parameters and accounts.
func NewInitSupplyInstruction(
	// Parameters:
	amount uint64,
	// Accounts:
	adminAuthority ag_solanago.PublicKey,
	safeSettings ag_solanago.PublicKey,
	tokenMint ag_solanago.PublicKey,
	whitelistedToken ag_solanago.PublicKey,
	sendFrom ag_solanago.PublicKey,
	sendTo ag_solanago.PublicKey,
	tokenProgram ag_solanago.PublicKey,
	associatedTokenProgram ag_solanago.PublicKey,
	systemProgram ag_solanago.PublicKey) *InitSupply {
	return NewInitSupplyInstructionBuilder().
		SetAmount(amount).
		SetAdminAuthorityAccount(adminAuthority).
		SetSafeSettingsAccount(safeSettings).
		SetTokenMintAccount(tokenMint).
		SetWhitelistedTokenAccount(whitelistedToken).
		SetSendFromAccount(sendFrom).
		SetSendToAccount(sendTo).
		SetTokenProgramAccount(tokenProgram).
		SetAssociatedTokenProgramAccount(associatedTokenProgram).
		SetSystemProgramAccount(systemProgram)
}
