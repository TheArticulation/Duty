package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgCreateDuty = "create_duty"
const TypeMsgUpdateDuty = "update_duty"
const TypeMsgDeleteDuty = "delete_duty"

var _ sdk.Msg = &MsgCreateDuty{}

func NewMsgCreateDuty(creator string, title string, description string, duration time.Duration) *MsgCreateDuty {
	return &MsgCreateDuty{
		Creator:     creator,
		Title:       title,
		Description: description,
		Duration:    duration,
	}
}

func (msg *MsgCreateDuty) Route() string {
	return RouterKey
}

func (msg *MsgCreateDuty) Type() string {
	return TypeMsgCreateDuty
}

func (msg *MsgCreateDuty) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreateDuty) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateDuty) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if msg.Title == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "title cannot be empty")
	}
	if msg.Duration <= 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "duration must be positive")
	}
	return nil
}

var _ sdk.Msg = &MsgUpdateDuty{}

func NewMsgUpdateDuty(creator string, id string, title string, description string) *MsgUpdateDuty {
	return &MsgUpdateDuty{
		Id:          id,
		Creator:     creator,
		Title:       title,
		Description: description,
	}
}

func (msg *MsgUpdateDuty) Route() string {
	return RouterKey
}

func (msg *MsgUpdateDuty) Type() string {
	return TypeMsgUpdateDuty
}

func (msg *MsgUpdateDuty) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateDuty) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateDuty) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if msg.Id == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "id cannot be empty")
	}
	if msg.Title == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "title cannot be empty")
	}
	return nil
}

var _ sdk.Msg = &MsgDeleteDuty{}

func NewMsgDeleteDuty(creator string, id string) *MsgDeleteDuty {
	return &MsgDeleteDuty{
		Id:      id,
		Creator: creator,
	}
}

func (msg *MsgDeleteDuty) Route() string {
	return RouterKey
}

func (msg *MsgDeleteDuty) Type() string {
	return TypeMsgDeleteDuty
}

func (msg *MsgDeleteDuty) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgDeleteDuty) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDeleteDuty) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if msg.Id == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "id cannot be empty")
	}
	return nil
}
