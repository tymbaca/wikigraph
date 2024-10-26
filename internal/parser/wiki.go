package parser

import (
	"context"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/PuerkitoBio/goquery"
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

func (w *WikiParser) ParseChilds(ctx context.Context, parentURL string) ([]string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, parentURL, nil)
	if err != nil {
		return nil, err
	}

	baseURL, err := url.Parse(parentURL)
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
		text := s.Text()
		matches := wikiLinkRegex.FindAllStringSubmatch(text, -1)
		for _, match := range matches {
			child := match[1]
			child = "https://" + baseURL.Host + child
			childs = append(childs, child)
		}
	})

	return childs, nil
}
