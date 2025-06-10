package repository

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/epic55/BankAppNew/internal/models"
)

type RepositoryInterface interface {
	GetByID(id int64) (*models.User, error)
	BuyingRepo(userId int) (*models.Account, error)
	UpdateAccount(updatedBalance, changesToAccountBalance float64, id int, AccountCurrency, typeofoperation2 string, date1 string) (string, error)
}

type repositoryStruct struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) RepositoryInterface {
	return &repositoryStruct{db: db}
}

func (r *repositoryStruct) GetByID(id int64) (*models.User, error) {
	row := r.db.QueryRow("SELECT id, name FROM users WHERE id=$1", id)
	var user models.User
	err := row.Scan(&user.ID, &user.Name)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repositoryStruct) BuyingRepo(userId int) (*models.Account, error) {
	queryStmt := `SELECT * FROM accounts WHERE id = $1 ;`
	results, err := r.db.Query(queryStmt, userId)
	if err != nil {
		//log.Println("failed to execute query", err)
		return nil, err
	}

	var account models.Account
	for results.Next() {
		err = results.Scan(&account.Id, &account.Name, &account.Account, &account.Balance, &account.Currency, &account.Date, &account.Blocked, &account.Defaultaccount)
		if err != nil {
			//log.Println("failed to scan", err)
			return nil, err
		}
	}

	return &account, nil
}

func (r *repositoryStruct) UpdateAccount(updatedBalance, changesToAccountBalance float64, id int, AccountCurrency, typeofoperation2 string, date1 string) (string, error) {
	queryStmt := `UPDATE accounts SET balance = $2, currency = $3, date = $4 WHERE id = $1 RETURNING id;`
	err := r.db.QueryRow(queryStmt, &id, &updatedBalance, &AccountCurrency, date1).Scan(&id)
	if err != nil {
		//log.Println("failed to execute query:", err)
		//w.WriteHeader(500)
		return "", err
	} else {
		fmt.Printf("Balance is %s on %.2f Result: %.2f\n", typeofoperation2, changesToAccountBalance, updatedBalance)
	}

	return "Balance is " + string(typeofoperation2) + " on " + strconv.FormatFloat(changesToAccountBalance, 'f', 2, 64), nil
}
