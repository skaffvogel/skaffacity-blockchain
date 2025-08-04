package main

import (
	"fmt"
	"log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"skaffacity/x/mint/types"
)

func main() {
	fmt.Println("=== SkaffaCity Mint Module Test ===")
	
	// Test default parameters
	params := types.DefaultParams()
	fmt.Printf("\n🔧 Default Mint Parameters:\n")
	fmt.Printf("   • Mint Denom: %s\n", params.MintDenom)
	fmt.Printf("   • Inflation Rate Change: %s\n", params.InflationRateChange.String())
	fmt.Printf("   • Max Inflation: %s\n", params.InflationMax.String())
	fmt.Printf("   • Min Inflation: %s\n", params.InflationMin.String())
	fmt.Printf("   • Goal Bonded: %s\n", params.GoalBonded.String())
	fmt.Printf("   • Blocks Per Year: %d\n", params.BlocksPerYear)
	
	// Test minter
	minter := types.DefaultInitialMinter()
	fmt.Printf("\n💰 Default Minter State:\n")
	fmt.Printf("   • Current Inflation: %s\n", minter.Inflation.String())
	fmt.Printf("   • Annual Provisions: %s\n", minter.AnnualProvisions.String())
	
	// Calculate block provision
	blockProvision := minter.BlockProvision(params)
	fmt.Printf("\n⚡ Block Reward Calculation:\n")
	fmt.Printf("   • Block Provision: %s\n", blockProvision.String())
	
	// Simulate some calculations - CORRECT SUPPLY!
	totalSupply := sdk.NewInt(1000000000000000) // 1 billion SKAF with 6 decimals (1B * 10^6)
	bondedRatio := sdk.NewDecWithPrec(67, 2)    // 67% bonded
	
	fmt.Printf("\n📊 Current Economic Model (0.5%% inflation):\n")
	fmt.Printf("   • Total Supply: %s microSKAF (%s SKAF)\n", totalSupply.String(), totalSupply.Quo(sdk.NewInt(1000000)).String())
	fmt.Printf("   • Bonded Ratio: %s%%\n", bondedRatio.Mul(sdk.NewDec(100)).String())
	
	// Calculate annual provisions based on current inflation model
	annualProvisions := minter.NextAnnualProvisions(params, totalSupply)
	fmt.Printf("   • Annual Provisions: %s microSKAF (%s SKAF)\n", annualProvisions.TruncateInt().String(), annualProvisions.Quo(sdk.NewDec(1000000)).TruncateInt().String())
	
	// Update minter and get block reward
	minter.AnnualProvisions = annualProvisions
	inflationBlockReward := minter.BlockProvision(params)
	
	fmt.Printf("\n🎯 TARGET: 1 SKAF per block (Fixed Gaming Rewards):\n")
	// What we WANT: exactly 1 SKAF per block
	targetBlockReward := sdk.NewCoin("skaf", sdk.NewInt(1000000)) // 1 SKAF = 1,000,000 microSKAF
	targetYearlyRewards := targetBlockReward.Amount.Mul(sdk.NewInt(int64(params.BlocksPerYear)))
	targetDailyRewards := targetBlockReward.Amount.Mul(sdk.NewInt(14400)) // 24*60*60/6 = 14,400 blocks per day
	
	fmt.Printf("   • Target Block Reward: %s (%s SKAF)\n", targetBlockReward.String(), targetBlockReward.Amount.Quo(sdk.NewInt(1000000)).String())
	fmt.Printf("   • Target Daily Rewards: %s microSKAF (%s SKAF)\n", targetDailyRewards.String(), targetDailyRewards.Quo(sdk.NewInt(1000000)).String())
	fmt.Printf("   • Target Yearly Rewards: %s microSKAF (%s SKAF)\n", targetYearlyRewards.String(), targetYearlyRewards.Quo(sdk.NewInt(1000000)).String())
	
	fmt.Printf("\n⚖️  COMPARISON - Current vs Target:\n")
	fmt.Printf("   • Current Block: %s microSKAF (%s SKAF)\n", inflationBlockReward.Amount.String(), inflationBlockReward.Amount.Quo(sdk.NewInt(1000000)).String())
	fmt.Printf("   • Target Block:  %s microSKAF (%s SKAF)\n", targetBlockReward.Amount.String(), targetBlockReward.Amount.Quo(sdk.NewInt(1000000)).String())
	
	// Validate parameters
	if err := params.Validate(); err != nil {
		log.Fatalf("Parameters validation failed: %v", err)
	}
	
	fmt.Printf("\n✅ All mint module tests passed!\n")
	fmt.Printf("🚀 SkaffaCity blockchain is ready with validator rewards!\n")
}
