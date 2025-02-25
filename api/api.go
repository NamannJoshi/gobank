package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type ApiFunc func(http.ResponseWriter, *http.Request) error
type ApiError struct {
	Error string
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func makeHTTPHandleFunc(f ApiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error:err.Error()})
		}
	}
}

type ApiServer struct {
	listenAddr string
	store Storage
}

func NewServerApi(listenAddr string, store Storage) *ApiServer {
	return &ApiServer{
		listenAddr: listenAddr,
		store: store,
	}
}

func (s *ApiServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleAccount))
	router.HandleFunc("/account/{id}", makeHTTPHandleFunc(s.handleGetAccountByID))

	http.ListenAndServe(s.listenAddr, router)
}

func (s *ApiServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAccounts(w,r)
	}
	if r.Method == "POST" {
		return s.handleCreateAccount(w,r)
	}
	if r.Method == "PUT" {
		return s.handleUpdateAccount(w, r);
	}
	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r);
	}
	return fmt.Errorf("method not allowed : %v", r.Method)
}
func (s *ApiServer) handleGetAccounts(w http.ResponseWriter, r *http.Request) error {
	// account := NewAccount("BB", "Gandak")
	// account1 := NewAccount("BB", "Gandak")
	// account2 := NewAccount("BB", "Gandak")

	// accounts := []*Account{account, account1, account2}

	// type AccountsResponse struct {
	// 	Id int `json:"id"`
  //   Data []*Account `json:"data"`
	// }

	// accountsResponse := AccountsResponse{Data: accounts}

	accounts, err := s.store.GetAllAccounts();
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, accounts)
}
func (s *ApiServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	createAccountReq := new(CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(createAccountReq); err != nil {
		return err
	}
	account := NewAccount(createAccountReq.FirstName, createAccountReq.LastName)
	if err := s.store.CreateAccount(account); err!=nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}
func (s *ApiServer) handleUpdateAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}
func (s *ApiServer) handleGetAccountByID(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	query := r.URL.Query()
	fmt.Println(vars)
	fmt.Println(query)
	return nil
}
func (s *ApiServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}