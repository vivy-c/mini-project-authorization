package controllers

import (
	"errors"
	"login-api/src/controllers/dtos"
	usermodels "login-api/src/models/user"
	"login-api/src/utils"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

var userFactory = usermodels.NewModel(db)

// onLogin
func LoginHandler(ec echo.Context) error {
	var login dtos.LoginRequest

	if err := ec.Bind(&login); err != nil {
		status := http.StatusBadRequest
		return utils.CreateResponse(ec, status, http.StatusText(status), err.Error())
	}

	uid, err := userFactory.AttemptLogin(login.Email, login.Password)
	if err != nil {
		status := http.StatusUnauthorized
		return utils.CreateResponse(ec, status, http.StatusText(status), err.Error())
	}

	token, err := utils.GenerateJwt(uid)
	if err != nil {
		status := http.StatusInternalServerError
		return utils.CreateResponse(ec, status, http.StatusText(status), err.Error())
	}

	utils.SetAuthCookie(ec, token)
	status := http.StatusCreated
	return utils.CreateResponse(ec, status, http.StatusText(status), map[string]string{
		"token": token,
	})
}

// onRegister
func RegisterHandler(ec echo.Context) error {
	var register dtos.RegisterRequest

	if err := ec.Bind(&register); err != nil {
		status := http.StatusBadRequest
		return utils.CreateResponse(ec, status, http.StatusText(status), err.Error())
	}

	if register.Password != register.PasswordConfirmation {
		status := http.StatusBadRequest
		return utils.CreateResponse(ec, status, http.StatusText(status), errors.New("password didn't confirmed"))
	}

	uid, err := userFactory.AttemptRegister(usermodels.User{
		Name:     register.Name,
		Email:    register.Email,
		Password: register.Password,
	})
	if err != nil {
		status := http.StatusUnauthorized
		return utils.CreateResponse(ec, status, http.StatusText(status), err.Error())
	}

	token, err := utils.GenerateJwt(uid)
	if err != nil {
		status := http.StatusInternalServerError
		return utils.CreateResponse(ec, status, http.StatusText(status), err.Error())
	}

	utils.SetAuthCookie(ec, token)
	status := http.StatusCreated
	return utils.CreateResponse(ec, status, http.StatusText(status), map[string]string{
		"token": token,
	})
}

// onAmendProfileByID
func AmendProfileByIDHandler(ec echo.Context) error {
	var profile dtos.AmendProfileInputRequest

	if err := ec.Bind(&profile); err != nil {
		status := http.StatusBadRequest
		return utils.CreateResponse(ec, status, http.StatusText(status), err.Error())
	}

	rows, err := userFactory.UpdateUserProfile(usermodels.User{
		Name:  profile.Name,
		Email: profile.Email,
	}, ec.Param("id"))

	if rows <= 0 {
		status := http.StatusNotFound
		return utils.CreateResponse(ec, status, http.StatusText(status), nil)
	}

	if err != nil {
		status := http.StatusInternalServerError
		return utils.CreateResponse(ec, status, http.StatusText(status), err.Error())
	}

	status := http.StatusNoContent
	return utils.CreateResponse(ec, status, http.StatusText(status), nil)
}

// onAmendPasswordByID
func AmendPasswordByIDHandler(ec echo.Context) error {
	var password dtos.AmendProfilePasswordRequest

	if err := ec.Bind(&password); err != nil {
		status := http.StatusBadRequest
		return utils.CreateResponse(ec, status, http.StatusText(status), err.Error())
	}

	if password.Password != password.PasswordConfirmation {
		status := http.StatusBadRequest
		return utils.CreateResponse(ec, status, http.StatusText(status), errors.New("password didn't confirmed"))
	}

	rows, err := userFactory.UpdateUserPassword(usermodels.User{
		Password: password.Password,
	}, ec.Param("id"))

	if rows <= 0 {
		status := http.StatusNotFound
		return utils.CreateResponse(ec, status, http.StatusText(status), nil)
	}

	if err != nil {

		status := http.StatusInternalServerError
		return utils.CreateResponse(ec, status, http.StatusText(status), err.Error())
	}

	status := http.StatusNoContent
	return utils.CreateResponse(ec, status, http.StatusText(status), nil)
}

// onRemoveUserByID
func RemoveUserByIDHandler(ec echo.Context) error {
	rows, err := userFactory.DeleteUser(ec.Param("id"))
	if rows <= 0 {
		status := http.StatusNotFound
		return utils.CreateResponse(ec, status, http.StatusText(status), nil)
	}

	if err != nil {
		status := http.StatusInternalServerError
		return utils.CreateResponse(ec, status, http.StatusText(status), err.Error())
	}

	status := http.StatusNoContent
	return utils.CreateResponse(ec, status, http.StatusText(status), nil)
}

// onGetUserDetailByID
func GetUserDetailByIDHandler(ec echo.Context) error {
	user, err := userFactory.GetUser(ec.Param("id"))

	if user.ID == uuid.Nil {
		status := http.StatusNotFound
		return utils.CreateResponse(ec, status, http.StatusText(status), nil)
	}

	if err != nil {
		status := http.StatusInternalServerError
		return utils.CreateResponse(ec, status, http.StatusText(status), err.Error())
	}

	status := http.StatusOK
	return utils.CreateResponse(ec, status, http.StatusText(status), dtos.UserResponse{
		ID:        user.ID.String(),
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}

// onGetAllUsers
func GetAllUsersHandler(ec echo.Context) error {
	users, err := userFactory.GetAllUsers()

	if err != nil {
		status := http.StatusInternalServerError
		return utils.CreateResponse(ec, status, http.StatusText(status), err.Error())
	}

	var responses []dtos.UserResponse
	for _, user := range users {
		responses = append(responses, dtos.UserResponse{
			ID:        user.ID.String(),
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}

	status := http.StatusOK
	return utils.CreateResponse(ec, status, http.StatusText(status), responses)
}
