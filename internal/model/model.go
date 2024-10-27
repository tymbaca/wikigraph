package model

type Graph map[int]Article

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
