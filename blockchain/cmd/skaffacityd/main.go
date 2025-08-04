package main

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/cosmos/cosmos-sdk/server"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/config"
	"github.com/cosmos/cosmos-sdk/client/debug"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	
	dbm "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	
	"skaffacity/app"
)

func main() {
	// Set address prefixes for SkaffaCity blockchain
	fmt.Println("🏙️ Starting SkaffaCity Blockchain...")
	fmt.Println("🔧 Configuring address prefixes (skaffa1...)...")
	app.SetConfig()
	fmt.Println("✅ Address configuration complete")
	
	rootCmd := NewRootCmd()
	
	if err := svrcmd.Execute(rootCmd, "", app.DefaultNodeHome); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error executing command: %v\n", err)
		os.Exit(1)
	}
}

func NewRootCmd() *cobra.Command {
	encodingConfig := app.MakeEncodingConfig()
	initClientCtx := client.Context{}.
		WithCodec(encodingConfig.Marshaler).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(types.AccountRetriever{}).
		WithBroadcastMode(flags.BroadcastSync).
		WithHomeDir(app.DefaultNodeHome)

	rootCmd := &cobra.Command{
		Use:   "skaffacityd",
		Short: "SkaffaCity Blockchain Daemon",
		Long: `SkaffaCity is a Cosmos SDK blockchain designed for gaming and NFTs.

This application provides:
- Native SKAF token with mint module
- NFT cosmetic system
- Player profiles and achievements  
- Social features (likes, follows)
- Moderated chat system`,
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			return client.SetCmdClientContextHandler(initClientCtx, cmd)
		},
	}

	initRootCmd(rootCmd, encodingConfig)

	return rootCmd
}

func initRootCmd(rootCmd *cobra.Command, encodingConfig app.EncodingConfig) {
	rootCmd.AddCommand(
		genutilcli.InitCmd(app.ModuleBasics, app.DefaultNodeHome),
		genutilcli.CollectGenTxsCmd(banktypes.GenesisBalancesIterator{}, app.DefaultNodeHome, genutiltypes.DefaultMessageValidator),
		genutilcli.GenTxCmd(app.ModuleBasics, encodingConfig.TxConfig, banktypes.GenesisBalancesIterator{}, app.DefaultNodeHome),
		genutilcli.ValidateGenesisCmd(app.ModuleBasics),
		genutilcli.AddGenesisAccountCmd(app.DefaultNodeHome),
		debug.Cmd(),
		config.Cmd(),
	)

	// Add server commands for blockchain node operations  
	server.AddCommands(rootCmd, app.DefaultNodeHome, newApp, nil, addModuleInitFlags)

	// Add key commands
	rootCmd.AddCommand(
		rpc.StatusCommand(),
		queryCommand(),
		txCommand(),
		keys.Commands(app.DefaultNodeHome),
	)
}

func newApp(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	appOpts servertypes.AppOptions,
) servertypes.Application {
	return app.NewSkaffaCityApp(
		logger,
		db,
		traceStore,
		true,                          // loadLatest
		map[int64]bool{},             // skipUpgradeHeights
		app.DefaultNodeHome,          // homePath
		uint(1),                      // invCheckPeriod
		app.MakeEncodingConfig(),     // encodingConfig
		appOpts,
	)
}

func appExport(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	height int64,
	forZeroHeight bool,
	jailAllowedAddrs []string,
	appOpts servertypes.AppOptions,
) (servertypes.ExportedApp, error) {
	skaffaApp := app.NewSkaffaCityApp(
		logger,
		db,
		traceStore,
		height != -1,                 // loadLatest
		map[int64]bool{},             // skipUpgradeHeights
		app.DefaultNodeHome,          // homePath
		uint(1),                      // invCheckPeriod
		app.MakeEncodingConfig(),     // encodingConfig
		appOpts,
	)
	
	return skaffaApp.ExportAppStateAndValidators(forZeroHeight, jailAllowedAddrs)
}

func addModuleInitFlags(startCmd *cobra.Command) {
	crisis.AddModuleInitFlags(startCmd)
}

func queryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "query",
		Aliases:                    []string{"q"},
		Short:                      "Querying subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		rpc.ValidatorCommand(),
		rpc.BlockCommand(),
	)

	app.ModuleBasics.AddQueryCommands(cmd)
	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")

	return cmd
}

func txCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "tx",
		Short:                      "Transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	app.ModuleBasics.AddTxCommands(cmd)
	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")

	return cmd
}
