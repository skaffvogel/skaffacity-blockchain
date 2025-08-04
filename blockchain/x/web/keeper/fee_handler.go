package keeper

import (
	"fmt"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
	"skaffacity/x/web/types"
)

// FeeHandler handles the distribution of transaction fees
type FeeHandler struct {
	bankKeeper  types.BankKeeper
	authKeeper  types.AccountKeeper
}

// NewFeeHandler creates a new fee handler
func NewFeeHandler(bankKeeper types.BankKeeper, authKeeper types.AccountKeeper) FeeHandler {
	return FeeHandler{
		bankKeeper: bankKeeper,
		authKeeper: authKeeper,
	}
}

// DistributeFees distributes transaction fees according to the configuration
func (fh FeeHandler) DistributeFees(ctx sdk.Context, webKeeper Keeper, feeCollector string, totalFees sdk.Coins) error {
	if totalFees.IsZero() {
		return nil
	}
	
	// Get fee distribution config
	webConfig := webKeeper.GetWebConfig(ctx)
	feeDistribution := webConfig.FeeDistribution
	
	// If fee distribution is disabled, all fees go to fee collector (normal behavior)
	if !feeDistribution.Enabled {
		return nil
	}
	
	// Validate configuration
	if err := feeDistribution.Validate(); err != nil {
		ctx.Logger().Error("Invalid fee distribution configuration", "error", err)
		return nil // Don't fail the transaction, just log and continue
	}
	
	// Calculate fee distribution
	developerFee, validatorFee := feeDistribution.CalculateFees(totalFees)
	
	ctx.Logger().Info("Distributing transaction fees",
		"total_fees", totalFees.String(),
		"developer_fee", developerFee.String(),
		"validator_fee", validatorFee.String(),
		"developer_address", feeDistribution.DeveloperAddress,
	)
	
	// Get fee collector account
	feeCollectorAddr := fh.authKeeper.GetModuleAddress(feeCollector)
	if feeCollectorAddr == nil {
		return fmt.Errorf("fee collector account not found: %s", feeCollector)
	}
	
	// Send developer fee to developer address
	if !developerFee.IsZero() {
		developerAddr, err := sdk.AccAddressFromBech32(feeDistribution.DeveloperAddress)
		if err != nil {
			ctx.Logger().Error("Invalid developer address", "address", feeDistribution.DeveloperAddress, "error", err)
		} else {
			err = fh.bankKeeper.SendCoinsFromModuleToAccount(ctx, feeCollector, developerAddr, developerFee)
			if err != nil {
				ctx.Logger().Error("Failed to send developer fee", "error", err)
			} else {
				ctx.Logger().Info("Developer fee sent successfully",
					"amount", developerFee.String(),
					"recipient", feeDistribution.DeveloperAddress,
				)
				
				// Emit event for developer fee
				ctx.EventManager().EmitEvent(
					sdk.NewEvent(
						"developer_fee_distribution",
						sdk.NewAttribute("developer_address", feeDistribution.DeveloperAddress),
						sdk.NewAttribute("amount", developerFee.String()),
						sdk.NewAttribute("percentage", fmt.Sprintf("%d", feeDistribution.DeveloperFeePercentage)),
					),
				)
			}
		}
	}
	
	// The remaining fees stay in the fee collector for validator rewards
	// This happens automatically as we only send out the developer portion
	
	// Emit event for fee distribution
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"fee_distribution",
			sdk.NewAttribute("total_fees", totalFees.String()),
			sdk.NewAttribute("developer_fee", developerFee.String()),
			sdk.NewAttribute("validator_fee", validatorFee.String()),
			sdk.NewAttribute("developer_percentage", fmt.Sprintf("%d", feeDistribution.DeveloperFeePercentage)),
			sdk.NewAttribute("validator_percentage", fmt.Sprintf("%d", feeDistribution.ValidatorFeePercentage)),
		),
	)
	
	return nil
}
