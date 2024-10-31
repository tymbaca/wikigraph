package parser

import (
	"context"
	"errors"

	"github.com/tymbaca/wikigraph/internal/model"
	"github.com/tymbaca/wikigraph/pkg/httpx"
)

const (
	EnRegion = "en"
	RuRegion = "ru"
)

type WikiApiParser struct {
	region string
	client httpx.Client
}

func NewWikiApiParser(region string, client httpx.Client) *WikiApiParser {
	return &WikiApiParser{
		region: region,
		client: client,
	}
}

func (w *WikiApiParser) Parse(ctx context.Context, url string) (model.ParsedArticle, error) {
	return model.ParsedArticle{}, errors.New("not implemented")
}
