package app

import (
    "encoding/json"
    
    sdk "github.com/cosmos/cosmos-sdk/types"
    banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
    
    minttypes "skaffacity/x/mint/types"
    nfttypes "skaffacity/x/nft/types"
    marketplacetypes "skaffacity/x/marketplace/types"
    govtypes "skaffacity/x/governance/types"
    stakingtypes "skaffacity/x/staking/types"
)

// GenesisState defines the genesis state for the entire application
type GenesisState map[string]json.RawMessage

func NewDefaultGenesisState() GenesisState {
    // Create SKAF token metadata
    skafDenom := &banktypes.Metadata{
        Description: "SkaffaCity native token for governance, staking, and marketplace transactions",
        Base:        "skaf",
        Display:     "SKAF",
        Name:        "SkaffaCity",
        Symbol:      "SKAF",
        DenomUnits: []*banktypes.DenomUnit{
            {
                Denom:    "uskaf", // micro SKAF (smallest unit)
                Exponent: 0,
                Aliases:  []string{"microskaf"},
            },
            {
                Denom:    "mskaf", // milli SKAF
                Exponent: 3,
                Aliases:  []string{"milliskaf"},
            },
            {
                Denom:    "skaf", // base SKAF
                Exponent: 6,
                Aliases:  []string{"SKAF"},
            },
        },
    }
    
    // Initial supply: 1 billion SKAF tokens
    initialSupply := sdk.NewCoins(
        sdk.NewCoin("skaf", sdk.NewInt(1000000000000000)), // 1B SKAF (with 6 decimals = 1B * 10^6)
    )
    
    // Genesis accounts with initial SKAF balances
    genesisAccounts := []banktypes.Balance{
        {
            Address: "skaffa1foundation000000000000000000000000", // Foundation treasury
            Coins:   sdk.NewCoins(sdk.NewCoin("skaf", sdk.NewInt(500000000000000))), // 500M SKAF
        },
        {
            Address: "skaffa1rewards00000000000000000000000000", // Rewards pool
            Coins:   sdk.NewCoins(sdk.NewCoin("skaf", sdk.NewInt(300000000000000))), // 300M SKAF
        },
        {
            Address: "skaffa1community0000000000000000000000000", // Community fund
            Coins:   sdk.NewCoins(sdk.NewCoin("skaf", sdk.NewInt(200000000000000))), // 200M SKAF
        },
    }
    
    // Bank genesis state with SKAF configuration
    bankGenesis := banktypes.GenesisState{
        Params: banktypes.DefaultParams(),
        Balances: genesisAccounts,
        Supply:   initialSupply,
        DenomMetadata: []banktypes.Metadata{*skafDenom},
        SendEnabled: []banktypes.SendEnabled{
            {
                Denom:   "skaf",
                Enabled: true,
            },
        },
    }
    
    bankGenesisJSON, _ := json.Marshal(bankGenesis)
    
    // Mint genesis state with SKAF parameters
    mintGenesis := minttypes.GenesisState{
        Minter: minttypes.DefaultInitialMinter(),
        Params: minttypes.DefaultParams(),
    }
    
    mintGenesisJSON, _ := json.Marshal(mintGenesis)
    
    return GenesisState{
        banktypes.ModuleName:        bankGenesisJSON,
        minttypes.ModuleName:        mintGenesisJSON,
        "auth":                      []byte(`{"params":{"max_memo_characters":"256","tx_sig_limit":"7","tx_size_cost_per_byte":"10","sig_verify_cost_ed25519":"590","sig_verify_cost_secp256k1":"1000"},"accounts":[]}`),
        nfttypes.ModuleName:         []byte(`{}`),
        marketplacetypes.ModuleName: []byte(`{}`),
        govtypes.ModuleName:         []byte(`{}`),
        stakingtypes.ModuleName:     []byte(`{}`),
    }
}
