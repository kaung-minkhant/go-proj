package main

import "github.com/google/uuid"

type CreateTransferRequestParams struct {
	ToAccount uuid.UUID `json:"to_account"`
	Amount    int64     `json:"amount"`
}
