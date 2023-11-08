package common

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

// Ok 成功
func Ok(c echo.Context, data any) error {
	return c.JSON(http.StatusOK, map[string]interface{}{"code": 200, "data": data, "msg": "success"})
}

// Fail 失败
func Fail(c echo.Context, code int, msg string) error {
	return c.JSON(http.StatusOK, map[string]interface{}{"code": code, "data": nil, "msg": msg})
}
