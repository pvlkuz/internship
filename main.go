package main

import (
	b64 "encoding/base64"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
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
		result = b64.StdEncoding.EncodeToString(f)
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

func ReverseHandler(w http.ResponseWriter, r *http.Request) {

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "could not read request body", http.StatusBadRequest)
	}
	_, err = w.Write([]byte(reverse(string(reqBody))))
	if err != nil {
		http.Error(w, "write output error", http.StatusBadRequest)
	}
}

func CaesarHandler(w http.ResponseWriter, r *http.Request) {
	shiftStr := r.FormValue("shift")
	shift, err := strconv.Atoi(shiftStr)
	if err != nil {
		http.Error(w, "No integer shift given", http.StatusBadRequest)
		return
	}
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "could not read request body", http.StatusBadRequest)
	}
	_, err = w.Write([]byte(caesar(string(reqBody), shift)))
	if err != nil {
		http.Error(w, "write output error", http.StatusBadRequest)
	}
}

func Base64Handler(w http.ResponseWriter, r *http.Request) {
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "could not read request body", http.StatusBadRequest)
	}
	_, err = w.Write([]byte(b64.StdEncoding.EncodeToString(reqBody)))
	if err != nil {
		http.Error(w, "write output error", http.StatusBadRequest)
	}
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
	var port string

	cmd := flag.NewFlagSet("transform", flag.ExitOnError)
	cmd.BoolVar(&config.StdIN, "input_stdin", true, "stnIN")
	cmd.StringVar(&config.FileIn, "input", "default", "file input")
	cmd.BoolVar(&config.StdOut, "output_std", true, "stdOUT")
	cmd.StringVar(&config.FileOut, "output", "default", "file output")
	cmd.IntVar(&useC, "c", 123321, "caesar")
	cmd.BoolVar(&useB, "base64", false, "base64")

	cmd1 := flag.NewFlagSet("serve", flag.ExitOnError)
	cmd1.StringVar(&port, "port", ":8080", "server port")

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
	case "serve":
		http.HandleFunc("/reverse", ReverseHandler)
		http.HandleFunc("/caesar", CaesarHandler)
		http.HandleFunc("/base64", Base64Handler)
		log.Fatal(http.ListenAndServe(port, nil))
	}

}
