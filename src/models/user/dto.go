package usermodels

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       uuid.UUID `json:"ID" gorm:"type:string;size:36;primaryKey"`
	Name     string    `json:"name"`
	Email    string    `json:"email" gorm:"unique"`
	Password string    `json:"password"`
}

type UserModel interface {
	GetAllUsers() (users []User, err error)
	GetUser(id string) (user User, err error)
	UpdateUserProfile(user User, uid string) (affected int, err error)
	UpdateUserPassword(user User, uid string) (affected int, err error)
	DeleteUser(id string) (affected int, err error)

	// auth
	AttemptLogin(email, password string) (uid string, err error)
	AttemptRegister(user User) (uid string, err error)
}
