package application

import (
	"expresApi/src/users/domain"

	"golang.org/x/crypto/bcrypt"
)

type CreateUser struct {
	db domain.IUser
}

func NewCreateUser(db domain.IUser) *CreateUser {
	return &CreateUser{db: db}
}

func (cu *CreateUser) Execute(userName string, email string, password string) error {
	// Generar hash de la contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Por defecto, los usuarios nuevos están activos (estado = true)
	err = cu.db.SaveUser(userName, email, string(hashedPassword), true)
	if err != nil {
		return err
	}

	return nil
}
