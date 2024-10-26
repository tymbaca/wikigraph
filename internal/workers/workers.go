// Package workers collects url graph from provided API and stores it into provided storage
package workers

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/tymbaca/wikigraph/internal/errs"
	"github.com/tymbaca/wikigraph/internal/logger"
)

func New(workers int, rate time.Duration, parser Parser, storage Storage) *Workers {
	return &Workers{
		workers: workers,
		rate:    rate,
		parser:  parser,
		storage: storage,
	}
}

type Workers struct {
	done    atomic.Int32
	workers int
	rate    time.Duration
	parser  Parser
	storage Storage
}

func (ws *Workers) Launch(ctx context.Context, initialUrls ...string) error {
	if err := ws.storage.AddPendingURLs(ctx, initialUrls...); err != nil {
		return err
	}

	if err := ws.storage.ResetInProgressURLs(ctx); err != nil {
		return err
	}

	ws.done.Store(0)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for id := range ws.workers {
		go ws.runWorker(ctx, id+1, ws.parseURL)
	}

	ws.waitAndShutdown(ctx, cancel)

	return nil
}

func (ws *Workers) waitAndShutdown(ctx context.Context, _ context.CancelFunc) {
	<-ctx.Done()
	logger.Info("context canceled, starting to shutdown")
	start := time.Now()
	defer func() {
		logger.Infof("shutdown completed, exiting, time elapced: %s", time.Since(start))
	}()
	// wait until all workers are done, or exit
	t := time.NewTimer(20 * time.Second)
	for {
		if ws.done.Load() == int32(ws.workers) {
			return
		}
		if ws.done.Load() > int32(ws.workers) {
			logger.Fatal("shutdown: done counter is bigger than worker count")
		}

		select {
		case <-t.C:
			logger.Warn("shotdown: wait time for workers exceeded, exiting")
			return
		case <-time.After(10 * time.Millisecond):
		}

	}

	// TODO shutdown if all urls are completed
}

func (ws *Workers) runWorker(ctx context.Context, workerID int, handler func(ctx context.Context, workerID int) error) {
	for {
		select {
		case <-ctx.Done():
			ws.done.Add(1)
			return
		case <-time.After(ws.rate):
			if err := handler(ctx, workerID); err != nil {
				logger.Errorf("%s", err)
			}
		}
	}
}

func (ws *Workers) retryFailedURL(ctx context.Context, workerID int) error {
	// url :=
	return errors.New("not implemented")
}

func (ws *Workers) parseURL(_ context.Context, workerID int) error {
	// detach from main context
	// ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	// defer cancel()
	ctx := context.Background()

	url, err := ws.storage.GetURLToProcess(ctx)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			logger.Infof("worker %d: no pending urls", workerID)
			return nil
		}

		return fmt.Errorf("worker %d: got error when getting pending url: %s", workerID, err)
	}

	logger.Infof("worker %d: processing url: %s", workerID, url)

	article, err := ws.parser.Parse(ctx, url)
	if err != nil {
		if failErr := ws.storage.SetFailed(ctx, url, err); failErr != nil {
			logger.Errorf("worker %d: can't set url status to failed, url: %s, err: %s", workerID, url, failErr)
		}
		return fmt.Errorf("worker %d: got error when parsing url: %s", workerID, err)
	}

	err = ws.storage.SaveParsedArticle(ctx, article)
	if err != nil {
		return fmt.Errorf("worker %d: got error when saving url childs to storage: %s", workerID, err)
	}

	return nil
}
