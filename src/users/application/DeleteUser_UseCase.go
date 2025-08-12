package application

import (
	"expresApi/src/users/domain"
)

type DeleteUser struct {
	db domain.IUser
}

func NewDeleteUser(db domain.IUser) *DeleteUser {	
	return &DeleteUser{db: db}
}

func (du *DeleteUser) Execute(id int32) error {
	return du.db.DeleteUser(id)
}