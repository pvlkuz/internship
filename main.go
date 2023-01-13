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

func main() {

	if len(os.Args) < 2 {
		fmt.Println("expected subcommand")
		return
	}

	switch os.Args[1] {
	case "transform":
		var in io.Reader
		var out io.Writer
		var config IoConfig
		var CaesarShift int
		var UseBase64, ioinput bool

		cmd := flag.NewFlagSet("transform", flag.ExitOnError)
		cmd.BoolVar(&config.StdIN, "input_stdin", true, "stnIN")
		cmd.StringVar(&config.FileIn, "input", "default", "file input")
		cmd.BoolVar(&config.StdOut, "output_std", true, "stdOUT")
		cmd.StringVar(&config.FileOut, "output", "default", "file output")
		cmd.IntVar(&CaesarShift, "caesar", 0, "caesar")
		cmd.BoolVar(&UseBase64, "base64", false, "base64")

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
			log.Print(fmt.Errorf("error in transfroming: %w", err))
			return
		}
	case "crud":
		m, err := migrate.New("file://./migration", "postgresql://postgres:password@database:5432/postgres?sslmode=disable")
		if err != nil {
			time.Sleep(2 * time.Second)
			m, err = migrate.New("file://./migration", "postgresql://postgres:password@database:5432/postgres?sslmode=disable")
			if err != nil {
				log.Fatalf("failed to migration init: %s", err.Error())
			}
		}
		err = m.Up()
		if err != nil {
			log.Fatalf("failed to migrate up: %s", err.Error())
		}

		db, err := database.NewDB()
		if err != nil {
			time.Sleep(2 * time.Second)
			db, err = database.NewDB()
			if err != nil {
				log.Fatalf("failed to initialize db: %s", err.Error())
			}

		}

		handler := crud_handler.NewHandler(*db)
		handler.RunServer()

	}

}
