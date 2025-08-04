package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

const (
	defaultPort = ":8090"
	deploymentFee = 100000000 // 100 SKAF in microSKAF
)

type WebServer struct {
	router *mux.Router
}

type TokenDeployment struct {
	Name              string  `json:"name"`
	Symbol            string  `json:"symbol"`
	TotalSupply       string  `json:"total_supply"`
	Decimals          int     `json:"decimals"`
	InitialRecipient  string  `json:"initial_recipient"`
	Description       string  `json:"description"`
	LogoURL           string  `json:"logo_url"`
	Mintable          bool    `json:"mintable"`
	Burnable          bool    `json:"burnable"`
	TransferFee       float64 `json:"transfer_fee"`
	DeployerAddress   string  `json:"deployer_address"`
}

type BlockchainStats struct {
	Height          string `json:"height"`
	ChainID         string `json:"chain_id"`
	TotalTokens     int    `json:"total_tokens"`
	TotalSupply     string `json:"total_supply"`
	ActiveValidators int   `json:"active_validators"`
	BlockTime       string `json:"block_time"`
}

type DeployedToken struct {
	Denom       string    `json:"denom"`
	Name        string    `json:"name"`
	Symbol      string    `json:"symbol"`
	Supply      string    `json:"supply"`
	Deployer    string    `json:"deployer"`
	DeployedAt  time.Time `json:"deployed_at"`
	TxHash      string    `json:"tx_hash"`
}

func NewWebServer() *WebServer {
	ws := &WebServer{
		router: mux.NewRouter(),
	}
	ws.setupRoutes()
	return ws
}

func (ws *WebServer) setupRoutes() {
	// Static files and templates
	ws.router.HandleFunc("/", ws.handleDashboard).Methods("GET")
	ws.router.HandleFunc("/dashboard", ws.handleDashboard).Methods("GET")
	ws.router.HandleFunc("/explorer", ws.handleExplorer).Methods("GET")
	ws.router.HandleFunc("/deploy", ws.handleDeployPage).Methods("GET")
	ws.router.HandleFunc("/admin", ws.handleAdmin).Methods("GET")

	// API endpoints
	api := ws.router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/stats", ws.handleStats).Methods("GET")
	api.HandleFunc("/tokens", ws.handleTokensList).Methods("GET")
	api.HandleFunc("/deploy", ws.handleTokenDeploy).Methods("POST")
	api.HandleFunc("/balance/{address}/{denom}", ws.handleBalance).Methods("GET")
	api.HandleFunc("/transactions", ws.handleTransactions).Methods("GET")
	api.HandleFunc("/validators", ws.handleValidators).Methods("GET")
	
	// Blockchain proxy endpoints
	proxy := ws.router.PathPrefix("/api/blockchain").Subrouter()
	proxy.HandleFunc("/status", ws.proxyToBlockchain).Methods("GET")
	proxy.HandleFunc("/blocks/{height}", ws.proxyToBlockchain).Methods("GET")
	proxy.HandleFunc("/txs/{hash}", ws.proxyToBlockchain).Methods("GET")
}

