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
	GetAllAccounts() ([]*Account, error)
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

func (s *PostgreStore) createAccountTable() error {
    query := `
        create table if not exists account (
            id serial primary key,
            firstName varchar(200),
            lastName varchar(200),
            number serial,
            balance decimal(10, 2),
            createdAt timestamp with time zone
        );
    `

    _, err := s.conn.Exec(context.Background(), query)
    return err
}

func (s *PostgreStore) CreateAccount(account *Account) error {
    query := `
			insert into account (firstName, lastName, number, balance, createdAt) values ($1, $2, $3, $4, $5);
    `
    args := []interface{}{
        account.FirstName,
        account.LastName,
        account.Number,
        account.Balance,
        account.CreatedAt,
    }

    _, err := s.conn.Exec(context.Background(), query, args...)
    return err
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
func(s *PostgreStore) GetAllAccounts() ([]*Account, error) {
	query := `
		select * from account;
	`
	rows, err := s.conn.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accounts := []*Account{}
	for rows.Next() {
		var account Account
		err := rows.Scan(
			&account.ID,
			&account.FirstName,
			&account.LastName,
			&account.Number,
			&account.Balance,
			&account.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, &account)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return accounts, nil
}