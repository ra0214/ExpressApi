package application

import (
	"expresApi/src/users/domain"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type LoginUser struct {
	db domain.IUser
}

func NewLoginUser(db domain.IUser) *LoginUser {
	return &LoginUser{
		db: db,
	}
}

func (lu *LoginUser) Execute(userName string, password string) (*domain.User, error) {
	// Obtener usuario solo con el username
	user, err := lu.db.GetUserByCredentials(userName)
	if err != nil {
		return nil, err
	}

	// Verificar la contraseña
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("credenciales inválidas")
	}

	return user, nil
}
