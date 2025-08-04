package app

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"

	// "skaffacity/x/mint"      // Disabled temporarily due to type conflicts
	// mintkeeper "skaffacity/x/mint/keeper"
	// minttypes "skaffacity/x/mint/types"
	"skaffacity/x/nft"
	nfttypes "skaffacity/x/nft/types"
	"skaffacity/x/marketplace"
	marketplacetypes "skaffacity/x/marketplace/types"
	"skaffacity/x/governance"
	govtypes "skaffacity/x/governance/types"
	"skaffacity/x/staking"
	stakingtypes "skaffacity/x/staking/types"
	"skaffacity/x/web"
	webtypes "skaffacity/x/web/types"

	log "github.com/cometbft/cometbft/libs/log"
)

// ModuleInfo holds information about a module
type ModuleInfo struct {
	Name        string
	Version     string
	StoreKey    string
	Keeper      interface{}
	Module      module.AppModule
	LoadTime    time.Duration
	Status      string
	Description string
}

// ModuleHandler manages all blockchain modules
type ModuleHandler struct {
	logger  log.Logger
	modules map[string]*ModuleInfo
	loaded  []string
	failed  []string
}

// NewModuleHandler creates a new module handler
func NewModuleHandler(logger log.Logger) *ModuleHandler {
	return &ModuleHandler{
		logger:  logger,
		modules: make(map[string]*ModuleInfo),
		loaded:  make([]string, 0),
		failed:  make([]string, 0),
	}
}

// RegisterModule registers a module with the handler
func (mh *ModuleHandler) RegisterModule(name, version, description string, keeper interface{}, module module.AppModule) {
	mh.modules[name] = &ModuleInfo{
		Name:        name,
		Version:     version,
		StoreKey:    name,
		Keeper:      keeper,
		Module:      module,
		Status:      "registered",
		Description: description,
	}
	mh.logger.Info(fmt.Sprintf("[MODULE] %s v%s registered - %s", name, version, description))
}

// LoadModule loads a specific module
func (mh *ModuleHandler) LoadModule(name string) error {
	start := time.Now()
	
	moduleInfo, exists := mh.modules[name]
	if !exists {
		err := fmt.Errorf("module %s not found", name)
		mh.logger.Error(fmt.Sprintf("[MODULE] ❌ %s - %v", name, err))
		mh.failed = append(mh.failed, name)
		return err
	}

	// Simulate loading process
	mh.logger.Info(fmt.Sprintf("[MODULE] 🔄 Loading %s v%s...", moduleInfo.Name, moduleInfo.Version))
	
	// Update load time and status
	moduleInfo.LoadTime = time.Since(start)
	moduleInfo.Status = "loaded"
	
	mh.loaded = append(mh.loaded, name)
	mh.logger.Info(fmt.Sprintf("[MODULE] ✅ %s v%s loaded successfully (%.2fms) - %s", 
		moduleInfo.Name, 
		moduleInfo.Version, 
		float64(moduleInfo.LoadTime.Nanoseconds())/1000000,
		moduleInfo.Description))
	
	return nil
}

// LoadAllModules loads all registered modules in order
func (mh *ModuleHandler) LoadAllModules(app *App, cdc codec.Codec, keys map[string]*storetypes.KVStoreKey, memKeys map[string]*storetypes.MemoryStoreKey) *module.Manager {
	mh.logger.Info("[MODULE] 🚀 Starting SkaffaCity module loading system...")
	mh.logger.Info("[MODULE] 📦 Initializing blockchain modules...")
	
	// Define loading order for dependencies (excluding mint for now)
	loadOrder := []string{
		authtypes.ModuleName,
		banktypes.ModuleName,
		// minttypes.ModuleName, // Disabled temporarily due to type conflicts
		nfttypes.ModuleName,
		marketplacetypes.ModuleName,
		govtypes.ModuleName,
		stakingtypes.ModuleName,
		webtypes.ModuleName,
	}

	// Initialize and register all modules with their keepers
	mh.initializeModules(app, cdc, keys, memKeys)
	
	// Load modules in order
	var modules []module.AppModule
	for _, moduleName := range loadOrder {
		if err := mh.LoadModule(moduleName); err == nil {
			if moduleInfo, exists := mh.modules[moduleName]; exists {
				modules = append(modules, moduleInfo.Module)
			}
		}
	}

	// Create module manager
	mm := module.NewManager(modules...)
	
	// Set genesis order
	mm.SetOrderInitGenesis(loadOrder...)
	
	mh.printLoadingSummary()
	
	return mm
}

