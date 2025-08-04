#!/bin/bash

# SkaffaCity Blockchain VPS Deployment Script
# This script deploys SkaffaCity blockchain to a VPS with fee distribution system

set -e

echo "🏙️ SkaffaCity Blockchain VPS Deployment"
echo "========================================"

# Configuration
CHAIN_ID="skaffacity-1"
MONIKER="skaffacity-node"
HOME_DIR="$HOME/.skaffacity"
BINARY_NAME="skaffacityd"
SERVICE_NAME="skaffacity"
FEE_DISTRIBUTION_DEV_ADDRESS=""  # Leave empty - will be created during deployment

# Ask user if they want to create a developer address now
CREATE_DEV_ADDRESS_NOW=true  # Set to false if you have an existing address

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if running as root
if [ "$EUID" -eq 0 ]; then
    print_error "Please do not run this script as root"
    exit 1
fi

# Check if developer address is set
if [ "$CREATE_DEV_ADDRESS_NOW" = true ] && [ -z "$FEE_DISTRIBUTION_DEV_ADDRESS" ]; then
    print_status "Developer address not set. Will create one during deployment."
    print_status "You can also set an existing address later with: skaffacityd tx web set-developer-address <your-address>"
elif [ -n "$FEE_DISTRIBUTION_DEV_ADDRESS" ]; then
    print_status "Using provided developer address: $FEE_DISTRIBUTION_DEV_ADDRESS"
else
    print_warning "No developer address configured. Fee distribution will be disabled initially."
    print_warning "You can set it later with: skaffacityd tx web set-developer-address <your-address>"
fi

print_status "Starting SkaffaCity blockchain deployment..."

# Debug information
print_status "Debug Information:"
echo "- Current user: $(whoami)"
echo "- Home directory: $HOME"
echo "- Blockchain home: $HOME_DIR"
echo "- Binary location: $(which $BINARY_NAME 2>/dev/null || echo 'Not found in PATH')"
echo "- Binary version: $($BINARY_NAME version 2>/dev/null || echo 'Cannot get version')"
echo ""

# 1. Update system
print_status "Updating system packages..."
sudo apt update && sudo apt upgrade -y

# 2. Install dependencies
print_status "Installing dependencies..."
sudo apt install -y curl wget git build-essential jq ufw fail2ban

# 3. Install Go if not present
if ! command -v go &> /dev/null; then
    print_status "Installing Go..."
    wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
    sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
    echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
    export PATH=$PATH:/usr/local/go/bin
    rm go1.21.0.linux-amd64.tar.gz
    print_success "Go installed successfully"
else
    print_success "Go is already installed"
fi

# 4. Build blockchain binary
print_status "Building SkaffaCity blockchain..."
cd blockchain
make build
sudo cp bin/skaffacityd /usr/local/bin/
sudo chmod +x /usr/local/bin/skaffacityd
cd ..

# Verify binary installation
if ! command -v $BINARY_NAME &> /dev/null; then
    print_error "Binary installation failed - $BINARY_NAME not found in PATH"
    exit 1
fi

# Test binary
print_status "Testing binary functionality..."
if ! $BINARY_NAME version &> /dev/null; then
    print_error "Binary is not working correctly"
    exit 1
fi

# Check for any existing global config that might interfere
if [ -f "$HOME/.skaffacityd" ] || [ -d "$HOME/.skaffacityd" ]; then
    print_warning "Found existing skaffacityd config in home directory. This might cause conflicts."
    print_status "Moving existing config to backup..."
    mv "$HOME/.skaffacityd" "$HOME/.skaffacityd.backup.$(date +%s)" 2>/dev/null || true
fi

print_success "Binary installed and verified: /usr/local/bin/skaffacityd"

# 5. Initialize chain
print_status "Initializing blockchain..."

# Clean any existing configuration that might conflict
if [ -d "$HOME_DIR" ]; then
    print_warning "Existing blockchain configuration found. Backing up and cleaning..."
    mv "$HOME_DIR" "${HOME_DIR}.backup.$(date +%s)" 2>/dev/null || true
fi

# Ensure home directory exists with correct permissions
mkdir -p $HOME_DIR
chmod 700 $HOME_DIR

# Initialize with explicit home directory
print_status "Running: $BINARY_NAME init $MONIKER --chain-id $CHAIN_ID --home $HOME_DIR"
$BINARY_NAME init $MONIKER --chain-id $CHAIN_ID --home $HOME_DIR --overwrite

# Verify initialization was successful
if [ ! -f "$HOME_DIR/config/node_key.json" ]; then
    print_error "Blockchain initialization failed - node_key.json not found at $HOME_DIR/config/node_key.json"
    print_error "Directory contents:"
    ls -la "$HOME_DIR/config/" 2>/dev/null || print_error "Config directory does not exist"
    exit 1
