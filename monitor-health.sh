#!/bin/bash

# SkaffaCity Blockchain Health Monitor
# Monitors node health, fee distribution, and earnings

set -e

BINARY_NAME="skaffacityd"
HOME_DIR="$HOME/skaffacity"
SERVICE_NAME="skaffacity"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m'

print_header() {
    echo -e "${BLUE}$1${NC}"
    echo "$(printf '=%.0s' {1..60})"
}

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

check_service_health() {
    print_header "🔧 Service Health Check"
    
    # Check if service is running
    if sudo systemctl is-active --quiet $SERVICE_NAME; then
        print_success "SkaffaCity service is running"
        
        # Get service status details
        UPTIME=$(sudo systemctl show $SERVICE_NAME --property=ActiveEnterTimestamp --value)
        print_status "Service uptime: $UPTIME"
        
        # Check if process is consuming reasonable resources
        PID=$(sudo systemctl show $SERVICE_NAME --property=MainPID --value)
        if [ "$PID" != "0" ]; then
            MEMORY=$(ps -p $PID -o rss= 2>/dev/null | awk '{print $1/1024 " MB"}' 2>/dev/null || echo "Unknown")
            CPU=$(ps -p $PID -o %cpu= 2>/dev/null || echo "Unknown")
            print_status "Memory usage: $MEMORY"
            print_status "CPU usage: $CPU%"
        fi
    else
        print_error "SkaffaCity service is not running!"
        print_status "Service status:"
        sudo systemctl status $SERVICE_NAME --no-pager
        return 1
    fi
    
    echo ""
}

check_blockchain_sync() {
    print_header "⛓️ Blockchain Sync Status"
    
    # Check if node is responding
    if ! curl -s localhost:26657/status >/dev/null 2>&1; then
        print_error "Node RPC not responding on port 26657"
        return 1
    fi
    
    # Get sync status
    STATUS=$($BINARY_NAME status --home $HOME_DIR 2>/dev/null || echo "{}")
    
    if [ "$STATUS" = "{}" ]; then
        print_error "Cannot get node status"
        return 1
    fi
    
    CATCHING_UP=$(echo $STATUS | jq -r '.SyncInfo.catching_up // true')
    LATEST_BLOCK=$(echo $STATUS | jq -r '.SyncInfo.latest_block_height // "0"')
    LATEST_TIME=$(echo $STATUS | jq -r '.SyncInfo.latest_block_time // "unknown"')
    
    if [ "$CATCHING_UP" = "false" ]; then
        print_success "Node is fully synchronized"
    else
        print_warning "Node is still catching up"
    fi
    
    print_status "Latest block height: $LATEST_BLOCK"
    print_status "Latest block time: $LATEST_TIME"
    
    # Check if blocks are being produced
    sleep 5
    NEW_STATUS=$($BINARY_NAME status --home $HOME_DIR 2>/dev/null || echo "{}")
    NEW_BLOCK=$(echo $NEW_STATUS | jq -r '.SyncInfo.latest_block_height // "0"')
    
    if [ "$NEW_BLOCK" -gt "$LATEST_BLOCK" ]; then
        print_success "Blocks are being produced (height increased from $LATEST_BLOCK to $NEW_BLOCK)"
    else
        print_warning "No new blocks in last 5 seconds"
    fi
    
    echo ""
}

