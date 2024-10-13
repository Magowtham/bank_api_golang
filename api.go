package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type apiHandler func(http.ResponseWriter, *http.Request) error

type apiError struct {
	error string
}

type APIServer struct {
	listenAddr string
	storage    Storage
}

// constructor
func NewAPIServer(listenAddr string, storage Storage) *APIServer {
	return &APIServer{
		listenAddr,
		storage,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/init", makeHTTPAPIHandler(s.handleInitDB)).Methods("GET")
	router.HandleFunc("/account", makeHTTPAPIHandler(s.handleCreateAccount)).Methods("POST")
	router.HandleFunc("/accounts", makeHTTPAPIHandler(s.handleGetAccounts)).Methods("GET")
	router.HandleFunc("/account/{id}", makeHTTPAPIHandler(s.handleGetAccount)).Methods("GET")
	router.HandleFunc("/account/{id}", makeHTTPAPIHandler(s.handleUpdateAccount)).Methods("PUT")
	router.HandleFunc("/account/{id}", makeHTTPAPIHandler(s.handleDeleteAccount)).Methods("DELETE")

	log.Printf("http server is listening on -> %s", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
}

func writeJson(w http.ResponseWriter, status int, message any) error {
	fmt.Println(message)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(message)
}

// currying
func makeHTTPAPIHandler(handler apiHandler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if error := handler(w, r); error != nil {
			writeJson(w, http.StatusBadRequest, apiError{error: error.Error()})
		}
	}
}

func (s *APIServer) handleInitDB(w http.ResponseWriter, r *http.Request) error {
	error := s.storage.InitDB()

	if error != nil {
		return error
	}

	return writeJson(w, http.StatusOK, "database successfully intialized")
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	var accountRequest AccountRequest
	error := json.NewDecoder(r.Body).Decode(&accountRequest)

	if error != nil {
		return error
	}

	account := NewAccount(
		accountRequest.FirstName,
		accountRequest.LastName,
		accountRequest.Email,
		accountRequest.PhoneNumber,
	)

	error = s.storage.CreateAccount(account)

	if error != nil {
		return error
	}

	return writeJson(w, http.StatusOK, "account created successfully")
}

func (s *APIServer) handleUpdateAccount(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)

	id_str := vars["id"]

	id_int, error := strconv.Atoi(id_str)

	if error != nil {
		return nil
	}

	var accountUpdateRequest AccountRequest

	error = json.NewDecoder(r.Body).Decode(&accountUpdateRequest)

	if error != nil {
		return error
	}

	error = s.storage.UpdateAccount(
		id_int,
		accountUpdateRequest.FirstName,
		accountUpdateRequest.LastName,
		accountUpdateRequest.Email,
		accountUpdateRequest.PhoneNumber,
	)

	if error != nil {
		return error
	}

	return writeJson(w, http.StatusOK, "account updated successfully")
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)

	id_str := vars["id"]

	id_int, error := strconv.Atoi(id_str)

	if error != nil {
		return error
	}

	error = s.storage.DeleteAccountByID(id_int)

	if error != nil {
		return error
	}

	return writeJson(w, http.StatusOK, "account deleted successfully")
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id_str := vars["id"]
	id_int, error := strconv.Atoi(id_str)

	if error != nil {
		return error
	}

	account, error := s.storage.GetAccountByID(id_int)

	if error != nil {
		return error
	}

	return writeJson(w, http.StatusOK, account)
}

func (s *APIServer) handleGetAccounts(w http.ResponseWriter, _ *http.Request) error {
	accounts, error := s.storage.GetAllAccounts()

	if error != nil {
		return error
	}

	return writeJson(w, http.StatusOK, accounts)
}
