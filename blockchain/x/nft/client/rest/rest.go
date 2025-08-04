package rest

import (
    "fmt"
    "net/http"
    "encoding/json"

    "github.com/cosmos/cosmos-sdk/client"
    "github.com/gorilla/mux"
    
    "skaffacity/x/nft/types"
)

func RegisterRoutes(clientCtx client.Context, r *mux.Router) {
    r.HandleFunc("/nft/nfts", listNFTsHandler(clientCtx)).Methods("GET")
    r.HandleFunc("/nft/nfts/{id}", getNFTHandler(clientCtx)).Methods("GET")
    r.HandleFunc("/nft/land", listLandHandler(clientCtx)).Methods("GET")
    r.HandleFunc("/nft/items", listItemsHandler(clientCtx)).Methods("GET")
    r.HandleFunc("/nft/badges", listBadgesHandler(clientCtx)).Methods("GET")
    r.HandleFunc("/nft/owner/{address}", listNFTsByOwnerHandler(clientCtx)).Methods("GET")
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

func listNFTsHandler(clientCtx client.Context) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        res, height, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/nfts", types.QuerierRoute), nil)
        if err != nil {
            writeErrorResponse(w, http.StatusInternalServerError, err.Error())
            return
        }

        clientCtx = clientCtx.WithHeight(height)
        writeResponse(w, clientCtx, res)
    }
}

func getNFTHandler(clientCtx client.Context) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        id := vars["id"]

        res, height, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/nft/%s", types.QuerierRoute, id), nil)
        if err != nil {
            writeErrorResponse(w, http.StatusNotFound, err.Error())
            return
        }

        clientCtx = clientCtx.WithHeight(height)
        writeResponse(w, clientCtx, res)
    }
}

func listLandHandler(clientCtx client.Context) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        res, height, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/land", types.QuerierRoute), nil)
        if err != nil {
            writeErrorResponse(w, http.StatusInternalServerError, err.Error())
            return
        }

        clientCtx = clientCtx.WithHeight(height)
        writeResponse(w, clientCtx, res)
    }
}

func listItemsHandler(clientCtx client.Context) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        res, height, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/items", types.QuerierRoute), nil)
        if err != nil {
            writeErrorResponse(w, http.StatusInternalServerError, err.Error())
            return
        }

        clientCtx = clientCtx.WithHeight(height)
        writeResponse(w, clientCtx, res)
    }
}

func listBadgesHandler(clientCtx client.Context) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        res, height, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/badges", types.QuerierRoute), nil)
        if err != nil {
            writeErrorResponse(w, http.StatusInternalServerError, err.Error())
            return
        }

        clientCtx = clientCtx.WithHeight(height)
        writeResponse(w, clientCtx, res)
    }
}

func listNFTsByOwnerHandler(clientCtx client.Context) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        address := vars["address"]

        res, height, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/nfts/%s", types.QuerierRoute, address), nil)
        if err != nil {
            writeErrorResponse(w, http.StatusNotFound, err.Error())
            return
        }

        clientCtx = clientCtx.WithHeight(height)
        writeResponse(w, clientCtx, res)
    }
}
