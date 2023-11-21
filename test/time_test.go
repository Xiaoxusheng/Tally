package test

import (
	"fmt"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	times := time.NewTicker(time.Second * 5)
	i := 0
	for {
		select {
		case <-times.C:
			i++
			fmt.Println("time out", i)
		}
	}

}
