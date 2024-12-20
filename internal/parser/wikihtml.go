package parser

import (
	"context"
	"fmt"
	"net/http"
	urllib "net/url"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/samber/lo"
	"github.com/tymbaca/wikigraph/internal/logger"
	"github.com/tymbaca/wikigraph/internal/model"
	"github.com/tymbaca/wikigraph/pkg/httpx"
	"golang.org/x/net/html"
)

func NewWikiHtmlParser(client httpx.Client) *WikiHtmlParser {
	return &WikiHtmlParser{
		client: client,
	}
}

type WikiHtmlParser struct {
	client httpx.Client
}

var (
	_wikiLinkRegex   = regexp.MustCompile(`href=["'](\/wiki.*?)["']`)
	_ignoredSuffixes = []string{
		".svg", ".png", ".gif", ".jpg", ".jpeg", ".webp",
		".SVG", ".PNG", ".GIF", ".JPG", ".JPEG", ".WEBP",
	}
)

func (w *WikiHtmlParser) Parse(ctx context.Context, url string) (model.ParsedArticle, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return model.ParsedArticle{}, err
	}

	// oldPath := req.URL.Path
	// req.URL.Path = req.URL.EscapedPath()
	// logger.Debugf("old path %s, new path %s", oldPath, req.URL.Path)

	parentURL, err := urllib.Parse(url)
	if err != nil {
		return model.ParsedArticle{}, err
	}

	resp, err := w.client.Do(req)
	if err != nil {
		return model.ParsedArticle{}, err
	}
	defer resp.Body.Close()

	root, err := html.Parse(resp.Body)
	if err != nil {
		return model.ParsedArticle{}, err
	}

	if resp.StatusCode >= 400 {
		_ = logToFile(path.Base(url), root)
		return model.ParsedArticle{}, fmt.Errorf("got error from wikipedia (%s), status code %d (%s)", url, resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	doc := goquery.NewDocumentFromNode(root)
	if err != nil {
		return model.ParsedArticle{}, err
	}

	nameSelection := doc.Find("#firstHeading")
	name := nameSelection.Text()
	// logger.Debug(name)
	if len(name) == 0 {
		logger.Warnf("can't find article name, url: %s", url)
		_ = logToFile(path.Base(url), root)

		name = path.Base(parentURL.Path)
	}

	childs := make([]string, 0, 100)
	doc.Find("#bodyContent").Not(".mw-references-wrap").Each(func(i int, s *goquery.Selection) {
		text, err := s.Html()
		if err != nil {
			return
		}
		matches := _wikiLinkRegex.FindAllStringSubmatch(text, -1)
		for _, match := range matches {
			childURL, err := urllib.Parse(match[1])
			if err != nil {
				continue
			}

			// leave only schema, domain and path
			childURL.Fragment = ""
			childURL.RawFragment = ""
			childURL.RawQuery = ""

			if !childURL.IsAbs() {
				childURL.Scheme = parentURL.Scheme
				childURL.Host = parentURL.Host
			}

			child := childURL.String()
			// child, err = urllib.PathUnescape(child)
			// if err != nil {
			// 	logger.Fatalf("can't unescape the path: %s", childURL.String())
			// 	return
			// }

			if hasSuffixes(child, _ignoredSuffixes...) {
				continue
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

func logToFile(name string, root *html.Node) error {
	f, err := os.Create("debug/" + name + ".html")
	if err != nil {
		return err
	}

	return html.Render(f, root)
}

func hasSuffixes(s string, suffixes ...string) bool {
	for _, suffix := range suffixes {
		if strings.HasSuffix(s, suffix) {
			return true
		}
	}

	return false
}
