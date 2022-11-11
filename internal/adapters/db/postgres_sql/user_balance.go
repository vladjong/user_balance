package postgressql

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/vladjong/user_balance/internal/entities"
)

const (
	CustomersTable   = "customers"
	TransactionTable = "transactions"
	AccountsTable    = "accounts"
	HistoryTable     = "history"
)

type userBalanceStorage struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *userBalanceStorage {
	return &userBalanceStorage{
		db: db,
	}
}

func (d *userBalanceStorage) PostCustomerBalance(customer entities.Customer) error {
	query := fmt.Sprintf(`INSERT INTO %s (id, balance)
							VALUES ($1, $2) ON CONFLICT (id)
							DO UPDATE SET (id, balance) = (EXCLUDED.id, EXCLUDED.balance + customers.balance)`, CustomersTable)
	_, err := d.db.Exec(query, customer.Id, customer.Balance)
	return err
}

func (d *userBalanceStorage) GetCustomerBalance(id int) (customer entities.Customer, err error) {
	query := fmt.Sprintf(`SELECT * FROM %s WHERE id = $1`, CustomersTable)
	var customers []entities.Customer
	if err := d.db.Select(&customers, query, id); err != nil {
		return customer, err
	}
	if len(customers) == 0 {
		return customer, errors.New("error: id don't exist")
	}
	return customers[0], nil
}

func (d *userBalanceStorage) PostReserveBalance(transaction entities.Transaction, customer entities.Customer, account entities.Acount) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	if err = d.updateCustomerBalance(customer); err != nil {
		tx.Rollback()
		return err
	}
	logrus.Info("tet")
	if err = d.postAccountBalance(account); err != nil {
		tx.Rollback()
		return err
	}
	transactionQuery := fmt.Sprintf(`INSERT INTO %s (customer_id, service_id, order_id, cost, transaction_datetime)
										VALUES ($1, $2, $3, $4, $5)`, TransactionTable)
	_, err = tx.Exec(transactionQuery, transaction.CustomeId, transaction.ServiceID, transaction.OrderID, transaction.Cost, transaction.TransactionDatiTime)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (d *userBalanceStorage) PostDeReservingBalance(customer entities.Customer, account entities.Acount, history entities.History) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	accountingQuery := fmt.Sprintf(`INSERT INTO %s (transaction_id, accounting_datetime, status_transaction)
										VALUES ($1, $2, $3)`, HistoryTable)
	_, err = tx.Exec(accountingQuery, history.TransactionId, history.AccountingDatetime, history.StatusTransaction)
	if err != nil {
		tx.Rollback()
		return err
	}
	if history.StatusTransaction {
		if err = d.updateAccountBalance(account); err != nil {
			tx.Rollback()
			return err
		}
	} else {
		if err = d.updateCustomerBalance(customer); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (d *userBalanceStorage) GetTransactionId(transaction entities.Transaction) (int, error) {
	var ids []int
	transactionQuery := fmt.Sprintf(`SELECT id FROM %s WHERE customer_id = $1 AND service_id = $2 AND order_id = $3 AND cost = $4`, TransactionTable)
	if err := d.db.Select(&ids, transactionQuery, transaction.CustomeId, transaction.ServiceID, transaction.OrderID, transaction.Cost); err != nil {
		return 0, err
	}
	return ids[0], nil
}

func (d *userBalanceStorage) updateCustomerBalance(customer entities.Customer) error {
	query := fmt.Sprintf("UPDATE %s SET balance = $1 WHERE id = $2", CustomersTable)
	_, err := d.db.Exec(query, customer.Balance, customer.Id)
	return err
}

func (d *userBalanceStorage) updateAccountBalance(account entities.Acount) error {
	query := fmt.Sprintf("UPDATE %s SET balance = balance - $1 WHERE id = $2", AccountsTable)
	_, err := d.db.Exec(query, account.Balance, account.CustomerId)
	return err
}

func (d *userBalanceStorage) postAccountBalance(account entities.Acount) error {
	query := fmt.Sprintf(`INSERT INTO %s (customer_id, balance)
							VALUES ($1, $2) ON CONFLICT (customer_id)
							DO UPDATE SET (customer_id, balance) = (EXCLUDED.customer_id, EXCLUDED.balance + accounts.balance)`, AccountsTable)
	_, err := d.db.Exec(query, account.CustomerId, account.Balance)
	return err
}
