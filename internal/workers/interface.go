package workers

import (
	"context"

	"github.com/tymbaca/wikigraph/internal/model"
)

type Parser interface {
	Parse(ctx context.Context, url string) (model.ParsedArticle, error)
}

type Storage interface {
	ResetInProgressURLs(ctx context.Context) error
	GetNotCompletedCount(ctx context.Context) (int, error)
	GetURLToProcess(ctx context.Context) (string, error)
	GetFailedURL(ctx context.Context) (string, error)
	AddPendingURLs(ctx context.Context, urls ...string) error
	SaveParsedArticle(ctx context.Context, article model.ParsedArticle) error
	SetFailed(ctx context.Context, url string, err error) error
}
