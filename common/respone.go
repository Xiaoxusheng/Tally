package common

import (
	"errors"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"os"
)

// Ok 成功
func Ok(c echo.Context, data any) error {
	return c.JSON(http.StatusOK, map[string]interface{}{"code": http.StatusOK, "data": data, "msg": "success"})
}

// Fail 失败
func Fail(c echo.Context, code int, msg string) error {
	return c.JSON(http.StatusOK, map[string]interface{}{"code": code, "data": nil, "msg": msg})
}

func Html(c echo.Context, code int, msg string) error {
	return c.HTML(code, msg)
}

// Picture 发送图片
func Picture(c echo.Context, file string) error {
	files, err := os.Open(file)
	if err != nil {
		return err
	}
	fileInfo, err := files.Stat()
	if err != nil {
		return err
	}
	if fileInfo.Size() > 5<<10<<10 {
		return errors.New("size is exceed")
	}

	_, err = io.Copy(c.Response(), files)
	if err != nil {
		return err
	}
	return nil
}
