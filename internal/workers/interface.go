package workers

import "context"

type Parser interface {
	ParseChilds(ctx context.Context, url string) ([]string, error)
}

type Storage interface {
	ResetInProgressURLs(ctx context.Context) error
	GetNotCompletedCount(ctx context.Context) (int, error)
	GetURLToProcess(ctx context.Context) (string, error)
	GetFailedURL(ctx context.Context) (string, error)
	AddPendingURLs(ctx context.Context, urls ...string) error
	SaveChildURLs(ctx context.Context, url string, childs []string) error
	SetFailed(ctx context.Context, url string, err error) error
}