func (ws *WebServer) handleDashboard(w http.ResponseWriter, r *http.Request) {
	tmpl := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>SkaffaCity Blockchain Dashboard</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            background: linear-gradient(135deg, #0d1117 0%, #161b22 50%, #21262d 100%);
            color: #c9d1d9;
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            min-height: 100vh;
        }
        .container { max-width: 1200px; margin: 0 auto; padding: 20px; }
        .header {
            text-align: center;
            padding: 40px 0;
            background: linear-gradient(90deg, #00d4ff, #8b5cf6);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
        }
        .header h1 { font-size: 3rem; font-weight: bold; }
        .header p { font-size: 1.2rem; margin-top: 10px; color: #8b949e; }
        .nav {
            display: flex;
            justify-content: center;
            gap: 20px;
            margin: 30px 0;
        }
        .nav-btn {
            background: linear-gradient(45deg, #00d4ff, #8b5cf6);
            color: white;
            padding: 12px 24px;
            text-decoration: none;
            border-radius: 8px;
            font-weight: bold;
            transition: transform 0.2s;
        }
        .nav-btn:hover { transform: translateY(-2px); box-shadow: 0 5px 15px rgba(0, 212, 255, 0.4); }
        .stats-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
            gap: 20px;
            margin: 40px 0;
        }
        .stat-card {
            background: rgba(33, 38, 45, 0.8);
            border: 1px solid #30363d;
            border-radius: 12px;
            padding: 20px;
            text-align: center;
            backdrop-filter: blur(10px);
        }
        .stat-value {
            font-size: 2rem;
            font-weight: bold;
            color: #00d4ff;
            margin-bottom: 8px;
        }
        .stat-label { color: #8b949e; }
        .features {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
            gap: 30px;
            margin: 40px 0;
        }
        .feature-card {
            background: rgba(33, 38, 45, 0.8);
            border: 1px solid #30363d;
            border-radius: 12px;
            padding: 30px;
            text-align: center;
        }
        .feature-icon {
            font-size: 3rem;
            margin-bottom: 15px;
            color: #8b5cf6;
        }
        .feature-title {
            font-size: 1.4rem;
            font-weight: bold;
            margin-bottom: 10px;
            color: #00d4ff;
        }
        .footer {
            text-align: center;
            margin-top: 60px;
            color: #8b949e;
            border-top: 1px solid #30363d;
            padding-top: 20px;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🏙️ SkaffaCity Blockchain</h1>
            <p>Gaming-Focused Blockchain with SKAF Tokens & NFTs</p>
        </div>

        <div class="nav">
            <a href="/dashboard" class="nav-btn">📊 Dashboard</a>
            <a href="/explorer" class="nav-btn">🔍 Explorer</a>
            <a href="/deploy" class="nav-btn">🚀 Deploy Token</a>
            <a href="/admin" class="nav-btn">⚙️ Admin</a>
        </div>

        <div class="stats-grid" id="stats">
            <div class="stat-card">
                <div class="stat-value" id="height">Loading...</div>
                <div class="stat-label">Block Height</div>
            </div>
            <div class="stat-card">
                <div class="stat-value" id="validators">Loading...</div>
                <div class="stat-label">Active Validators</div>
            </div>
            <div class="stat-card">
                <div class="stat-value" id="tokens">Loading...</div>
                <div class="stat-label">Custom Tokens</div>
            </div>
            <div class="stat-card">
                <div class="stat-value" id="supply">Loading...</div>
                <div class="stat-label">SKAF Supply</div>
            </div>
        </div>

        <div class="features">
            <div class="feature-card">
                <div class="feature-icon">🎮</div>
                <div class="feature-title">Gaming Focused</div>
                <p>Built specifically for gaming applications with 1 SKAF per block rewards and low fees.</p>
            </div>
            <div class="feature-card">
                <div class="feature-icon">🪙</div>
                <div class="feature-title">Token Factory</div>
                <p>Deploy your own custom tokens for just 100 SKAF with full features like minting and burning.</p>
            </div>
            <div class="feature-card">
                <div class="feature-icon">🎨</div>
                <div class="feature-title">NFT System</div>
                <p>Comprehensive NFT marketplace for cosmetics, achievements, and in-game assets.</p>
            </div>
            <div class="feature-card">
                <div class="feature-icon">⚡</div>
                <div class="feature-title">Fast & Secure</div>
                <p>Built on Cosmos SDK with instant finality and enterprise-grade security.</p>
            </div>
        </div>

        <div class="footer">
            <p>SkaffaCity Blockchain v1.0.0 | Powered by Cosmos SDK</p>
        </div>
    </div>

    <script>
        async function loadStats() {
            try {
                const response = await fetch('/api/stats');
                const stats = await response.json();
                
                document.getElementById('height').textContent = stats.height;
                document.getElementById('validators').textContent = stats.active_validators;
                document.getElementById('tokens').textContent = stats.total_tokens;
                document.getElementById('supply').textContent = stats.total_supply + ' SKAF';
            } catch (error) {
                console.error('Failed to load stats:', error);
            }
        }

        loadStats();
        setInterval(loadStats, 30000); // Update every 30 seconds
    </script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(tmpl))
}

func (ws *WebServer) handleDeployPage(w http.ResponseWriter, r *http.Request) {
	tmpl := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Deploy Token - SkaffaCity</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            background: linear-gradient(135deg, #0d1117 0%, #161b22 50%, #21262d 100%);
            color: #c9d1d9;
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            min-height: 100vh;
        }
        .container { max-width: 800px; margin: 0 auto; padding: 20px; }
        .header {
            text-align: center;
            padding: 40px 0;
            background: linear-gradient(90deg, #00d4ff, #8b5cf6);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
        }
        .nav {
            display: flex;
            justify-content: center;
            gap: 20px;
            margin: 30px 0;
        }
        .nav-btn {
            background: linear-gradient(45deg, #00d4ff, #8b5cf6);
            color: white;
            padding: 12px 24px;
            text-decoration: none;
            border-radius: 8px;
            font-weight: bold;
        }
        .form-container {
            background: rgba(33, 38, 45, 0.8);
            border: 1px solid #30363d;
            border-radius: 12px;
            padding: 30px;
            margin: 20px 0;
        }
        .form-group {
            margin-bottom: 20px;
        }
        .form-group label {
            display: block;
            margin-bottom: 5px;
            color: #00d4ff;
            font-weight: bold;
        }
        .form-group input, .form-group textarea, .form-group select {
            width: 100%;
            padding: 12px;
            background: rgba(13, 17, 23, 0.8);
            border: 1px solid #30363d;
            border-radius: 6px;
            color: #c9d1d9;
            font-size: 16px;
        }
        .checkbox-group {
            display: flex;
            align-items: center;
            gap: 10px;
        }
        .checkbox-group input[type="checkbox"] {
            width: auto;
        }
        .deploy-btn {
            width: 100%;
            padding: 15px;
            background: linear-gradient(45deg, #00d4ff, #8b5cf6);
            color: white;
            border: none;
            border-radius: 8px;
            font-size: 18px;
            font-weight: bold;
            cursor: pointer;
            transition: transform 0.2s;
        }
        .deploy-btn:hover {
            transform: translateY(-2px);
            box-shadow: 0 5px 15px rgba(0, 212, 255, 0.4);
        }
        .fee-info {
            background: rgba(139, 92, 246, 0.1);
            border: 1px solid #8b5cf6;
            border-radius: 8px;
            padding: 15px;
            margin-bottom: 20px;
            text-align: center;
        }
        .result {
            margin-top: 20px;
            padding: 15px;
            border-radius: 8px;
        }
        .success {
            background: rgba(34, 197, 94, 0.1);
            border: 1px solid #22c55e;
            color: #22c55e;
        }
        .error {
            background: rgba(239, 68, 68, 0.1);
            border: 1px solid #ef4444;
            color: #ef4444;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🚀 Deploy Custom Token</h1>
            <p>Create your own token on SkaffaCity Blockchain</p>
        </div>

        <div class="nav">
            <a href="/dashboard" class="nav-btn">📊 Dashboard</a>
            <a href="/explorer" class="nav-btn">🔍 Explorer</a>
            <a href="/deploy" class="nav-btn">🚀 Deploy Token</a>
            <a href="/admin" class="nav-btn">⚙️ Admin</a>
        </div>

        <div class="fee-info">
            <h3>💰 Deployment Fee: 100 SKAF</h3>
            <p>One-time fee to deploy your custom token with all features included</p>
        </div>

        <div class="form-container">
            <form id="deployForm">
                <div class="form-group">
                    <label>Token Name *</label>
                    <input type="text" name="name" placeholder="e.g. GameCoin" required>
                </div>

                <div class="form-group">
                    <label>Symbol *</label>
                    <input type="text" name="symbol" placeholder="e.g. GAME" required>
                </div>

                <div class="form-group">
                    <label>Total Supply *</label>
                    <input type="number" name="total_supply" placeholder="1000000" required>
                </div>

                <div class="form-group">
                    <label>Decimals</label>
                    <select name="decimals">
                        <option value="6">6 (Standard)</option>
                        <option value="18">18 (Ethereum Style)</option>
                        <option value="8">8 (Bitcoin Style)</option>
                        <option value="0">0 (Whole Numbers)</option>
                    </select>
                </div>

                <div class="form-group">
                    <label>Initial Recipient Address *</label>
                    <input type="text" name="initial_recipient" placeholder="skaffa1..." required>
                </div>

                <div class="form-group">
                    <label>Description</label>
                    <textarea name="description" placeholder="Describe your token's purpose and utility" rows="3"></textarea>
                </div>

                <div class="form-group">
                    <label>Logo URL</label>
                    <input type="url" name="logo_url" placeholder="https://example.com/logo.png">
                </div>

                <div class="form-group">
                    <div class="checkbox-group">
                        <input type="checkbox" name="mintable" id="mintable">
                        <label for="mintable">Mintable (Allow creating more tokens later)</label>
                    </div>
                </div>

                <div class="form-group">
                    <div class="checkbox-group">
                        <input type="checkbox" name="burnable" id="burnable">
                        <label for="burnable">Burnable (Allow destroying tokens)</label>
                    </div>
                </div>

                <div class="form-group">
                    <label>Transfer Fee % (to deployer)</label>
                    <input type="number" name="transfer_fee" step="0.01" min="0" max="10" placeholder="0.00">
                </div>

                <div class="form-group">
                    <label>Your Address (Deployer) *</label>
                    <input type="text" name="deployer_address" placeholder="skaffa1..." required>
                </div>

                <button type="submit" class="deploy-btn">Deploy Token for 100 SKAF</button>
            </form>

            <div id="result"></div>
        </div>
    </div>

    <script>
        document.getElementById('deployForm').addEventListener('submit', async (e) => {
            e.preventDefault();
            
            const formData = new FormData(e.target);
            const data = Object.fromEntries(formData.entries());
            
            // Convert checkboxes to booleans
            data.mintable = formData.has('mintable');
            data.burnable = formData.has('burnable');
            data.transfer_fee = parseFloat(data.transfer_fee) || 0;
            
            const resultDiv = document.getElementById('result');
            resultDiv.innerHTML = '<p>Deploying token...</p>';
            
            try {
                const response = await fetch('/api/deploy', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(data)
                });
                
                const result = await response.json();
                
                if (response.ok) {
                    resultDiv.className = 'result success';
                    resultDiv.innerHTML = ` + "`" + `
                        <h3>✅ Token Deployed Successfully!</h3>
                        <p><strong>Transaction Hash:</strong> ${result.tx_hash}</p>
                        <p><strong>Token Denom:</strong> ${result.denom}</p>
                        <p>Your token is now live on SkaffaCity Blockchain!</p>
                    ` + "`" + `;
                } else {
                    resultDiv.className = 'result error';
                    resultDiv.innerHTML = ` + "`" + `
                        <h3>❌ Deployment Failed</h3>
                        <p>${result.error}</p>
                    ` + "`" + `;
                }
            } catch (error) {
                resultDiv.className = 'result error';
                resultDiv.innerHTML = ` + "`" + `
                    <h3>❌ Error</h3>
                    <p>Failed to communicate with blockchain: ${error.message}</p>
                ` + "`" + `;
            }
        });
    </script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(tmpl))
}

func (ws *WebServer) handleStats(w http.ResponseWriter, r *http.Request) {
	// Mock data for now - replace with actual blockchain queries
	stats := BlockchainStats{
		Height:           "123456",
		ChainID:          "skaffacity-1",
		TotalTokens:      25,
		TotalSupply:      "1,000,000,000",
		ActiveValidators: 4,
		BlockTime:        "6.2s",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func (ws *WebServer) handleTokenDeploy(w http.ResponseWriter, r *http.Request) {
	var deployment TokenDeployment
	if err := json.NewDecoder(r.Body).Decode(&deployment); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if deployment.Name == "" || deployment.Symbol == "" || 
	   deployment.TotalSupply == "" || deployment.DeployerAddress == "" ||
	   deployment.InitialRecipient == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Here you would implement the actual token deployment
	// For now, return a mock response
	response := map[string]interface{}{
		"success": true,
		"tx_hash": "A1B2C3D4E5F6...",
		"denom":   fmt.Sprintf("factory/%s/%s", deployment.DeployerAddress, deployment.Symbol),
		"message": "Token deployed successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (ws *WebServer) handleExplorer(w http.ResponseWriter, r *http.Request) {
	// Block explorer implementation would go here
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("<h1>🔍 Block Explorer - Coming Soon</h1>"))
}

func (ws *WebServer) handleAdmin(w http.ResponseWriter, r *http.Request) {
	// Admin panel implementation would go here
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("<h1>⚙️ Admin Panel - Coming Soon</h1>"))
}

func (ws *WebServer) handleTokensList(w http.ResponseWriter, r *http.Request) {
	// Mock deployed tokens data
	tokens := []DeployedToken{
		{
			Denom:      "factory/skaffa1.../GAME",
			Name:       "GameCoin",
			Symbol:     "GAME",
			Supply:     "1,000,000",
			Deployer:   "skaffa1...",
			DeployedAt: time.Now().Add(-24 * time.Hour),
			TxHash:     "ABC123...",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokens)
}

func (ws *WebServer) handleBalance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]
	denom := vars["denom"]
	
	// Mock balance data
	balance := map[string]interface{}{
		"address": address,
		"denom":   denom,
		"amount":  "1000000",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(balance)
}

func (ws *WebServer) handleTransactions(w http.ResponseWriter, r *http.Request) {
	// Mock transaction data
	txs := []map[string]interface{}{
		{
			"hash":   "ABC123...",
			"height": "123456",
			"type":   "token_deploy",
			"fee":    "100000000uskaf",
			"time":   time.Now().Format(time.RFC3339),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(txs)
}

func (ws *WebServer) handleValidators(w http.ResponseWriter, r *http.Request) {
	// Mock validator data
	validators := []map[string]interface{}{
		{
			"operator_address": "skaffavaloper1...",
			"moniker":         "SkaffaCity-Validator-1",
			"status":          "BOND_STATUS_BONDED",
			"tokens":          "1000000000",
			"commission_rate": "0.05",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(validators)
}

func (ws *WebServer) proxyToBlockchain(w http.ResponseWriter, r *http.Request) {
	// Proxy requests to the actual blockchain RPC/API
	// This would forward requests to localhost:26657 or localhost:1317
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Blockchain proxy endpoint - implementation pending",
	})
}

func main() {
	fmt.Println("🏙️ Starting SkaffaCity Web Module...")
	fmt.Printf("🌐 Server starting on %s\n", defaultPort)
	fmt.Println("📊 Dashboard: http://localhost" + defaultPort)
	fmt.Println("🚀 Token Deployer: http://localhost" + defaultPort + "/deploy")
	fmt.Println("🔍 Block Explorer: http://localhost" + defaultPort + "/explorer")
	fmt.Println("⚙️ Admin Panel: http://localhost" + defaultPort + "/admin")

	server := NewWebServer()

	// Setup CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})

	handler := c.Handler(server.router)

	log.Fatal(http.ListenAndServe(defaultPort, handler))
}
