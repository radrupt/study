package main

import (
	"fmt"
	"strconv"
	"time"
)

func PrintA_Z(nCh, azCh chan int) {
	for {
		for i := 0; i < 26; i++ {
			<-azCh
			fmt.Print(string(rune(i + 65)))
			time.Sleep(time.Duration(time.Second * 1))
			nCh <- 1
		}
	}
}

func PrintNum(nCh, azCh chan int) {
	i := 1
	for {
		<-nCh
		fmt.Print(strconv.Itoa(i))
		time.Sleep(time.Duration(time.Second * 1))
		i++
		azCh <- 1
	}
}

func main() {
	// 65~90, 0~25, mod 26
	nCh := make(chan int)
	azCh := make(chan int)
	go PrintA_Z(nCh, azCh)
	go PrintNum(nCh, azCh)
	nCh <- 1
	time.Sleep(time.Duration(time.Second * 100))
}
