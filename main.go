package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

import _ "embed"

//go:embed recent-publications.sparql
var query string

func main() {
	client := &http.Client{}
	req, err := http.NewRequest(
		"GET",
		"https://query.wikidata.org/sparql?query="+url.QueryEscape(query),
		nil,
	)
	// req.Header.Add("Content-Type", "application/sparql-query")
	req.Header.Add("Accept", "application/sparql-results+json")
	req.Header.Add("User-Agent", "book-alert/0.1.0 (github.com/jwoudenberg/book-alert); book-alert@jasperwoudenberg.com)")
	if err != nil {
		log.Println(err)
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	if resp.StatusCode != 200 {
		log.Println("Unexpected response:", resp.Status)
		return
	}
	xml, err := rawSparqlToRawFeed(body)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(string(xml))
}

func rawSparqlToRawFeed(bytes []byte) ([]byte, error) {
	sparql, err := sparqlFromJson(bytes)
	if err != nil {
		return nil, err
	}
	feed := sparqlToFeed(sparql)
	xml, err := renderXml(feed)
	if err != nil {
		return nil, err
	}
	return xml, nil
}

func sparqlFromJson(bytes []byte) (sparqlResults, error) {
	var sparql sparqlResults
	err := json.Unmarshal(bytes, &sparql)
	if err != nil {
		return sparql, err
	}
	return sparql, nil
}

type sparqlResults struct {
	Results struct {
		Bindings []struct {
			Work            sparqlValue
			AuthorLabel     sparqlValue
			Title           sparqlValue
			PublicationDate sparqlValue
		}
	}
}

type sparqlValue struct {
	Value string
}

type xmlFeed struct {
	XMLName xml.Name `xml:"feed"`
	Ns      string   `xml:"xmlns,attr"`
	Id      string   `xml:"id"`
	Title   string   `xml:"title"`
	Link    xmlLink  `xml:"link"`
	Updated string   `xml:"updated"`
	Entries []xmlEntry
}

type xmlEntry struct {
	XMLName xml.Name `xml:"entry"`
	Title   string   `xml:"title"`
	Link    xmlLink  `xml:"link"`
	Id      string   `xml:"id"`
	Updated string   `xml:"updated"`
	Author  string   `xml:"author>name"`
}

type xmlLink struct {
	XMLName xml.Name `xml:"link"`
	Rel     string   `xml:"rel,attr"`
	Href    string   `xml:"href,attr"`
}

func renderXml(feed xmlFeed) ([]byte, error) {
	feedContents, err := xml.MarshalIndent(feed, "", "  ")
	if err != nil {
		return nil, err
	}
	xmlFeed := append([]byte(xml.Header), feedContents...)
	return xmlFeed, nil
}

func sparqlToFeed(sparql sparqlResults) xmlFeed {
	books := sparql.Results.Bindings
	var updated string
	if len(books) > 0 {
		updated = books[0].PublicationDate.Value
	} else {
		updated = "1970-01-01T00:00:00Z"
	}
	feed := xmlFeed{
		Ns:      "http://www.w3.org/2005/Atom",
		Title:   "Book Alert",
		Id:      "https://github.com/jwoudenberg/book-alert",
		Link:    xmlLink{Rel: "self", Href: "https://github.com/jwoudenberg/book-alert"},
		Updated: updated,
	}
	for _, book := range books {
		entry := xmlEntry{
			Title:   book.Title.Value,
			Link:    xmlLink{Rel: "alternate", Href: book.Work.Value},
			Id:      book.Work.Value,
			Updated: book.PublicationDate.Value,
			Author:  book.AuthorLabel.Value,
		}
		feed.Entries = append(feed.Entries, entry)
	}
	return feed
}
