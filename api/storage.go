package api

import (
	"context"
	"fmt"
	"log"
	"strconv"

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
			insert into account (firstName, lastName, number, balance, createdAt) values (@firstName, @lastName, @number, @balance, @createdAt);
    `
    args := pgx.NamedArgs{
			"firstName": account.FirstName,
			"lastName": account.LastName,
			"number": account.Number,
			"balance": account.Balance,
			"createdAt": account.CreatedAt,
    }

    _, err := s.conn.Exec(context.Background(), query, args)
    return err
}

func (s *PostgreStore) UpdateAccount(account *Account, accountId int) error {
    query := "UPDATE account SET "
    namedArgs := pgx.NamedArgs{}
    updates := make(map[string]interface{})
    argCount := 1

    if account.FirstName != "" {
        updates["firstName"] = account.FirstName
    }
    if account.LastName != "" {
        updates["lastName"] = account.LastName
    }
    if account.Number != 0 {
        updates["number"] = account.Number
    }
    if account.Balance != 0 {
        updates["balance"] = account.Balance
    }

    if len(updates) == 0 {
        return nil 
    }

    for key, value := range updates {
        query += key + " = @" + key + ", "
        namedArgs[key] = value
        argCount++
    }

    query = query[:len(query)-2] // Remove trailing comma and space
    query += " WHERE id = $" + strconv.Itoa(argCount)
    namedArgs["id"] = accountId

    _, err := s.conn.Exec(context.Background(), query, namedArgs)
    return err
}

func(s *PostgreStore) DeleteAccount(accountId int) (error) {
	query := `
		delete from account where id = @accountId;
	`
	args := pgx.NamedArgs{
		"accountId": strconv.Itoa(accountId),
	}
	_, err := s.conn.Exec(context.Background(), query, args)
	if err != nil {
		log.Fatalf("error while deletion in DB: %v", err)
	}
	return err
}

func(s *PostgreStore) GetAccountByID(accountId int) (*Account, error) {
	query := `
		select * from account where id = @accountId;
	`
	args := pgx.NamedArgs{
		"accountId": strconv.Itoa(accountId),
	}
	row := s.conn.QueryRow(context.Background(), query, args)

	var res Account 
	err := row.Scan(&res.ID, &res.FirstName, &res.LastName, &res.Number, &res.Balance, &res.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("this is gone brrrrr")
		}
	return nil, err
	}
	
	return &res, nil
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

	//Short way to new pgx version
	// accounts, err := pgx.CollectRows(rows, pgx.RowToStructByName[*Account])
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