// Code generated by https://github.com/gagliardetto/anchor-go. DO NOT EDIT.

package bridge

import (
	"errors"
	ag_binary "github.com/gagliardetto/binary"
	ag_solanago "github.com/gagliardetto/solana-go"
	ag_format "github.com/gagliardetto/solana-go/text/format"
	ag_treeout "github.com/gagliardetto/treeout"
)

// ExecuteTransfer is the `executeTransfer` instruction.
type ExecuteTransfer struct {
	BatchNonce      *ag_binary.Uint128
	IsLastTransfer  *bool
	TransferRequest *TransferRequest

	// [0] = [WRITE, SIGNER] relayer
	//
	// [1] = [WRITE] bridgeSettings
	//
	// [2] = [WRITE] relayerSettings
	//
	// [3] = [WRITE] safeSettings
	//
	// [4] = [WRITE] executedBatch
	//
	// [5] = [WRITE] executedTransfer
	//
	// [6] = [WRITE] whitelistedToken
	//
	// [7] = [WRITE] tokenMint
	//
	// [8] = [WRITE] safeSettingsAta
	//
	// [9] = [WRITE] sendTo
	//
	// [10] = [] recipient
	//
	// [11] = [] tokensSafeProgram
	//
	// [12] = [] tokenProgram
	//
	// [13] = [] associatedTokenProgram
	//
	// [14] = [] systemProgram
	ag_solanago.AccountMetaSlice `bin:"-"`
}

// NewExecuteTransferInstructionBuilder creates a new `ExecuteTransfer` instruction builder.
func NewExecuteTransferInstructionBuilder() *ExecuteTransfer {
	nd := &ExecuteTransfer{
		AccountMetaSlice: make(ag_solanago.AccountMetaSlice, 15),
	}
	return nd
}

// SetBatchNonce sets the "batchNonce" parameter.
func (inst *ExecuteTransfer) SetBatchNonce(batchNonce ag_binary.Uint128) *ExecuteTransfer {
	inst.BatchNonce = &batchNonce
	return inst
}

// SetIsLastTransfer sets the "isLastTransfer" parameter.
func (inst *ExecuteTransfer) SetIsLastTransfer(isLastTransfer bool) *ExecuteTransfer {
	inst.IsLastTransfer = &isLastTransfer
	return inst
}

// SetTransferRequest sets the "transferRequest" parameter.
func (inst *ExecuteTransfer) SetTransferRequest(transferRequest TransferRequest) *ExecuteTransfer {
	inst.TransferRequest = &transferRequest
	return inst
}

// SetRelayerAccount sets the "relayer" account.
func (inst *ExecuteTransfer) SetRelayerAccount(relayer ag_solanago.PublicKey) *ExecuteTransfer {
	inst.AccountMetaSlice[0] = ag_solanago.Meta(relayer).WRITE().SIGNER()
	return inst
}

// GetRelayerAccount gets the "relayer" account.
func (inst *ExecuteTransfer) GetRelayerAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(0)
}

// SetBridgeSettingsAccount sets the "bridgeSettings" account.
func (inst *ExecuteTransfer) SetBridgeSettingsAccount(bridgeSettings ag_solanago.PublicKey) *ExecuteTransfer {
	inst.AccountMetaSlice[1] = ag_solanago.Meta(bridgeSettings).WRITE()
	return inst
}

// GetBridgeSettingsAccount gets the "bridgeSettings" account.
func (inst *ExecuteTransfer) GetBridgeSettingsAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(1)
}

// SetRelayerSettingsAccount sets the "relayerSettings" account.
func (inst *ExecuteTransfer) SetRelayerSettingsAccount(relayerSettings ag_solanago.PublicKey) *ExecuteTransfer {
	inst.AccountMetaSlice[2] = ag_solanago.Meta(relayerSettings).WRITE()
	return inst
}

// GetRelayerSettingsAccount gets the "relayerSettings" account.
func (inst *ExecuteTransfer) GetRelayerSettingsAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(2)
}

// SetSafeSettingsAccount sets the "safeSettings" account.
func (inst *ExecuteTransfer) SetSafeSettingsAccount(safeSettings ag_solanago.PublicKey) *ExecuteTransfer {
	inst.AccountMetaSlice[3] = ag_solanago.Meta(safeSettings).WRITE()
	return inst
}

// GetSafeSettingsAccount gets the "safeSettings" account.
func (inst *ExecuteTransfer) GetSafeSettingsAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(3)
}

