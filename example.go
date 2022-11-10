// go mod init /Users/pavlo.kuzmin/Documents/Internship/1

package main

import (
	"fmt"
	"os"
)

func reverse(s *string) string {
	rns := []rune(*s)
	for i, j := 0, len(rns)-1; i < j; i, j = i+1, j-1 {
		rns[i], rns[j] = rns[j], rns[i]
	}
	return string(rns)
}

func main() {
	//fmt.Fprintln(os.Stdout, "hello, world")

	arg := os.Args[1]
	fmt.Println(reverse(&arg))
}
