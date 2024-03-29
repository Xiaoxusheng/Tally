package middleware

import (
	"Tally/common"
	"Tally/config"
	"Tally/global"
	"Tally/utils"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"strings"
)

type MyCustomClaims struct {
	Identity string `json:"identity"`
	jwt.RegisteredClaims
}

func ParseToken() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Token from another example.  This token is expired
			tokenString := c.Request().Header.Get("Authorization")
			if tokenString == "" {
				return common.Fail(c, global.VerifyCode, "token 不能为空")
			}
			t := strings.Split(tokenString, " ")

			claims := new(utils.MyCustomClaims)
			token, err := jwt.ParseWithClaims(t[len(t)-1], claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(config.Config.Jwt.Key), nil
			})
			_, errs := global.Global.Redis.Get(global.Global.Ctx, claims.Identity).Result()
			if errs != nil {
				fmt.Println(errs)
				return common.Fail(c, global.VerifyCode, "token过期或退出登录")
			}

			c.Set("identity", claims.Identity)
			fmt.Println("id", claims.Identity)

			switch {
			case token.Valid:
				fmt.Println("You look nice today")
			case errors.Is(err, jwt.ErrTokenMalformed):
				fmt.Println("That's not even a token")
				return common.Fail(c, global.VerifyCode, "token 错误")
			case errors.Is(err, jwt.ErrTokenSignatureInvalid):
				// Invalid signature
				fmt.Println("Invalid signature")
				return common.Fail(c, global.VerifyCode, "Invalid signature")
			case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
				// Token is either expired or not active yet
				fmt.Println("Timing is everything")
				return common.Fail(c, global.VerifyCode, "Timing is everything")
			default:
				fmt.Println("Couldn't handle this token:", err)
				return common.Fail(c, global.VerifyCode, "Couldn't handle this token")
			}
			return next(c)
		}
	}

}
