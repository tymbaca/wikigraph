package parser

import (
	"context"
	"net/http"
	urllib "net/url"
	"regexp"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/samber/lo"
	"github.com/tymbaca/wikigraph/internal/logger"
	"github.com/tymbaca/wikigraph/internal/model"
)

func NewWikiParser() *WikiParser {
	return &WikiParser{
		client: &http.Client{Timeout: 1 * time.Minute},
	}
}

type WikiParser struct {
	client *http.Client
}

var wikiLinkRegex = regexp.MustCompile(`href="(\/wiki.*?)"`)

func (w *WikiParser) Parse(ctx context.Context, url string) (model.ParsedArticle, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return model.ParsedArticle{}, err
	}

	parentURL, err := urllib.Parse(url)
	if err != nil {
		return model.ParsedArticle{}, err
	}

	resp, err := w.client.Do(req)
	if err != nil {
		return model.ParsedArticle{}, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return model.ParsedArticle{}, err
	}

	name := doc.Find("#firstHeading > span").Text()
	if len(name) == 0 {
		name = parentURL.Path
	}

	childs := make([]string, 0, 100)
	doc.Find("#bodyContent").Not(".mw-references-wrap").Each(func(i int, s *goquery.Selection) {
		text, err := s.Html()
		if err != nil {
			return
		}
		matches := wikiLinkRegex.FindAllStringSubmatch(text, -1)
		for _, match := range matches {
			childURL, err := urllib.Parse(match[1])
			if err != nil {
				return
			}

			// leave only schema, domain and path
			childURL.Fragment = ""
			childURL.RawFragment = ""
			childURL.RawQuery = ""

			var child string
			if !childURL.IsAbs() {
				childURL.Scheme = parentURL.Scheme
				childURL.Host = parentURL.Host
			}

			child = childURL.String()
			child, err = urllib.PathUnescape(child) // TODO: do we really need this?
			if err != nil {
				logger.Fatalf("can't unescape the path: %s", childURL.String())
				return
			}

			childs = append(childs, child)
		}
	})

	return model.ParsedArticle{
		Name:      name,
		URL:       url,
		ChildURLs: lo.Uniq(childs),
	}, nil
}
