package tests

import (
	"login-api/src/controllers/dtos"
	usermodels "login-api/src/models/user"
	"login-api/src/utils"
	"net/http"
	"testing"
	"time"

	"github.com/gavv/httpexpect"
	"github.com/google/uuid"
)

type UserData struct {
	ID        uuid.UUID
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

var sampleUserData = UserData{
	ID:        uuid.Must(uuid.NewRandom()),
	Name:      "test",
	Email:     "test@example.com",
	Password:  "thestrongestpassword",
	CreatedAt: time.Now(),
	UpdatedAt: time.Now(),
}

func (s *TestSuite) seedUser() *TestSuite {
	user := usermodels.User{
		ID:       sampleUserData.ID,
		Name:     sampleUserData.Name,
		Email:    sampleUserData.Email,
		Password: utils.MustHashed(utils.CreateHash(sampleUserData.Password)),
	}
	s.DB.Create(&user)
	return s
}

// auth
func TestRegister(t *testing.T) {
	svr := new(TestSuite).SetupEnvironment().SetupSuite().StartServer()

	t.Run("should successfully registered", func(t *testing.T) {
		test := httpexpect.New(t, svr.Server.URL)
		test.POST("/api/register").WithJSON(dtos.RegisterRequest{
			Name:                 sampleUserData.Name,
			Email:                sampleUserData.Email,
			Password:             sampleUserData.Password,
			PasswordConfirmation: sampleUserData.Password,
		}).Expect().Status(http.StatusCreated)
		svr.TearDownTest()
	})

	t.Run("should reject different password input", func(t *testing.T) {
		test := httpexpect.New(t, svr.Server.URL)
		test.POST("/api/register").WithJSON(dtos.RegisterRequest{
			Name:                 sampleUserData.Name,
			Email:                sampleUserData.Email,
			Password:             sampleUserData.Password,
			PasswordConfirmation: "sampleUserData.Password",
		}).Expect().Status(http.StatusBadRequest)
		svr.TearDownTest()
	})
}

func TestLogin(t *testing.T) {
	svr := new(TestSuite).SetupEnvironment().SetupSuite().StartServer().seedUser()
	defer svr.TearDownTest()

	t.Run("should successfully logged in", func(t *testing.T) {
		test := httpexpect.New(t, svr.Server.URL)
		test.POST("/api/login").WithJSON(dtos.LoginRequest{
			Email:    sampleUserData.Email,
			Password: sampleUserData.Password,
		}).Expect().Status(http.StatusCreated)
	})

	t.Run("should got unauthorized error", func(t *testing.T) {
		test := httpexpect.New(t, svr.Server.URL)
		test.POST("/api/login").WithJSON(dtos.LoginRequest{
			Email:    sampleUserData.Email,
			Password: "sampleUserData.Password",
		}).Expect().Status(http.StatusUnauthorized)
	})
}

// user functionality
func TestGetAllUsers(t *testing.T) {
	svr := new(TestSuite).SetupEnvironment().SetupSuite().StartServer().seedUser()
	defer svr.TearDownTest()

	loginFirst := func(uid string) (token string) {
		token, _ = utils.GenerateJwt(uid)
		return token
	}

	t.Run("should get all users", func(t *testing.T) {
		token := loginFirst(sampleUserData.ID.String())
		test := httpexpect.New(t, svr.Server.URL)
		test.GET("/api/users").WithCookie("token", token).Expect().Status(http.StatusOK)
	})

	t.Run("should got unauthenticated error", func(t *testing.T) {
		test := httpexpect.New(t, svr.Server.URL)
		test.GET("/api/users").Expect().Status(http.StatusUnauthorized)
	})
}

func TestGetUserDetailByID(t *testing.T) {
	svr := new(TestSuite).SetupEnvironment().SetupSuite().StartServer().seedUser()
	defer svr.TearDownTest()

	loginFirst := func(uid string) (token string) {
		token, _ = utils.GenerateJwt(uid)
		return token
	}

	t.Run("should get user detail by id", func(t *testing.T) {
		token := loginFirst(sampleUserData.ID.String())
		test := httpexpect.New(t, svr.Server.URL)
		test.GET("/api/users/d/"+sampleUserData.ID.String()).WithCookie("token", token).Expect().Status(http.StatusOK)
	})

	t.Run("should got unauthenticated error", func(t *testing.T) {
		test := httpexpect.New(t, svr.Server.URL)
		test.GET("/api/users/d/" + sampleUserData.ID.String()).Expect().Status(http.StatusUnauthorized)
	})

	t.Run("should got not found error", func(t *testing.T) {
		token := loginFirst(sampleUserData.ID.String())
		test := httpexpect.New(t, svr.Server.URL)
		test.GET("/api/users/d/"+uuid.Must(uuid.NewRandom()).String()).WithCookie("token", token).Expect().Status(http.StatusNotFound)
	})
}

func TestDeleteUserByID(t *testing.T) {
	svr := new(TestSuite).SetupEnvironment().SetupSuite().StartServer().seedUser()
	defer svr.TearDownTest()

	loginFirst := func(uid string) (token string) {
		token, _ = utils.GenerateJwt(uid)
		return token
	}

	t.Run("should delete user by id", func(t *testing.T) {
		token := loginFirst(sampleUserData.ID.String())
		test := httpexpect.New(t, svr.Server.URL)
		test.DELETE("/api/users/d/"+sampleUserData.ID.String()).WithCookie("token", token).Expect().Status(http.StatusNoContent)
	})

	t.Run("should got unauthenticated error", func(t *testing.T) {
		test := httpexpect.New(t, svr.Server.URL)
		test.DELETE("/api/users/d/" + sampleUserData.ID.String()).Expect().Status(http.StatusUnauthorized)
	})

	t.Run("should got not found error", func(t *testing.T) {
		token := loginFirst(sampleUserData.ID.String())
		test := httpexpect.New(t, svr.Server.URL)
		test.DELETE("/api/users/d/"+uuid.Must(uuid.NewRandom()).String()).WithCookie("token", token).Expect().Status(http.StatusNotFound)
	})
}

func TestAmendUserProfileByID(t *testing.T) {
	svr := new(TestSuite).SetupEnvironment().SetupSuite().StartServer().seedUser()
	defer svr.TearDownTest()

	loginFirst := func(uid string) (token string) {
		token, _ = utils.GenerateJwt(uid)
		return token
	}

	t.Run("should update user by id", func(t *testing.T) {
		token := loginFirst(sampleUserData.ID.String())
		test := httpexpect.New(t, svr.Server.URL)
		test.PUT("/api/users/p/"+sampleUserData.ID.String()).WithCookie("token", token).WithJSON(dtos.AmendProfileInputRequest{
			Name:  sampleUserData.Name,
			Email: sampleUserData.Email,
		}).Expect().Status(http.StatusNoContent)
	})

	t.Run("should got unauthenticated error", func(t *testing.T) {
		test := httpexpect.New(t, svr.Server.URL)
		test.PUT("/api/users/p/" + sampleUserData.ID.String()).WithJSON(dtos.AmendProfileInputRequest{
			Name:  sampleUserData.Name,
			Email: sampleUserData.Email,
		}).Expect().Status(http.StatusUnauthorized)
	})

	t.Run("should got not found error", func(t *testing.T) {
		token := loginFirst(sampleUserData.ID.String())
		test := httpexpect.New(t, svr.Server.URL)
		test.PUT("/api/users/p/"+uuid.Must(uuid.NewRandom()).String()).WithJSON(dtos.AmendProfileInputRequest{
			Name:  sampleUserData.Name,
			Email: sampleUserData.Email,
		}).WithCookie("token", token).Expect().Status(http.StatusNotFound)
	})
}

func TestAmendUserSecurityByID(t *testing.T) {
	svr := new(TestSuite).SetupEnvironment().SetupSuite().StartServer().seedUser()
	defer svr.TearDownTest()

	loginFirst := func(uid string) (token string) {
		token, _ = utils.GenerateJwt(uid)
		return token
	}

	t.Run("should update user password by id", func(t *testing.T) {
		token := loginFirst(sampleUserData.ID.String())
		test := httpexpect.New(t, svr.Server.URL)
		test.PUT("/api/users/s/"+sampleUserData.ID.String()).WithCookie("token", token).WithJSON(dtos.AmendProfilePasswordRequest{
			Password:             sampleUserData.Password,
			PasswordConfirmation: sampleUserData.Password,
		}).Expect().Status(http.StatusNoContent)
	})

	t.Run("should got unauthenticated error", func(t *testing.T) {
		test := httpexpect.New(t, svr.Server.URL)
		test.PUT("/api/users/s/" + sampleUserData.ID.String()).WithJSON(dtos.AmendProfilePasswordRequest{
			Password:             sampleUserData.Password,
			PasswordConfirmation: sampleUserData.Password,
		}).Expect().Status(http.StatusUnauthorized)
	})

	t.Run("should got not found error", func(t *testing.T) {
		token := loginFirst(sampleUserData.ID.String())
		test := httpexpect.New(t, svr.Server.URL)
		test.PUT("/api/users/s/"+uuid.Must(uuid.NewRandom()).String()).WithJSON(dtos.AmendProfilePasswordRequest{
			Password:             sampleUserData.Password,
			PasswordConfirmation: sampleUserData.Password,
		}).WithCookie("token", token).Expect().Status(http.StatusNotFound)
	})
}
