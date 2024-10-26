package parser

import (
	"context"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/samber/lo"
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

func (w *WikiParser) ParseChilds(ctx context.Context, parent string) ([]string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, parent, nil)
	if err != nil {
		return nil, err
	}

	parentURL, err := url.Parse(parent)
	if err != nil {
		return nil, err
	}

	resp, err := w.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	childs := make([]string, 0, 100)
	doc.Find("#bodyContent").Not(".mw-references-wrap").Each(func(i int, s *goquery.Selection) {
		text, err := s.Html()
		if err != nil {
			return
		}
		matches := wikiLinkRegex.FindAllStringSubmatch(text, -1)
		for _, match := range matches {
			childURL, err := url.Parse(match[1])
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

			childs = append(childs, child)
		}
	})

	return lo.Uniq(childs), nil
}
