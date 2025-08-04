package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

const ModuleName = "governance"

var (
	ErrProposalNotFound     = sdkerrors.Register(ModuleName, 1, "proposal not found")
	ErrInvalidProposal      = sdkerrors.Register(ModuleName, 2, "invalid proposal")
	ErrVotingPeriodEnded    = sdkerrors.Register(ModuleName, 3, "voting period has ended")
	ErrInsufficientStake    = sdkerrors.Register(ModuleName, 4, "insufficient stake")
	ErrAlreadyVoted         = sdkerrors.Register(ModuleName, 5, "already voted")
	ErrInvalidVoteOption    = sdkerrors.Register(ModuleName, 6, "invalid vote option")
)
