package main

import (
	"math/rand/v2"
	"strconv"
	"time"
)

type AccountRequest struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}

type Account struct {
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	Email         string    `json:"email"`
	PhoneNumber   string    `json:"phone_number"`
	AccountNumber string    `json:"account_number"`
	CreatedAt     time.Time `json:"created_at"`
}

func NewAccount(firstName, lastName, email, phoneNumber string) *Account {
	return &Account{
		FirstName:     firstName,
		LastName:      lastName,
		Email:         email,
		PhoneNumber:   phoneNumber,
		AccountNumber: strconv.Itoa(rand.Int()),
		CreatedAt:     time.Now().UTC(),
	}
}
