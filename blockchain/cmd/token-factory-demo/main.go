package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "token-factory-demo",
		Short: "SkaffaCity Token Factory Demo",
		Long:  "Demo application for creating and managing tokens on SkaffaCity blockchain",
		Run: func(cmd *cobra.Command, args []string) {
			printWelcome()
		},
	}

	rootCmd.AddCommand(
		createTokenCmd(),
		mintTokenCmd(),
		burnTokenCmd(),
		transferTokenCmd(),
		queryTokenCmd(),
		queryBalanceCmd(),
		helpCmd(),
	)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func printWelcome() {
	fmt.Println("🏗️  SkaffaCity Token Factory Demo")
	fmt.Println("==================================")
	fmt.Println()
	fmt.Println("Welcome to the SkaffaCity Token Factory!")
	fmt.Println("This demo shows how to create and manage custom tokens.")
	fmt.Println()
	fmt.Println("Available commands:")
	fmt.Println("  create-token   - Create a new token")
	fmt.Println("  mint-token     - Mint additional tokens")
	fmt.Println("  burn-token     - Burn tokens")
	fmt.Println("  transfer-token - Transfer tokens")
	fmt.Println("  query-token    - Query token info")
	fmt.Println("  query-balance  - Query token balance")
	fmt.Println("  help           - Show usage examples")
	fmt.Println()
	fmt.Println("Use --help with any command for more details.")
	fmt.Println()
}

func createTokenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-token [name] [symbol] [description] [decimals] [max-supply] [initial-supply] [mintable] [burnable]",
		Short: "Create a new token",
		Args:  cobra.ExactArgs(8),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("🚀 Creating new token: %s (%s)\n", args[0], args[1])
			fmt.Printf("   Description: %s\n", args[2])
			fmt.Printf("   Decimals: %s\n", args[3])
			fmt.Printf("   Max Supply: %s\n", args[4])
			fmt.Printf("   Initial Supply: %s\n", args[5])
			fmt.Printf("   Mintable: %s\n", args[6])
			fmt.Printf("   Burnable: %s\n", args[7])
			
			// Validate parameters
			decimals, err := strconv.Atoi(args[3])
			if err != nil || decimals < 0 || decimals > 18 {
				return fmt.Errorf("❌ Invalid decimals: must be 0-18")
			}
			
			// Generate mock denom
			mockCreator := "skaf1abc123def456ghi789jkl012mno345pqr678stu"
			mockDenom := fmt.Sprintf("factory/%s/a1b2c3d4", mockCreator)
			
			fmt.Printf("✅ Token creation parameters validated!\n")
			fmt.Printf("   Generated denom: %s\n", mockDenom)
			fmt.Printf("   Creator: %s\n", mockCreator)
			fmt.Printf("\n📝 Next steps:\n")
			fmt.Printf("   1. Use 'skaffacityd tx tokenfactory create-token' to create on blockchain\n")
			fmt.Printf("   2. Broadcast the transaction with proper key and chain-id\n")
			fmt.Printf("   3. Query the token using the generated denom\n")
			
			return nil
		},
	}
	return cmd
}

func mintTokenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mint-token [denom] [amount] [recipient]",
		Short: "Mint tokens to an address",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("💰 Minting tokens...\n")
			fmt.Printf("   Token: %s\n", args[0])
			fmt.Printf("   Amount: %s\n", args[1])
			fmt.Printf("   Recipient: %s\n", args[2])
			
			fmt.Printf("✅ Mint transaction prepared!\n")
			fmt.Printf("   Note: Only token creator can mint (if mintable=true)\n")
			
			return nil
		},
	}
	return cmd
}

func burnTokenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "burn-token [denom] [amount]",
		Short: "Burn tokens from your address",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("🔥 Burning tokens...\n")
			fmt.Printf("   Token: %s\n", args[0])
			fmt.Printf("   Amount: %s\n", args[1])
			
			fmt.Printf("✅ Burn transaction prepared!\n")
			fmt.Printf("   Note: Token must be burnable and you must have sufficient balance\n")
			
			return nil
		},
	}
	return cmd
}

func transferTokenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer-token [recipient] [denom] [amount]",
		Short: "Transfer tokens to another address",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("📤 Transferring tokens...\n")
			fmt.Printf("   Recipient: %s\n", args[0])
			fmt.Printf("   Token: %s\n", args[1])
			fmt.Printf("   Amount: %s\n", args[2])
			
			fmt.Printf("✅ Transfer transaction prepared!\n")
			fmt.Printf("   Note: You must have sufficient balance for the transfer\n")
			
			return nil
		},
	}
	return cmd
}

func queryTokenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query-token [denom]",
		Short: "Query token information",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("🔍 Querying token information...\n")
			fmt.Printf("   Token: %s\n", args[0])
			
			// Mock token data for demo
			token := map[string]interface{}{
				"creator":        "skaf1abc123def456ghi789...",
				"denom":          args[0],
				"name":           "Demo Token",
				"symbol":         "DEMO",
				"description":    "A demonstration token on SkaffaCity",
				"decimals":       6,
				"max_supply":     "1000000000",
				"current_supply": "100000000",
				"mintable":       true,
				"burnable":       true,
				"created_at":     12345,
			}
			
			jsonData, _ := json.MarshalIndent(token, "", "  ")
			fmt.Printf("📋 Token Information:\n%s\n", string(jsonData))
			
			return nil
		},
	}
	return cmd
}

func queryBalanceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query-balance [address] [denom]",
		Short: "Query token balance",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("💳 Querying balance...\n")
			fmt.Printf("   Address: %s\n", args[0])
			fmt.Printf("   Token: %s\n", args[1])
			
			// Mock balance data for demo
			balance := map[string]interface{}{
				"address": args[0],
				"denom":   args[1],
				"balance": "1000000",
			}
			
			jsonData, _ := json.MarshalIndent(balance, "", "  ")
			fmt.Printf("💰 Balance Information:\n%s\n", string(jsonData))
			
			return nil
		},
	}
	return cmd
}

func helpCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "help",
		Short: "Show usage examples",
		Run: func(cmd *cobra.Command, args []string) {
			printUsageInfo()
		},
	}
	return cmd
}

func printUsageInfo() {
	fmt.Println("\n📚 SkaffaCity Token Factory Usage Guide:")
	fmt.Println("=====================================")
	fmt.Println()
	fmt.Println("🏗️  CREATE TOKEN:")
	fmt.Println("   ./token-factory-demo create-token \"My Token\" \"MTK\" \"My custom token\" 18 1000000000 100000000 true true")
	fmt.Println()
	fmt.Println("💰 MINT TOKENS:")
	fmt.Println("   ./token-factory-demo mint-token factory/skaf1abc.../12345678 1000000 skaf1recipient...")
	fmt.Println()
	fmt.Println("🔥 BURN TOKENS:")
	fmt.Println("   ./token-factory-demo burn-token factory/skaf1abc.../12345678 500000")
	fmt.Println()
	fmt.Println("📤 TRANSFER TOKENS:")
	fmt.Println("   ./token-factory-demo transfer-token skaf1recipient... factory/skaf1abc.../12345678 100000")
	fmt.Println()
	fmt.Println("🔍 QUERY TOKEN:")
	fmt.Println("   ./token-factory-demo query-token factory/skaf1abc.../12345678")
	fmt.Println()
	fmt.Println("💳 QUERY BALANCE:")
	fmt.Println("   ./token-factory-demo query-balance skaf1abc... factory/skaf1abc.../12345678")
	fmt.Println()
	fmt.Println("📝 Note: This is a demo application. For actual blockchain interaction,")
	fmt.Println("    use 'skaffacityd tx tokenfactory' and 'skaffacityd query tokenfactory'")
	fmt.Println()
}
