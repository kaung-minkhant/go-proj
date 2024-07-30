package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (apiServer *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		return fmt.Errorf("method not supported")
	}
	var body AccountLoginRequestParams
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return fmt.Errorf("cannot parse credentials")
	}
	if uuidString := body.Number.String(); uuidString == "" {
		return fmt.Errorf("cannot parse credentials")
	}
	if body.Password == "" {
		return fmt.Errorf("invalid credentials")
	}
	account, err := apiServer.store.GetAccountByAccountNumber(body.Number)
	if err != nil {
		ErrorLog("Cannot find account with account number %v", err)
		return fmt.Errorf("cannot find accout with account number")
	}
	if err := ComparePassword(body.Password, account.EncryptedPassword); err != nil {
		return fmt.Errorf("invalid credentials")
	}

	token, err := makeJWTToken(account)
	if err != nil {
		return err
	}
	fmt.Println("JWT token", token)

	return WriteJSON(w, http.StatusAccepted, AccountLoginResponse{
		AccessToken: token,
	})
}

func (apiServer *APIServer) handleAccountAdmin(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return apiServer.handleGetAccounts(w, r)
	}
	if r.Method == "POST" {
		return apiServer.handleCreateAccount(w, r)
	}
	return fmt.Errorf("method not allowed %v", r.Method)
}

func (apiServer *APIServer) handlerAccountWithID(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return apiServer.handleGetAccountById(w, r)
	}
	if r.Method == "DELETE" {
		return apiServer.handleDeleteAccount(w, r)
	}
	return fmt.Errorf("method not allowed %v", r.Method)
}

func GetAccountIDFromUrl(r *http.Request) (uuid.UUID, error) {
	idString, ok := mux.Vars(r)["id"]
	if !ok {
		return uuid.UUID{}, fmt.Errorf("account id required")
	}
	id, err := uuid.Parse(idString)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("Account ID parsing failed with error %v", err)
	}
	return id, nil
}

func (apiServer *APIServer) handleGetAccountById(w http.ResponseWriter, r *http.Request) error {
	id, err := GetAccountIDFromUrl(r)
	if err != nil {
		return err
	}
	account, err := apiServer.store.GetAccountById(id)
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, account)
}

func (apiServer *APIServer) handleGetAccounts(w http.ResponseWriter, r *http.Request) error {
	accounts, err := apiServer.store.GetAccounts()
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, accounts)
}

func (apiServer *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	var body CreateAccountRequestParams
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return err
	}
	newAccount, err := NewAccount(body.FirstName, body.LastName, body.Password)
	if err != nil {
		ErrorLog("Error while making account %v", err)
		return err
	}

	createdAccount, err := apiServer.store.CreateAccount(newAccount)
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusCreated, createdAccount)
}

func (apiServer *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := GetAccountIDFromUrl(r)
	if err != nil {
		return err
	}
	err = apiServer.store.DeleteAccount(id)
	if err != nil {
		return fmt.Errorf("cannot delete account with error %v", err)
	}
	return WriteJSON(w, http.StatusAccepted, struct{}{})
}

func (apiServer *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		return fmt.Errorf("method not supported %s", r.Method)
	}
	var params CreateTransferRequestParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		return err
	}
	defer r.Body.Close()

	return WriteJSON(w, http.StatusOK, params)
}
