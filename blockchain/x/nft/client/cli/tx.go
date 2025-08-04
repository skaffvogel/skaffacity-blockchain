package cli

import (
    "github.com/spf13/cobra"
    "github.com/cosmos/cosmos-sdk/client"
    "github.com/cosmos/cosmos-sdk/client/flags"
    "github.com/cosmos/cosmos-sdk/client/tx"
)

// GetTxCmd returns the transaction commands for the NFT module
func GetTxCmd() *cobra.Command {
    cmd := &cobra.Command{
        Use:                        "nft",
        Short:                      "NFT transaction subcommands",
        DisableFlagParsing:        true,
        SuggestionsMinimumDistance: 2,
        RunE:                       client.ValidateCmd,
    }

    cmd.AddCommand(
        GetCmdMintNFT(),
        GetCmdTransferNFT(),
        GetCmdAttachToItem(),
    )

    return cmd
}

func GetCmdMintNFT() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "mint [type] [name] [description] [recipient]",
        Short: "Mint a new NFT",
        Args:  cobra.ExactArgs(4),
        RunE: func(cmd *cobra.Command, args []string) error {
            clientCtx, err := client.GetClientTxContext(cmd)
            if err != nil {
                return err
            }

            // msg := types.NewMsgMintNFT(
            //     clientCtx.GetFromAddress().String(),
            //     args[0], // type
            //     args[1], // name
            //     args[2], // description
            //     args[3], // recipient
            // )

            return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), nil)
        },
    }

    flags.AddTxFlagsToCmd(cmd)
    return cmd
}

func GetCmdTransferNFT() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "transfer [nft-id] [recipient]",
        Short: "Transfer an NFT to another address",
        Args:  cobra.ExactArgs(2),
        RunE: func(cmd *cobra.Command, args []string) error {
            clientCtx, err := client.GetClientTxContext(cmd)
            if err != nil {
                return err
            }

            // msg := types.NewMsgTransferNFT(
            //     clientCtx.GetFromAddress().String(),
            //     args[0], // nft-id
            //     args[1], // recipient
            // )

            return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), nil)
        },
    }

    flags.AddTxFlagsToCmd(cmd)
    return cmd
}

func GetCmdAttachToItem() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "attach [item-id] [attachment-id]",
        Short: "Attach an NFT to an item",
        Args:  cobra.ExactArgs(2),
        RunE: func(cmd *cobra.Command, args []string) error {
            clientCtx, err := client.GetClientTxContext(cmd)
            if err != nil {
                return err
            }

            // msg := types.NewMsgAttachToItem(
            //     clientCtx.GetFromAddress().String(),
            //     args[0], // item-id
            //     args[1], // attachment-id
            // )

            return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), nil)
        },
    }

    flags.AddTxFlagsToCmd(cmd)
    return cmd
}
