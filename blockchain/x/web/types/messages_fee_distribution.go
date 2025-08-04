package types

import (
	"encoding/json"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// MsgSetDeveloperAddress sets the developer address for fee distribution
type MsgSetDeveloperAddress struct {
	Creator        string `json:"creator"`
	DeveloperAddress string `json:"developer_address"`
}

// NewMsgSetDeveloperAddress creates a new MsgSetDeveloperAddress
func NewMsgSetDeveloperAddress(creator, developerAddress string) *MsgSetDeveloperAddress {
	return &MsgSetDeveloperAddress{
		Creator:          creator,
		DeveloperAddress: developerAddress,
	}
}

// Route returns the route of MsgSetDeveloperAddress
func (msg *MsgSetDeveloperAddress) Route() string {
	return RouterKey
}

// Type returns the type of MsgSetDeveloperAddress
func (msg *MsgSetDeveloperAddress) Type() string {
	return "set_developer_address"
}

// GetSigners returns the signers of MsgSetDeveloperAddress
func (msg *MsgSetDeveloperAddress) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// GetSignBytes returns the sign bytes of MsgSetDeveloperAddress
func (msg *MsgSetDeveloperAddress) GetSignBytes() []byte {
	bz, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(bz)
}

// ValidateBasic validates the basic fields of MsgSetDeveloperAddress
func (msg *MsgSetDeveloperAddress) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	
	_, err = sdk.AccAddressFromBech32(msg.DeveloperAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid developer address (%s)", err)
	}
	
	return nil
}

// MsgEnableFeeDistribution enables or disables fee distribution
type MsgEnableFeeDistribution struct {
	Creator string `json:"creator"`
	Enabled bool   `json:"enabled"`
}

// NewMsgEnableFeeDistribution creates a new MsgEnableFeeDistribution
func NewMsgEnableFeeDistribution(creator string, enabled bool) *MsgEnableFeeDistribution {
	return &MsgEnableFeeDistribution{
		Creator: creator,
		Enabled: enabled,
	}
}

// Route returns the route of MsgEnableFeeDistribution
func (msg *MsgEnableFeeDistribution) Route() string {
	return RouterKey
}

// Type returns the type of MsgEnableFeeDistribution
func (msg *MsgEnableFeeDistribution) Type() string {
	return "enable_fee_distribution"
}

// GetSigners returns the signers of MsgEnableFeeDistribution
func (msg *MsgEnableFeeDistribution) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// GetSignBytes returns the sign bytes of MsgEnableFeeDistribution
func (msg *MsgEnableFeeDistribution) GetSignBytes() []byte {
	bz, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(bz)
}

// ValidateBasic validates the basic fields of MsgEnableFeeDistribution
func (msg *MsgEnableFeeDistribution) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	
	return nil
}

// Response types

// MsgSetDeveloperAddressResponse is the response for MsgSetDeveloperAddress
type MsgSetDeveloperAddressResponse struct{}

// MsgEnableFeeDistributionResponse is the response for MsgEnableFeeDistribution
type MsgEnableFeeDistributionResponse struct{}
