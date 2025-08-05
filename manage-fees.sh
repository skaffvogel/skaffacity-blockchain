#!/bin/bash

# SkaffaCity Fee Distribution Management Script
# Manage your fee distribution settings easily

set -e

BINARY_NAME="skaffacityd"
HOME_DIR="$HOME/skaffacity"

# Colors
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

show_help() {
    echo "🏙️ SkaffaCity Fee Distribution Manager"
    echo "======================================"
    echo ""
    echo "Usage: $0 [COMMAND] [OPTIONS]"
    echo ""
    echo "Commands:"
    echo "  status                    - Show current fee distribution configuration"
    echo "  enable                    - Enable fee distribution"
    echo "  disable                   - Disable fee distribution"
    echo "  set-address <address>     - Set developer address for fee collection"
    echo "  earnings                  - Show developer fee earnings"
    echo "  monitor                   - Monitor fee distribution in real-time"
    echo "  help                      - Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 status"
    echo "  $0 set-address skaffa1abc123..."
    echo "  $0 enable"
    echo "  $0 earnings"
    echo ""
}

get_fee_config() {
    $BINARY_NAME query web web-config --home $HOME_DIR --output json 2>/dev/null || echo "{}"
}

show_status() {
    print_status "Getting fee distribution configuration..."
    
    CONFIG=$(get_fee_config)
    
    if [ "$CONFIG" = "{}" ]; then
        print_error "Could not retrieve configuration. Is the node running?"
        return 1
    fi
    
    ENABLED=$(echo $CONFIG | jq -r '.web_config.fee_distribution.enabled // false')
    DEV_ADDR=$(echo $CONFIG | jq -r '.web_config.fee_distribution.developer_address // "not set"')
    DEV_PCT=$(echo $CONFIG | jq -r '.web_config.fee_distribution.developer_fee_percentage // "0"')
    VAL_PCT=$(echo $CONFIG | jq -r '.web_config.fee_distribution.validator_fee_percentage // "0"')
    
    echo ""
    echo "💰 Fee Distribution Status"
    echo "========================="
    echo "Status: $([ "$ENABLED" = "true" ] && echo -e "${GREEN}ENABLED${NC}" || echo -e "${RED}DISABLED${NC}")"
    echo "Developer Address: $DEV_ADDR"
    echo "Developer Fee: $(echo "scale=2; $DEV_PCT / 100" | bc)%"
    echo "Validator Fee: $(echo "scale=2; $VAL_PCT / 100" | bc)%"
    echo ""
    
    if [ "$ENABLED" = "true" ] && [ "$DEV_ADDR" != "not set" ]; then
        print_success "Fee distribution is active and earning!"
    elif [ "$ENABLED" = "true" ] && [ "$DEV_ADDR" = "not set" ]; then
        print_warning "Fee distribution enabled but no developer address set"
    else
        print_warning "Fee distribution is disabled"
    fi
}

set_developer_address() {
    local address=$1
    
    if [ -z "$address" ]; then
        print_error "Please provide a developer address"
        echo "Usage: $0 set-address <skaffa1...>"
        return 1
    fi
    
    # Validate address format
    if [[ ! $address =~ ^skaffa1[a-z0-9]{38}$ ]]; then
        print_error "Invalid address format. Must start with 'skaffa1' and be 45 characters long"
        return 1
    fi
    
    print_status "Setting developer address to: $address"
    
    # Create transaction
    $BINARY_NAME tx web set-developer-address $address \
        --from validator \
        --chain-id skaffacity-1 \
        --home $HOME_DIR \
        --gas auto \
        --gas-adjustment 1.3 \
        --fees 1000token \
        --yes
    
    print_success "Developer address updated successfully!"
    print_status "Waiting for transaction to be processed..."
    sleep 5
    show_status
}

enable_fee_distribution() {
    print_status "Enabling fee distribution..."
    
    $BINARY_NAME tx web enable-fee-distribution true \
        --from validator \
        --chain-id skaffacity-1 \
        --home $HOME_DIR \
        --gas auto \
        --gas-adjustment 1.3 \
        --fees 1000token \
        --yes
    
    print_success "Fee distribution enabled!"
    print_status "Waiting for transaction to be processed..."
    sleep 5
    show_status
}

