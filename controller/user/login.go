package user

import (
	"Tally/common"
	"github.com/labstack/echo/v4"
)

func Login(c echo.Context) error {

	return common.Ok(c, "ping")
}
