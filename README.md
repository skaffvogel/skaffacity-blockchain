# SkaffaCity Blockchain VPS Deployment

Complete production-ready deployment package for SkaffaCity blockchain with integrated fee distribution system.

## 🚀 Quick Deployment

### Prerequisites
- Ubuntu 20.04+ VPS
- Minimum 4GB RAM, 2 CPU cores
- 50GB+ storage
- Root or sudo access

### One-Command Deployment

```bash
# Clone the repository
git clone https://github.com/skaffvogel/skaffacity-blockchain.git
cd skaffacity-blockchain

# Make deployment script executable
chmod +x deploy-vps.sh

# Run deployment (update developer address in script first!)
./deploy-vps.sh
```

## 📋 What Gets Deployed

### Blockchain Components
- **SkaffaCity Binary**: Full blockchain node with all modules
- **Module Handler**: Centralized module loading with detailed logging
- **Fee Distribution**: Automatic 90/10 fee split system
- **Custom Modules**: NFT, Marketplace, Governance, Staking, Web

### Infrastructure
- **Systemd Service**: Auto-restart blockchain on boot/crash
- **Firewall Configuration**: Secure ports (26656, 26657, 1317, 9090)
- **Monitoring**: Service health checks and logging
- **Account Management**: Validator account creation

### Fee Distribution System
- **Automatic Distribution**: Every block distributes fees
- **90% to Validators**: Keeps network incentivized
- **10% to Developer**: Your revenue stream
- **Configurable**: Change addresses and percentages
- **Event Logging**: Track all fee distributions

## 💰 Fee Distribution Configuration

### Update Developer Address
```bash
# Edit deploy-vps.sh before deployment
FEE_DISTRIBUTION_DEV_ADDRESS="skaffa1your-actual-address-here"

# Or update after deployment
./manage-fees.sh set-address skaffa1your-actual-address-here
```

### Enable Fee Distribution
```bash
./manage-fees.sh enable
```

### Monitor Earnings
```bash
# Check current status
./manage-fees.sh status

# View earnings
./manage-fees.sh earnings

# Real-time monitoring
./manage-fees.sh monitor
```

## 🌐 Network Endpoints

After deployment, your blockchain will be accessible at:

- **RPC**: `http://YOUR_VPS_IP:26657`
- **API**: `http://YOUR_VPS_IP:1317`
- **gRPC**: `YOUR_VPS_IP:9090`
- **P2P**: `YOUR_VPS_IP:26656`

## 🔧 Management Commands

### Service Management
```bash
# Check status
sudo systemctl status skaffacity

# View logs
sudo journalctl -u skaffacity -f

# Restart service
sudo systemctl restart skaffacity

# Stop service
sudo systemctl stop skaffacity
```

### Blockchain Operations
```bash
# Check node status
skaffacityd status --home ~/.skaffacity

# View accounts
skaffacityd keys list --home ~/.skaffacity

# Check balances
skaffacityd query bank balances $(skaffacityd keys show validator -a --home ~/.skaffacity) --home ~/.skaffacity
```

### Fee Distribution Management
```bash
# Show current configuration
./manage-fees.sh status

# Update developer address
./manage-fees.sh set-address skaffa1your-address

# Enable/disable fee distribution
./manage-fees.sh enable
./manage-fees.sh disable

# Monitor earnings in real-time
./manage-fees.sh monitor
```

## 📊 Revenue Model

### Fee Structure
- **Transaction Fees**: Set by users, collected automatically
- **Developer Share**: 10% of all transaction fees
- **Validator Share**: 90% of all transaction fees
- **Distribution**: Automatic every block

### Earnings Calculation
```
Daily Earnings = (Daily Transactions × Average Fee × 10%)
Monthly Earnings = Daily Earnings × 30
```

### Example Revenue
- 1000 transactions/day × 0.001 tokens fee × 10% = 0.1 tokens/day
- At $1/token = $0.10/day = $3/month
- Scale with network growth!

## 🔒 Security Features

### Network Security
- **Firewall**: UFW enabled with specific ports
- **Fail2ban**: Protection against brute force attacks
- **User Privileges**: Non-root service execution

### Blockchain Security
- **Validator Keys**: Secure key generation and storage
- **Address Validation**: Validates all fee distribution addresses
- **Error Handling**: Graceful failure handling for fee distribution

## 📈 Monitoring & Analytics

### Built-in Monitoring
- **Module Loading**: Detailed startup logging with [MODULE] prefix
- **Fee Distribution**: Event logging for all fee transfers
- **Service Health**: Systemd status monitoring
- **Network Stats**: Block height, transaction count

### Log Locations
- **Blockchain Logs**: `sudo journalctl -u skaffacity -f`
- **Module Loading**: Visible during startup
- **Fee Distribution**: In blockchain event logs

## 🚀 Scaling & Growth

### Network Growth Strategy
1. **Deploy Blockchain**: Get your node running
2. **Enable Fee Distribution**: Start earning from day 1
3. **Promote Usage**: Drive transaction volume
4. **Monitor Earnings**: Track revenue growth
5. **Reinvest**: Scale infrastructure as needed

### Performance Optimization
- **Multiple Validators**: Add more validators as network grows
- **Load Balancing**: Scale API endpoints
- **Database Optimization**: Tune for high transaction volume

## 🛠️ Troubleshooting

### Common Issues

#### Service Won't Start
```bash
# Check logs
sudo journalctl -u skaffacity -f

# Verify binary
which skaffacityd

# Check permissions
ls -la /usr/local/bin/skaffacityd
```

#### Fee Distribution Not Working
```bash
# Check configuration
./manage-fees.sh status

# Verify developer address
./manage-fees.sh earnings

# Enable if disabled
./manage-fees.sh enable
```

#### Network Connectivity
```bash
# Check ports
sudo ufw status

# Test endpoints
curl http://localhost:26657/status
curl http://localhost:1317/cosmos/base/tendermint/v1beta1/node_info
```

### Support Resources
- **GitHub Issues**: Report bugs and feature requests
- **Documentation**: Check module-specific docs
- **Community**: Join SkaffaCity Discord/Telegram

## 📝 File Structure

```
skaffacity-deployment/
├── blockchain/                 # Complete blockchain source
│   ├── app/                   # Application logic & module handler
│   ├── x/                     # Custom modules (NFT, Marketplace, etc.)
│   ├── cmd/                   # Binary commands
│   └── Makefile              # Build configuration
├── deploy-vps.sh             # Main deployment script
├── manage-fees.sh            # Fee distribution management
└── README.md                 # This file
```

## 🎯 Next Steps

1. **Deploy**: Run the deployment script on your VPS
2. **Configure**: Set your developer address for fee collection
3. **Test**: Send transactions and verify fee distribution
4. **Monitor**: Use monitoring tools to track earnings
5. **Scale**: Promote your blockchain and grow transaction volume

## 💡 Pro Tips

- **Backup Keys**: Always backup your validator keys
- **Monitor Regularly**: Check node health and earnings daily
- **Promote Usage**: Drive transaction volume to increase earnings
- **Stay Updated**: Keep blockchain software updated
- **Plan for Growth**: Scale infrastructure as network grows

---

🏙️ **SkaffaCity Blockchain** - Building the future of decentralized cities with automatic revenue generation!

💰 **Start earning transaction fees today!**
