package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	l "github.com/sirupsen/logrus"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	l.Infoln("Artifactory upload service is listening at port 8080...")

	e.GET("/api/artifactory/upload", HandleGetUpload)

	e.Logger.Fatal(e.Start(":8080"))
}

func HandleGetUpload(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"msg": "Test endpoint"})
}