// SetExecutedBatchAccount sets the "executedBatch" account.
func (inst *ExecuteTransfer) SetExecutedBatchAccount(executedBatch ag_solanago.PublicKey) *ExecuteTransfer {
	inst.AccountMetaSlice[4] = ag_solanago.Meta(executedBatch).WRITE()
	return inst
}

// GetExecutedBatchAccount gets the "executedBatch" account.
func (inst *ExecuteTransfer) GetExecutedBatchAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(4)
}

// SetExecutedTransferAccount sets the "executedTransfer" account.
func (inst *ExecuteTransfer) SetExecutedTransferAccount(executedTransfer ag_solanago.PublicKey) *ExecuteTransfer {
	inst.AccountMetaSlice[5] = ag_solanago.Meta(executedTransfer).WRITE()
	return inst
}

// GetExecutedTransferAccount gets the "executedTransfer" account.
func (inst *ExecuteTransfer) GetExecutedTransferAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(5)
}

// SetWhitelistedTokenAccount sets the "whitelistedToken" account.
func (inst *ExecuteTransfer) SetWhitelistedTokenAccount(whitelistedToken ag_solanago.PublicKey) *ExecuteTransfer {
	inst.AccountMetaSlice[6] = ag_solanago.Meta(whitelistedToken).WRITE()
	return inst
}

// GetWhitelistedTokenAccount gets the "whitelistedToken" account.
func (inst *ExecuteTransfer) GetWhitelistedTokenAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(6)
}

// SetTokenMintAccount sets the "tokenMint" account.
func (inst *ExecuteTransfer) SetTokenMintAccount(tokenMint ag_solanago.PublicKey) *ExecuteTransfer {
	inst.AccountMetaSlice[7] = ag_solanago.Meta(tokenMint).WRITE()
	return inst
}

// GetTokenMintAccount gets the "tokenMint" account.
func (inst *ExecuteTransfer) GetTokenMintAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(7)
}

// SetSafeSettingsAtaAccount sets the "safeSettingsAta" account.
func (inst *ExecuteTransfer) SetSafeSettingsAtaAccount(safeSettingsAta ag_solanago.PublicKey) *ExecuteTransfer {
	inst.AccountMetaSlice[8] = ag_solanago.Meta(safeSettingsAta).WRITE()
	return inst
}

// GetSafeSettingsAtaAccount gets the "safeSettingsAta" account.
func (inst *ExecuteTransfer) GetSafeSettingsAtaAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(8)
}

// SetSendToAccount sets the "sendTo" account.
func (inst *ExecuteTransfer) SetSendToAccount(sendTo ag_solanago.PublicKey) *ExecuteTransfer {
	inst.AccountMetaSlice[9] = ag_solanago.Meta(sendTo).WRITE()
	return inst
}

// GetSendToAccount gets the "sendTo" account.
func (inst *ExecuteTransfer) GetSendToAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(9)
}

// SetRecipientAccount sets the "recipient" account.
func (inst *ExecuteTransfer) SetRecipientAccount(recipient ag_solanago.PublicKey) *ExecuteTransfer {
	inst.AccountMetaSlice[10] = ag_solanago.Meta(recipient)
	return inst
}

// GetRecipientAccount gets the "recipient" account.
func (inst *ExecuteTransfer) GetRecipientAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(10)
}

// SetTokensSafeProgramAccount sets the "tokensSafeProgram" account.
func (inst *ExecuteTransfer) SetTokensSafeProgramAccount(tokensSafeProgram ag_solanago.PublicKey) *ExecuteTransfer {
	inst.AccountMetaSlice[11] = ag_solanago.Meta(tokensSafeProgram)
	return inst
}

// GetTokensSafeProgramAccount gets the "tokensSafeProgram" account.
func (inst *ExecuteTransfer) GetTokensSafeProgramAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(11)
}

// SetTokenProgramAccount sets the "tokenProgram" account.
func (inst *ExecuteTransfer) SetTokenProgramAccount(tokenProgram ag_solanago.PublicKey) *ExecuteTransfer {
	inst.AccountMetaSlice[12] = ag_solanago.Meta(tokenProgram)
	return inst
}

// GetTokenProgramAccount gets the "tokenProgram" account.
func (inst *ExecuteTransfer) GetTokenProgramAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(12)
}

// SetAssociatedTokenProgramAccount sets the "associatedTokenProgram" account.
func (inst *ExecuteTransfer) SetAssociatedTokenProgramAccount(associatedTokenProgram ag_solanago.PublicKey) *ExecuteTransfer {
	inst.AccountMetaSlice[13] = ag_solanago.Meta(associatedTokenProgram)
	return inst
}

