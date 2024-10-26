package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/tymbaca/wikigraph/internal/errs"
	"github.com/tymbaca/wikigraph/internal/logger"
)

type storage struct {
	db *sql.DB
}

func New(db *sql.DB) *storage {
	return &storage{
		db: db,
	}
}

func (s *storage) ResetInProgressURLs(ctx context.Context) error {
	_, err := squirrel.Update(_articleTable).Set("status", _pending).Where(squirrel.Eq{"status": []status{_inProgress, _failed}}).RunWith(s.db).ExecContext(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *storage) GetNotCompletedCount(ctx context.Context) (int, error) {
	qb := squirrel.Select("COUNT(*)").From(_articleTable).Where("status != ?", _completed)

	var count int
	if err := qb.RunWith(s.db).QueryRowContext(ctx).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func (s *storage) GetURLToProcess(ctx context.Context) (string, error) {
	subq, args := squirrel.Select("id").From(_articleTable).
		Where("status = ?", _pending).
		OrderBy("RANDOM()").
		Limit(1).MustSql()
	// logger.Debug(subq)
	qb := squirrel.Update(_articleTable).Set("status", _inProgress).Where("id = ("+subq+")", args...)
	q, args := qb.MustSql()
	logger.Debugf("sql: %s | args: %v", q, args)

	var url string
	err := qb.RunWith(s.db).QueryRowContext(ctx).Scan(&url)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errs.ErrNotFound
		}
		return "", err
	}

	return url, nil
}

func (s *storage) GetFailedURL(ctx context.Context) (string, error) {
	qb := squirrel.Select("url").From(_articleTable).Where("status = ?", _failed)

	var url string
	err := qb.RunWith(s.db).QueryRowContext(ctx).Scan(&url)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errs.ErrNotFound
		}
		return "", err
	}

	return url, nil
}

func (s *storage) AddPendingURLs(ctx context.Context, urls ...string) error {
	if len(urls) == 0 {
		return nil
	}

	qb := squirrel.Insert(_articleTable).Columns("url")
	for _, url := range urls {
		qb = qb.Values(url)
	}

	qb = qb.Suffix("ON CONFLICT (url) DO NOTHING")
	_, err := qb.RunWith(s.db).ExecContext(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *storage) SaveChildURLs(ctx context.Context, parent string, childs []string) error {
	return s.inTx(ctx, func(ctx context.Context, tx *sql.Tx) error {
		// Set status to COMPLETED
		_, err := squirrel.Update(_articleTable).Set("status", _completed).
			Where(squirrel.Eq{"url": parent}).
			RunWith(tx).
			ExecContext(ctx)
		if err != nil {
			return err
		}

		// We don't need to create relations nor new urls if there is no childs
		if len(childs) == 0 {
			return nil
		}

		// Insert urls and get IDs
		iqb := squirrel.Insert(_articleTable).Columns("url")
		for _, child := range childs {
			iqb = iqb.Values(child)
		}
		iqb = iqb.Suffix("ON CONFLICT (url) DO UPDATE SET url = EXCLUDED.url RETURNING id")

		rows, err := iqb.RunWith(tx).QueryContext(ctx)
		if err != nil {
			return err
		}
		defer rows.Close()

		var childIDs []int
		for rows.Next() {
			var id int
			if err := rows.Scan(&id); err != nil {
				return err
			}

			childIDs = append(childIDs, id)
		}
		if err := rows.Err(); err != nil {
			return err
		}

		parentID, err := s.getIDInTx(ctx, tx, parent)
		if err != nil {
			return err
		}

		// Insert relations
		iqb = squirrel.Insert("relation").Columns("from_id", "to_id")
		for _, childID := range childIDs {
			iqb = iqb.Values(parentID, childID)
		}

		if _, err := squirrel.ExecContextWith(ctx, tx, iqb); err != nil {
			return err
		}

		return nil
	})
}

func (s *storage) SetFailed(ctx context.Context, url string, failErr error) error {
	_, err := squirrel.Update(_articleTable).
		Set("status", _failed).
		Set("error", failErr.Error()).
		Where(squirrel.Eq{"url": url}).
		RunWith(s.db).Exec()
	if err != nil {
		return err
	}

	return nil
}

func (s *storage) getIDInTx(ctx context.Context, tx *sql.Tx, url string) (int, error) {
	qb := squirrel.Select("id").From(_articleTable).Where("url = ?", url)

	var id int
	if err := qb.RunWith(tx).QueryRowContext(ctx).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (s *storage) inTx(ctx context.Context, fn func(ctx context.Context, tx *sql.Tx) error) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	if err = fn(ctx, tx); err != nil {
		must(tx.Rollback())
		return err
	}

	must(tx.Commit())

	return nil
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
