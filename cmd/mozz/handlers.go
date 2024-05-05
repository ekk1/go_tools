package main

import (
	"encoding/json"
	"fmt"
	"go_utils/utils"
	"io"
	"net/http"
)

type Account struct {
	Name           string `json:"name"`
	Group          string `json:"group"`
	InitialBalance int64  `json:"initial_balance"`
	Balance        int64  `json:"balance"`
	IsCredit       bool   `json:"is_credit"`
	BillingDate    string `json:"billing_date"`
	Records        []*Record
}

type Record struct {
	Amount      int64  `json:"amount"`
	Type        string `json:"type"`
	FromAccount string `json:"from_account"`
	ToAccount   string `json:"to_account"`
	Time        int64  `json:"time"`
}

func handleAccounts(w http.ResponseWriter, r *http.Request) {
	// if r.Header.Get("Content-Type") != "application/json" {
	// 	http.Error(w, "Invalid content type", http.StatusBadRequest)
	// 	return
	// }
	utils.LogPrintDebug(fmt.Sprintf("Request from %s: %s %s", r.RemoteAddr, r.Method, r.URL.Path))
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET, POST")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	switch r.Method {
	case http.MethodGet:
		accData, err := json.Marshal(Accounts)
		if err != nil {
			http.Error(w, "Failed to marshal accounts", http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.Write(accData)
	case http.MethodPost:
		acc := &Account{}
		reqBody, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			utils.LogPrintError(err)
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(reqBody, acc)
		if err != nil {
			http.Error(w, "Failed to unmarshal account", http.StatusBadRequest)
			return
		}
		if acc.Name == "" {
			http.Error(w, "Invalid account name", http.StatusBadRequest)
			return
		}
		for _, a := range Accounts {
			if a.Name == acc.Name {
				http.Error(w, "Account already exists", http.StatusBadRequest)
				return
			}
		}
		if acc.Group == "" {
			acc.Group = "default"
		}
		if acc.InitialBalance < 0 {
			http.Error(w, "Invalid initial balance", http.StatusBadRequest)
			return
		}
		if acc.IsCredit && acc.BillingDate == "" {
			http.Error(w, "Invalid billing date", http.StatusBadRequest)
			return
		}
		acc.Balance = acc.InitialBalance
		Accounts = append(Accounts, acc)
		kv.Set(AccountPrefix+acc.Name, string(reqBody))
		if err := kv.Save(); err != nil {
			http.Error(w, "Failed to save account", http.StatusInternalServerError)
			return
		}
		w.Write([]byte("OK"))
	case http.MethodOptions:
		w.WriteHeader(http.StatusOK)
		return
	default:
		http.Error(w, "Invalid method", http.StatusBadRequest)
	}

}
