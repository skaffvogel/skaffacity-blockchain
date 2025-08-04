package app

import (
    "os"
    
    "github.com/cosmos/cosmos-sdk/baseapp"
    "github.com/cosmos/cosmos-sdk/types/module"
    // "github.com/cosmos/cosmos-sdk/x/bank"   // Used in moduleHandler
    bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
    banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
    // "github.com/cosmos/cosmos-sdk/x/auth"   // Used in moduleHandler
    authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
    authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
    "github.com/cosmos/cosmos-sdk/codec"
    "github.com/cosmos/cosmos-sdk/codec/types"
    paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
    storetypes "github.com/cosmos/cosmos-sdk/store/types"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/cosmos/cosmos-sdk/server/api"
    "github.com/cosmos/cosmos-sdk/server/config"
    servertypes "github.com/cosmos/cosmos-sdk/server/types"
    "github.com/cosmos/cosmos-sdk/client"
    dbm "github.com/cometbft/cometbft-db"
    log "github.com/cometbft/cometbft/libs/log"
    tmtypes "github.com/cometbft/cometbft/types"
    
    // "skaffacity/x/mint"        // TODO: uncomment when mint keeper is fixed
    // mintkeeper "skaffacity/x/mint/keeper"
    // minttypes "skaffacity/x/mint/types"
    
    // "skaffacity/x/nft"         // Used in moduleHandler
    nftkeeper "skaffacity/x/nft/keeper"
    nfttypes "skaffacity/x/nft/types"
    // "skaffacity/x/marketplace" // Used in moduleHandler
    marketplacekeeper "skaffacity/x/marketplace/keeper"
    marketplacetypes "skaffacity/x/marketplace/types"
    // "skaffacity/x/governance"  // Used in moduleHandler
    governancekeeper "skaffacity/x/governance/keeper"
    govtypes "skaffacity/x/governance/types"
    // "skaffacity/x/staking"     // Used in moduleHandler
    stakingkeeper "skaffacity/x/staking/keeper"
    stakingtypes "skaffacity/x/staking/types"
    // "skaffacity/x/web"         // Used in moduleHandler
    webkeeper "skaffacity/x/web/keeper"
    webtypes "skaffacity/x/web/types"
)

const (
    AppName = "skaffacity"
    Version = "1.0.0"
    
    // Bech32 address prefixes
    Bech32PrefixAccAddr  = "skaffa"
    Bech32PrefixAccPub   = "skaffapub"
    Bech32PrefixValAddr  = "skaffavaloper"
    Bech32PrefixValPub   = "skaffavaloperpub"
    Bech32PrefixConsAddr = "skaffavalcons"
    Bech32PrefixConsPub  = "skaffavalconspub"
)

var (
    DefaultNodeHome = os.ExpandEnv("$HOME/skaffacity")
)

// SetConfig sets the global SDK configuration for address prefixes
func SetConfig() {
    config := sdk.GetConfig()
    config.SetBech32PrefixForAccount(Bech32PrefixAccAddr, Bech32PrefixAccPub)
    config.SetBech32PrefixForValidator(Bech32PrefixValAddr, Bech32PrefixValPub)
    config.SetBech32PrefixForConsensusNode(Bech32PrefixConsAddr, Bech32PrefixConsPub)
    config.Seal()
}

type App struct {
    *baseapp.BaseApp
    
    // Store keys
    keys    map[string]*storetypes.KVStoreKey
    tkeys   map[string]*storetypes.TransientStoreKey
    memKeys map[string]*storetypes.MemoryStoreKey
    
    // Keepers
    AccountKeeper authkeeper.AccountKeeper
    BankKeeper    bankkeeper.Keeper
    // MintKeeper    mintkeeper.Keeper  // Disabled temporarily
    NFTKeeper     nftkeeper.Keeper
    MarketKeeper  marketplacekeeper.Keeper
    GovKeeper     governancekeeper.Keeper
    StakingKeeper stakingkeeper.Keeper
    WebKeeper     webkeeper.Keeper
    
    // Module Manager
    mm *module.Manager
    
    // Module Handler - centralized module management
    moduleHandler *ModuleHandler
    
    // Encoding config
    cdc               codec.Codec
    legacyAmino       *codec.LegacyAmino
    interfaceRegistry types.InterfaceRegistry
}

