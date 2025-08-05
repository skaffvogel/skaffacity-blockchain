#!/bin/bash

# Create app.toml for SkaffaCity blockchain
# Use this script if app.toml is missing on your VPS

set -e

print_status() {
    echo "🔄 $1"
}

print_success() {
    echo "✅ $1"
}

print_error() {
    echo "❌ $1"
}

# Check if we're in the right directory
if [ ! -f "config/config.toml" ]; then
    print_error "config/config.toml not found. Please run this from the ~/.skaffacity directory"
    exit 1
fi

# Create app.toml
print_status "Creating app.toml configuration file..."

cat > "config/app.toml" << 'EOF'
# This is a TOML config file for SkaffaCity blockchain.
# For more information, see the documentation at https://docs.cosmos.network/

###############################################################################
###                           Base Configuration                            ###
###############################################################################

# Defines the minimum gas prices to accept for transactions. Any transaction with
# a gas price lower than this value will be rejected.
minimum-gas-prices = "0.001uskaf"

# default: the last 100 states are kept in addition to every 500th state; pruning at 10 block intervals
# nothing: all historic states will be saved, nothing will be deleted (i.e. archival node)
# everything: all saved states will be deleted, storing only the current and previous state; pruning at 10 block intervals
# custom: allow pruning options to be manually specified through 'pruning-keep-recent', 'pruning-keep-every', and 'pruning-interval'
pruning = "default"

# These are applied if and only if the pruning strategy is custom.
pruning-keep-recent = "0"
pruning-keep-every = "0"
pruning-interval = "0"

# HaltHeight contains a non-zero block height at which a node will gracefully
# halt and shutdown that can be used to assist upgrades and testing.
halt-height = 0

# HaltTime contains a non-zero minimum block time (in Unix seconds) at which
# a node will gracefully halt and shutdown that can be used to assist upgrades
# and testing.
halt-time = 0

# MinRetainBlocks defines the minimum block height offset from the current
# block being committed, such that all blocks past this offset are pruned
# from Tendermint.
min-retain-blocks = 0

# InterBlockCache enables inter-block caching.
inter-block-cache = true

# IndexEvents defines the set of events in the form {eventType}.{attributeKey},
# which informs Tendermint what to index. If empty, all events will be indexed.
index-events = []

# IavlCacheSize set the size of the iavl tree cache.
iavl-cache-size = 781250

# IavlDisableFastNode enables or disables the fast node feature of IAVL.
iavl-disable-fastnode = false

###############################################################################
###                         Telemetry Configuration                         ###
###############################################################################

[telemetry]

# Prefixed with keys to separate services.
service-name = ""

# Enabled enables the application telemetry functionality. When enabled,
# an in-memory sink is also enabled by default. Operators may also enabled
# other sinks such as Prometheus.
enabled = false

# Enable prefixing gauge values with hostname.
enable-hostname = false

# Enable adding hostname to labels.
enable-hostname-label = false

# Enable adding service to labels.
enable-service-label = false

# PrometheusRetentionTime, when positive, enables a Prometheus metrics sink.
prometheus-retention-time = 0

# GlobalLabels defines a global set of name/value label tuples applied to all
# metrics emitted using the wrapper functions defined in telemetry package.
global-labels = [
]

###############################################################################
###                           API Configuration                             ###
###############################################################################

[api]

# Enable defines if the API server should be enabled.
enable = true

# Swagger defines if swagger documentation should automatically be registered.
swagger = true

# Address defines the API server to listen on.
address = "tcp://0.0.0.0:1317"

# MaxOpenConnections defines the number of maximum open connections.
max-open-connections = 1000

# RPCReadTimeout defines the Tendermint RPC read timeout (in seconds).
rpc-read-timeout = 10

# RPCWriteTimeout defines the Tendermint RPC write timeout (in seconds).
rpc-write-timeout = 0

# RPCMaxBodyBytes defines the Tendermint maximum response body (in bytes).
rpc-max-body-bytes = 1000000

# EnableUnsafeCORS defines if CORS should be enabled (unsafe - use it at your own risk).
enabled-unsafe-cors = true

###############################################################################
###                           Rosetta Configuration                         ###
###############################################################################

[rosetta]

# Enable defines if the Rosetta API server should be enabled.
enable = false

# Address defines the Rosetta API server to listen on.
address = ":8080"

# Network defines the name of the blockchain that will be returned by Rosetta.
blockchain = "skaffacity"

# Network defines the name of the network that will be returned by Rosetta.
network = "mainnet"

# Retries defines the number of retries when connecting to the node before failing.
retries = 3

# Offline defines if Rosetta server should run in offline mode.
offline = false

###############################################################################
###                           gRPC Configuration                            ###
###############################################################################

[grpc]

# Enable defines if the gRPC server should be enabled.
enable = true

# Address defines the gRPC server address to bind to.
address = "0.0.0.0:9090"

###############################################################################
###                        gRPC Web Configuration                           ###
###############################################################################

[grpc-web]

# GRPCWebEnable defines if the gRPC-web should be enabled.
# NOTE: gRPC must also be enabled, otherwise, this configuration is a no-op.
enable = true

# Address defines the gRPC-web server address to bind to.
address = "0.0.0.0:9091"

# EnableUnsafeCORS defines if CORS should be enabled (unsafe - use it at your own risk).
enable-unsafe-cors = true

###############################################################################
###                        State Sync Configuration                         ###
###############################################################################

[state-sync]

# State sync rapidly bootstraps a new node by discovering, fetching, and restoring a state machine
# snapshot from peers instead of fetching and replaying historical blocks. Requires some peers in
# the network to take and serve state machine snapshots. State sync is not attempted if the node
# has any local state (LastBlockHeight > 0). The node will have a truncated block history,
# starting from the height of the snapshot.
enable = false

# RPC servers (comma-separated) for light client verification of the synced state machine and
# retrieval of state data for node bootstrapping. Also needs a trusted height and corresponding
# header hash obtained from a trusted source, and a period during which validators can be trusted.
rpc-servers = ""

# Trust height and corresponding header hash from a trusted source.
trust-height = 0
trust-hash = ""

# Time period during which validators can be trusted.
trust-period = "112h0m0s"

# Number of snapshots to keep on disk. 0 to disable snapshots.
snapshots = 3

# Snapshot format. Set to 0 to disable snapshots.
snapshot-format = 1

# Snapshot interval. 0 to disable snapshots.
snapshot-interval = 0

# Snapshot keep recent. 0 to keep all snapshots.
snapshot-keep-recent = 2

EOF

print_success "app.toml created successfully!"
print_status "File created at: config/app.toml"

# Verify the file was created
if [ -f "config/app.toml" ]; then
    print_success "✅ app.toml exists and is ready to use"
    print_status "Key configurations:"
    echo "   - API enabled on port 1317"
    echo "   - gRPC enabled on port 9090"
    echo "   - CORS enabled for web access"
    echo "   - Minimum gas price: 0.001uskaf"
else
    print_error "❌ Failed to create app.toml"
    exit 1
fi

print_status "You can now restart your SkaffaCity blockchain service!"
