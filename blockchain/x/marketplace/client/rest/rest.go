package rest

import (
    "fmt"
    "net/http"
    "encoding/json"

    "github.com/cosmos/cosmos-sdk/client"
    "github.com/gorilla/mux"
    
    "skaffacity/x/marketplace/types"
)

func RegisterRoutes(clientCtx client.Context, r *mux.Router) {
    r.HandleFunc("/marketplace/listings", listListingsHandler(clientCtx)).Methods("GET")
    r.HandleFunc("/marketplace/listings/{id}", getListingHandler(clientCtx)).Methods("GET")
    r.HandleFunc("/marketplace/listings/type/{type}", listListingsByTypeHandler(clientCtx)).Methods("GET")
    r.HandleFunc("/marketplace/listings/owner/{address}", listListingsByOwnerHandler(clientCtx)).Methods("GET")
    r.HandleFunc("/marketplace/stats", getMarketStatsHandler(clientCtx)).Methods("GET")
}

func writeErrorResponse(w http.ResponseWriter, status int, err string) {
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(map[string]string{"error": err})
}

func writeResponse(w http.ResponseWriter, clientCtx client.Context, res []byte) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(res)
}

func listListingsHandler(clientCtx client.Context) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        res, height, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/listings", types.QuerierRoute), nil)
        if err != nil {
            writeErrorResponse(w, http.StatusInternalServerError, err.Error())
            return
        }

        clientCtx = clientCtx.WithHeight(height)
        writeResponse(w, clientCtx, res)
    }
}

func getListingHandler(clientCtx client.Context) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        id := vars["id"]

        res, height, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/listing/%s", types.QuerierRoute, id), nil)
        if err != nil {
            writeErrorResponse(w, http.StatusNotFound, err.Error())
            return
        }

        clientCtx = clientCtx.WithHeight(height)
        writeResponse(w, clientCtx, res)
    }
}

func listListingsByTypeHandler(clientCtx client.Context) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        nftType := vars["type"]

        res, height, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/listings/type/%s", types.QuerierRoute, nftType), nil)
        if err != nil {
            writeErrorResponse(w, http.StatusInternalServerError, err.Error())
            return
        }

        clientCtx = clientCtx.WithHeight(height)
        writeResponse(w, clientCtx, res)
    }
}

func listListingsByOwnerHandler(clientCtx client.Context) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        address := vars["address"]

        res, height, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/listings/owner/%s", types.QuerierRoute, address), nil)
        if err != nil {
            writeErrorResponse(w, http.StatusNotFound, err.Error())
            return
        }

        clientCtx = clientCtx.WithHeight(height)
        writeResponse(w, clientCtx, res)
    }
}

func getMarketStatsHandler(clientCtx client.Context) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        res, height, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/stats", types.QuerierRoute), nil)
        if err != nil {
            writeErrorResponse(w, http.StatusInternalServerError, err.Error())
            return
        }

        clientCtx = clientCtx.WithHeight(height)
        writeResponse(w, clientCtx, res)
    }
}
