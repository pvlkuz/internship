package main

import (
	"flag"
	"fmt"
	"os"
)

func caesar(s string, shift int) string {
	rns := []rune(s)
	for i := 0; i < len(rns); i++ {
		r := int(rns[i]) + shift
		if r > 'z' {
			rns[i] = rune(r - 26)
		} else if r < 'a' {
			rns[i] = rune(r + 26)
		} else {
			rns[i] = rune(r)
		}
	}
	return string(rns)
}

func reverse(s string) string {
	rns := []rune(s)
	for i, j := 0, len(rns)-1; i < j; i, j = i+1, j-1 {
		rns[i], rns[j] = rns[j], rns[i]
	}
	return string(rns)
}

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {

	transform := flag.NewFlagSet("transform", flag.ExitOnError)
	var useOut, useIN string
	var useC int
	transform.StringVar(&useOut, "f", "default", "file output")
	transform.StringVar(&useIN, "fi", "default", "file input")
	transform.IntVar(&useC, "c", 0, "caesar")
	html := flag.NewFlagSet("html", flag.ExitOnError)
	html.Int("port", 0, "port")

	if len(os.Args) < 2 {
		fmt.Println("expected subcommand")
		os.Exit(1)
	}

	switch os.Args[1] {

	case "transform":

		err := transform.Parse(os.Args[2:])
		checkError(err)
		if useOut != "default" {

			fmt.Println("subcommand 'transform' with flag 'file output'")
			f, err := os.Create(useOut)
			checkError(err)
			_, err2 := f.WriteString(reverse(transform.Args()[0]))
			checkError(err2)

		} else if useIN != "default" {

			fmt.Println("subcommand 'transform' with flag 'file input'")
			dat, err := os.ReadFile(useIN)
			checkError(err)
			fmt.Println(reverse(string(dat)))

		} else if useC != 0 {

			fmt.Println("subcommand 'transform' with flag 'caesar'")
			fmt.Println(caesar(transform.Args()[0], useC))

		} else if useC == 0 {

			fmt.Println("subcommand 'transform' with flag 'caesar' but shift=0")
			fmt.Println(transform.Args()[0])

		} else {

			fmt.Println("subcommand 'transform'")
			fmt.Println(reverse(transform.Args()[0]))

		}

	case "html":

		err := html.Parse(os.Args[2:])
		checkError(err)
		fmt.Println("subcommand 'html'")

	default:

		fmt.Println("expected 'transform' or 'html' subcommands")
		os.Exit(1)

	}
}
