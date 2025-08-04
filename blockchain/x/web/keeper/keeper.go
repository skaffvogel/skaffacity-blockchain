package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"skaffacity/x/web/types"
)

type (
	Keeper struct {
		cdc        codec.BinaryCodec
		storeKey   storetypes.StoreKey
		memKey     storetypes.StoreKey
		paramstore paramtypes.Subspace  // Back to original
		bankKeeper types.BankKeeper
		authKeeper types.AccountKeeper
		feeHandler FeeHandler
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	bankKeeper types.BankKeeper,
	authKeeper types.AccountKeeper,
) *Keeper {
	// set KeyTable if it has not already been set (use string comparison to check for zero value)
	if ps.Name() != "" && !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	keeper := &Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		memKey:     memKey,
		paramstore: ps,
		bankKeeper: bankKeeper,
		authKeeper: authKeeper,
	}
	
	// Initialize fee handler
	keeper.feeHandler = NewFeeHandler(bankKeeper, authKeeper)

	return keeper
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// SetWebConfig sets the web configuration
func (k Keeper) SetWebConfig(ctx sdk.Context, webConfig types.WebConfig) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&webConfig)
	store.Set(types.KeyPrefix(types.WebConfigKey), b)
}

// RemoveWebConfig removes the web configuration
func (k Keeper) RemoveWebConfig(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.KeyPrefix(types.WebConfigKey))
}

// GetWebConfig returns the web configuration or default if not found
func (k Keeper) GetWebConfig(ctx sdk.Context) types.WebConfig {
	webConfig, found := k.getWebConfig(ctx)
	if !found {
		// Return default config if not found
		return types.DefaultWebConfig()
	}
	return webConfig
}

// getWebConfig is the internal method that returns found status
func (k Keeper) getWebConfig(ctx sdk.Context) (val types.WebConfig, found bool) {
	store := ctx.KVStore(k.storeKey)
	
	b := store.Get(types.KeyPrefix(types.WebConfigKey))
	if b == nil {
		return val, false
	}
	
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// Fee Distribution Methods

// DistributeFees distributes transaction fees according to configuration
func (k Keeper) DistributeFees(ctx sdk.Context, feeCollector string, totalFees sdk.Coins) error {
	return k.feeHandler.DistributeFees(ctx, k, feeCollector, totalFees)
}

// GetFeeDistributionConfig returns the current fee distribution configuration
func (k Keeper) GetFeeDistributionConfig(ctx sdk.Context) types.FeeDistribution {
	webConfig := k.GetWebConfig(ctx)
	return webConfig.FeeDistribution
}

// SetDeveloperAddress updates the developer address for fee distribution
func (k Keeper) SetDeveloperAddress(ctx sdk.Context, address string) error {
	// Validate address
	if _, err := sdk.AccAddressFromBech32(address); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid developer address: %s", err.Error())
	}
	
	// Get current config
	webConfig := k.GetWebConfig(ctx)
	
	// Update developer address
	webConfig.FeeDistribution.DeveloperAddress = address
	
	// Validate new configuration
	if err := webConfig.FeeDistribution.Validate(); err != nil {
		return err
	}
	
	// Save updated config
	k.SetWebConfig(ctx, webConfig)
	
	ctx.Logger().Info("Developer address updated", "new_address", address)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"developer_address_updated",
			sdk.NewAttribute("new_address", address),
		),
	)
	
	return nil
}

// EnableFeeDistribution enables or disables fee distribution
func (k Keeper) EnableFeeDistribution(ctx sdk.Context, enabled bool) error {
	// Get current config
	webConfig := k.GetWebConfig(ctx)
	
	// Update enabled status
	webConfig.FeeDistribution.Enabled = enabled
	
	// If enabling, validate the configuration
	if enabled {
		if err := webConfig.FeeDistribution.Validate(); err != nil {
			return err
		}
	}
	
	// Save updated config
	k.SetWebConfig(ctx, webConfig)
	
	ctx.Logger().Info("Fee distribution status updated", "enabled", enabled)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"fee_distribution_status_updated",
			sdk.NewAttribute("enabled", fmt.Sprintf("%t", enabled)),
		),
	)
	
	return nil
}
