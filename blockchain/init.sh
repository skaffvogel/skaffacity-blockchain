#!/bin/sh

# Initialize chain
skaffacityd init test --chain-id skaffacity-testnet-1

# Create keys
skaffacityd keys add validator --keyring-backend test
skaffacityd keys add user1 --keyring-backend test
skaffacityd keys add user2 --keyring-backend test

# Add genesis accounts
skaffacityd add-genesis-account $(skaffacityd keys show validator -a --keyring-backend test) 100000000skaf
skaffacityd add-genesis-account $(skaffacityd keys show user1 -a --keyring-backend test) 100000000skaf
skaffacityd add-genesis-account $(skaffacityd keys show user2 -a --keyring-backend test) 100000000skaf

# Create validator transaction
skaffacityd gentx validator 70000000skaf --chain-id skaffacity-testnet-1 --keyring-backend test

# Collect genesis transactions
skaffacityd collect-gentxs

# Start the chain
skaffacityd start
