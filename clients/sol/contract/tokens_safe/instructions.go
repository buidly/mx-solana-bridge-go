// Code generated by https://github.com/gagliardetto/anchor-go. DO NOT EDIT.

package tokens_safe

import (
	"bytes"
	"fmt"
	ag_spew "github.com/davecgh/go-spew/spew"
	ag_binary "github.com/gagliardetto/binary"
	ag_solanago "github.com/gagliardetto/solana-go"
	ag_text "github.com/gagliardetto/solana-go/text"
	ag_treeout "github.com/gagliardetto/treeout"
)

var ProgramID ag_solanago.PublicKey

func SetProgramID(pubkey ag_solanago.PublicKey) {
	ProgramID = pubkey
	ag_solanago.RegisterInstructionDecoder(ProgramID, registryDecodeInstruction)
}

const ProgramName = "TokensSafe"

func init() {
	if !ProgramID.IsZero() {
		ag_solanago.RegisterInstructionDecoder(ProgramID, registryDecodeInstruction)
	}
}

var (
	Instruction_InitializeSettings = ag_binary.TypeID([8]byte{71, 239, 156, 98, 109, 81, 123, 78})

	Instruction_Pause = ag_binary.TypeID([8]byte{211, 22, 221, 251, 74, 121, 193, 47})

	Instruction_Unpause = ag_binary.TypeID([8]byte{169, 144, 4, 38, 10, 141, 188, 255})

	Instruction_TransferAdmin = ag_binary.TypeID([8]byte{42, 242, 66, 106, 228, 10, 111, 156})

	Instruction_RenounceAdmin = ag_binary.TypeID([8]byte{223, 213, 55, 194, 0, 108, 225, 137})

	Instruction_TransferBridge = ag_binary.TypeID([8]byte{23, 249, 126, 68, 26, 27, 39, 180})

	Instruction_WhitelistToken = ag_binary.TypeID([8]byte{6, 141, 83, 167, 31, 6, 2, 224})

	Instruction_RemoveTokenFromWhitelist = ag_binary.TypeID([8]byte{201, 48, 92, 116, 169, 35, 188, 193})

	Instruction_SetBatchSlotLimit = ag_binary.TypeID([8]byte{16, 130, 93, 4, 89, 56, 251, 187})

	Instruction_SetBatchSize = ag_binary.TypeID([8]byte{49, 113, 32, 31, 165, 75, 186, 177})

	Instruction_SetTokenMinLimit = ag_binary.TypeID([8]byte{144, 167, 123, 251, 97, 120, 83, 131})

	Instruction_SetTokenMaxLimit = ag_binary.TypeID([8]byte{42, 214, 71, 228, 81, 225, 175, 208})

	Instruction_Deposit = ag_binary.TypeID([8]byte{242, 35, 198, 137, 82, 225, 242, 182})

	Instruction_InitSupply = ag_binary.TypeID([8]byte{211, 0, 180, 127, 125, 56, 66, 4})

	Instruction_BridgeTransfer = ag_binary.TypeID([8]byte{78, 60, 231, 120, 171, 69, 0, 50})

	Instruction_RecoverLostFunds = ag_binary.TypeID([8]byte{162, 223, 39, 58, 212, 115, 100, 154})
)

// InstructionIDToName returns the name of the instruction given its ID.
func InstructionIDToName(id ag_binary.TypeID) string {
	switch id {
	case Instruction_InitializeSettings:
		return "InitializeSettings"
	case Instruction_Pause:
		return "Pause"
	case Instruction_Unpause:
		return "Unpause"
	case Instruction_TransferAdmin:
		return "TransferAdmin"
	case Instruction_RenounceAdmin:
		return "RenounceAdmin"
	case Instruction_TransferBridge:
		return "TransferBridge"
	case Instruction_WhitelistToken:
		return "WhitelistToken"
	case Instruction_RemoveTokenFromWhitelist:
		return "RemoveTokenFromWhitelist"
	case Instruction_SetBatchSlotLimit:
		return "SetBatchSlotLimit"
	case Instruction_SetBatchSize:
		return "SetBatchSize"
	case Instruction_SetTokenMinLimit:
		return "SetTokenMinLimit"
	case Instruction_SetTokenMaxLimit:
		return "SetTokenMaxLimit"
	case Instruction_Deposit:
		return "DepositInstr"
	case Instruction_InitSupply:
		return "InitSupply"
	case Instruction_BridgeTransfer:
		return "BridgeTransfer"
	case Instruction_RecoverLostFunds:
		return "RecoverLostFunds"
	default:
		return ""
	}
}

