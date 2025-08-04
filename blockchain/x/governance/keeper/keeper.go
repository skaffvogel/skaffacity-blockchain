package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	"skaffacity/x/governance/types"
)

type Keeper struct {
	storeKey      storetypes.StoreKey
	cdc           codec.BinaryCodec
	stakingKeeper types.StakingKeeper
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	stakingKeeper types.StakingKeeper,
) *Keeper {
	return &Keeper{
		storeKey:      storeKey,
		cdc:           cdc,
		stakingKeeper: stakingKeeper,
	}
}

// Basic stub methods
func (k Keeper) CreateProposal(ctx sdk.Context, title, description, proposer string) (uint64, error) {
	return 1, nil
}

func (k Keeper) Vote(ctx sdk.Context, proposalID uint64, voter, voteOption string) error {
	return nil
}

func (k Keeper) GetProposal(ctx sdk.Context, proposalID uint64) (types.Proposal, bool) {
	return types.Proposal{}, false
}

func (k Keeper) IsProposalActive(ctx sdk.Context, proposalID uint64) bool {
	return false
}
