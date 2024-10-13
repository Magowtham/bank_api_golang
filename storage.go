package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Storage interface {
	InitDB() error
	CreateAccount(*Account) error
	UpdateAccount(int, string, string, string, string) error
	DeleteAccountByID(int) error
	GetAccountByID(int) (*Account, error)
	GetAllAccounts() ([]Account, error)
}

type PostgresStorage struct {
	db *sql.DB
}

func NewDataBase() *PostgresStorage {
	database_url := os.Getenv("DATABASE_URL")

	db, err := sql.Open("pgx", database_url)

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	if err != nil {
		log.Fatalln("Database configuration error ->", err.Error())
	}

	if err = db.Ping(); err != nil {
		log.Fatalln("Cannot ping to database because of error ->", err.Error())
	}

	log.Println("connected to database")

	return &PostgresStorage{
		db,
	}
}

func (storage *PostgresStorage) InitDB() error {

	query := `CREATE TABLE IF NOT EXISTS bank (
		id SERIAL PRIMARY KEY,
		first_name VARCHAR(100) NOT NULL,
		last_name VARCHAR(100) NOT NULL,
		email VARCHAR(50) UNIQUE NOT NULL,
		phone_number VARCHAR(50) UNIQUE NOT NULL,
		account_number VARCHAR(100) UNIQUE NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`

	_, err := storage.db.Exec(query)

	return err
}

func (storage *PostgresStorage) CreateAccount(acc *Account) error {

	query := `INSERT INTO bank (
		first_name,
		last_name,
		email,
		phone_number,
		account_number
		) VALUES ($1,$2,$3,$4,$5)`

	_, err := storage.db.Exec(query,
		acc.FirstName,
		acc.LastName,
		acc.Email,
		acc.PhoneNumber,
		acc.AccountNumber)

	return err
}

func (storage *PostgresStorage) UpdateAccount(id int, firstName, lastName, email, phoneNumber string) error {
	query := `UPDATE bank SET
		first_name=$1,
		last_name=$2,
		email=$3,
		phone_number=$4,
		WHERE id=$5`

	_, err := storage.db.Exec(query, firstName, lastName, email, phoneNumber, id)

	return err
}

func (storage *PostgresStorage) DeleteAccountByID(id int) error {
	query := `DELETE FROM bank WHERE id=$1`

	_, err := storage.db.Exec(query, id)

	return err
}

func (storage *PostgresStorage) GetAccountByID(id int) (*Account, error) {

	var account Account

	query := `SELECT first_name,last_name,email,phone_number,account_number FROM bank WHERE id=$1`

	row := storage.db.QueryRow(query, id)

	err := row.Scan(
		&account.FirstName,
		&account.LastName,
		&account.Email,
		&account.PhoneNumber,
		&account.AccountNumber,
	)

	if err != nil {
		return nil, err
	}

	return &account, nil

}

func (storage *PostgresStorage) GetAllAccounts() ([]Account, error) {
	var accounts []Account
	var account Account

	query := `SELECT first_name,last_name,email,phone_number,account_number FROM bank`

	rows, err := storage.db.Query(query)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err := rows.Scan(&account.FirstName, &account.LastName, &account.Email, &account.PhoneNumber, &account.AccountNumber)

		if err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}
