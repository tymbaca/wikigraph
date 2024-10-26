package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
	"github.com/tymbaca/wikigraph/internal/logger"
	"github.com/tymbaca/wikigraph/internal/parser"
	"github.com/tymbaca/wikigraph/internal/storage"
	"github.com/tymbaca/wikigraph/internal/workers"
	"github.com/tymbaca/wikigraph/migrations"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	db, err := sql.Open("sqlite3", "./example.db") // TODO flag
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Test the connection
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	goose.SetDialect("sqlite3")
	goose.SetBaseFS(migrations.FS)
	err = goose.Up(db, ".")
	if err != nil {
		log.Fatal(err)
	}

	logger.Info("Connected to the database!")

	parser := parser.NewWikiParser()
	storage := storage.New(db)
	workers := workers.New(10, time.Duration(50*time.Millisecond), parser, storage)

	workers.Launch(ctx, `https://ru.wikipedia.org/wiki/%d0%9c%d0%b8%d1%84`)
}

type fakeAPI struct{}

func (a fakeAPI) ParseChilds(ctx context.Context, url string) ([]string, error) {
	var childs []string
	time.Sleep(500 + time.Duration(gofakeit.IntN(500))*time.Millisecond)

	if gofakeit.Int()%20 == 0 {
		return nil, errors.New("bad error from api")
	}

	for range gofakeit.IntN(20) {
		childs = append(childs, gofakeit.URL())
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	return childs, nil
}
