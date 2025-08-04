#!/bin/bash

# SkaffaCity Developer Address Generator
# Creates a developer address for fee distribution before deployment

set -e

echo "🔑 SkaffaCity Developer Address Generator"
echo "========================================"

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

BINARY_NAME="skaffacityd"

# Check if binary exists
if [ ! -f "./blockchain/bin/$BINARY_NAME" ]; then
    print_error "SkaffaCity binary not found!"
    print_status "Please build the blockchain first:"
    print_status "cd blockchain && make build"
    exit 1
fi

print_status "Creating developer address for fee distribution..."
echo ""

# Generate a new address
print_status "Generating new SkaffaCity address..."
ADDRESS_INFO=$(./blockchain/bin/$BINARY_NAME keys add developer --dry-run 2>/dev/null)

# Extract address from output
ADDRESS=$(echo "$ADDRESS_INFO" | grep "address:" | awk '{print $2}')
MNEMONIC=$(echo "$ADDRESS_INFO" | tail -n 1)

print_success "Developer address generated!"
echo ""
echo "📋 Your Developer Address Information:"
echo "======================================"
echo "Address: $ADDRESS"
echo ""
echo "🔐 Mnemonic Phrase (KEEP THIS SAFE!):"
echo "$MNEMONIC"
echo ""

print_warning "IMPORTANT SECURITY NOTES:"
echo "• Write down the mnemonic phrase above"
echo "• Store it in a secure location (offline)"
echo "• Never share your mnemonic with anyone"
echo "• This is the ONLY way to recover your account"
echo ""

print_status "To use this address in deployment:"
echo "1. Edit deploy-vps.sh"
echo "2. Set: FEE_DISTRIBUTION_DEV_ADDRESS=\"$ADDRESS\""
echo "3. Set: CREATE_DEV_ADDRESS_NOW=false"
echo "4. Run deployment: ./deploy-vps.sh"
echo ""

print_status "Or copy this configuration:"
echo "==============================="
echo "FEE_DISTRIBUTION_DEV_ADDRESS=\"$ADDRESS\""
echo "CREATE_DEV_ADDRESS_NOW=false"
echo "==============================="
echo ""

# Ask if user wants to update deploy-vps.sh automatically
read -p "Do you want to automatically update deploy-vps.sh with this address? (y/n): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    print_status "Updating deploy-vps.sh..."
    
    # Update the configuration in deploy-vps.sh
    sed -i "s/FEE_DISTRIBUTION_DEV_ADDRESS=\"\"/FEE_DISTRIBUTION_DEV_ADDRESS=\"$ADDRESS\"/" deploy-vps.sh
    sed -i "s/CREATE_DEV_ADDRESS_NOW=true/CREATE_DEV_ADDRESS_NOW=false/" deploy-vps.sh
    
    print_success "deploy-vps.sh updated successfully!"
    print_status "Your deployment is now configured with address: $ADDRESS"
    print_status "Run './deploy-vps.sh' to deploy with fee distribution enabled"
else
    print_status "Manual configuration required."
    print_status "Update deploy-vps.sh with the configuration shown above."
fi

echo ""
print_success "Developer address generation completed!"
print_warning "Remember to save your mnemonic phrase securely!"

# Save address info to file
cat > developer-address-info.txt <<EOF
SkaffaCity Developer Address Information
Generated: $(date)
========================================

Address: $ADDRESS

Mnemonic Phrase:
$MNEMONIC

Configuration for deploy-vps.sh:
FEE_DISTRIBUTION_DEV_ADDRESS="$ADDRESS"
CREATE_DEV_ADDRESS_NOW=false

SECURITY WARNING:
Keep this file secure and delete it after saving the mnemonic phrase elsewhere!
EOF

print_status "Address information saved to: developer-address-info.txt"
print_warning "Delete this file after saving the mnemonic phrase securely!"
echo ""