type Instruction struct {
	ag_binary.BaseVariant
}

func (inst *Instruction) EncodeToTree(parent ag_treeout.Branches) {
	if enToTree, ok := inst.Impl.(ag_text.EncodableToTree); ok {
		enToTree.EncodeToTree(parent)
	} else {
		parent.Child(ag_spew.Sdump(inst))
	}
}

var InstructionImplDef = ag_binary.NewVariantDefinition(
	ag_binary.AnchorTypeIDEncoding,
	[]ag_binary.VariantType{
		{
			"initialize_settings", (*InitializeSettings)(nil),
		},
		{
			"pause", (*Pause)(nil),
		},
		{
			"unpause", (*Unpause)(nil),
		},
		{
			"transfer_admin", (*TransferAdmin)(nil),
		},
		{
			"renounce_admin", (*RenounceAdmin)(nil),
		},
		{
			"transfer_bridge", (*TransferBridge)(nil),
		},
		{
			"whitelist_token", (*WhitelistToken)(nil),
		},
		{
			"remove_token_from_whitelist", (*RemoveTokenFromWhitelist)(nil),
		},
		{
			"set_batch_slot_limit", (*SetBatchSlotLimit)(nil),
		},
		{
			"set_batch_size", (*SetBatchSize)(nil),
		},
		{
			"set_token_min_limit", (*SetTokenMinLimit)(nil),
		},
		{
			"set_token_max_limit", (*SetTokenMaxLimit)(nil),
		},
		{
			"deposit", (*DepositInstr)(nil),
		},
		{
			"init_supply", (*InitSupply)(nil),
		},
		{
			"bridge_transfer", (*BridgeTransfer)(nil),
		},
		{
			"recover_lost_funds", (*RecoverLostFunds)(nil),
		},
	},
)

func (inst *Instruction) ProgramID() ag_solanago.PublicKey {
	return ProgramID
}

func (inst *Instruction) Accounts() (out []*ag_solanago.AccountMeta) {
	return inst.Impl.(ag_solanago.AccountsGettable).GetAccounts()
}

func (inst *Instruction) Data() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := ag_binary.NewBorshEncoder(buf).Encode(inst); err != nil {
		return nil, fmt.Errorf("unable to encode instruction: %w", err)
	}
	return buf.Bytes(), nil
}

func (inst *Instruction) TextEncode(encoder *ag_text.Encoder, option *ag_text.Option) error {
	return encoder.Encode(inst.Impl, option)
}

func (inst *Instruction) UnmarshalWithDecoder(decoder *ag_binary.Decoder) error {
	return inst.BaseVariant.UnmarshalBinaryVariant(decoder, InstructionImplDef)
}

func (inst *Instruction) MarshalWithEncoder(encoder *ag_binary.Encoder) error {
	err := encoder.WriteBytes(inst.TypeID.Bytes(), false)
	if err != nil {
		return fmt.Errorf("unable to write variant type: %w", err)
	}
	return encoder.Encode(inst.Impl)
}

func registryDecodeInstruction(accounts []*ag_solanago.AccountMeta, data []byte) (interface{}, error) {
	inst, err := DecodeInstruction(accounts, data)
	if err != nil {
		return nil, err
	}
	return inst, nil
}

func DecodeInstruction(accounts []*ag_solanago.AccountMeta, data []byte) (*Instruction, error) {
	inst := new(Instruction)
	if err := ag_binary.NewBorshDecoder(data).Decode(inst); err != nil {
		return nil, fmt.Errorf("unable to decode instruction: %w", err)
	}
	if v, ok := inst.Impl.(ag_solanago.AccountsSettable); ok {
		err := v.SetAccounts(accounts)
		if err != nil {
			return nil, fmt.Errorf("unable to set accounts for instruction: %w", err)
		}
	}
	return inst, nil
}
