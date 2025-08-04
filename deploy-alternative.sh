#!/bin/bash

# SkaffaCity Blockchain Alternative VPS Deployment Script
# This script uses a different approach to avoid home directory conflicts

set -e

echo "🏙️ SkaffaCity Blockchain Alternative VPS Deployment"
echo "=================================================="

# Configuration
CHAIN_ID="skaffacity-1"
MONIKER="skaffacity-node"
# Use the default path that the binary expects
HOME_DIR="/home/$(whoami)/skaffacity"  # Changed to avoid conflicts
BINARY_NAME="skaffacityd"
SERVICE_NAME="skaffacity"
FEE_DISTRIBUTION_DEV_ADDRESS=""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m'

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

print_status "Starting alternative SkaffaCity blockchain deployment..."

# Debug information
print_status "System Information:"
echo "- Current user: $(whoami)"
echo "- Working directory: $(pwd)"
echo "- Blockchain home will be: $HOME_DIR"
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
    print_error "Binary installation failed"
    exit 1
fi

print_success "Binary installed successfully"

# 5. Clean start - remove any existing configuration
print_status "Preparing clean blockchain environment..."
if [ -d "$HOME_DIR" ]; then
    print_warning "Removing existing blockchain configuration..."
    rm -rf "$HOME_DIR"
fi

# Create the directory structure manually
mkdir -p "$HOME_DIR"
mkdir -p "$HOME_DIR/config"
mkdir -p "$HOME_DIR/data"
chmod -R 700 "$HOME_DIR"

# 6. Try initialization without explicit home flag first
print_status "Attempting blockchain initialization..."

# Method 1: Let the binary use its default path, then move files
print_status "Method 1: Using binary default path..."
cd "$HOME_DIR"

# Initialize in current directory
if $BINARY_NAME init $MONIKER --chain-id $CHAIN_ID; then
    print_success "Initialization successful using default method"
else
    print_warning "Default method failed, trying alternative..."
    
    # Method 2: Try with explicit home
    if $BINARY_NAME init $MONIKER --chain-id $CHAIN_ID --home "$HOME_DIR"; then
        print_success "Initialization successful with explicit home"
    else
        print_error "Both initialization methods failed"
        print_status "Trying manual configuration setup..."
        
        # Method 3: Manual configuration creation
        print_status "Creating configuration files manually..."
        
        # Create basic config files
        cat > "$HOME_DIR/config/config.toml" << 'EOF'
# This is a TOML config file for SkaffaCity
proxy_app = "tcp://127.0.0.1:26658"
moniker = "skaffacity-node"

[rpc]
laddr = "tcp://0.0.0.0:26657"

[p2p]
laddr = "tcp://0.0.0.0:26656"
external_address = ""
EOF

        cat > "$HOME_DIR/config/app.toml" << 'EOF'
minimum-gas-prices = "0.001token"

[api]
enable = true
address = "tcp://0.0.0.0:1317"

[grpc]
enable = true
address = "0.0.0.0:9090"
EOF

        # Create a basic genesis file
        cat > "$HOME_DIR/config/genesis.json" << EOF
{
  "genesis_time": "$(date -u +%Y-%m-%dT%H:%M:%S.%6NZ)",
  "chain_id": "$CHAIN_ID",
  "initial_height": "1",
  "consensus_params": {
    "block": {
      "max_bytes": "22020096",
      "max_gas": "-1",
      "time_iota_ms": "1000"
    },
    "evidence": {
      "max_age_num_blocks": "100000",
      "max_age_duration": "172800000000000",
      "max_bytes": "1048576"
    },
    "validator": {
      "pub_key_types": [
        "ed25519"
      ]
    },
    "version": {}
  },
  "app_hash": "",
  "app_state": {
    "auth": {
      "params": {
        "max_memo_characters": "256",
        "tx_sig_limit": "7",
        "tx_size_cost_per_byte": "10",
        "sig_verify_cost_ed25519": "590",
        "sig_verify_cost_secp256k1": "1000"
      },
      "accounts": []
    },
    "bank": {
      "params": {
        "send_enabled": [],
        "default_send_enabled": true
      },
      "balances": [],
      "supply": [],
      "denom_metadata": []
    },
    "staking": {
      "params": {
        "unbonding_time": "1814400s",
        "max_validators": 100,
        "max_entries": 7,
        "historical_entries": 10000,
        "bond_denom": "token"
      },
      "last_total_power": "0",
      "last_validator_powers": [],
      "validators": [],
      "delegations": [],
      "unbonding_delegations": [],
      "redelegations": [],
      "exported": false
    }
  }
}
EOF

        # Generate node key manually
        $BINARY_NAME tendermint show-node-id --home "$HOME_DIR" 2>/dev/null || {
            print_status "Generating node keys manually..."
            openssl rand -hex 32 > "$HOME_DIR/config/node_key.json.tmp"
            echo '{"priv_key":{"type":"tendermint/PrivKeyEd25519","value":"'$(cat "$HOME_DIR/config/node_key.json.tmp")'"}}' > "$HOME_DIR/config/node_key.json"
            rm "$HOME_DIR/config/node_key.json.tmp"
        }
        
        print_success "Manual configuration created"
    fi
