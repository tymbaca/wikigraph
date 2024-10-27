package main

import (
	"context"
	"database/sql"
	"log"

	rl "github.com/gen2brain/raylib-go/raylib"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
	"github.com/tymbaca/wikigraph/internal/graphviz"
	"github.com/tymbaca/wikigraph/internal/logger"
	"github.com/tymbaca/wikigraph/internal/storage"
	"github.com/tymbaca/wikigraph/migrations"
)

func main() {
	ctx := context.Background()
	// ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	// defer cancel()

	db, err := sql.Open("sqlite3", "./wolof-wiki.db") // TODO flag
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

	db.SetMaxOpenConns(1)

	// _, err = db.Exec("DELETE FROM url; DELETE FROM relation;")
	// if err != nil {
	// 	panic(err)
	// }

	logger.Info("Connected to the database!")

	storage := storage.New(db)

	logger.Info("Getting the graph..")
	graph, err := storage.GetGraph(ctx)
	if err != nil {
		panic(err)
	}

	logger.Info("Generating the layout..")
	layout := graphviz.ForceDirLayout(graph)
	logger.Debugf("%d", len(layout))

	// Render stuff
	rl.InitWindow(1000, 1000, "viz")

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		for id := range graph {
			pos := layout[id]
			logger.Debugf("Drawing %d node in graph, pos: %v", id, pos)
			rl.DrawCircle(int32(pos.X), int32(pos.Y), 3, rl.White)
		}

		rl.EndDrawing()
	}
}
