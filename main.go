package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"main/handler"
	"main/transformer"
	"net/http"
	"os"
)

func transform(in io.Reader, out io.Writer, CaesarShift int, Base64Use bool) error {
	var result string
	var err error
	var tr transformer.Transformer
	switch {
	case Base64Use:
		tr = transformer.NewBase64Transformer()
	case CaesarShift != 0:
		tr = transformer.NewCaesarTransformer(CaesarShift)
	default:
		tr = transformer.NewReverseTransformer()
	}
	result, err = tr.Transform(in)
	if err != nil {
		return fmt.Errorf("TRANSFORMER error: %w", err)
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
	var CaesarShift int
	var Base64Use bool
	var port string

	cmd := flag.NewFlagSet("transform", flag.ExitOnError)
	cmd.BoolVar(&config.StdIN, "input_stdin", true, "stnIN")
	cmd.StringVar(&config.FileIn, "input", "default", "file input")
	cmd.BoolVar(&config.StdOut, "output_std", true, "stdOUT")
	cmd.StringVar(&config.FileOut, "output", "default", "file output")
	cmd.IntVar(&CaesarShift, "caesar", 0, "caesar")
	cmd.BoolVar(&Base64Use, "base64", false, "base64")

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
		err = transform(in, out, CaesarShift, Base64Use)
		if err != nil {
			log.Print(fmt.Errorf("error in transfroming: %w", err))
			return
		}
	case "serve":
		http.HandleFunc("/reverse", handler.ReverseHandler)
		http.HandleFunc("/caesar", handler.CaesarHandler)
		http.HandleFunc("/base64", handler.Base64Handler)
		log.Fatal(http.ListenAndServe(port, nil))
	}

}
