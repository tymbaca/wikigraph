package workers

import (
	"context"

	"github.com/tymbaca/wikigraph/internal/model"
)

type Parser interface {
	Parse(ctx context.Context, url string) (model.ParsedArticle, error)
}

type Storage interface {
	GetURLToProcess(ctx context.Context) (string, error)
	SaveParsedArticle(ctx context.Context, article model.ParsedArticle) error

	AddPendingURLs(ctx context.Context, urls ...string) error
	ResetInProgressURLs(ctx context.Context) error
	SetFailed(ctx context.Context, url string, err error) error
	GetNotCompletedCount(ctx context.Context) (int, error)

	GetGraph(ctx context.Context) (model.Graph, error)
}