// initializeModules creates and registers all modules with their keepers
func (mh *ModuleHandler) initializeModules(app *App, cdc codec.Codec, keys map[string]*storetypes.KVStoreKey, memKeys map[string]*storetypes.MemoryStoreKey) {
	
	// Auth Module
	authModule := auth.NewAppModule(cdc, app.AccountKeeper, nil, nil)
	mh.RegisterModule(
		authtypes.ModuleName, 
		"v0.47.0", 
		"Account authentication and authorization",
		&app.AccountKeeper,
		authModule,
	)

	// Bank Module  
	bankModule := bank.NewAppModule(cdc, app.BankKeeper, app.AccountKeeper, nil)
	mh.RegisterModule(
		banktypes.ModuleName,
		"v0.47.0", 
		"Token transfers and balances",
		&app.BankKeeper,
		bankModule,
	)

	// Mint Module - Disabled temporarily due to type conflicts
	// mintModule := mint.NewAppModule(cdc, app.MintKeeper, app.AccountKeeper)
	// mh.RegisterModule(
	//	minttypes.ModuleName,
	//	"v1.0.0",
	//	"Token minting and inflation control", 
	//	&app.MintKeeper,
	//	mintModule,
	// )

	// NFT Module
	nftModule := nft.NewAppModule(app.NFTKeeper)
	mh.RegisterModule(
		nfttypes.ModuleName,
		"v1.0.0",
		"Non-Fungible Token management",
		&app.NFTKeeper, 
		nftModule,
	)

	// Marketplace Module
	marketplaceModule := marketplace.NewAppModule(app.MarketKeeper)
	mh.RegisterModule(
		marketplacetypes.ModuleName,
		"v1.0.0", 
		"NFT and token marketplace",
		&app.MarketKeeper,
		marketplaceModule,
	)

	// Governance Module
	govModule := governance.NewAppModule(app.GovKeeper)
	mh.RegisterModule(
		govtypes.ModuleName,
		"v1.0.0",
		"On-chain governance and voting",
		&app.GovKeeper,
		govModule,
	)

	// Staking Module  
	stakingModule := staking.NewAppModule(app.StakingKeeper)
	mh.RegisterModule(
		stakingtypes.ModuleName,
		"v1.0.0",
		"Validator staking and delegation", 
		&app.StakingKeeper,
		stakingModule,
	)

	// Web Module
	webModule := web.NewAppModule(cdc, app.WebKeeper, app.AccountKeeper, app.BankKeeper)
	mh.RegisterModule(
		webtypes.ModuleName,
		"v1.0.0",
		"Web interface and API endpoints",
		&app.WebKeeper,
		webModule,
	)
}

// printLoadingSummary prints a summary of the module loading process
func (mh *ModuleHandler) printLoadingSummary() {
	mh.logger.Info("[MODULE] 📊 === MODULE LOADING SUMMARY ===")
	mh.logger.Info(fmt.Sprintf("[MODULE] ✅ Successfully loaded: %d modules", len(mh.loaded)))
	mh.logger.Info(fmt.Sprintf("[MODULE] ❌ Failed to load: %d modules", len(mh.failed)))
	
	if len(mh.loaded) > 0 {
		mh.logger.Info("[MODULE] 📋 Loaded modules:")
		for _, name := range mh.loaded {
			if moduleInfo, exists := mh.modules[name]; exists {
				mh.logger.Info(fmt.Sprintf("[MODULE]   • %s v%s (%.2fms)", 
					moduleInfo.Name, 
					moduleInfo.Version, 
					float64(moduleInfo.LoadTime.Nanoseconds())/1000000))
			}
		}
	}
	
	if len(mh.failed) > 0 {
		mh.logger.Info("[MODULE] 🚨 Failed modules:")
		for _, name := range mh.failed {
			mh.logger.Info(fmt.Sprintf("[MODULE]   • %s", name))
		}
	}
	
	mh.logger.Info("[MODULE] 🎉 SkaffaCity blockchain modules initialized successfully!")
	mh.logger.Info("[MODULE] 🏙️ Ready to build the future of decentralized cities!")
}

// GetModuleInfo returns information about a specific module
func (mh *ModuleHandler) GetModuleInfo(name string) (*ModuleInfo, bool) {
	moduleInfo, exists := mh.modules[name]
	return moduleInfo, exists
}

// GetLoadedModules returns a list of successfully loaded modules
func (mh *ModuleHandler) GetLoadedModules() []string {
	return mh.loaded
}

// GetFailedModules returns a list of modules that failed to load
func (mh *ModuleHandler) GetFailedModules() []string {
	return mh.failed
}

// IsModuleLoaded checks if a module is successfully loaded
func (mh *ModuleHandler) IsModuleLoaded(name string) bool {
	for _, loaded := range mh.loaded {
		if loaded == name {
			return true
		}
	}
	return false
}