fi

print_success "Blockchain initialized successfully at $HOME_DIR"

# 6. Create accounts
print_status "Creating validator account..."
print_status "Running: $BINARY_NAME keys add validator --home $HOME_DIR"
$BINARY_NAME keys add validator --home $HOME_DIR

# 6.1. Create developer account if requested
if [ "$CREATE_DEV_ADDRESS_NOW" = true ] && [ -z "$FEE_DISTRIBUTION_DEV_ADDRESS" ]; then
    print_status "Creating developer account for fee distribution..."
    print_status "Running: $BINARY_NAME keys add developer --home $HOME_DIR"
    $BINARY_NAME keys add developer --home $HOME_DIR
    FEE_DISTRIBUTION_DEV_ADDRESS=$($BINARY_NAME keys show developer -a --home $HOME_DIR)
    print_success "Developer account created: $FEE_DISTRIBUTION_DEV_ADDRESS"
    print_warning "IMPORTANT: Save your developer account mnemonic phrase!"
    print_warning "You can export it with: $BINARY_NAME keys export developer --home $HOME_DIR"
fi

# 7. Create genesis account
print_status "Setting up genesis..."
VALIDATOR_ADDR=$($BINARY_NAME keys show validator -a --home $HOME_DIR)
$BINARY_NAME genesis add-genesis-account $VALIDATOR_ADDR 1000000000000token --home $HOME_DIR

# 8. Create genesis transaction
print_status "Creating genesis transaction..."
$BINARY_NAME genesis gentx validator 100000000token --chain-id $CHAIN_ID --home $HOME_DIR

# 9. Collect genesis transactions
$BINARY_NAME genesis collect-gentxs --home $HOME_DIR

# 10. Configure fee distribution in genesis
print_status "Configuring fee distribution system..."
GENESIS_FILE="$HOME_DIR/config/genesis.json"

# Update genesis with fee distribution config (if developer address is set)
if [ -n "$FEE_DISTRIBUTION_DEV_ADDRESS" ]; then
    print_status "Setting up fee distribution with developer address: $FEE_DISTRIBUTION_DEV_ADDRESS"
    
    # Create temporary genesis with fee distribution
    jq --arg dev_addr "$FEE_DISTRIBUTION_DEV_ADDRESS" '
    .app_state.web.web_config.fee_distribution = {
        "enabled": true,
        "developer_address": $dev_addr,
        "developer_fee_percentage": "1000",
        "validator_fee_percentage": "9000"
    }' $GENESIS_FILE > ${GENESIS_FILE}.tmp && mv ${GENESIS_FILE}.tmp $GENESIS_FILE
    
    print_success "Fee distribution configured: 10% to developer, 90% to validators"
else
    print_warning "Fee distribution disabled - no developer address set"
    print_status "You can enable it later with:"
    print_status "1. Create address: $BINARY_NAME keys add developer --home $HOME_DIR"
    print_status "2. Set address: $BINARY_NAME tx web set-developer-address <address>"
    print_status "3. Enable: $BINARY_NAME tx web enable-fee-distribution true"
fi

# 11. Configure node
print_status "Configuring node settings..."
CONFIG_FILE="$HOME_DIR/config/config.toml"
APP_FILE="$HOME_DIR/config/app.toml"

# Update config.toml
sed -i 's/minimum-gas-prices = ""/minimum-gas-prices = "0.001token"/' $CONFIG_FILE
sed -i 's/enable = false/enable = true/' $CONFIG_FILE

# Update app.toml
sed -i 's/enable = false/enable = true/' $APP_FILE
sed -i 's/address = "localhost:9090"/address = "0.0.0.0:9090"/' $APP_FILE

# 12. Setup systemd service
print_status "Creating systemd service..."

# Ensure correct ownership of home directory
sudo chown -R $USER:$USER $HOME_DIR

sudo tee /etc/systemd/system/$SERVICE_NAME.service > /dev/null <<EOF
[Unit]
Description=SkaffaCity Blockchain Node
After=network-online.target

[Service]
User=$USER
Group=$USER
WorkingDirectory=$HOME_DIR
ExecStart=/usr/local/bin/$BINARY_NAME start --home $HOME_DIR
Restart=on-failure
RestartSec=3
LimitNOFILE=65535
Environment=HOME=/home/$USER

[Install]
WantedBy=multi-user.target
EOF

# 13. Setup firewall
print_status "Configuring firewall..."
sudo ufw allow 22/tcp     # SSH
sudo ufw allow 26656/tcp  # P2P
sudo ufw allow 26657/tcp  # RPC
sudo ufw allow 1317/tcp   # API
sudo ufw allow 9090/tcp   # gRPC
sudo ufw --force enable

