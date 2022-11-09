// go mod init /Users/pavlo.kuzmin/Documents/Internship/1

package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Fprintln(os.Stdout, "hello, world")
	fmt.Println("HELLO, WORLD")
}
