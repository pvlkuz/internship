package main

import (
	b64 "encoding/base64"
	"flag"
	"fmt"
	"io"
	"log"
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

func transform(in io.Reader, out io.Writer, C int, B bool) error {
	f, err := io.ReadAll(in)
	if err != nil {
		return fmt.Errorf("read input error: %w", err)
	}

	var result string
	if B {
		result = b64.URLEncoding.EncodeToString(f)
	} else if C != 123321 {
		result = caesar(string(f), C)
	} else {
		result = reverse(string(f))
	}

	_, err = out.Write([]byte(result))
	if err != nil {
		return fmt.Errorf("write output error: %w", err)
	}

	return nil
}

type Config struct {
	StdIN, StdOut   bool
	FileIn, FileOut string
}

func main() {

	var in io.Reader
	var out io.Writer
	var config Config
	var useC int
	var useB bool

	cmd := flag.NewFlagSet("transform", flag.ExitOnError)
	cmd.BoolVar(&config.StdIN, "input_stdin", true, "stnIN")
	cmd.StringVar(&config.FileIn, "input", "default", "file input")
	cmd.BoolVar(&config.StdOut, "output_std", true, "stdOUT")
	cmd.StringVar(&config.FileOut, "output", "default", "file output")
	cmd.IntVar(&useC, "c", 123321, "caesar")
	cmd.BoolVar(&useB, "base64", false, "base64")

	if len(os.Args) < 2 {
		fmt.Println("expected subcommand")
		return
	}

	err := cmd.Parse(os.Args[2:])
	if err != nil {
		log.Print(fmt.Errorf("error in persing flags: %w", err))
		return
	}

	if config.FileIn != "default" {
		f, err := os.Open(config.FileIn)
		if err != nil {
			log.Print(fmt.Errorf("error in opening input file: %w", err))
		}
		defer f.Close()
		in = f
	} else {
		in = os.Stdin
	}

	if config.FileOut != "default" {
		f, err := os.Create(config.FileOut)
		if err != nil {
			log.Print(fmt.Errorf("error in creating output file: %w", err))
		}
		out = f
	} else {
		out = os.Stdout
	}

	switch os.Args[1] {
	case "transform":
		err = transform(in, out, useC, useB)
		if err != nil {
			log.Print(fmt.Errorf("error in transfroming: %w", err))
			return
		}

	}

}
