package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
	"github.com/tymbaca/wikigraph/internal/exporter"
	"github.com/tymbaca/wikigraph/internal/logger"
	"github.com/tymbaca/wikigraph/internal/parser"
	"github.com/tymbaca/wikigraph/internal/storage"
	"github.com/tymbaca/wikigraph/internal/workers"
	"github.com/tymbaca/wikigraph/migrations"
	"github.com/tymbaca/wikigraph/pkg/httpx"
)

const _help = `Wikigraph is a tool for parsing wikipedia graph.

Usage:

        wikigraph parse <db_path> [initial_links...]
        wikigraph export <db_path> <export_path>

parse
        Parses wiki to the specified database file.

        <db_path> 
                Path to sqlite database file. If it exists and it already
                has some data, then it will be used to continue the progress.
                If it doesn't exists, it will be created.

        [initial_links...]
                Optional wikipedia links that will be inserted into the parsing 
                queue at the start of a program.

export
        Exports parsed graph into the CSV table that containes 2 columns: 
        parent article name -> child article name. 

        <db_path> 
                Path to sqlite database file from where graph will be exported.

        <export_path> 
                Path where resulting graph (in csv) will be exported.
`

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if len(os.Args) < 3 {
		fmt.Println(_help)
		return
	}

	switch os.Args[1] {
	case "parse":
		parseCmd(ctx)
		return
	case "export":
		if len(os.Args) != 4 {
			fmt.Println(_help)
			return
		}
		exportCmd(ctx)
		return
	default:
		fmt.Println(_help)
		return
	}
}

func parseCmd(ctx context.Context) {
	dbPath := os.Args[2]
	db := connectToDB(dbPath)
	defer db.Close()

	client := httpx.NewRateLimitingClient(&http.Client{Timeout: 1 * time.Minute}, 20, 1)
	parser := parser.NewWikiHtmlParser(client)
	storage := storage.New(db)
	workers := workers.New(10, time.Duration(50*time.Millisecond), parser, storage)

	var initialLinks []string
	if len(os.Args) >= 4 {
		initialLinks = os.Args[3:]
	}

	if err := workers.Launch(ctx, initialLinks...); err != nil {
		logger.Fatal(err.Error())
	}
}

func exportCmd(ctx context.Context) {
	dbPath := os.Args[2]
	if !exists(dbPath) {
		fmt.Printf("database file doesn't exist: %s\n", dbPath)
		os.Exit(1)
	}

	db := connectToDB(dbPath)
	defer db.Close()
	_ = db

	exportPath := os.Args[3]

	storage := storage.New(db)

	f, err := os.Create(exportPath)
	if err != nil {
		logger.Fatalf("can't create export file: %s", err)
	}

	if err := exporter.Export(ctx, f, storage); err != nil {
		logger.Fatalf("can't export the graph to csv: %s", err)
	}
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}

	return true
}

func connectToDB(path string) *sql.DB {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("can't ping database: %s", err)
	}

	// Migrate
	goose.SetDialect("sqlite3")
	goose.SetBaseFS(migrations.FS)
	err = goose.Up(db, ".")
	if err != nil {
		log.Fatalf("can't migrate database: %s", err)
	}

	// To avoid sqlite "database is locked" error
	db.SetMaxOpenConns(1)

	logger.Debug("Connected to the database!")

	return db
}
