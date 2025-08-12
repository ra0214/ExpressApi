package application

import (
	"expresApi/src/users/domain"
)

type ViewUser struct {
	db domain.IUser
}

func NewViewUser(db domain.IUser) *ViewUser {
	return &ViewUser{db: db}
}

func (vu *ViewUser) Execute() ([]domain.User, error) {
	return vu.db.GetAll()
}