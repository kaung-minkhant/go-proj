package main

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) ([]Account, error)
	GetAccounts() ([]Account, error)
	DeleteAccount(uuid.UUID) error
	UpdateAccount(*Account) error
	GetAccountById(uuid.UUID) (*Account, error)
	GetAccountByAccountNumber(uuid.UUID) (*Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "postgres://root:root@localhost/gobank?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		// log.Fatal("Could not open DB with error", err)
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	var postgresStore PostgresStore
	postgresStore.db = db
	return &postgresStore, nil
}

func (store *PostgresStore) Init() error {
	return store.CreateAccountTable()
}
func (store *PostgresStore) CreateAccountTable() error {
	query := `CREATE TABLE IF NOT EXISTS accounts (
		id UUID PRIMARY KEY UNIQUE,
		first_name VARCHAR(50) NOT NULL,
		last_name VARCHAR(50),
		number UUID NOT NULL UNIQUE,
		balance BIGINT,
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NOT NULL,
		password TEXT NOT NULL
	);`
	_, err := store.db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (store *PostgresStore) GetAccounts() ([]Account, error) {
	query := `
		SELECT * FROM accounts;
	`
	rows, err := store.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var accounts []Account
	for rows.Next() {
		account, err := scanIntoAccount(rows)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, *account)
	}
	return accounts, nil
}
func (store *PostgresStore) GetAccountById(id uuid.UUID) (*Account, error) {
	query := `
		SELECT * FROM accounts
		WHERE id = $1;	
	`
	rows, err := store.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		account, err := scanIntoAccount(rows)
		if err != nil {
			return nil, err
		}
		return account, nil
		// accounts = append(accounts, *account)
	}
	return nil, fmt.Errorf("cannot find account")
}
func (store *PostgresStore) CreateAccount(account *Account) ([]Account, error) {
	query := `
		INSERT INTO accounts
		(id, first_name, last_name, number, balance, created_at, updated_at, password)
		VALUES ($1, $2, $3, $4, $5, NOW() AT TIME ZONE 'utc', NOW() AT TIME ZONE 'utc', $6)
		RETURNING *;
	`
	rows, err := store.db.Query(query, account.ID, account.FirstName, account.LastName, account.Number, account.Balance, account.EncryptedPassword)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var createdAccounts = []Account{}
	if rows.Next() {
		account, err := scanIntoAccount(rows)
		if err != nil {
			return nil, err
		}
		createdAccounts = append(createdAccounts, *account)
	}
	return createdAccounts, nil
}
func (store *PostgresStore) DeleteAccount(id uuid.UUID) error {
	query := `
		DELETE FROM accounts
		WHERE id = $1;
	`
	_, err := store.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
func (store *PostgresStore) UpdateAccount(account *Account) error {
	return nil
}

func scanIntoAccount(rows *sql.Rows) (*Account, error) {
	var account = new(Account)
	if err := rows.Scan(&account.ID, &account.FirstName, &account.LastName, &account.Number, &account.Balance, &account.CreatedAt, &account.UpdatedAt, &account.EncryptedPassword); err != nil {
		return nil, err
	}
	return account, nil
}

func (store *PostgresStore) GetAccountByAccountNumber(accountNumber uuid.UUID) (*Account, error) {
	query := `
		SELECT * FROM accounts
		WHERE number = $1;	
	`
	rows, err := store.db.Query(query, accountNumber)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		account, err := scanIntoAccount(rows)
		if err != nil {
			return nil, err
		}
		return account, nil
	}
	return nil, fmt.Errorf("cannot find account")
}
