package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

import _ "embed"

//go:embed recent-publications.sparql
var query string

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
	var data sparqlResults
	err2 := json.Unmarshal(body, &data)
	if err2 != nil {
		log.Println(err2)
		return
	}

	fmt.Printf("%+v\n", data.Results.Bindings)
}
