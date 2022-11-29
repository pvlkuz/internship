package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Transformer interface {
	Transform(in io.Reader) (string, error)
}

type CaesarTransformer struct {
	Shift int
}

func NewCaesarTransformer(shift int) *CaesarTransformer {
	return &CaesarTransformer{Shift: shift}
}

func (t *CaesarTransformer) Transform(in io.Reader) (string, error) {
	var result string
	f, err := io.ReadAll(in)
	if err != nil {
		return result, err
	}

	rns := []rune(string(f))
	for i := 0; i < len(rns); i++ {
		r := int(rns[i]) + t.Shift
		if r > 'z' {
			rns[i] = rune(r - 26)
		} else if r < 'a' {
			rns[i] = rune(r + 26)
		} else {
			rns[i] = rune(r)
		}
	}
	result = string(rns)

	return result, nil
}

type ReverseTransformer struct{}

func NewReverseTransformer() *ReverseTransformer {
	return &ReverseTransformer{}
}

func (t *ReverseTransformer) Transform(in io.Reader) (string, error) {
	var result string
	f, err := io.ReadAll(in)
	if err != nil {
		return result, err
	}

	rns := []rune(string(f))
	for i, j := 0, len(rns)-1; i < j; i, j = i+1, j-1 {
		rns[i], rns[j] = rns[j], rns[i]
	}
	result = string(rns)
	return result, nil
}

type Base64Transformer struct{}

func NewBase64Transformer() *Base64Transformer {
	return &Base64Transformer{}
}

func (t *Base64Transformer) Transform(in io.Reader) (string, error) {
	var result string
	f, err := io.ReadAll(in)
	if err != nil {
		return result, err
	}
	result = base64.StdEncoding.EncodeToString(f)
	return result, nil
}

func transform(in io.Reader, out io.Writer, C int, B bool) error {
	var result string
	var err error
	if B {
		result, err = NewBase64Transformer().Transform(in)
	} else if C != 123321 {
		result, err = NewCaesarTransformer(C).Transform(in)
	} else {
		result, err = NewReverseTransformer().Transform(in)
	}
	if err != nil {
		return fmt.Errorf("TRANSFORMER error: %w", err)
	}

	_, err = out.Write([]byte(result))
	if err != nil {
		return fmt.Errorf("write output error: %w", err)
	}

	return nil
}

func ReverseHandler(w http.ResponseWriter, r *http.Request) {
	result, err := NewReverseTransformer().Transform(r.Body)
	if err != nil {
		http.Error(w, "Server ReverseHandler Transformer error", http.StatusBadRequest)
		return
	}
	_, err = w.Write([]byte(result))
	if err != nil {
		http.Error(w, "server write output error", http.StatusBadRequest)
		return
	}
}

func CaesarHandler(w http.ResponseWriter, r *http.Request) {
	shiftStr := r.FormValue("shift")
	shift, err := strconv.Atoi(shiftStr)
	if err != nil {
		http.Error(w, "No integer shift given", http.StatusBadRequest)
		return
	}
	result, err := NewCaesarTransformer(shift).Transform(r.Body)
	if err != nil {
		http.Error(w, "Server CaesarHandler Transformer error", http.StatusBadRequest)
		return
	}
	_, err = w.Write([]byte(result))
	if err != nil {
		http.Error(w, "server write output error", http.StatusBadRequest)
		return
	}
}

func Base64Handler(w http.ResponseWriter, r *http.Request) {
	result, err := NewBase64Transformer().Transform(r.Body)
	if err != nil {
		http.Error(w, "Server Bae64Handler Transformer error", http.StatusBadRequest)
		return
	}

	_, err = w.Write([]byte(result))
	if err != nil {
		http.Error(w, "server write output error", http.StatusBadRequest)
		return
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
