package middleware

import (
	"Tally/common"
	"Tally/config"
	"Tally/global"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"net/http"
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
			tokenString, ok := c.Get("Authorization").(string)
			if !ok {
				return common.Fail(c, http.StatusOK, "Authorization header not found")
			}
			t := strings.Split(tokenString, " ")
			claims := new(MyCustomClaims)
			token, err := jwt.ParseWithClaims(t[len(t)-1], claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(config.Config.Jwt.Key), nil
			})
			_, err = global.Global.Redis.Get(global.Global.Ctx, claims.Identity).Result()
			if err != nil {
				return common.Fail(c, http.StatusOK, "token 过期")
			}

			switch {
			case token.Valid:
				fmt.Println("You look nice today")
			case errors.Is(err, jwt.ErrTokenMalformed):
				fmt.Println("That's not even a token")
			case errors.Is(err, jwt.ErrTokenSignatureInvalid):
				// Invalid signature
				fmt.Println("Invalid signature")
			case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
				// Token is either expired or not active yet
				fmt.Println("Timing is everything")
			default:
				fmt.Println("Couldn't handle this token:", err)
			}
			return nil

		}
	}

}

/*return func(c echo.Context) error {
	// Token from another example.  This token is expired
	tokenString, ok := c.Get("Authorization").(string)
	if !ok {
		return common.Fail(c, http.StatusOK, "Authorization header not found")
	}
	t := strings.Split(tokenString, " ")
	claims := new(MyCustomClaims)
	token, err := jwt.ParseWithClaims(t[len(t)-1], claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Config.Jwt.Key), nil
	})
	_, err = global.Global.Redis.Get(global.Global.Ctx, claims.Identity).Result()
	if err != nil {
		return common.Fail(c, http.StatusOK, "token 过期")
	}

	switch {
	case token.Valid:
		fmt.Println("You look nice today")
	case errors.Is(err, jwt.ErrTokenMalformed):
		fmt.Println("That's not even a token")
	case errors.Is(err, jwt.ErrTokenSignatureInvalid):
		// Invalid signature
		fmt.Println("Invalid signature")
	case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
		// Token is either expired or not active yet
		fmt.Println("Timing is everything")
	default:
		fmt.Println("Couldn't handle this token:", err)
	}
	return nil
}*/
