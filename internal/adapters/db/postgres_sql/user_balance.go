package postgressql

import "github.com/jmoiron/sqlx"

type userBalanceStorage struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *userBalanceStorage {
	return &userBalanceStorage{
		db: db,
	}
}
