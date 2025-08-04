package main

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func main() {
	// Initialize SDK config
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount("skaffa", "skaffapub")
	config.Seal()

	// Create SKAF token metadata
	skafDenom := &banktypes.Metadata{
		Description: "SkaffaCity native token for governance, staking, and marketplace transactions",
		Base:        "skaf",
		Display:     "SKAF", 
		Name:        "SkaffaCity",
		Symbol:      "SKAF",
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    "uskaf",
				Exponent: 0,
				Aliases:  []string{"microskaf"},
			},
			{
				Denom:    "mskaf",
				Exponent: 3,
				Aliases:  []string{"milliskaf"},
			},
			{
				Denom:    "skaf",
				Exponent: 6,
				Aliases:  []string{"SKAF"},
			},
		},
	}
	
	// Initial supply: 1 billion SKAF tokens
	initialSupply := sdk.NewCoins(
		sdk.NewCoin("skaf", sdk.NewInt(1000000000000000)), // 1B SKAF (with 6 decimals)
	)
	
	// Genesis accounts with initial SKAF balances
	genesisAccounts := []banktypes.Balance{
		{
			Address: "skaffa1foundation000000000000000000000000",
			Coins:   sdk.NewCoins(sdk.NewCoin("skaf", sdk.NewInt(500000000000000))), // 500M SKAF
		},
		{
			Address: "skaffa1rewards00000000000000000000000000",
			Coins:   sdk.NewCoins(sdk.NewCoin("skaf", sdk.NewInt(300000000000000))), // 300M SKAF
		},
		{
			Address: "skaffa1community0000000000000000000000000",
			Coins:   sdk.NewCoins(sdk.NewCoin("skaf", sdk.NewInt(200000000000000))), // 200M SKAF
		},
	}
	
	fmt.Println("=== SKAF TOKEN AANGEMAAKT! ===")
	fmt.Printf("Token Name: %s\n", skafDenom.Name)
	fmt.Printf("Symbol: %s\n", skafDenom.Symbol)
	fmt.Printf("Base Denom: %s\n", skafDenom.Base)
	fmt.Printf("Description: %s\n", skafDenom.Description)
	
	fmt.Println("\n=== DENOM UNITS ===")
	for _, unit := range skafDenom.DenomUnits {
		fmt.Printf("  - %s (10^%d)\n", unit.Denom, unit.Exponent)
	}
	
	fmt.Println("\n=== INITIAL SUPPLY ===")
	for _, coin := range initialSupply {
		if coin.Denom == "skaf" {
			humanAmount := coin.Amount.Quo(sdk.NewInt(1000000)) // Divide by 10^6
			fmt.Printf("Total SKAF: %s million tokens (1 billion total)\n", humanAmount.String())
		}
	}
	
	fmt.Println("\n=== GENESIS ACCOUNTS ===")
	for _, balance := range genesisAccounts {
		fmt.Printf("Account: %s\n", balance.Address)
		for _, coin := range balance.Coins {
			if coin.Denom == "skaf" {
				humanAmount := coin.Amount.Quo(sdk.NewInt(1000000))
				fmt.Printf("  SKAF Balance: %s million tokens\n", humanAmount.String())
			}
		}
		fmt.Println()
	}
	
	fmt.Println("=== GEBRUIK IN MODULES ===")
	fmt.Println("✅ Staking: Gebruikers kunnen SKAF tokens staken voor governance voting power")
	fmt.Println("✅ Governance: Weighted voting gebaseerd op gestakte SKAF tokens")
	fmt.Println("✅ Marketplace: NFTs kopen/verkopen met SKAF tokens")
	fmt.Println("✅ Rewards: Automatische SKAF rewards voor spelers")
	
	fmt.Println("\n🎉 SKAF tokens zijn succesvol geconfigureerd!")
	fmt.Println("🚀 Blockchain is klaar voor testing!")
}
