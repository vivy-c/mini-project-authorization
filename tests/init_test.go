package tests

import (
	"login-api/src/configs"
	usermodels "login-api/src/models/user"
	"login-api/src/routes"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type TestSuite struct {
	suite.Suite
	DB        *gorm.DB
	Echo      *echo.Echo
	Server    *httptest.Server
	ConnStr   string
	JWTSecret string
}

func (s *TestSuite) SetupEnvironment() *TestSuite {
	config, _ := configs.LoadServerConfig(".")
	s.ConnStr = config.ConnectionString
	s.JWTSecret = config.JWTsecret
	return s
}

func (s *TestSuite) SetupSuite() *TestSuite {
	s.DB, _ = gorm.Open(mysql.Open(s.ConnStr), &gorm.Config{})

	s.DB.AutoMigrate(&usermodels.User{})
	return s
}

func (s *TestSuite) TearDownTest() *TestSuite {
	s.DB.Exec("truncate table users")
	return s
}

func (s *TestSuite) TearDownSuite() {
	require.NoError(s.T(), s.DB.Migrator().DropTable("users"))
}

func (s *TestSuite) StartServer() *TestSuite {
	s.Server = httptest.NewServer(routes.New())
	return s
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
