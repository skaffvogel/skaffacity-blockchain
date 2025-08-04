package ante

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	webkeeper "skaffacity/x/web/keeper"
)

// FeeDistributionDecorator distributes transaction fees according to configuration
type FeeDistributionDecorator struct {
	accountKeeper ante.AccountKeeper
	bankKeeper    authtypes.BankKeeper
	webKeeper     webkeeper.Keeper
	feegrantKeeper ante.FeegrantKeeper
}

// NewFeeDistributionDecorator creates a new fee distribution decorator
func NewFeeDistributionDecorator(ak ante.AccountKeeper, bk authtypes.BankKeeper, wk webkeeper.Keeper, fk ante.FeegrantKeeper) FeeDistributionDecorator {
	return FeeDistributionDecorator{
		accountKeeper:  ak,
		bankKeeper:     bk,
		webKeeper:      wk,
		feegrantKeeper: fk,
	}
}

// AnteHandle implements the AnteDecorator interface
func (fdd FeeDistributionDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	// Only process fees if not simulating
	if !simulate {
		feeTx, ok := tx.(sdk.FeeTx)
		if !ok {
			return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
		}

		// Get fees from transaction
		fees := feeTx.GetFee()
		
		// If there are fees, distribute them
		if !fees.IsZero() {
			// Distribute fees through web keeper
			err := fdd.webKeeper.DistributeFees(ctx, authtypes.FeeCollectorName, fees)
			if err != nil {
				ctx.Logger().Error("Failed to distribute fees", "error", err)
				// Don't fail the transaction if fee distribution fails, just log
			}
		}
	}
	
	// Continue with the next decorator
	return next(ctx, tx, simulate)
}

// PostHandler for fee distribution (alternative approach)
type FeeDistributionPostHandler struct {
	webKeeper webkeeper.Keeper
}

// NewFeeDistributionPostHandler creates a new fee distribution post handler
func NewFeeDistributionPostHandler(wk webkeeper.Keeper) FeeDistributionPostHandler {
	return FeeDistributionPostHandler{
		webKeeper: wk,
	}
}

// PostHandle implements the PostHandler interface
func (fdph FeeDistributionPostHandler) PostHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, success bool, next sdk.PostHandler) (newCtx sdk.Context, err error) {
	// Only distribute fees for successful transactions that are not simulations
	if success && !simulate {
		feeTx, ok := tx.(sdk.FeeTx)
		if !ok {
			return next(ctx, tx, simulate, success)
		}

		// Get fees from transaction
		fees := feeTx.GetFee()
		
		// If there are fees, distribute them
		if !fees.IsZero() {
			ctx.Logger().Info("Distributing transaction fees", 
				"fees", fees.String(),
				"tx_hash", fmt.Sprintf("%X", ctx.TxBytes()),
			)
			
			// Distribute fees through web keeper
			err := fdph.webKeeper.DistributeFees(ctx, authtypes.FeeCollectorName, fees)
			if err != nil {
				ctx.Logger().Error("Failed to distribute fees", "error", err)
				// Don't fail the transaction if fee distribution fails, just log
			}
		}
	}
	
	// Continue with the next post handler
	return next(ctx, tx, simulate, success)
}
