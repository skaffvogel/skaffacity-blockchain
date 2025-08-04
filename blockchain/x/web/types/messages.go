package types

import (
	"context"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgUpdateWebConfig = "update_web_config"

var _ sdk.Msg = &MsgUpdateWebConfig{}

// MsgUpdateWebConfig defines a message to update web configuration
type MsgUpdateWebConfig struct {
	Creator string    `json:"creator"`
	Config  WebConfig `json:"config"`
}

// ProtoMessage implements the proto.Message interface for MsgUpdateWebConfig.
func (msg *MsgUpdateWebConfig) ProtoMessage() {}

// Reset implements the proto.Message interface for MsgUpdateWebConfig.
func (msg *MsgUpdateWebConfig) Reset() { *msg = MsgUpdateWebConfig{} }

// String implements the proto.Message interface for MsgUpdateWebConfig.
func (msg *MsgUpdateWebConfig) String() string {
	return 	fmt.Sprintf("MsgUpdateWebConfig{Creator: %s, Config: %s}",
		msg.Creator, msg.Config.String())
}

// MsgUpdateWebConfigResponse defines the response for MsgUpdateWebConfig
type MsgUpdateWebConfigResponse struct{}

// ProtoMessage implements proto.Message interface
func (m *MsgUpdateWebConfigResponse) ProtoMessage() {}

// Reset implements proto.Message interface
func (m *MsgUpdateWebConfigResponse) Reset() {
	*m = MsgUpdateWebConfigResponse{}
}

// String implements proto.Message interface
func (m *MsgUpdateWebConfigResponse) String() string {
	return "MsgUpdateWebConfigResponse{}"
}

// MsgServer interface
type MsgServer interface {
	UpdateWebConfig(context.Context, *MsgUpdateWebConfig) (*MsgUpdateWebConfigResponse, error)
}

func NewMsgUpdateWebConfig(creator string, config WebConfig) *MsgUpdateWebConfig {
	return &MsgUpdateWebConfig{
		Creator: creator,
		Config:  config,
	}
}

func (msg *MsgUpdateWebConfig) Route() string {
	return RouterKey
}

func (msg *MsgUpdateWebConfig) Type() string {
	return TypeMsgUpdateWebConfig
}

func (msg *MsgUpdateWebConfig) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateWebConfig) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateWebConfig) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	
	if msg.Config.Port == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "port cannot be zero")
	}
	
	if msg.Config.Host == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "host cannot be empty")
	}
	
	return nil
}
