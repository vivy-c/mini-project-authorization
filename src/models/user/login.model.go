package usermodels

import (
	"errors"
	"login-api/src/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type models struct {
	DB *gorm.DB
}

// AttemptLogin implements UserModel
func (m *models) AttemptLogin(email string, password string) (uid string, err error) {
	var user User

	if err = m.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return "", errors.New("user not found")
	}

	if !utils.ValidateHash(password, user.Password) {
		return "", errors.New("password not match")
	}

	return user.ID.String(), nil
}

// AttemptRegister implements UserModel
func (m *models) AttemptRegister(user User) (uid string, err error) {
	user = User{
		ID:       uuid.Must(uuid.NewRandom()),
		Name:     user.Name,
		Email:    user.Email,
		Password: utils.MustHashed(utils.CreateHash(user.Password)),
	}
	return user.ID.String(), m.DB.Create(&user).Error
}

// DeleteUser implements UserModel
func (m *models) DeleteUser(id string) (affected int, err error) {
	q := m.DB.Where("ID = ?", id).Delete(new(User))
	return int(q.RowsAffected), q.Error
}

// GetAllUsers implements UserModel
func (m *models) GetAllUsers() (users []User, err error) {
	err = m.DB.Find(&users).Error
	return users, err
}

// GetUser implements UserModel
func (m *models) GetUser(id string) (user User, err error) {
	err = m.DB.Where("ID = ?", id).Find(&user).Error
	return user, err
}

// UpdateUserPassword implements UserModel
func (m *models) UpdateUserPassword(user User, uid string) (affected int, err error) {
	q := m.DB.Where("ID = ?", uid).Updates(&user)
	return int(q.RowsAffected), q.Error
}

// UpdateUserProfile implements UserModel
func (m *models) UpdateUserProfile(user User, uid string) (affected int, err error) {
	q := m.DB.Where("ID = ?", uid).Updates(&user)
	return int(q.RowsAffected), q.Error
}

func NewModel(DB *gorm.DB) UserModel {
	return &models{DB}
}