// GetAssociatedTokenProgramAccount gets the "associatedTokenProgram" account.
func (inst *ExecuteTransfer) GetAssociatedTokenProgramAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(13)
}

// SetSystemProgramAccount sets the "systemProgram" account.
func (inst *ExecuteTransfer) SetSystemProgramAccount(systemProgram ag_solanago.PublicKey) *ExecuteTransfer {
	inst.AccountMetaSlice[14] = ag_solanago.Meta(systemProgram)
	return inst
}

// GetSystemProgramAccount gets the "systemProgram" account.
func (inst *ExecuteTransfer) GetSystemProgramAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(14)
}

func (inst ExecuteTransfer) Build() *Instruction {
	return &Instruction{BaseVariant: ag_binary.BaseVariant{
		Impl:   inst,
		TypeID: Instruction_ExecuteTransfer,
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst ExecuteTransfer) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *ExecuteTransfer) Validate() error {
	// Check whether all (required) parameters are set:
	{
		if inst.BatchNonce == nil {
			return errors.New("BatchNonce parameter is not set")
		}
		if inst.IsLastTransfer == nil {
			return errors.New("IsLastTransfer parameter is not set")
		}
		if inst.TransferRequest == nil {
			return errors.New("TransferRequest parameter is not set")
		}
	}

	// Check whether all (required) accounts are set:
	{
		if inst.AccountMetaSlice[0] == nil {
			return errors.New("accounts.Relayer is not set")
		}
		if inst.AccountMetaSlice[1] == nil {
			return errors.New("accounts.BridgeSettings is not set")
		}
		if inst.AccountMetaSlice[2] == nil {
			return errors.New("accounts.RelayerSettings is not set")
		}
		if inst.AccountMetaSlice[3] == nil {
			return errors.New("accounts.SafeSettings is not set")
		}
		if inst.AccountMetaSlice[4] == nil {
			return errors.New("accounts.ExecutedBatch is not set")
		}
		if inst.AccountMetaSlice[5] == nil {
			return errors.New("accounts.ExecutedTransfer is not set")
		}
		if inst.AccountMetaSlice[6] == nil {
			return errors.New("accounts.WhitelistedToken is not set")
		}
		if inst.AccountMetaSlice[7] == nil {
			return errors.New("accounts.TokenMint is not set")
		}
		if inst.AccountMetaSlice[8] == nil {
			return errors.New("accounts.SafeSettingsAta is not set")
		}
		if inst.AccountMetaSlice[9] == nil {
			return errors.New("accounts.SendTo is not set")
		}
		if inst.AccountMetaSlice[10] == nil {
			return errors.New("accounts.Recipient is not set")
		}
		if inst.AccountMetaSlice[11] == nil {
			return errors.New("accounts.TokensSafeProgram is not set")
		}
		if inst.AccountMetaSlice[12] == nil {
			return errors.New("accounts.TokenProgram is not set")
		}
		if inst.AccountMetaSlice[13] == nil {
			return errors.New("accounts.AssociatedTokenProgram is not set")
		}
		if inst.AccountMetaSlice[14] == nil {
			return errors.New("accounts.SystemProgram is not set")
		}
	}
	return nil
}

func (inst *ExecuteTransfer) EncodeToTree(parent ag_treeout.Branches) {
	parent.Child(ag_format.Program(ProgramName, ProgramID)).
		//
		ParentFunc(func(programBranch ag_treeout.Branches) {
			programBranch.Child(ag_format.Instruction("ExecuteTransfer")).
				//
				ParentFunc(func(instructionBranch ag_treeout.Branches) {

					// Parameters of the instruction:
					instructionBranch.Child("Params[len=3]").ParentFunc(func(paramsBranch ag_treeout.Branches) {
						paramsBranch.Child(ag_format.Param("     BatchNonce", *inst.BatchNonce))
						paramsBranch.Child(ag_format.Param(" IsLastTransfer", *inst.IsLastTransfer))
						paramsBranch.Child(ag_format.Param("TransferRequest", *inst.TransferRequest))
					})

					// Accounts of the instruction:
					instructionBranch.Child("Accounts[len=15]").ParentFunc(func(accountsBranch ag_treeout.Branches) {
						accountsBranch.Child(ag_format.Meta("               relayer", inst.AccountMetaSlice.Get(0)))
						accountsBranch.Child(ag_format.Meta("        bridgeSettings", inst.AccountMetaSlice.Get(1)))
						accountsBranch.Child(ag_format.Meta("       relayerSettings", inst.AccountMetaSlice.Get(2)))
						accountsBranch.Child(ag_format.Meta("          safeSettings", inst.AccountMetaSlice.Get(3)))
						accountsBranch.Child(ag_format.Meta("         executedBatch", inst.AccountMetaSlice.Get(4)))
						accountsBranch.Child(ag_format.Meta("      executedTransfer", inst.AccountMetaSlice.Get(5)))
						accountsBranch.Child(ag_format.Meta("      whitelistedToken", inst.AccountMetaSlice.Get(6)))
						accountsBranch.Child(ag_format.Meta("             tokenMint", inst.AccountMetaSlice.Get(7)))
						accountsBranch.Child(ag_format.Meta("       safeSettingsAta", inst.AccountMetaSlice.Get(8)))
						accountsBranch.Child(ag_format.Meta("                sendTo", inst.AccountMetaSlice.Get(9)))
						accountsBranch.Child(ag_format.Meta("             recipient", inst.AccountMetaSlice.Get(10)))
						accountsBranch.Child(ag_format.Meta("     tokensSafeProgram", inst.AccountMetaSlice.Get(11)))
						accountsBranch.Child(ag_format.Meta("          tokenProgram", inst.AccountMetaSlice.Get(12)))
						accountsBranch.Child(ag_format.Meta("associatedTokenProgram", inst.AccountMetaSlice.Get(13)))
						accountsBranch.Child(ag_format.Meta("         systemProgram", inst.AccountMetaSlice.Get(14)))
					})
				})
		})
}

func (obj ExecuteTransfer) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	// Serialize `BatchNonce` param:
	err = encoder.Encode(obj.BatchNonce)
	if err != nil {
		return err
	}
	// Serialize `IsLastTransfer` param:
	err = encoder.Encode(obj.IsLastTransfer)
	if err != nil {
		return err
	}
	// Serialize `TransferRequest` param:
	err = encoder.Encode(obj.TransferRequest)
	if err != nil {
		return err
	}
	return nil
}
func (obj *ExecuteTransfer) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	// Deserialize `BatchNonce`:
	err = decoder.Decode(&obj.BatchNonce)
	if err != nil {
		return err
	}
	// Deserialize `IsLastTransfer`:
	err = decoder.Decode(&obj.IsLastTransfer)
	if err != nil {
		return err
	}
	// Deserialize `TransferRequest`:
	err = decoder.Decode(&obj.TransferRequest)
	if err != nil {
		return err
	}
	return nil
}

