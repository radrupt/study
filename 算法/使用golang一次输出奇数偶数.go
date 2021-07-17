package main

import (
	"fmt"
	"sync"
)

func main() {
	jiCh := make(chan int)
	ouCh := make(chan int)
	done := make(chan int)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for {
			select {
			case <-done:
				wg.Done()
				return
			default:
				n := <-jiCh
				fmt.Println(n)
				ouCh <- n + 1
			}
		}
	}()
	go func() {
		for {
			select {
			case <-done:
				wg.Done()
				return
			default:
				n := <-ouCh
				fmt.Println(n)
				if n > 20 {
					close(done)
				}
				jiCh <- n + 1
			}
		}
	}()

	jiCh <- 1
	wg.Wait()
}
