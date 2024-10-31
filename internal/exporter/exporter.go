package exporter

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"

	"github.com/schollz/progressbar/v3"
	"github.com/tymbaca/wikigraph/internal/logger"
	"github.com/tymbaca/wikigraph/internal/model"
)

type storage interface {
	GetGraph(ctx context.Context) (model.Graph, error)
}

func Export(ctx context.Context, w io.Writer, storage storage) error {
	csvWriter := csv.NewWriter(w)

	bar := progressbar.Default(-1, "Extracting graph from DB... (it might take a while)")
	graph, err := storage.GetGraph(ctx)
	if err != nil {
		logger.Fatalf("can't get the graph: %s", err)
	}
	bar.Finish()

	bar = progressbar.Default(int64(len(graph)), "Exporting...")

	csvWriter.Write([]string{"from", "to"})

	for _, article := range graph {
		for _, childID := range article.Childs {
			child, ok := graph[childID]
			if !ok {
				continue
			}

			if err := csvWriter.Write([]string{article.Name, child.Name}); err != nil {
				return fmt.Errorf("can't write to csv file: %w", err)
			}
		}
		bar.Add(1)
	}

	csvWriter.Flush()
	if err := csvWriter.Error(); err != nil {
		return fmt.Errorf("can't write to csv file (at flush): %w", err)
	}

	return nil
}
