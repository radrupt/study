package main

import "fmt"

func RevertS(s string) string {
	sRune := []rune(s)
	l := len(sRune)
	for i := 0; i < l/2; i++ {
		sRune[i], sRune[l-i-1] = sRune[l-i-1], sRune[i]
	}
	return string(sRune)
}

func main() {
	s := "232323sdsd"
	fmt.Println(RevertS(s))
}