check_fee_distribution() {
    print_header "💰 Fee Distribution Status"
    
    # Get fee distribution configuration
    CONFIG=$($BINARY_NAME query web web-config --home $HOME_DIR --output json 2>/dev/null || echo "{}")
    
    if [ "$CONFIG" = "{}" ]; then
        print_error "Cannot retrieve fee distribution configuration"
        return 1
    fi
    
    ENABLED=$(echo $CONFIG | jq -r '.web_config.fee_distribution.enabled // false')
    DEV_ADDR=$(echo $CONFIG | jq -r '.web_config.fee_distribution.developer_address // "not set"')
    DEV_PCT=$(echo $CONFIG | jq -r '.web_config.fee_distribution.developer_fee_percentage // "0"')
    
    if [ "$ENABLED" = "true" ]; then
        print_success "Fee distribution is enabled"
        print_status "Developer address: $DEV_ADDR"
        print_status "Developer fee percentage: $(echo "scale=2; $DEV_PCT / 100" | bc)%"
        
        # Check developer address balance if set
        if [ "$DEV_ADDR" != "not set" ] && [ "$DEV_ADDR" != "null" ]; then
            BALANCE=$($BINARY_NAME query bank balances $DEV_ADDR --home $HOME_DIR --output json 2>/dev/null || echo '{"balances":[]}')
            TOKENS=$(echo $BALANCE | jq -r '.balances[] | select(.denom=="token") | .amount // "0"')
            
            if [ "$TOKENS" != "0" ]; then
                TOKENS_FORMATTED=$(echo "scale=6; $TOKENS / 1000000" | bc)
                print_success "Developer balance: $TOKENS_FORMATTED tokens"
            else
                print_status "Developer balance: 0 tokens"
            fi
        else
            print_warning "Developer address not configured"
        fi
    else
        print_warning "Fee distribution is disabled"
    fi
    
    echo ""
}

check_network_activity() {
    print_header "📊 Network Activity"
    
    # Get recent block data
    LATEST_BLOCK=$($BINARY_NAME status --home $HOME_DIR 2>/dev/null | jq -r '.SyncInfo.latest_block_height // "0"')
    
    if [ "$LATEST_BLOCK" = "0" ]; then
        print_error "Cannot get latest block height"
        return 1
    fi
    
    # Check transactions in recent blocks
    TOTAL_TXS=0
    BLOCK_COUNT=10
    
    print_status "Checking last $BLOCK_COUNT blocks for transaction activity..."
    
    for i in $(seq 1 $BLOCK_COUNT); do
        BLOCK_HEIGHT=$((LATEST_BLOCK - i + 1))
        if [ $BLOCK_HEIGHT -gt 0 ]; then
            BLOCK_RESULT=$($BINARY_NAME query block $BLOCK_HEIGHT --home $HOME_DIR --output json 2>/dev/null || echo '{"block":{"data":{"txs":[]}}}')
            TX_COUNT=$(echo $BLOCK_RESULT | jq -r '.block.data.txs | length')
            TOTAL_TXS=$((TOTAL_TXS + TX_COUNT))
            
            if [ $TX_COUNT -gt 0 ]; then
                print_status "Block $BLOCK_HEIGHT: $TX_COUNT transactions"
            fi
        fi
    done
    
    if [ $TOTAL_TXS -gt 0 ]; then
        print_success "Total transactions in last $BLOCK_COUNT blocks: $TOTAL_TXS"
        AVG_TXS=$(echo "scale=2; $TOTAL_TXS / $BLOCK_COUNT" | bc)
        print_status "Average transactions per block: $AVG_TXS"
    else
        print_warning "No transactions found in recent blocks"
    fi
    
    echo ""
}

check_endpoints() {
    print_header "🌐 Network Endpoints"
    
    # Get public IP
    PUBLIC_IP=$(curl -s ifconfig.me 2>/dev/null || echo "unknown")
    print_status "Public IP: $PUBLIC_IP"
    
    # Check RPC endpoint
    if curl -s localhost:26657/status >/dev/null 2>&1; then
        print_success "RPC endpoint (26657) is responding"
        echo "  - Local: http://localhost:26657"
        echo "  - Public: http://$PUBLIC_IP:26657"
    else
        print_error "RPC endpoint (26657) is not responding"
    fi
    
    # Check API endpoint
    if curl -s localhost:1317/cosmos/base/tendermint/v1beta1/node_info >/dev/null 2>&1; then
        print_success "API endpoint (1317) is responding"
        echo "  - Local: http://localhost:1317"
        echo "  - Public: http://$PUBLIC_IP:1317"
    else
        print_error "API endpoint (1317) is not responding"
    fi
    
    # Check if ports are open
    print_status "Checking firewall status..."
    UFW_STATUS=$(sudo ufw status | grep -E "(26656|26657|1317|9090)" || echo "No rules found")
    echo "$UFW_STATUS"
    
    echo ""
}

