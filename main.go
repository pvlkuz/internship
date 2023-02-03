package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"main/crud_handler"
	database "main/data-base"
	"main/transformer"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type IoConfig struct {
	StdIN, StdOut   bool
	FileIn, FileOut string
}

const connStr = "postgresql://postgres:password@database:5432/postgres?sslmode=disable"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("expected command, type help for details")
		return
	}

	switch os.Args[1] {
	default:
		fmt.Println("type help for details")

	case "help":
		fmt.Println(" ")
		fmt.Println("Usage:  ./bin/main COMMAND [OPTIONS]")
		fmt.Println(" ")
		fmt.Println("Commands:")
		fmt.Println("	transform \t Transform string - reversing it if no other option provided (by default input from std.in and output to std.out)")
		fmt.Println("	crud \t\t Start a server listening on port 8080, and connecting to db on port 5432 (use docker-compose to start app and database together)")
		fmt.Println(" ")
		fmt.Println("Options:")
		fmt.Println("	-input \t\t Path to input file")
		fmt.Println("	-output \t Path to output file")
		fmt.Println("	-caesar \t Run Caesar cipher, provide shift number (different from 0)")
		fmt.Println("	-base64 \t Run Base64 cipher")
		fmt.Println(" ")

	case "transform":
		var in io.Reader
		var out io.Writer
		var config IoConfig
		var CaesarShift int
		var UseBase64, ioinput bool

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

		if config.FileIn != "default" {
			f, err := os.Open(config.FileIn)
			if err != nil {
				log.Print(fmt.Errorf("error in opening input file: %w", err))
			}
			defer f.Close()
			in = f
		} else {
			in = os.Stdin
			ioinput = true
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

		err = transformer.BasicTransform(in, out, CaesarShift, UseBase64, ioinput)
		if err != nil {
			log.Print(fmt.Errorf("error in transforming: %w", err))
			return
		}
	case "crud":
		m, err := migrate.New("file://./migration", connStr)
		if err != nil {
			time.Sleep(2 * time.Second)
			m, err = migrate.New("file://./migration", connStr)
			if err != nil {
				log.Print(fmt.Errorf("failed to migration init: %s", err.Error()))
				return
			}
		}
		err = m.Up()
		if err != nil {
			log.Print(fmt.Errorf("failed to migrate up: %s", err.Error()))
			return
		}

		db, err := database.NewDB(connStr)
		if err != nil {
			time.Sleep(2 * time.Second)
			db, err = database.NewDB(connStr)
			if err != nil {
				// log.Fatalf("failed to initialize db: %s", err.Error())
				log.Print(fmt.Errorf("failed to initialize db: %s", err.Error()))
				return
			}
		}
		handler := crud_handler.NewHandler(db)
		handler.RunServer()
	}
}