fi

cd - > /dev/null

# Verify we have the required files
required_files=("$HOME_DIR/config/config.toml" "$HOME_DIR/config/app.toml" "$HOME_DIR/config/genesis.json")
for file in "${required_files[@]}"; do
    if [ ! -f "$file" ]; then
        print_error "Required file missing: $file"
        exit 1
    fi
done

print_success "Blockchain configuration ready at $HOME_DIR"

# 7. Create validator account
print_status "Creating validator account..."
if $BINARY_NAME keys add validator --home "$HOME_DIR" --keyring-backend test; then
    VALIDATOR_ADDR=$($BINARY_NAME keys show validator -a --home "$HOME_DIR" --keyring-backend test)
    print_success "Validator account created: $VALIDATOR_ADDR"
else
    print_error "Failed to create validator account"
    exit 1
fi

# 8. Create developer account
print_status "Creating developer account..."
if $BINARY_NAME keys add developer --home "$HOME_DIR" --keyring-backend test; then
    FEE_DISTRIBUTION_DEV_ADDRESS=$($BINARY_NAME keys show developer -a --home "$HOME_DIR" --keyring-backend test)
    print_success "Developer account created: $FEE_DISTRIBUTION_DEV_ADDRESS"
else
    print_warning "Could not create developer account automatically"
fi

# 9. Setup systemd service with correct paths
print_status "Creating systemd service..."
sudo tee /etc/systemd/system/$SERVICE_NAME.service > /dev/null <<EOF
[Unit]
Description=SkaffaCity Blockchain Node
After=network-online.target

[Service]
Type=exec
User=$(whoami)
Group=$(whoami)
WorkingDirectory=$HOME_DIR
ExecStart=/usr/local/bin/$BINARY_NAME start --home $HOME_DIR
Restart=on-failure
RestartSec=3
LimitNOFILE=65535
Environment=HOME=$HOME_DIR

[Install]
WantedBy=multi-user.target
EOF

# 10. Setup firewall
print_status "Configuring firewall..."
sudo ufw allow 22/tcp
sudo ufw allow 26656/tcp
sudo ufw allow 26657/tcp
sudo ufw allow 1317/tcp
sudo ufw allow 9090/tcp
sudo ufw --force enable

# 11. Start the service
print_status "Starting blockchain service..."
sudo systemctl daemon-reload
sudo systemctl enable $SERVICE_NAME

# Try to start and check if it works
if sudo systemctl start $SERVICE_NAME; then
    sleep 5
    if sudo systemctl is-active --quiet $SERVICE_NAME; then
        print_success "🎉 SkaffaCity blockchain is running!"
    else
        print_warning "Service started but may have issues. Checking logs..."
        sudo journalctl -u $SERVICE_NAME --no-pager -n 10
    fi
else
    print_error "Failed to start service"
    sudo journalctl -u $SERVICE_NAME --no-pager -n 10
fi

# 12. Display information
echo ""
echo "🌐 Your blockchain should be accessible at:"
echo "- RPC: http://$(curl -s ifconfig.me 2>/dev/null || echo 'YOUR_IP'):26657"
echo "- API: http://$(curl -s ifconfig.me 2>/dev/null || echo 'YOUR_IP'):1317"
echo ""
echo "🔧 Management commands:"
echo "- Check status: sudo systemctl status $SERVICE_NAME"
echo "- View logs: sudo journalctl -u $SERVICE_NAME -f"
echo "- Restart: sudo systemctl restart $SERVICE_NAME"
echo ""
echo "🔑 Account information:"
echo "- Validator: $VALIDATOR_ADDR"
if [ -n "$FEE_DISTRIBUTION_DEV_ADDRESS" ]; then
    echo "- Developer: $FEE_DISTRIBUTION_DEV_ADDRESS"
fi
echo "- List keys: $BINARY_NAME keys list --home $HOME_DIR --keyring-backend test"
echo ""
print_success "Deployment completed!"
