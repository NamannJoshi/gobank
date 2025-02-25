package api

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account, int) error
	GetAccountByID(int) (*Account, error)
}

type PostgreStore struct {
	conn *pgx.Conn
}

func NewPostgreStore() (*PostgreStore, error) {
	connString := "postgres://postgres:mukeshakun@localhost:5432/test"
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		return nil, err
	}

	return &PostgreStore{
		conn: conn,
	}, nil
}

func(s *PostgreStore) Init() error {
	return s.createAccountTable()
}

func(s *PostgreStore) createAccountTable() error {
	query := `
		create table if not exists account (
			id serial primary key,
			firstName varchar(200),
			lastName varchar(200),
			number serial,
			balance decimal(10, 2)
		);
	`

	_, err := s.conn.Exec(context.Background(), query)
	return err
}

func(s *PostgreStore) CreateAccount(account *Account) (error) {
	return nil
}
func(s *PostgreStore) UpdateAccount(account *Account, i int) (error) {
	return nil
}
func(s *PostgreStore) DeleteAccount(i int) (error) {
	return nil
}
func(s *PostgreStore) GetAccountByID(int) (*Account, error) {
	return &Account{}, nil
}