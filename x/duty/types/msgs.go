package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type DutyMetadata struct {
	// ECDSA secp256k1 public key used to sign Hyperlane checkpoints (hex or 0x…)
	CheckpointPubKey string `json:"checkpoint_pub_key"`
	// Public location for signatures (e.g., s3://bucket/prefix or https://…)
	CheckpointStorageURI string `json:"checkpoint_storage_uri"`
}

const (
	TypeMsgSetDutyMetadata = "set_duty_metadata"
)

type MsgSetDutyMetadata struct {
	// signer is the consensus validator operator address (valoper…)
	Signer   string       `json:"signer"`
	Metadata DutyMetadata `json:"metadata"`
}

func (m MsgSetDutyMetadata) Route() string { return RouterKey }
func (m MsgSetDutyMetadata) Type() string  { return TypeMsgSetDutyMetadata }
func (m MsgSetDutyMetadata) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.ValAddressFromBech32(m.Signer)
	return []sdk.AccAddress{sdk.AccAddress(addr.Bytes())}
}
func (m MsgSetDutyMetadata) ValidateBasic() error {
	if _, err := sdk.ValAddressFromBech32(m.Signer); err != nil {
		return sdkerrors.Wrap(err, "invalid valoper")
	}
	if len(m.Metadata.CheckpointPubKey) == 0 || len(m.Metadata.CheckpointStorageURI) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "missing metadata")
	}
	return nil
}
