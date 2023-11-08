package utils

import (
	"Tally/config"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var mySigningKey = []byte(config.Config.Jwt.Key)

type MyCustomClaims struct {
	Identity string `json:"identity"`
	jwt.RegisteredClaims
}

func GetToken(identity string) {
	// Create claims with multiple fields populated
	claims := MyCustomClaims{
		"bar",
		jwt.RegisteredClaims{
			// A usual scenario is to set the expiration time relative to the current time
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.Config.Jwt.Time * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	fmt.Printf("foo: %v\n", claims.Identity)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	fmt.Println(ss, err)
}
