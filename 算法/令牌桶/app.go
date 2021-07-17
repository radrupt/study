package main

import (
	"fmt"
	"time"

	"golang.org/x/time/rate"
)

func MakeToken(internel int, bucket int) {

}

func main() {
	// 定期发放令牌
	limiter := rate.NewLimiter(rate.Every(100*time.Millisecond), 1)
	index := 0
	go func() {
		for {
			if limiter.Allow() { // do something
				index++
				fmt.Println(index)
			}
		}
	}()
	time.Sleep(time.Duration(time.Second * 10))
	// 消费
}