generate_summary() {
    print_header "📋 Health Summary"
    
    # Overall health score
    HEALTH_SCORE=0
    MAX_SCORE=5
    
    # Check service
    if sudo systemctl is-active --quiet $SERVICE_NAME; then
        HEALTH_SCORE=$((HEALTH_SCORE + 1))
    fi
    
    # Check RPC
    if curl -s localhost:26657/status >/dev/null 2>&1; then
        HEALTH_SCORE=$((HEALTH_SCORE + 1))
    fi
    
    # Check sync status
    STATUS=$($BINARY_NAME status --home $HOME_DIR 2>/dev/null || echo "{}")
    CATCHING_UP=$(echo $STATUS | jq -r '.SyncInfo.catching_up // true')
    if [ "$CATCHING_UP" = "false" ]; then
        HEALTH_SCORE=$((HEALTH_SCORE + 1))
    fi
    
    # Check fee distribution
    CONFIG=$($BINARY_NAME query web web-config --home $HOME_DIR --output json 2>/dev/null || echo "{}")
    ENABLED=$(echo $CONFIG | jq -r '.web_config.fee_distribution.enabled // false')
    if [ "$ENABLED" = "true" ]; then
        HEALTH_SCORE=$((HEALTH_SCORE + 1))
    fi
    
    # Check API
    if curl -s localhost:1317/cosmos/base/tendermint/v1beta1/node_info >/dev/null 2>&1; then
        HEALTH_SCORE=$((HEALTH_SCORE + 1))
    fi
    
    # Display health score
    HEALTH_PERCENTAGE=$((HEALTH_SCORE * 100 / MAX_SCORE))
    
    if [ $HEALTH_PERCENTAGE -ge 80 ]; then
        print_success "Overall Health: $HEALTH_PERCENTAGE% ($HEALTH_SCORE/$MAX_SCORE) - Excellent! 🎉"
    elif [ $HEALTH_PERCENTAGE -ge 60 ]; then
        print_warning "Overall Health: $HEALTH_PERCENTAGE% ($HEALTH_SCORE/$MAX_SCORE) - Good, some issues to address"
    else
        print_error "Overall Health: $HEALTH_PERCENTAGE% ($HEALTH_SCORE/$MAX_SCORE) - Needs attention!"
    fi
    
    echo ""
    print_status "Recommendations:"
    
    if ! sudo systemctl is-active --quiet $SERVICE_NAME; then
        echo "  • Start the SkaffaCity service: sudo systemctl start $SERVICE_NAME"
    fi
    
    if [ "$CATCHING_UP" = "true" ]; then
        echo "  • Wait for blockchain to fully synchronize"
    fi
    
    if [ "$ENABLED" = "false" ]; then
        echo "  • Enable fee distribution: ./manage-fees.sh enable"
    fi
    
    DEV_ADDR=$(echo $CONFIG | jq -r '.web_config.fee_distribution.developer_address // "not set"')
    if [ "$DEV_ADDR" = "not set" ] || [ "$DEV_ADDR" = "null" ]; then
        echo "  • Set developer address: ./manage-fees.sh set-address <your-address>"
    fi
    
    echo ""
}

# Main execution
echo "🏙️ SkaffaCity Blockchain Health Monitor"
echo "========================================"
echo "Timestamp: $(date)"
echo ""

check_service_health
check_blockchain_sync
check_fee_distribution
check_network_activity
check_endpoints
generate_summary

print_success "Health check completed!"

# Save report to file
REPORT_FILE="health-report-$(date +%Y%m%d-%H%M%S).txt"
{
    echo "SkaffaCity Blockchain Health Report"
    echo "Generated: $(date)"
    echo "=================================="
    echo ""
} > $REPORT_FILE

# Rerun checks and append to file
{
    check_service_health 2>&1
    check_blockchain_sync 2>&1
    check_fee_distribution 2>&1
    check_network_activity 2>&1
    check_endpoints 2>&1
    generate_summary 2>&1
} >> $REPORT_FILE 2>&1

print_status "Health report saved to: $REPORT_FILE"
