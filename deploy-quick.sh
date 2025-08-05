#!/bin/bash

# SkaffaCity Blockchain - Quick Fix Script
# Use this when you only need to fix configuration issues

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

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_header() {
    echo -e "${BLUE}================================${NC}"
    echo -e "${BLUE}🏙️  SkaffaCity Quick Fix${NC}"
    echo -e "${BLUE}================================${NC}"
}

HOME_DIR="$HOME/skaffacity"

print_header

# Stop service if running
print_status "Stopping any running services..."
sudo systemctl stop skaffacityd 2>/dev/null || true
sudo pkill skaffacityd 2>/dev/null || true
sleep 2

# Fix app.toml if missing
if [ ! -f "$HOME_DIR/config/app.toml" ]; then
    print_status "Creating missing app.toml..."
    cat > "$HOME_DIR/config/app.toml" << 'EOF'
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
else
    print_status "app.toml exists, skipping creation"
fi

# Fix systemd service if missing
if [ ! -f "/etc/systemd/system/skaffacityd.service" ]; then
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
else
    print_status "Systemd service exists, reloading..."
    sudo systemctl daemon-reload
fi

# Start the service
print_status "Starting SkaffaCity blockchain..."
sudo systemctl start skaffacityd
sleep 3

# Check status
if sudo systemctl is-active --quiet skaffacityd; then
    print_success "✅ SkaffaCity blockchain is running!"
    print_status "Recent logs:"
    sudo journalctl -u skaffacityd -n 10 --no-pager
else
    print_error "❌ Service failed to start"
    print_status "Error logs:"
    sudo journalctl -u skaffacityd -n 20 --no-pager
fi

echo ""
print_status "Commands to monitor:"
echo "  sudo systemctl status skaffacityd"
echo "  sudo journalctl -u skaffacityd -f"
echo "  curl http://localhost:26657/status"
