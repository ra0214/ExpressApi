package domain

type IUser interface {
	SaveUser(userName string, email string, password string, estado bool) error
	DeleteUser(id int32) error
	UpdateUser(id int32, userName string, email string, password string) error
	GetAll() ([]User, error)
	GetUserByCredentials(userName string) (*User, error)
	ToggleUserStatus(id int32) (bool, error)
}

type User struct {
	ID        int32  `json:"id"`
	UserName  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password,omitempty"`
	Estado    bool   `json:"estado"`
	CreatedAt string `json:"created_at"`
}

func NewUser(userName string, email string, password string, estado bool) *User {
	return &User{
		UserName: userName,
		Email:    email,
		Password: password,
		Estado:   estado,
	}
}

func (u *User) SetUserName(userName string) {
	u.UserName = userName
}
