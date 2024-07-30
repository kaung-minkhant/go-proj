package main

import (
	"time"

	"github.com/google/uuid"
)

type CreateAccountRequestParams struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
}

type AccountLoginResponse struct {
	AccessToken string `json:"access_token"`
}

type AccountLoginRequestParams struct {
	Number   uuid.UUID `json:"number"`
	Password string    `json:"password"`
}

type Account struct {
	ID                uuid.UUID `json:"id"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	EncryptedPassword string    `json:"-"`
	Number            uuid.UUID `json:"number"`
	Balance           int64     `json:"balance"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func NewAccount(firstName, lastName, password string) (*Account, error) {
	encryptedPass, err := EncryptPassword(password)
	if err != nil {
		return nil, err
	}
	return &Account{
		ID:                uuid.New(),
		FirstName:         firstName,
		LastName:          lastName,
		Number:            uuid.New(),
		EncryptedPassword: string(encryptedPass),
	}, nil
}
