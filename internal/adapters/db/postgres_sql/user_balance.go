package postgressql

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/vladjong/user_balance/internal/entities"
)

const (
	CustomersTable   = "customers"
	TransactionTable = "transaction"
	AccountringTable = "accounting"
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
							DO UPDATE SET (id, balance) = (EXCLUDED.id, EXCLUDED.balance + customers.balance)
							RETURNING id`, CustomersTable)
	row := d.db.QueryRow(query, customer.Id, customer.Balance)
	var id int
	if err := row.Scan(&id); err != nil {
		return err
	}
	return nil
}

func (d *userBalanceStorage) GetCustomerBalance(id int) (customer entities.Customer, err error) {
	query := fmt.Sprintf(`SELECT * FROM %s WHERE id = $1`, CustomersTable)
	var customers []entities.Customer
	if err := d.db.Select(&customers, query, id); err != nil {
		return customer, err
	}
	if len(customers) == 0 {
		return customer, fmt.Errorf("error: id don't exist")
	}
	return customers[0], nil
}
