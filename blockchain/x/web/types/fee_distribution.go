package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"gopkg.in/yaml.v2"
)

// FeeDistribution defines the fee distribution configuration
type FeeDistribution struct {
	// DeveloperAddress is the address that receives the developer fee
	DeveloperAddress string `json:"developer_address" yaml:"developer_address"`
	
	// DeveloperFeePercentage is the percentage of fees going to developer (in basis points, 1000 = 10%)
	DeveloperFeePercentage uint64 `json:"developer_fee_percentage" yaml:"developer_fee_percentage"`
	
	// ValidatorFeePercentage is the percentage of fees going to validators (in basis points, 9000 = 90%)
	ValidatorFeePercentage uint64 `json:"validator_fee_percentage" yaml:"validator_fee_percentage"`
	
	// Enabled determines if fee distribution is active
	Enabled bool `json:"enabled" yaml:"enabled"`
}

// DefaultFeeDistribution returns the default fee distribution configuration
func DefaultFeeDistribution() FeeDistribution {
	return FeeDistribution{
		DeveloperAddress:       "", // Will be set in config
		DeveloperFeePercentage: 1000, // 10%
		ValidatorFeePercentage: 9000, // 90%
		Enabled:                true,
	}
}

// Validate performs basic validation of fee distribution configuration
func (fd FeeDistribution) Validate() error {
	if fd.Enabled {
		if fd.DeveloperAddress == "" {
			return sdkerrors.ErrInvalidAddress.Wrap("developer address cannot be empty when fee distribution is enabled")
		}
		
		// Validate address format
		if _, err := sdk.AccAddressFromBech32(fd.DeveloperAddress); err != nil {
			return sdkerrors.ErrInvalidAddress.Wrapf("invalid developer address: %s", err.Error())
		}
		
		// Validate percentages add up to 100% (10000 basis points)
		total := fd.DeveloperFeePercentage + fd.ValidatorFeePercentage
		if total != 10000 {
			return sdkerrors.ErrInvalidRequest.Wrapf("fee percentages must add up to 10000 (100%%), got %d", total)
		}
	}
	
	return nil
}

// String implements the Stringer interface
func (fd FeeDistribution) String() string {
	out, _ := yaml.Marshal(fd)
	return string(out)
}

// CalculateFees calculates the distribution of fees
func (fd FeeDistribution) CalculateFees(totalFees sdk.Coins) (developerFee, validatorFee sdk.Coins) {
	if !fd.Enabled || totalFees.IsZero() {
		return sdk.NewCoins(), totalFees
	}
	
	var developerCoins []sdk.Coin
	var validatorCoins []sdk.Coin
	
	for _, coin := range totalFees {
		// Calculate developer fee (10%)
		devAmount := coin.Amount.MulRaw(int64(fd.DeveloperFeePercentage)).QuoRaw(10000)
		
		// Calculate validator fee (90%)
		valAmount := coin.Amount.Sub(devAmount)
		
		if devAmount.IsPositive() {
			developerCoins = append(developerCoins, sdk.NewCoin(coin.Denom, devAmount))
		}
		
		if valAmount.IsPositive() {
			validatorCoins = append(validatorCoins, sdk.NewCoin(coin.Denom, valAmount))
		}
	}
	
	return sdk.NewCoins(developerCoins...), sdk.NewCoins(validatorCoins...)
}
