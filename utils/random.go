package utils

import (
	"math/rand"
	"time"
)

func GetRandom() int64 {
	rand.NewSource(time.Now().UnixNano())
	n := rand.Int63n(24)
	if n == 0 {
		n = rand.Int63n(24)
	}
	return n
}