// NewExecuteTransferInstruction declares a new ExecuteTransfer instruction with the provided parameters and accounts.
func NewExecuteTransferInstruction(
	// Parameters:
	batchNonce ag_binary.Uint128,
	isLastTransfer bool,
	transferRequest TransferRequest,
	// Accounts:
	relayer ag_solanago.PublicKey,
	bridgeSettings ag_solanago.PublicKey,
	relayerSettings ag_solanago.PublicKey,
	safeSettings ag_solanago.PublicKey,
	executedBatch ag_solanago.PublicKey,
	executedTransfer ag_solanago.PublicKey,
	whitelistedToken ag_solanago.PublicKey,
	tokenMint ag_solanago.PublicKey,
	safeSettingsAta ag_solanago.PublicKey,
	sendTo ag_solanago.PublicKey,
	recipient ag_solanago.PublicKey,
	tokensSafeProgram ag_solanago.PublicKey,
	tokenProgram ag_solanago.PublicKey,
	associatedTokenProgram ag_solanago.PublicKey,
	systemProgram ag_solanago.PublicKey) *ExecuteTransfer {
	return NewExecuteTransferInstructionBuilder().
		SetBatchNonce(batchNonce).
		SetIsLastTransfer(isLastTransfer).
		SetTransferRequest(transferRequest).
		SetRelayerAccount(relayer).
		SetBridgeSettingsAccount(bridgeSettings).
		SetRelayerSettingsAccount(relayerSettings).
		SetSafeSettingsAccount(safeSettings).
		SetExecutedBatchAccount(executedBatch).
		SetExecutedTransferAccount(executedTransfer).
		SetWhitelistedTokenAccount(whitelistedToken).
		SetTokenMintAccount(tokenMint).
		SetSafeSettingsAtaAccount(safeSettingsAta).
		SetSendToAccount(sendTo).
		SetRecipientAccount(recipient).
		SetTokensSafeProgramAccount(tokensSafeProgram).
		SetTokenProgramAccount(tokenProgram).
		SetAssociatedTokenProgramAccount(associatedTokenProgram).
		SetSystemProgramAccount(systemProgram)
}
