package helpers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	l "github.com/sirupsen/logrus"
)

type APIError struct {
	StatusCode int    `json:"statusCode"`
	Msg        string `json:"msg"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("api error: %d - %s", e.StatusCode, e.Msg)
}

func NewAPIError(statusCode int, err error) APIError {
	return APIError{
		StatusCode: statusCode,
		Msg:        err.Error(),
	}
}

type EchoAPIFunc func(c echo.Context) error

func EchoErrorWrapper(h EchoAPIFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := h(c); err != nil {
			if apiErr, ok := err.(APIError); ok {
				l.Errorf(apiErr.Error())
				return WriteJSON(c, apiErr.StatusCode, apiErr)
			} else {
				errResp := map[string]any{
					"statusCode": http.StatusInternalServerError,
					"msg":        "internal server error",
				}
				l.Errorf(err.Error())
				return WriteJSON(c, http.StatusInternalServerError, errResp)
			}
		}
		return nil
	}
}

func WriteJSON(c echo.Context, statusCode int, v any) error {
	return c.JSON(statusCode, v)
}

func MethodNotAllowed() error {
	return NewAPIError(http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
}

func InvalidForm(err error) error {
	return NewAPIError(http.StatusBadRequest, fmt.Errorf("invalid form data: %s", err.Error()))
}

func InvalidFileSize(maxSize int64) error {
	return NewAPIError(http.StatusBadRequest, fmt.Errorf("file size exceeds %d byte(s)", maxSize))
}

func InvalidFileKey() error {
	return NewAPIError(http.StatusBadRequest, fmt.Errorf("the key should be 'file' for single file uploading and 'files' for multiple files uploading"))
}

func Unauthorized() error {
	return NewAPIError(http.StatusUnauthorized, fmt.Errorf("Unauthorized"))
}
