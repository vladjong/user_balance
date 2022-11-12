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
	query := fmt.Sprintf(`INSERT INTO %s (id, balance)
							VALUES ($1, $2) ON CONFLICT (id)
							DO UPDATE SET (id, balance) = (EXCLUDED.id, EXCLUDED.balance + customers.balance)`, CustomersTable)
	if _, err := tx.Exec(query, customer.Id, customer.Balance); err != nil {
		tx.Rollback()
		return err
	}
	var id int
	transactionQuery := fmt.Sprintf(`INSERT INTO %s (customer_id, service_id, order_id, cost, transaction_datetime)
										VALUES ($1, $2, $3, $4, $5) RETURNING id`, TransactionTable)
	row := tx.QueryRow(transactionQuery, transaction.CustomeId, transaction.ServiceID, transaction.OrderID, transaction.Cost, transaction.TransactionDatiTime)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return err
	}
	historyQuery := fmt.Sprintf(`INSERT INTO %s (transaction_id, accounting_datetime, status_transaction)
									VALUES ($1, $2, $3)`, HistoryTable)
	if _, err := tx.Exec(historyQuery, id, transaction.TransactionDatiTime, true); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
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
	updateCustomerBalance := fmt.Sprintf("UPDATE %s SET balance = $1 WHERE id = $2", CustomersTable)
	if _, err := tx.Exec(updateCustomerBalance, customer.Balance, customer.Id); err != nil {
		tx.Rollback()
		return err
	}
	query := fmt.Sprintf(`INSERT INTO %s (customer_id, balance)
							VALUES ($1, $2) ON CONFLICT (customer_id)
							DO UPDATE SET (customer_id, balance) = (EXCLUDED.customer_id, EXCLUDED.balance + accounts.balance)`, AccountsTable)
	if _, err := tx.Exec(query, transaction.CustomeId, transaction.Cost); err != nil {
		tx.Rollback()
		return err
	}
	var id int
	transactionQuery := fmt.Sprintf(`INSERT INTO %s (customer_id, service_id, order_id, cost, transaction_datetime)
										VALUES ($1, $2, $3, $4, $5) RETURNING id`, TransactionTable)
	row := tx.QueryRow(transactionQuery, transaction.CustomeId, transaction.ServiceID, transaction.OrderID, transaction.Cost, transaction.TransactionDatiTime)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return err
	}
	expectTransactionQuery := fmt.Sprintf(`INSERT INTO %s (transaction_id)
											VALUES ($1)`, ExpectedTransaction)
	_, err = tx.Exec(expectTransactionQuery, id)
	if err != nil {
		tx.Rollback()
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
	searchTransactionId := fmt.Sprintf(`SELECT e.transaction_id
										FROM %s AS e
											JOIN %s t ON e.transaction_id = t.id
										WHERE t.customer_id = $1 AND t.service_id = $2 AND t.order_id = $3 AND t.cost = $4
										`, ExpectedTransaction, TransactionTable)
	if err := d.db.Select(&id, searchTransactionId, transaction.CustomeId, transaction.ServiceID, transaction.OrderID, transaction.Cost); err != nil {
		return err
	}
	if len(id) == 0 {
		return errors.New("error: this id don't exist")
	}
	history.TransactionId = id[0]
	deleteTransactionQuery := fmt.Sprintf(`DELETE FROM %s WHERE transaction_id = $1`, ExpectedTransaction)
	if _, err := tx.Exec(deleteTransactionQuery, history.TransactionId); err != nil {
		tx.Rollback()
		return err
	}
	historyQuery := fmt.Sprintf(`INSERT INTO %s (transaction_id, accounting_datetime, status_transaction)
									VALUES ($1, $2, $3)`, HistoryTable)
	if _, err := tx.Exec(historyQuery, history.TransactionId, history.AccountingDatetime, history.StatusTransaction); err != nil {
		tx.Rollback()
		return err
	}
	updateAccountBalance := fmt.Sprintf(`UPDATE %s SET balance = balance - $1 WHERE customer_id = $2`, AccountsTable)
	if _, err := tx.Exec(updateAccountBalance, transaction.Cost, transaction.CustomeId); err != nil {
		tx.Rollback()
		return err
	}
	if !history.StatusTransaction {
		updateCustomerBalance := fmt.Sprintf(`UPDATE %s SET balance = balance + $1 WHERE id = $2`, CustomersTable)
		if _, err := tx.Exec(updateCustomerBalance, transaction.Cost, transaction.CustomeId); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (d *userBalanceStorage) GetHistoryReport(date time.Time) (report []entities.Report, err error) {
	query := fmt.Sprintf(`SELECT ROW_NUMBER() OVER(ORDER BY name) AS id, name, SUM(cost) AS all_sum
							FROM %s
							WHERE $1 <= accounting_datetime AND $1::timestamp + INTERVAL '1' MONTH > accounting_datetime
							GROUP BY name`, ReportView)
	if err := d.db.Select(&report, query, date); err != nil {
		return report, err
	}
	if len(report) == 0 {
		return report, errors.New("error: history report empty")
	}
	return report, nil
}

func (d *userBalanceStorage) GetCustomerReport(id int, date time.Time) (report []entities.CustomerReport, err error) {
	query := fmt.Sprintf(`SELECT ROW_NUMBER() OVER(ORDER BY date DESC, sum DESC) AS id, service_name, order_name, sum, date
							FROM %s
							WHERE $1 <= date
							AND $1::timestamp + INTERVAL '1' MONTH > date
							AND customer_id = $2
							ORDER BY date DESC, sum DESC`, CustomerReportView)
	if err := d.db.Select(&report, query, date, id); err != nil {
		return report, err
	}
	if len(report) == 0 {
		return report, errors.New("error: customer history report empty")
	}
	return report, nil
}
