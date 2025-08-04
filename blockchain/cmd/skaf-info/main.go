package main

import (
	"fmt"
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"skaffacity/app"
)

func main() {
	// Initialize SDK config
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount("skaffa", "skaffapub")
	config.Seal()

	// Create genesis state with SKAF tokens
	genesis := app.NewDefaultGenesisState()
	
	// Pretty print the bank genesis to show SKAF tokens
	var bankGenesis banktypes.GenesisState
	if err := json.Unmarshal(genesis["bank"], &bankGenesis); err == nil {
		fmt.Println("=== SKAF TOKEN CONFIGURATION ===")
		
		// Show token metadata
		for _, metadata := range bankGenesis.DenomMetadata {
			fmt.Printf("Token Name: %s\n", metadata.Name)
			fmt.Printf("Symbol: %s\n", metadata.Symbol)
			fmt.Printf("Base Denom: %s\n", metadata.Base)
			fmt.Printf("Description: %s\n", metadata.Description)
			
			fmt.Println("\nDenom Units:")
			for _, unit := range metadata.DenomUnits {
				fmt.Printf("  - %s (10^%d)\n", unit.Denom, unit.Exponent)
			}
		}
		
		// Show initial supply
		fmt.Println("\n=== INITIAL SUPPLY ===")
		for _, coin := range bankGenesis.Supply {
			// Convert to human readable (divide by 10^6 for SKAF)
			if coin.Denom == "skaf" {
				humanAmount := coin.Amount.Quo(sdk.NewInt(1000000)) // Divide by 10^6
				fmt.Printf("Total SKAF: %s (1 billion tokens)\n", humanAmount.String())
			}
		}
		
		// Show genesis accounts
		fmt.Println("\n=== GENESIS ACCOUNTS ===")
		for _, balance := range bankGenesis.Balances {
			fmt.Printf("Account: %s\n", balance.Address)
			for _, coin := range balance.Coins {
				if coin.Denom == "skaf" {
					humanAmount := coin.Amount.Quo(sdk.NewInt(1000000))
					fmt.Printf("  SKAF Balance: %s million tokens\n", humanAmount.String())
				}
			}
			fmt.Println()
		}
	}
}
