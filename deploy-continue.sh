#!/bin/bash

# SkaffaCity Blockchain VPS Deployment - Continue from Configuration Step
# Use this script to continue deployment from where it usually fails
# This skips the build process and starts from system configuration

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_status() {
    echo -e "${BLUE}🔄 $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️ $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_header() {
    echo -e "${BLUE}===========================================${NC}"
    echo -e "${BLUE}🏙️  SkaffaCity Blockchain - Continue Deploy${NC}"
    echo -e "${BLUE}===========================================${NC}"
}

# Configuration variables
HOME_DIR="$HOME/skaffacity"
CHAIN_ID="skaffacity-1"
MONIKER="skaffacity-validator"
KEYRING_BACKEND="test"
FEE_DISTRIBUTION_DEV_ADDRESS=""
CREATE_DEV_ADDRESS_NOW=true

print_header

# Check if binary exists
if ! command -v skaffacityd &> /dev/null; then
    print_error "skaffacityd binary not found! Please run the full deployment script first."
    print_status "Run: ./deploy-vps.sh or ./deploy-alternative.sh"
    exit 1
fi

print_success "Binary found: $(which skaffacityd)"

# Step 1: Clean existing configuration
print_status "Cleaning existing configuration..."
sudo systemctl stop skaffacityd 2>/dev/null || true
sudo pkill skaffacityd 2>/dev/null || true
rm -rf "$HOME_DIR"
print_success "Configuration cleaned"

# Step 2: Initialize node
print_status "Initializing blockchain node..."
skaffacityd init "$MONIKER" --chain-id "$CHAIN_ID" --home "$HOME_DIR"
print_success "Node initialized"

# Step 3: Create app.toml if missing
print_status "Creating app.toml configuration..."
cat > "$HOME_DIR/config/app.toml" << 'EOF'
# SkaffaCity Blockchain Configuration
minimum-gas-prices = "0.001uskaf"
pruning = "default"
pruning-keep-recent = "0"
pruning-keep-every = "0"
pruning-interval = "0"
halt-height = 0
halt-time = 0
min-retain-blocks = 0
inter-block-cache = true
index-events = []
iavl-cache-size = 781250
iavl-disable-fastnode = true

[telemetry]
service-name = ""
enabled = false
enable-hostname = false
enable-hostname-label = false
enable-service-label = false
prometheus-retention-time = 0
global-labels = []

[api]
enable = true
swagger = false
address = "tcp://0.0.0.0:1317"
max-open-connections = 1000
rpc-read-timeout = 10
rpc-write-timeout = 0
rpc-max-body-bytes = 1000000
enabled-unsafe-cors = true

[grpc]
enable = true
address = "0.0.0.0:9090"
max-recv-msg-size = "10485760"
max-send-msg-size = "2147483647"

[grpc-web]
enable = true
address = "0.0.0.0:9091"
enable-unsafe-cors = false

[state-sync]
snapshot-interval = 0
snapshot-keep-recent = 2

[store]
streamers = []

[streamers]
[streamers.file]
keys = ["*", ]
write_dir = ""
prefix = ""
output-metadata = "true"
stop-node-on-error = "true"
fsync = "false"
EOF
print_success "app.toml created"

# Step 4: Update config.toml
print_status "Configuring config.toml..."
CONFIG_FILE="$HOME_DIR/config/config.toml"
if [ -f "$CONFIG_FILE" ]; then
    sed -i 's/laddr = "tcp:\/\/127.0.0.1:26657"/laddr = "tcp:\/\/0.0.0.0:26657"/' "$CONFIG_FILE"
    sed -i 's/cors_allowed_origins = \[\]/cors_allowed_origins = ["*"]/' "$CONFIG_FILE"
    sed -i 's/allow_duplicate_ip = false/allow_duplicate_ip = true/' "$CONFIG_FILE"
    print_success "config.toml updated"
else
    print_warning "config.toml not found, skipping updates"
fi

# Step 5: Create validator account
print_status "Creating validator account..."
if ! skaffacityd keys show validator --home "$HOME_DIR" --keyring-backend "$KEYRING_BACKEND" 2>/dev/null; then
    skaffacityd keys add validator --home "$HOME_DIR" --keyring-backend "$KEYRING_BACKEND"
    print_success "Validator account created"
else
    print_warning "Validator account already exists"
fi

# Step 6: Create developer account for fee distribution
if [ "$CREATE_DEV_ADDRESS_NOW" = true ]; then
    print_status "Creating developer account for fee distribution..."
    if ! skaffacityd keys show developer --home "$HOME_DIR" --keyring-backend "$KEYRING_BACKEND" 2>/dev/null; then
        skaffacityd keys add developer --home "$HOME_DIR" --keyring-backend "$KEYRING_BACKEND"
        FEE_DISTRIBUTION_DEV_ADDRESS=$(skaffacityd keys show developer -a --home "$HOME_DIR" --keyring-backend "$KEYRING_BACKEND")
        print_success "Developer account created: $FEE_DISTRIBUTION_DEV_ADDRESS"
    else
        FEE_DISTRIBUTION_DEV_ADDRESS=$(skaffacityd keys show developer -a --home "$HOME_DIR" --keyring-backend "$KEYRING_BACKEND")
        print_warning "Developer account already exists: $FEE_DISTRIBUTION_DEV_ADDRESS"
    fi
