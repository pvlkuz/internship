package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"main/cache"
	database "main/data-base"
	"main/handler"
	"main/service"
	"main/transformer"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type IoConfig struct {
	FileIn, FileOut string
}

const connStr = "postgresql://postgres:password@database:5432/postgres?sslmode=disable"

//nolint:forbidigo
func printHelp() {
	fmt.Println(` `)
	fmt.Println(`Usage:  ./bin/main COMMAND [OPTIONS]`)
	fmt.Println(` `)
	fmt.Println(`Commands:`)
	fmt.Println(`	transform \t Transform string - reversing it if no other option provided`)
	fmt.Println(`			(by default input from std.in and output to std.out)`)
	fmt.Println(` `)
	fmt.Println(`	crud \t\t Start a server listening on port 8080, and connecting to db on port 5432`)
	fmt.Println(` `)
	fmt.Println(`Options:`)
	fmt.Println(`	-input \t\t Path to input file`)
	fmt.Println(`	-output \t Path to output file`)
	fmt.Println(`	-caesar \t Run Caesar cipher, provide shift number (different from 0)`)
	fmt.Println(`	-base64 \t Run Base64 cipher`)
	fmt.Println(` `)
}

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
		printHelp()
	case "transform":
		transform()
	case "crud":
		m, err := migrate.New("file://./migration", connStr)
		if err != nil {
			time.Sleep(2 * time.Second)

			m, err = migrate.New("file://./migration", connStr)
			if err != nil {
				log.Print(fmt.Errorf("failed to migration init: %w", err))
				return
			}
		}

		err = m.Up()
		if err != nil {
			log.Print(fmt.Errorf("failed to migrate up: %w", err))
			return
		}

		db, err := database.NewDB(connStr)
		if err != nil {
			log.Print(fmt.Errorf("failed to initialize db: %w", err))
			return
		}

		cache := cache.NewLruCache(10)
		service := service.NewService(db, cache)
		myhandler := handler.NewHandler(service)

		err = myhandler.RunServer()
		if err != nil {
			log.Print(err)
		}
	}
}
