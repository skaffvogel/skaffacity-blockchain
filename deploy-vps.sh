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
FEE_DISTRIBUTION_DEV_ADDRESS="skaffa1your-developer-address-here"  # UPDATE THIS!

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
if [ "$FEE_DISTRIBUTION_DEV_ADDRESS" = "skaffa1your-developer-address-here" ]; then
    print_warning "Please update FEE_DISTRIBUTION_DEV_ADDRESS in this script with your actual address"
    print_warning "You can set it later with: skaffacityd tx web set-developer-address <your-address>"
fi

print_status "Starting SkaffaCity blockchain deployment..."

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
make build
sudo cp bin/skaffacityd /usr/local/bin/
sudo chmod +x /usr/local/bin/skaffacityd
print_success "Binary installed to /usr/local/bin/skaffacityd"

# 5. Initialize chain
print_status "Initializing blockchain..."
$BINARY_NAME init $MONIKER --chain-id $CHAIN_ID --home $HOME_DIR

# 6. Create accounts
print_status "Creating validator account..."
$BINARY_NAME keys add validator --home $HOME_DIR

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
if [ "$FEE_DISTRIBUTION_DEV_ADDRESS" != "skaffa1your-developer-address-here" ]; then
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
    print_warning "Fee distribution disabled - update developer address later"
fi

# 11. Configure node
print_status "Configuring node settings..."
CONFIG_FILE="$HOME_DIR/config/config.toml"
APP_FILE="$HOME_DIR/config/app.toml"

# Update config.toml
sed -i 's/timeout_commit = "5s"/timeout_commit = "1s"/' $CONFIG_FILE
sed -i 's/cors_allowed_origins = \[\]/cors_allowed_origins = ["*"]/' $CONFIG_FILE

# Update app.toml
sed -i 's/enable = false/enable = true/' $APP_FILE
sed -i 's/swagger = false/swagger = true/' $APP_FILE

# 12. Setup systemd service
print_status "Creating systemd service..."
sudo tee /etc/systemd/system/$SERVICE_NAME.service > /dev/null <<EOF
[Unit]
Description=SkaffaCity Blockchain Node
After=network-online.target

[Service]
User=$USER
ExecStart=/usr/local/bin/$BINARY_NAME start --home $HOME_DIR
Restart=on-failure
RestartSec=3
LimitNOFILE=65535

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

# 14. Enable and start service
print_status "Starting SkaffaCity blockchain service..."
sudo systemctl daemon-reload
sudo systemctl enable $SERVICE_NAME
sudo systemctl start $SERVICE_NAME

# 15. Wait for service to start
print_status "Waiting for service to start..."
sleep 10

# 16. Check service status
if sudo systemctl is-active --quiet $SERVICE_NAME; then
    print_success "SkaffaCity blockchain service is running!"
else
    print_error "Service failed to start. Check logs with: sudo journalctl -u $SERVICE_NAME -f"
    exit 1
fi

# 17. Display important information
echo ""
echo "🎉 SkaffaCity Blockchain Deployment Complete!"
echo "============================================="
echo ""
echo "📋 Deployment Information:"
echo "- Chain ID: $CHAIN_ID"
echo "- Home Directory: $HOME_DIR"
echo "- Binary Location: /usr/local/bin/$BINARY_NAME"
echo "- Service Name: $SERVICE_NAME"
echo ""
echo "🌐 Network Endpoints:"
echo "- RPC: http://$(curl -s ifconfig.me):26657"
echo "- API: http://$(curl -s ifconfig.me):1317"
echo "- gRPC: $(curl -s ifconfig.me):9090"
echo ""
echo "💰 Fee Distribution:"
if [ "$FEE_DISTRIBUTION_DEV_ADDRESS" != "skaffa1your-developer-address-here" ]; then
    echo "- Status: ENABLED"
    echo "- Developer Address: $FEE_DISTRIBUTION_DEV_ADDRESS"
    echo "- Developer Fee: 10%"
    echo "- Validator Fee: 90%"
else
    echo "- Status: DISABLED"
    echo "- To enable: skaffacityd tx web set-developer-address <your-address>"
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
echo "- View keys: $BINARY_NAME keys list --home $HOME_DIR"
echo ""
echo "📈 Monitoring:"
echo "- Node status: $BINARY_NAME status --home $HOME_DIR"
echo "- Balance: $BINARY_NAME query bank balances $VALIDATOR_ADDR --home $HOME_DIR"
echo ""
echo "🚀 Your SkaffaCity blockchain is now running!"
echo "💼 Start earning fees from transactions automatically!"

# 18. Show next steps
echo ""
echo "📝 Next Steps:"
echo "1. Save your validator key: $BINARY_NAME keys export validator --home $HOME_DIR"
echo "2. Update developer address: $BINARY_NAME tx web set-developer-address <your-address>"
echo "3. Enable fee distribution: $BINARY_NAME tx web enable-fee-distribution true"
echo "4. Monitor your earnings and node health"
echo ""
print_success "Deployment completed successfully! 🎉"
