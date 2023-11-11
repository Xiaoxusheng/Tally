package utils

import (
	"github.com/labstack/echo/v4"
	"math/rand"
	"time"
)

func GetRandom(m int64) int64 {
	rand.NewSource(time.Now().UnixNano())
	n := rand.Int63n(m)
	if n == 0 {
		n = rand.Int63n(24)
	}
	return n
}

func GetIdentity(c echo.Context, s string) string {
	str, ok := c.Get(s).(string)
	if !ok {
		return ""
	}
	return str
}
