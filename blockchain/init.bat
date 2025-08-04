@echo off

REM Initialize chain
skaffacityd init test --chain-id skaffacity-testnet-1

REM Create keys
skaffacityd keys add validator --keyring-backend test
skaffacityd keys add user1 --keyring-backend test
skaffacityd keys add user2 --keyring-backend test

REM Add genesis accounts
for /f "tokens=*" %%a in ('skaffacityd keys show validator -a --keyring-backend test') do set VAL_ADDR=%%a
skaffacityd add-genesis-account %VAL_ADDR% 100000000skaf

for /f "tokens=*" %%a in ('skaffacityd keys show user1 -a --keyring-backend test') do set USER1_ADDR=%%a
skaffacityd add-genesis-account %USER1_ADDR% 100000000skaf

for /f "tokens=*" %%a in ('skaffacityd keys show user2 -a --keyring-backend test') do set USER2_ADDR=%%a
skaffacityd add-genesis-account %USER2_ADDR% 100000000skaf

REM Create validator transaction
skaffacityd gentx validator 70000000skaf --chain-id skaffacity-testnet-1 --keyring-backend test

REM Collect genesis transactions
skaffacityd collect-gentxs

REM Start the chain
skaffacityd start
