package test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestRandom(t *testing.T) {
	rand.NewSource(time.Now().UnixNano())
	fmt.Println(rand.Int63n(12))
	fmt.Println(time.Now().Unix(), time.Unix(1699703569, 0).String())

}
