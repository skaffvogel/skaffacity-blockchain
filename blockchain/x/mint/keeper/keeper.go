package keeper

import (
	"github.com/cometbft/cometbft/libs/log"
	
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"skaffacity/x/mint/types"
)

// Keeper of the mint store
type Keeper struct {
	cdc              codec.BinaryCodec
	storeKey         storetypes.StoreKey
	paramstore       paramtypes.Subspace
	accountKeeper    types.AccountKeeper
	bankKeeper       types.BankKeeper
	stakingKeeper    types.StakingKeeper
	feeCollectorName string
}

// NewKeeper creates a new mint Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec,
	key storetypes.StoreKey,
	ak types.AccountKeeper,
	bk types.BankKeeper,
) Keeper {
	return Keeper{
		cdc:              cdc,
		storeKey:         key,
		paramstore:       paramtypes.Subspace{}, // Use empty subspace for now
		accountKeeper:    ak,
		bankKeeper:       bk,
		stakingKeeper:    nil, // Will be set later if needed
		feeCollectorName: "fee_collector",
	}
}

// NewKeeperFull creates a new mint Keeper instance with full parameters
func NewKeeperFull(
	cdc codec.BinaryCodec,
	key storetypes.StoreKey,
	paramstore paramtypes.Subspace,
	ak types.AccountKeeper,
	bk types.BankKeeper,
	sk types.StakingKeeper,
	feeCollectorName string,
) Keeper {
	return Keeper{
		cdc:              cdc,
		storeKey:         key,
		paramstore:       paramstore,
		accountKeeper:    ak,
		bankKeeper:       bk,
		stakingKeeper:    sk,
		feeCollectorName: feeCollectorName,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

// GetParams returns the total set of minting parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramstore.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the total set of minting parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}

// GetMinter gets the minter.
func (k Keeper) GetMinter(ctx sdk.Context) (minter types.Minter) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.MinterKey)
	if b == nil {
		panic("stored minter should not have been nil")
	}

	k.cdc.MustUnmarshal(b, &minter)
	return
}

// SetMinter sets the minter.
func (k Keeper) SetMinter(ctx sdk.Context, minter types.Minter) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&minter)
	store.Set(types.MinterKey, b)
}

// MintCoins implements an alias call to the underlying supply keeper's
// MintCoins to be used in BeginBlocker.
func (k Keeper) MintCoins(ctx sdk.Context, newCoins sdk.Coins) error {
	if newCoins.Empty() {
		// skip as no coins need to be minted
		return nil
	}

	return k.bankKeeper.MintCoins(ctx, types.ModuleName, newCoins)
}

// AddCollectedFees implements an alias call to the underlying supply keeper's
// AddCollectedFees to be used in BeginBlocker.
func (k Keeper) AddCollectedFees(ctx sdk.Context, fees sdk.Coins) error {
	return k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, k.feeCollectorName, fees)
}

// BlockProvision returns the provisions for a block based on the annual
// provisions rate.
func (k Keeper) BlockProvision(ctx sdk.Context) sdk.Coin {
	minter := k.GetMinter(ctx)
	params := k.GetParams(ctx)
	return minter.BlockProvision(params)
}

// GetProportions gets the balance of the `MintedDenom` from minted coins and returns coins according to the `AllocationRatio`.
func (k Keeper) GetProportions(ctx sdk.Context, mintedCoin sdk.Coin, ratio sdk.Dec) sdk.Coin {
	return sdk.NewCoin(mintedCoin.Denom, sdk.NewDecFromInt(mintedCoin.Amount).Mul(ratio).TruncateInt())
}
