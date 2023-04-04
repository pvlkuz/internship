package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"main/cache"
	database "main/data-base"
	"main/httpserver"
	"main/service"
	"main/transformer"
	"os"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type IoConfig struct {
	FileIn, FileOut string
}

const connStr = "postgresql://postgres:password@database:5432/postgres?sslmode=disable"
const help = `
Usage:  ./bin/main COMMAND [OPTIONS]

Commands:
	transform \t Transform string - reversing it if no other option provided
	(by default input from std.in and output to std.out)

	crud \t\t Start a server listening on port 8080, and connecting to db on port 5432

Options:
	-input \t\t Path to input file
	-output \t Path to output file
	-caesar \t Run Caesar cipher, provide shift number (different from 0)
	-base64 \t Run Base64 cipher \n`

func ioSetup(config IoConfig) (io.Reader, io.Writer, bool) {
	var (
		in      io.Reader
		out     io.Writer
		ioInput bool
	)

	if config.FileIn == "default" {
		in = os.Stdin
		ioInput = true
	} else {
		f, err := os.Open(config.FileIn)
		if err != nil {
			log.Print(fmt.Errorf("error in opening input file: %w", err))
		}
		defer f.Close()
		in = f
	}

	if config.FileOut == "default" {
		out = os.Stdout
	} else {
		f, err := os.Create(config.FileOut)
		if err != nil {
			log.Print(fmt.Errorf("error in creating output file: %w", err))
		}

		out = f
	}

	return in, out, ioInput
}

func transform() {
	var (
		in                 io.Reader
		out                io.Writer
		config             IoConfig
		CaesarShift        int
		UseBase64, ioInput bool
	)

	cmd := flag.NewFlagSet("transform", flag.ExitOnError)
	cmd.StringVar(&config.FileIn, "input", "default", "Path to file input")
	cmd.StringVar(&config.FileOut, "output", "default", "Path to file output")
	cmd.IntVar(&CaesarShift, "caesar", 0, "Run Caesar cipher with provided shift")
	cmd.BoolVar(&UseBase64, "base64", false, "Run Base64 ")

	err := cmd.Parse(os.Args[2:])
	if err != nil {
		log.Print(fmt.Errorf("error in persing flags: %w", err))
		return
	}

	in, out, ioInput = ioSetup(config)

	err = transformer.BasicTransform(in, out, CaesarShift, UseBase64, ioInput)
	if err != nil {
		log.Print(fmt.Errorf("error in transforming: %w", err))
		return
	}
}

//nolint:forbidigo
func main() {
	if len(os.Args) < 2 {
		fmt.Println("expected command, type help for details")
		return
	}

	switch os.Args[1] {
	default:
		fmt.Println("type help for details")

	case "help":
		fmt.Println(help)
	case "transform":
		transform()
	case "crud":
		db, err := database.NewDB(connStr)
		if err != nil {
			log.Print(fmt.Errorf("failed to initialize db: %w", err))
			return
		}

		err = db.Migr(connStr)
		if err != nil {
			log.Print(fmt.Errorf("failed to migrate up: %w", err))
			return
		}

		cache := cache.NewLruCache(10)
		service := service.NewService(db, cache)
		myhandler := httpserver.NewHandler(service)

		err = myhandler.RunServer()
		if err != nil {
			log.Print(err)
		}
	}
}
