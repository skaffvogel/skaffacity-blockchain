package types

import (
    sdk "github.com/cosmos/cosmos-sdk/types"
    "time"
)

// VoteOption defines vote options as string
type VoteOption string

const (
    VoteYes     VoteOption = "yes"
    VoteNo      VoteOption = "no"
    VoteAbstain VoteOption = "abstain"
)

// VotingParams defines the parameters for voting
type VotingParams struct {
    VotingPeriod    time.Duration `json:"voting_period"`
    QuorumThreshold sdk.Dec       `json:"quorum_threshold"`
    MinStakeToVote  sdk.Int       `json:"min_stake_to_vote"`
}
