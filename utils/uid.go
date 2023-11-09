package utils

import (
	"crypto/md5"
	"fmt"
	uuid "github.com/satori/go.uuid"
)

func GetUidV5(u string) string {
	return uuid.NewV5(uuid.NewV4(), u).String()
}

func GetUidV4() string {
	return uuid.NewV4().String()
}

func Md5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}