fi

# Step 7: Add genesis account
print_status "Adding genesis account..."
VALIDATOR_ADDRESS=$(skaffacityd keys show validator -a --home "$HOME_DIR" --keyring-backend "$KEYRING_BACKEND")
skaffacityd add-genesis-account "$VALIDATOR_ADDRESS" 1000000000000uskaf --home "$HOME_DIR"
print_success "Genesis account added"

# Step 8: Create genesis transaction
print_status "Creating genesis transaction..."
skaffacityd gentx validator 500000000000uskaf \
    --chain-id "$CHAIN_ID" \
    --home "$HOME_DIR" \
    --keyring-backend "$KEYRING_BACKEND"
print_success "Genesis transaction created"

# Step 9: Collect genesis transactions
print_status "Collecting genesis transactions..."
skaffacityd collect-gentxs --home "$HOME_DIR"
print_success "Genesis transactions collected"

# Step 10: Validate genesis
print_status "Validating genesis..."
skaffacityd validate-genesis --home "$HOME_DIR"
print_success "Genesis validated"

# Step 11: Configure firewall
print_status "Configuring firewall..."
sudo ufw allow 22/tcp 2>/dev/null || true
sudo ufw allow 26656/tcp 2>/dev/null || true
sudo ufw allow 26657/tcp 2>/dev/null || true
sudo ufw allow 1317/tcp 2>/dev/null || true
sudo ufw allow 9090/tcp 2>/dev/null || true
sudo ufw --force enable 2>/dev/null || true
print_success "Firewall configured"

# Step 12: Create systemd service
print_status "Creating systemd service..."
sudo tee /etc/systemd/system/skaffacityd.service > /dev/null <<EOF
[Unit]
Description=SkaffaCity Blockchain Node
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
User=$USER
WorkingDirectory=$HOME
ExecStart=/usr/local/bin/skaffacityd start --home $HOME_DIR --minimum-gas-prices="0.001uskaf"
Restart=always
RestartSec=3
LimitNOFILE=65535
Environment=DAEMON_HOME=$HOME_DIR
Environment=DAEMON_NAME=skaffacityd
StandardOutput=journal
StandardError=journal
SyslogIdentifier=skaffacityd

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sudo systemctl enable skaffacityd
print_success "Systemd service created"

# Step 13: Start the service
print_status "Starting SkaffaCity blockchain service..."
sudo systemctl start skaffacityd
sleep 5

# Step 14: Check service status
print_status "Checking service status..."
if sudo systemctl is-active --quiet skaffacityd; then
    print_success "Service is running!"
    
    # Show service status
    echo ""
    print_status "Service Status:"
    sudo systemctl status skaffacityd --no-pager -l
    
    echo ""
    print_status "Recent logs:"
    sudo journalctl -u skaffacityd -n 20 --no-pager
    
else
    print_error "Service failed to start!"
    print_status "Checking logs for errors..."
    sudo journalctl -u skaffacityd -n 50 --no-pager
    exit 1
fi

# Step 15: Display final information
echo ""
print_success "🎉 SkaffaCity Blockchain deployment completed!"
echo ""
echo "=========================================="
echo "📋 DEPLOYMENT SUMMARY"
echo "=========================================="
echo "🔗 Chain ID: $CHAIN_ID"
echo "🏠 Home Directory: $HOME_DIR"
echo "👤 Validator Address: $VALIDATOR_ADDRESS"
if [ -n "$FEE_DISTRIBUTION_DEV_ADDRESS" ]; then
echo "💰 Developer Address: $FEE_DISTRIBUTION_DEV_ADDRESS"
fi
echo ""
echo "🌐 Network Endpoints:"
echo "  RPC: http://$(curl -s ifconfig.me):26657"
echo "  API: http://$(curl -s ifconfig.me):1317"
echo "  gRPC: $(curl -s ifconfig.me):9090"
echo ""
echo "🔧 Management Commands:"
echo "  Status: sudo systemctl status skaffacityd"
echo "  Logs: sudo journalctl -u skaffacityd -f"
echo "  Restart: sudo systemctl restart skaffacityd"
echo ""
echo "💡 Next Steps:"
echo "  1. Monitor logs: sudo journalctl -u skaffacityd -f"
echo "  2. Test endpoints: curl http://localhost:26657/status"
echo "  3. Check validator: skaffacityd status --home $HOME_DIR"
if [ -n "$FEE_DISTRIBUTION_DEV_ADDRESS" ]; then
echo "  4. Monitor earnings: Check developer address balance"
fi
echo "=========================================="
print_success "Ready to receive transactions! 🚀"