# 14. Validate configuration before starting service
print_status "Validating blockchain configuration..."

# Check if all required files exist
required_files=(
    "$HOME_DIR/config/node_key.json"
    "$HOME_DIR/config/priv_validator_key.json"
    "$HOME_DIR/config/genesis.json"
    "$HOME_DIR/config/config.toml"
    "$HOME_DIR/config/app.toml"
)

for file in "${required_files[@]}"; do
    if [ ! -f "$file" ]; then
        print_error "Required configuration file missing: $file"
        exit 1
    fi
done

# Test if the blockchain can validate its configuration
print_status "Testing blockchain configuration..."
if ! $BINARY_NAME validate-genesis $HOME_DIR/config/genesis.json &> /dev/null; then
    print_error "Genesis file validation failed"
    exit 1
fi

print_success "Blockchain configuration validated successfully"

# 15. Enable and start service
print_status "Starting SkaffaCity blockchain service..."
sudo systemctl daemon-reload
sudo systemctl enable $SERVICE_NAME
sudo systemctl start $SERVICE_NAME

# 16. Wait for service to start
print_status "Waiting for service to start..."
sleep 10

# 17. Check service status
if sudo systemctl is-active --quiet $SERVICE_NAME; then
    print_success "SkaffaCity blockchain service is running!"
else
    print_error "Service failed to start. Check logs with: sudo journalctl -u $SERVICE_NAME -f"
    print_status "Showing recent logs:"
    sudo journalctl -u $SERVICE_NAME --no-pager -n 20
    exit 1
fi

# 18. Display deployment information
print_success "🎉 SkaffaCity Blockchain Deployment Complete!"
echo ""
echo "🌐 Your blockchain is accessible at:"
echo "- RPC: http://$(curl -s ifconfig.me):26657"
echo "- API: http://$(curl -s ifconfig.me):1317"
echo "- gRPC: $(curl -s ifconfig.me):9090"
echo ""
echo "💰 Fee Distribution:"
if [ -n "$FEE_DISTRIBUTION_DEV_ADDRESS" ]; then
    echo "- Status: ENABLED"
    echo "- Developer Address: $FEE_DISTRIBUTION_DEV_ADDRESS"
    echo "- Developer Fee: 10%"
    echo "- Validator Fee: 90%"
else
    echo "- Status: DISABLED"
    echo "- To enable:"
    echo "  1. Create developer account: $BINARY_NAME keys add developer --home $HOME_DIR"
    echo "  2. Set address: $BINARY_NAME tx web set-developer-address <address>"
    echo "  3. Enable distribution: $BINARY_NAME tx web enable-fee-distribution true"
fi
echo ""
echo "🔧 Useful Commands:"
echo "- Check status: sudo systemctl status $SERVICE_NAME"
echo "- View logs: sudo journalctl -u $SERVICE_NAME -f"
echo "- Restart: sudo systemctl restart $SERVICE_NAME"
echo "- Stop: sudo systemctl stop $SERVICE_NAME"
echo ""
echo "🔑 Account Information:"
echo "- Validator address: $VALIDATOR_ADDR"
if [ -n "$FEE_DISTRIBUTION_DEV_ADDRESS" ]; then
    echo "- Developer address: $FEE_DISTRIBUTION_DEV_ADDRESS"
fi
echo "- View keys: $BINARY_NAME keys list --home $HOME_DIR"
echo ""
echo "📈 Monitoring:"
echo "- Node status: $BINARY_NAME status --home $HOME_DIR"
echo "- Balance: $BINARY_NAME query bank balances $VALIDATOR_ADDR --home $HOME_DIR"
echo ""
echo "🚀 Your SkaffaCity blockchain is now running!"
echo "💼 Start earning fees from transactions automatically!"

# 19. Show next steps
echo ""
echo "📝 Next Steps:"
echo "1. Save your validator key: $BINARY_NAME keys export validator --home $HOME_DIR"
if [ -n "$FEE_DISTRIBUTION_DEV_ADDRESS" ]; then
    echo "2. Save your developer key: $BINARY_NAME keys export developer --home $HOME_DIR"
    echo "3. Fee distribution is already enabled and earning!"
    echo "4. Monitor your earnings: ./manage-fees.sh earnings"
else
    echo "2. Create developer address: $BINARY_NAME keys add developer --home $HOME_DIR"
    echo "3. Set developer address: $BINARY_NAME tx web set-developer-address <your-address>"
    echo "4. Enable fee distribution: $BINARY_NAME tx web enable-fee-distribution true"
fi
echo "5. Monitor your node health: ./monitor-health.sh"
echo ""
print_success "Deployment completed successfully! 🎉"
