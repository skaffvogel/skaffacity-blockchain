package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewMinter returns a new Minter object with the given inflation and annual
// provisions values.
func NewMinter(inflation, annualProvisions sdk.Dec) Minter {
	return Minter{
		Inflation:        inflation,
		AnnualProvisions: annualProvisions,
	}
}

// InitialMinter returns an initial Minter object with a given inflation value.
func InitialMinter(inflation sdk.Dec) Minter {
	return NewMinter(
		inflation,
		sdk.NewDec(0),
	)
}

// DefaultInitialMinter returns a default initial Minter object for a new chain
func DefaultInitialMinter() Minter {
	return InitialMinter(
		sdk.NewDecWithPrec(5, 3), // 0.5% inflation for gaming blockchain
	)
}

// Minter represents the minting state.
type Minter struct {
	Inflation        sdk.Dec `json:"inflation" yaml:"inflation"`                 // current annual inflation rate
	AnnualProvisions sdk.Dec `json:"annual_provisions" yaml:"annual_provisions"` // current annual expected provisions
}

// String implements fmt.Stringer
func (m Minter) String() string {
	return fmt.Sprintf(`Minter:
  Inflation:        %s
  Annual Provisions: %s`,
		m.Inflation, m.AnnualProvisions)
}

// BlockProvision returns the provisions for a block based on the annual
// provisions rate.
func (m Minter) BlockProvision(params Params) sdk.Coin {
	// For gaming blockchain, we want FIXED 1 SKAF per block
	// This is more predictable for gaming rewards than percentage-based inflation
	fixedRewardPerBlock := sdk.NewInt(1000000) // 1 SKAF = 1,000,000 microSKAF
	return sdk.NewCoin(params.MintDenom, fixedRewardPerBlock)
}

// NextInflationRate returns the new inflation rate for the next hour.
func (m Minter) NextInflationRate(params Params, bondedRatio sdk.Dec) sdk.Dec {
	// For gaming blockchain, we want stable low inflation
	// Keep it simple at 0.5% annually
	return sdk.NewDecWithPrec(5, 3) // 0.5%
}

// NextAnnualProvisions returns the annual provisions based on current total
// supply and inflation rate.
func (m Minter) NextAnnualProvisions(_ Params, totalSupply sdk.Int) sdk.Dec {
	return m.Inflation.MulInt(totalSupply)
}
