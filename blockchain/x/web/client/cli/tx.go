package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"skaffacity/x/web/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdUpdateWebConfig())

	return cmd
}

func CmdUpdateWebConfig() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-config [enabled] [port] [host] [api-endpoint] [ws-endpoint] [theme]",
		Short: "Update web configuration",
		Args:  cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			enabled, err := strconv.ParseBool(args[0])
			if err != nil {
				return err
			}

			port, err := strconv.ParseUint(args[1], 10, 32)
			if err != nil {
				return err
			}

			host := args[2]
			apiEndpoint := args[3]
			wsEndpoint := args[4]
			theme := args[5]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			config := types.WebConfig{
				Enabled:           enabled,
				Port:              uint32(port),
				Host:              host,
				ApiEndpoint:       apiEndpoint,
				WebsocketEndpoint: wsEndpoint,
				Theme:             theme,
				Features: []string{
					"tokenfactory",
					"nft",
					"marketplace",
					"governance",
					"staking",
				},
			}

			msg := types.NewMsgUpdateWebConfig(
				clientCtx.GetFromAddress().String(),
				config,
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
