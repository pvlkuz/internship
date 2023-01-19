package crud_handler

import (
	"fmt"
	"io"
	"log"
	"time"

	database "main/data-base"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Test_NewRecordHandlerFunc(t *testing.T) {

	db, err := database.NewDB()
	if err != nil {
		time.Sleep(2 * time.Second)
		db, err = database.NewDB()
		if err != nil {
			log.Fatalf("failed to initialize db: %s", err.Error())
		}

	}

	m, err := migrate.New("file://.././migration", "postgresql://postgres:password@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		time.Sleep(2 * time.Second)
		m, err = migrate.New("file://.././migration", "postgresql://postgres:password@localhost:5432/postgres?sslmode=disable")
		if err != nil {
			log.Fatalf("failed to migration init: %s", err.Error())
		}
	}
	err = m.Up()
	if err != nil {
		log.Fatalf("failed to migrate up: %s", err.Error())
	}

	reverseTs := httptest.NewServer(http.HandlerFunc(NewHandler(*db).NewRecord))
	defer reverseTs.Close()
	in := fmt.Sprintf("{%s:%s, %s:%s}", "'type'", "'base64'", "'input'", "'Man'")
	//t.Fatalf("in =%s", in)

	res, err := http.Post(reverseTs.URL, "application/json", strings.NewReader(in))
	if err != nil {
		log.Fatal(err)
	}

	result, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatalf("reading error")
	}

	if string(result) != "gfedcba" {
		t.Fatalf("result=%s", string(result))
	}
}
