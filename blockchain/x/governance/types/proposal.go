package types

import (
	"time"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Proposal represents a governance proposal
type Proposal struct {
	ID             uint64    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Title          string    `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	Description    string    `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	Proposer       string    `protobuf:"bytes,4,opt,name=proposer,proto3" json:"proposer,omitempty"`
	Status         string    `protobuf:"bytes,5,opt,name=status,proto3" json:"status,omitempty"`
	SubmitTime     time.Time `protobuf:"bytes,6,opt,name=submit_time,json=submitTime,proto3,stdtime" json:"submit_time"`
	VotingEndTime  time.Time `protobuf:"bytes,7,opt,name=voting_end_time,json=votingEndTime,proto3,stdtime" json:"voting_end_time"`
	YesVotes       sdk.Dec   `protobuf:"bytes,8,opt,name=yes_votes,json=yesVotes,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"yes_votes"`
	NoVotes        sdk.Dec   `protobuf:"bytes,9,opt,name=no_votes,json=noVotes,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"no_votes"`
	AbstainVotes   sdk.Dec   `protobuf:"bytes,10,opt,name=abstain_votes,json=abstainVotes,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"abstain_votes"`
}

// Vote represents a vote on a proposal
type Vote struct {
	ProposalID uint64    `protobuf:"varint,1,opt,name=proposal_id,json=proposalId,proto3" json:"proposal_id,omitempty"`
	Voter      string    `protobuf:"bytes,2,opt,name=voter,proto3" json:"voter,omitempty"`
	Option     string    `protobuf:"bytes,3,opt,name=option,proto3" json:"option,omitempty"`
	Weight     sdk.Dec   `protobuf:"bytes,4,opt,name=weight,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"weight"`
	Timestamp  time.Time `protobuf:"bytes,5,opt,name=timestamp,proto3,stdtime" json:"timestamp"`
}

// Proposal status constants
const (
	StatusSubmitted    = "submitted"
	StatusVotingPeriod = "voting_period"
	StatusPassed       = "passed"
	StatusRejected     = "rejected"
	StatusFailed       = "failed"
)

// Default voting period
var DefaultVotingPeriod = time.Hour * 24 * 7 // 7 days

// ProtoMessage implements the proto.Message interface for Proposal.
func (p *Proposal) ProtoMessage() {}

// Reset implements the proto.Message interface for Proposal.
func (p *Proposal) Reset() {}

// String implements the fmt.Stringer interface for Proposal.
func (p *Proposal) String() string { return "Proposal stub" }