disable_fee_distribution() {
    print_status "Disabling fee distribution..."
    
    $BINARY_NAME tx web enable-fee-distribution false \
        --from validator \
        --chain-id skaffacity-1 \
        --home $HOME_DIR \
        --gas auto \
        --gas-adjustment 1.3 \
        --fees 1000token \
        --yes
    
    print_success "Fee distribution disabled!"
    print_status "Waiting for transaction to be processed..."
    sleep 5
    show_status
}

show_earnings() {
    print_status "Checking developer fee earnings..."
    
    CONFIG=$(get_fee_config)
    DEV_ADDR=$(echo $CONFIG | jq -r '.web_config.fee_distribution.developer_address // ""')
    
    if [ -z "$DEV_ADDR" ] || [ "$DEV_ADDR" = "not set" ]; then
        print_error "No developer address configured"
        return 1
    fi
    
    print_status "Developer address: $DEV_ADDR"
    
    # Get balance
    BALANCE=$($BINARY_NAME query bank balances $DEV_ADDR --home $HOME_DIR --output json 2>/dev/null || echo '{"balances":[]}')
    
    echo ""
    echo "💰 Developer Fee Earnings"
    echo "========================"
    echo "Address: $DEV_ADDR"
    echo ""
    
    TOKENS=$(echo $BALANCE | jq -r '.balances[] | select(.denom=="token") | .amount // "0"')
    
    if [ "$TOKENS" != "0" ]; then
        echo "Token Balance: $(echo "scale=6; $TOKENS / 1000000" | bc) tokens"
        print_success "You have earned fees! 🎉"
    else
        echo "Token Balance: 0 tokens"
        print_status "No fees earned yet. Keep promoting your blockchain!"
    fi
    
    echo ""
    echo "📊 Estimated Daily Earnings:"
    echo "Based on current activity, track your progress over time"
}

monitor_fees() {
    print_status "Starting real-time fee distribution monitoring..."
    print_status "Press Ctrl+C to stop monitoring"
    echo ""
    
    while true; do
        # Get latest block
        LATEST_BLOCK=$($BINARY_NAME status --home $HOME_DIR 2>/dev/null | jq -r '.SyncInfo.latest_block_height // "0"')
        
        # Show current status
        clear
        echo "🏙️ SkaffaCity Fee Distribution Monitor"
        echo "======================================"
        echo "Latest Block: $LATEST_BLOCK"
        echo "Time: $(date)"
        echo ""
        
        show_status
        
        echo ""
        echo "📊 Recent Activity:"
        
        # Show recent transactions (last 10 blocks)
        for i in {1..5}; do
            BLOCK_HEIGHT=$((LATEST_BLOCK - i))
            if [ $BLOCK_HEIGHT -gt 0 ]; then
                BLOCK_RESULT=$($BINARY_NAME query block $BLOCK_HEIGHT --home $HOME_DIR 2>/dev/null || echo '{"block":{"data":{"txs":[]}}}')
                TX_COUNT=$(echo $BLOCK_RESULT | jq -r '.block.data.txs | length')
                if [ $TX_COUNT -gt 0 ]; then
                    echo "Block $BLOCK_HEIGHT: $TX_COUNT transactions"
                fi
            fi
        done
        
        echo ""
        echo "⏱️  Refreshing in 10 seconds... (Ctrl+C to stop)"
        sleep 10
    done
}

# Main command handling
case "$1" in
    "status")
        show_status
        ;;
    "enable")
        enable_fee_distribution
        ;;
    "disable")
        disable_fee_distribution
        ;;
    "set-address")
        set_developer_address "$2"
        ;;
    "earnings")
        show_earnings
        ;;
    "monitor")
        monitor_fees
        ;;
    "help"|"")
        show_help
        ;;
    *)
        print_error "Unknown command: $1"
        echo ""
        show_help
        exit 1
        ;;
esac
