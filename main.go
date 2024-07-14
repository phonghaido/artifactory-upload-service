package main

import (
	"crypto/subtle"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/phonghaido/artifactory-upload-service/configs"
	"github.com/phonghaido/artifactory-upload-service/db"
	"github.com/phonghaido/artifactory-upload-service/handlers"
	"github.com/phonghaido/artifactory-upload-service/helpers"
	l "github.com/sirupsen/logrus"
)

var userDB *db.PostgreSQL

func main() {
	config, err := configs.GetConfig()
	if err != nil {
		l.Fatal(err)
	}

	userDB = db.NewPostgreSQL(config.PostgresConnStr)
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.BasicAuth(func(username, password string, ctx echo.Context) (bool, error) {
		user, _ := userDB.SearchUser(username, password)

		if subtle.ConstantTimeCompare([]byte(username), []byte(user.Username)) == 1 && subtle.ConstantTimeCompare([]byte(password), []byte(user.Password)) == 1 {
			return true, nil
		}
		return false, nil
	}))

	l.Infoln("Artifactory upload service is listening at port 8080...")

	e.POST("/api/artifactory/upload/file", helpers.EchoErrorWrapper(handlers.HandlePostUploadFile))
	e.POST("/api/artifactory/upload/files", helpers.EchoErrorWrapper(handlers.HandlePostUploadFiles))

	e.Logger.Fatal(e.Start(":8080"))
}
