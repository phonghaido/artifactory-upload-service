package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/phonghaido/artifactory-upload-service/handlers"
	"github.com/phonghaido/artifactory-upload-service/helpers"
	l "github.com/sirupsen/logrus"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	l.Infoln("Artifactory upload service is listening at port 8080...")

	e.POST("/api/artifactory/upload/file", helpers.EchoErrorWrapper(handlers.HandlePostUploadFile))
	e.POST("/api/artifactory/upload/files", helpers.EchoErrorWrapper(handlers.HandlePostUploadFiles))

	e.Logger.Fatal(e.Start(":8080"))
}
