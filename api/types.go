package api

import (
	"math/rand/v2"
	"time"
)

type CreateAccountRequest struct {
	FirstName string	`json:"firstName"`
	LastName string `json:"lastName"`
}

type UpdateAccountRequest struct {
	FirstName *string	`json:"firstName"`
	LastName *string `json:"lastName"`
	Number *int `json:"number"`
	Balance *int `json:"balance"`
}

type TransferRequest struct {
	ToAccount int `json:"toAccount"`
	Amount int `json:"amount"`
}

type Account struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Number    int    `json:"number"`
	Balance   int    `json:"balance"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewAccount(firstName, lastName string) *Account {
	return &Account{
		ID:        rand.IntN(1000),
		FirstName: firstName,
		LastName:  lastName,
		Number:    rand.IntN(1000000),
		CreatedAt: time.Now().UTC(),
	}
}

func UpdateAccount(firstName *string, lastName *string, number *int, balance *int) (*Account) {
	account := &Account{}

	if firstName != nil {
		account.FirstName = *firstName
	}
	if lastName != nil {
		account.LastName = *lastName
	}
	if number != nil {
		account.Number = *number
	}
	if balance != nil {
		account.Balance = *balance
	}
	return account
}