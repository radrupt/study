package main

import (
	"context"
	"fmt"
	"time"
	"token/uber"

	"golang.org/x/time/rate"
)

func MakeToken(internel int, bucket int) {

}

func main() {
	// 定期发放令牌
	limiter := rate.NewLimiter(rate.Every(100*time.Millisecond), 1) // 每100毫秒生成1个，桶容量是10个，即1秒可以将桶放满
	index := 0
	fmt.Println(time.Now())
	c, _ := context.WithDeadline(context.Background(), time.Now().Add(time.Duration(time.Second*5)))
	go func() {
		for {
			if err := limiter.WaitN(c, 1); err != nil { // do something
				fmt.Println(err)
				return
			}
			index++
			fmt.Println(index)
		}
	}()
	time.Sleep(time.Duration(time.Second * 10))
	fmt.Println(time.Now())
	uber.Uber()
	// 消费
}
