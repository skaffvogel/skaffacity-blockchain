package cli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
)

func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "marketplace",
		Short:                      "Marketplace transaction subcommands",
		DisableFlagParsing:        true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		GetCmdCreateListing(),
		GetCmdBuyItem(),
		GetCmdCancelListing(),
	)

	return cmd
}

func GetCmdCreateListing() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [nft-id] [price]",
		Short: "Create a new marketplace listing",
		Args:  cobra.ExactArgs(2),


		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func GetCmdBuyItem() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "buy [listing-id]",
		Short: "Buy an item from the marketplace",
		Args:  cobra.ExactArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func GetCmdCancelListing() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel [listing-id]",
		Short: "Cancel a marketplace listing",
		Args:  cobra.ExactArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