func NewApp(
    logger log.Logger,
    db dbm.DB,
    traceStore interface{},
    loadLatest bool,
    skipUpgradeHeights map[int64]bool,
    homePath string,
    invCheckPeriod uint,
    encodingConfig interface{},
    appOpts interface{},
    baseAppOptions ...interface{},
) *App {
    // Create codec
    interfaceRegistry := types.NewInterfaceRegistry()
    cdc := codec.NewProtoCodec(interfaceRegistry)
    legacyAmino := codec.NewLegacyAmino()
    
    // Create base app
    bApp := baseapp.NewBaseApp(AppName, logger, db, nil)
    
    // Store keys
    keys := sdk.NewKVStoreKeys(
        authtypes.StoreKey,
        banktypes.StoreKey,
        // minttypes.StoreKey, // Disabled temporarily
        nfttypes.StoreKey,
        marketplacetypes.StoreKey,
        govtypes.StoreKey,
        stakingtypes.StoreKey,
        webtypes.StoreKey,
    )
    
    memKeys := sdk.NewMemoryStoreKeys(webtypes.MemStoreKey)
    
    app := &App{
        BaseApp:           bApp,
        cdc:               cdc,
        legacyAmino:       legacyAmino,
        interfaceRegistry: interfaceRegistry,
        keys:              keys,
        memKeys:           memKeys,
    }
    
    // Initialize module handler
    logger.Info("[SKAFFACITY] 🏙️ Initializing SkaffaCity Blockchain with Module Handler...")
    app.moduleHandler = NewModuleHandler(logger)
    
    // Initialize account keeper
    app.AccountKeeper = authkeeper.NewAccountKeeper(
        cdc,
        keys[authtypes.StoreKey],
        authtypes.ProtoBaseAccount,
        nil, // maccPerms - module account permissions
        sdk.Bech32MainPrefix,
        authtypes.NewModuleAddress("gov").String(),
    )
    
    // Initialize bank keeper
    app.BankKeeper = bankkeeper.NewBaseKeeper(
        cdc,
        keys[banktypes.StoreKey],
        app.AccountKeeper,
        nil, // blocked addresses
        authtypes.NewModuleAddress("gov").String(),
    )
    
    // Initialize mint keeper - using a basic implementation for now
    // app.MintKeeper = *mintkeeper.NewKeeper(
    //     cdc,
    //     keys[minttypes.StoreKey],
    //     app.AccountKeeper,
    //     app.BankKeeper,
    // )
    
    // Initialize custom keepers with proper dependencies
    app.NFTKeeper = *nftkeeper.NewKeeper(
        cdc,
        keys[nfttypes.StoreKey],
        app.BankKeeper, // Add bank keeper dependency
    )
    
    app.StakingKeeper = *stakingkeeper.NewKeeper(
        cdc,
        keys[stakingtypes.StoreKey],
        app.BankKeeper,
    )
    
    app.MarketKeeper = *marketplacekeeper.NewKeeper(
        cdc,
        keys[marketplacetypes.StoreKey],
        app.BankKeeper,
        &app.NFTKeeper,
    )
    
    app.GovKeeper = *governancekeeper.NewKeeper(
        cdc,
        keys[govtypes.StoreKey],
        &app.StakingKeeper,
    )
    
    app.WebKeeper = *webkeeper.NewKeeper(
        cdc,
        keys[webtypes.StoreKey],
        memKeys[webtypes.MemStoreKey],
        paramtypes.Subspace{}, // Use zero value instead of nil
        app.BankKeeper,
        app.AccountKeeper,
    )
    
    // Use module handler to load all modules with proper initialization
    app.mm = app.moduleHandler.LoadAllModules(app, cdc, keys, memKeys)
    
    // Mount stores
    app.MountKVStores(keys)
    app.MountMemoryStores(memKeys)
    
    if loadLatest {
        if err := app.LoadLatestVersion(); err != nil {
            logger.Error("[SKAFFACITY] ❌ Failed to load latest version", "error", err)
            panic(err)
        }
        logger.Info("[SKAFFACITY] ✅ Latest version loaded successfully")
    }
    
    logger.Info("[SKAFFACITY] 🎉 SkaffaCity blockchain application initialized successfully!")
    logger.Info("[SKAFFACITY] 🔑 Address prefix: skaffa1...")
    
    return app
}

// RegisterAPIRoutes registers all application module routes with the provided API server.
func (app *App) RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig) {
    // For now, we'll implement a minimal version
    // In a full implementation, you would register module-specific API routes here
}

// RegisterNodeService registers the node gRPC service with the client context.
func (app *App) RegisterNodeService(clientCtx client.Context) {
    // For now, we'll implement a minimal version
    // In a full implementation, you would register node services here
}

// RegisterTendermintService registers the tendermint queries service.
func (app *App) RegisterTendermintService(clientCtx client.Context) {
    // For now, we'll implement a minimal version
    // In a full implementation, you would register tendermint services here
}

// RegisterTxService registers the tx service.
func (app *App) RegisterTxService(clientCtx client.Context) {
    // For now, we'll implement a minimal version
    // In a full implementation, you would register tx services here
}

// ExportAppStateAndValidators exports the state of the application for a genesis file.
func (app *App) ExportAppStateAndValidators(forZeroHeight bool, jailAllowedAddrs []string) (servertypes.ExportedApp, error) {
    // For now, return a basic implementation
    // In a full implementation, you would export the actual state
    return servertypes.ExportedApp{
        AppState:    []byte("{}"),
        Validators:  []tmtypes.GenesisValidator{},
        Height:      0,
        ConsensusParams: nil,
    }, nil
}

// NewSkaffaCityApp creates a new SkaffaCity blockchain application
func NewSkaffaCityApp(
    logger log.Logger,
    db dbm.DB,
    traceStore interface{},
    loadLatest bool,
    skipUpgradeHeights map[int64]bool,
    homePath string,
    invCheckPeriod uint,
    encodingConfig interface{},
    appOpts interface{},
    baseAppOptions ...interface{},
) *App {
    return NewApp(
        logger,
        db,
        traceStore,
        loadLatest,
        skipUpgradeHeights,
        homePath,
        invCheckPeriod,
        encodingConfig,
        appOpts,
        baseAppOptions...,
    )
}

// GetModuleHandler returns the module handler instance
func (app *App) GetModuleHandler() *ModuleHandler {
    return app.moduleHandler
}

// GetLoadedModules returns a list of successfully loaded modules
func (app *App) GetLoadedModules() []string {
    if app.moduleHandler != nil {
        return app.moduleHandler.GetLoadedModules()
    }
    return []string{}
}

// IsModuleLoaded checks if a specific module is loaded
func (app *App) IsModuleLoaded(moduleName string) bool {
    if app.moduleHandler != nil {
        return app.moduleHandler.IsModuleLoaded(moduleName)
    }
    return false
}
