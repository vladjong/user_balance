package usecase

import (
	"github.com/vladjong/user_balance/internal/adapters/db"
)

type userBalanseUseCase struct {
	storage db.UserBalanse
}

func New(storage db.UserBalanse) *userBalanseUseCase {
	return &userBalanseUseCase{
		storage: storage,
	}
}
