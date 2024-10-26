package model

type Article struct {
	ID     int
	Name   string
	URL    string
	Childs []int
}

type ParsedArticle struct {
	Name      string
	URL       string
	ChildURLs []string
}
