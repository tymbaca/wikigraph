package storage

import (
	"context"
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/tymbaca/wikigraph/internal/model"
)

func (s *storage) GetGraph(ctx context.Context) (model.Graph, error) {
	var graph model.Graph
	err := s.inTx(ctx, func(ctx context.Context, tx *sql.Tx) error {
		articles, err := s.listAllArticles(ctx, tx)
		if err != nil {
			return err
		}

		relations, err := s.listAllRelations(ctx, tx)
		if err != nil {
			return err
		}

		for id, a := range articles {
			a.Childs = relations[id]
			articles[id] = a
		}

		graph = model.Graph(articles)

		return nil
	})
	if err != nil {
		return model.Graph{}, err
	}

	return graph, nil
}

func (s *storage) listAllArticles(ctx context.Context, tx *sql.Tx) (map[int]model.Article, error) {
	rows, err := squirrel.Select("id", "name", "url").Where("status = ?", _completed).From(_articleTable).RunWith(tx).QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []model.Article
	for rows.Next() {
		var a model.Article
		if err := rows.Scan(&a.ID, &a.Name, &a.URL); err != nil {
			return nil, err
		}

		articles = append(articles, a)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	result := make(map[int]model.Article, len(articles))
	for _, a := range articles {
		result[a.ID] = a
	}

	return result, nil
}

func (s *storage) listAllRelations(ctx context.Context, tx *sql.Tx) (map[int][]int, error) {
	rows, err := squirrel.Select("from_id", "to_id").From(_relationTable).RunWith(tx).QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pairs [][2]int
	for rows.Next() {
		var from, to int
		if err := rows.Scan(&from, &to); err != nil {
			return nil, err
		}

		pairs = append(pairs, [2]int{from, to})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	relations := make(map[int][]int)
	for _, pair := range pairs {
		from, to := pair[0], pair[1]
		relations[from] = append(relations[from], to)
	}

	return relations, nil
}
