package postgressql

import (
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/vladjong/user_balance/internal/entities"
)

const (
	CustomersTable      = "customers"
	TransactionTable    = "transactions"
	AccountsTable       = "accounts"
	HistoryTable        = "history"
	ExpectedTransaction = "expected_transactions"
	ReportView          = "history_report"
	CustomerReportView  = "customer_report"
)

type userBalanceStorage struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *userBalanceStorage {
	return &userBalanceStorage{
		db: db,
	}
}

func (d *userBalanceStorage) PostCustomerBalance(customer entities.Customer, transaction entities.Transaction) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	query := `INSERT INTO customers (id, balance)
				VALUES ($1, $2) ON CONFLICT (id)
				DO UPDATE SET (id, balance) = (EXCLUDED.id, EXCLUDED.balance + customers.balance)`
	if _, err := tx.Exec(query, customer.Id, customer.Balance); err != nil {
		if rb := tx.Rollback(); rb != nil {
			return rb
		}
		return err
	}
	var id int
	transactionQuery := `INSERT INTO transactions (customer_id, service_id, order_id, cost, transaction_datetime)
							VALUES ($1, $2, $3, $4, $5) RETURNING id`
	row := tx.QueryRow(transactionQuery, transaction.CustomeId, transaction.ServiceID, transaction.OrderID, transaction.Cost, transaction.TransactionDatiTime)
	if err := row.Scan(&id); err != nil {
		if rb := tx.Rollback(); rb != nil {
			return rb
		}
		return err
	}
	historyQuery := `INSERT INTO history (transaction_id, accounting_datetime, status_transaction)
						VALUES ($1, $2, $3)`
	if _, err := tx.Exec(historyQuery, id, transaction.TransactionDatiTime, true); err != nil {
		if rb := tx.Rollback(); rb != nil {
			return rb
		}
		return err
	}
	return tx.Commit()
}

func (d *userBalanceStorage) GetCustomerBalance(id int) (customer entities.Customer, err error) {
	query := `SELECT * FROM customers WHERE id = $1`
	var customers []entities.Customer
	if err := d.db.Select(&customers, query, id); err != nil {
		return customer, err
	}
	if len(customers) == 0 {
		return customer, errors.New("error: id don't exist")
	}
	return customers[0], nil
}

func (d *userBalanceStorage) PostReserveBalance(transaction entities.Transaction) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	customer, err := d.GetCustomerBalance(transaction.CustomeId)
	if err != nil {
		return err
	}
	if customer.Balance.LessThan(transaction.Cost) {
		return errors.New("error: customer balance less than transaction cost")
	}
	customer.Balance = customer.Balance.Sub(transaction.Cost)
	updateCustomerBalance := `UPDATE customers SET balance = $1 WHERE id = $2`
	if _, err := tx.Exec(updateCustomerBalance, customer.Balance, customer.Id); err != nil {
		if rb := tx.Rollback(); rb != nil {
			return rb
		}
		return err
	}
	query := `INSERT INTO accounts (customer_id, balance)
				VALUES ($1, $2) ON CONFLICT (customer_id)
				DO UPDATE SET (customer_id, balance) = (EXCLUDED.customer_id, EXCLUDED.balance + accounts.balance)`
	if _, err := tx.Exec(query, transaction.CustomeId, transaction.Cost); err != nil {
		if rb := tx.Rollback(); rb != nil {
			return rb
		}
		return err
	}
	var id int
	transactionQuery := `INSERT INTO transactions (customer_id, service_id, order_id, cost, transaction_datetime)
							VALUES ($1, $2, $3, $4, $5) RETURNING id`
	row := tx.QueryRow(transactionQuery, transaction.CustomeId, transaction.ServiceID, transaction.OrderID, transaction.Cost, transaction.TransactionDatiTime)
	if err := row.Scan(&id); err != nil {
		if rb := tx.Rollback(); rb != nil {
			return rb
		}
		return err
	}
	expectTransactionQuery := `INSERT INTO expected_transactions (transaction_id) VALUES ($1)`
	_, err = tx.Exec(expectTransactionQuery, id)
	if err != nil {
		if rb := tx.Rollback(); rb != nil {
			return rb
		}
		return err
	}
	return tx.Commit()
}

func (d *userBalanceStorage) PostDeReservingBalance(transaction entities.Transaction, history entities.History) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	var id []int
	searchTransactionId := `SELECT e.transaction_id
							FROM expected_transactions AS e
								JOIN transactions t ON e.transaction_id = t.id
							WHERE t.customer_id = $1 AND t.service_id = $2 AND t.order_id = $3 AND t.cost = $4`
	if err := d.db.Select(&id, searchTransactionId, transaction.CustomeId, transaction.ServiceID, transaction.OrderID, transaction.Cost); err != nil {
		return err
	}
	if len(id) == 0 {
		return errors.New("error: this id don't exist")
	}
	history.TransactionId = id[0]
	deleteTransactionQuery := `DELETE FROM expected_transactions WHERE transaction_id = $1`
	if _, err := tx.Exec(deleteTransactionQuery, history.TransactionId); err != nil {
		if rb := tx.Rollback(); rb != nil {
			return rb
		}
		return err
	}
	historyQuery := `INSERT INTO history (transaction_id, accounting_datetime, status_transaction) VALUES ($1, $2, $3)`
	if _, err := tx.Exec(historyQuery, history.TransactionId, history.AccountingDatetime, history.StatusTransaction); err != nil {
		if rb := tx.Rollback(); rb != nil {
			return rb
		}
		return err
	}
	updateAccountBalance := `UPDATE accounts SET balance = balance - $1 WHERE customer_id = $2`
	if _, err := tx.Exec(updateAccountBalance, transaction.Cost, transaction.CustomeId); err != nil {
		if rb := tx.Rollback(); rb != nil {
			return rb
		}
		return err
	}
	if !history.StatusTransaction {
		updateCustomerBalance := `UPDATE customers SET balance = balance + $1 WHERE id = $2`
		if _, err := tx.Exec(updateCustomerBalance, transaction.Cost, transaction.CustomeId); err != nil {
			if rb := tx.Rollback(); rb != nil {
				return rb
			}
			return err
		}
	}
	return tx.Commit()
}

func (d *userBalanceStorage) GetHistoryReport(date time.Time) (report []entities.Report, err error) {
	query := `SELECT ROW_NUMBER() OVER(ORDER BY name) AS id, name, SUM(cost) AS all_sum
				FROM history_report
				WHERE $1 <= accounting_datetime AND $1::timestamp + INTERVAL '1' MONTH > accounting_datetime
				GROUP BY name`
	if err := d.db.Select(&report, query, date); err != nil {
		return report, err
	}
	return report, nil
}

func (d *userBalanceStorage) GetCustomerReport(id int, date time.Time) (report []entities.CustomerReport, err error) {
	query := `SELECT ROW_NUMBER() OVER(ORDER BY date DESC, sum DESC) AS id, service_name, order_name, sum, status_transaction, date
				FROM customer_report
				WHERE $1 <= date
				AND $1::timestamp + INTERVAL '1' MONTH > date
				AND customer_id = $2
				ORDER BY date DESC, sum DESC`
	if err := d.db.Select(&report, query, date, id); err != nil {
		return report, err
	}
	if report == nil {
		empty := fmt.Sprintf("don't have customer id: %d history report in %s", id, date.String())
		return nil, errors.New(empty)
	}
	return report, nil
}
