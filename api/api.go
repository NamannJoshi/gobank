package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

type ApiFunc func(http.ResponseWriter, *http.Request) error
type ApiError struct {
	Error string `json:"error"`
}
type ApiServer struct {
	listenAddr string
	store Storage
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
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

func createJWT(account *Account) (string, error) {
	claims := &jwt.MapClaims{
		"expiredAt": 15000,
		"accountNumber": account.Number,
	}
	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50TnVtYmVyIjo3NjA3MzIsImV4cGlyZWRBdCI6MTUwMDB9.V9rDh9oNUYrALshSK-_Pf8R1koc4mLsYWcwCrgrFIms

func permissionDenied(w http.ResponseWriter) error {
	return WriteJSON(w, http.StatusForbidden, ApiError{Error: "permission denied"})
}

func withJWTauth(handlerFunc http.HandlerFunc, s Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("calling JWT auth middleware")
		
		tokenString := r.Header.Get("x-jwt-token")
		token, err := validateJWT(tokenString)
		if err != nil {
			permissionDenied(w)
			return
		}
		if !token.Valid {
			permissionDenied(w)
			return
		}
		userId := mux.Vars(r)["id"]

		claims := token.Claims.(jwt.MapClaims)
		fmt.Println(claims)

		i, err := strconv.Atoi(userId)
		if err != nil {
			return 
		}
		account, err := s.GetAccountByID(i)
		if err != nil {
			return
		}
		if account.Number != int(claims["accountNumber"].(float64)) {
			permissionDenied(w)
			return
		}

		handlerFunc(w, r)
	}
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")	

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
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
	router.HandleFunc("/account/{id}", withJWTauth(makeHTTPHandleFunc(s.handleAccountById), s.store))

	// can use method to specify method in routes

	router.HandleFunc("/transfer", makeHTTPHandleFunc(s.handleTransfer)).Methods("GET")

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
	return fmt.Errorf("method not allowed : %v", r.Method)
}

func (s *ApiServer) handleAccountById(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAccountByID(w, r)
	}
	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}
	// if r.Method == "PUT" {
	// 	return s.handleUpdateAccount(w, r)
	// }
	return fmt.Errorf("this method not available on this API: %v", r.Method)
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

	tokenString, err := createJWT(account)
	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
	}

	fmt.Println(tokenString)

	return WriteJSON(w, http.StatusOK, account)
}

func (s *ApiServer) handleUpdateAccount(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)["id"]
	param ,err := strconv.Atoi(vars)
	if err != nil {
		log.Fatalf("error while var conversion in deletion: %v", err)
	}
	updateAccountReq := new(UpdateAccountRequest)

	if err := json.NewDecoder(r.Body).Decode(updateAccountReq); err != nil {
		log.Fatalf("error while decoding update req: %v", err)
	}

	account := UpdateAccount(updateAccountReq.FirstName, updateAccountReq.LastName, updateAccountReq.Number, updateAccountReq.Balance)

	errar := s.store.UpdateAccount(account, param)
	return WriteJSON(w, http.StatusOK, errar)
}

func(s *ApiServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	transferReq := new(TransferRequest)
	if err := json.NewDecoder(r.Body).Decode(transferReq); err != nil {
		return err
	}
	defer r.Body.Close()

	return WriteJSON(w, http.StatusOK, transferReq)
}

func (s *ApiServer) handleGetAccountByID(w http.ResponseWriter, r *http.Request) error {
	fmt.Println("Hey its working!")
	vars := mux.Vars(r)["id"]
	// query := r.URL.Query()
	// fmt.Println(query)
	
	param, err := strconv.Atoi(vars)
	if err != nil {
		log.Fatalf("error while converting var: %v", err)
	}
	
	account, errer := s.store.GetAccountByID(param)
	if errer != nil {
		errf := fmt.Sprintf("Account %d not found", param)
		return WriteJSON(w, http.StatusNotFound, ApiError{Error: errf})
	}
	
	return WriteJSON(w, http.StatusOK, account)
}
func (s *ApiServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)["id"]
	param ,err := strconv.Atoi(vars)
	if err != nil {
		log.Fatalf("error while var conversion in deletion: %v", err)
	}

	if errar := s.store.DeleteAccount(param); err != nil {
		return WriteJSON(w, http.StatusBadRequest, errar)
	}
	return WriteJSON(w, http.StatusOK, map[string]int{"deleted": param})

}