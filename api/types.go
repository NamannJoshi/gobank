package api

import (
	"math/rand/v2"
	"time"
)

type CreateAccountRequest struct {
	FirstName string	`json:"firstName"`
	LastName string `json:"lastName"`
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