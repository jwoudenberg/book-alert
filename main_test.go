package main

import (
	"testing"
)

func TestPlaceholder(t *testing.T) {
	inputJson := `{
  "head" : {
    "vars" : [ "work", "authorLabel", "title", "publicationDate" ]
  },
  "results" : {
    "bindings" : [ {
      "work" : {
        "type" : "uri",
        "value" : "http://www.wikidata.org/entity/Q107123470"
      },
      "authorLabel" : {
        "xml:lang" : "en",
        "type" : "literal",
        "value" : "Martha Wells"
      },
      "title" : {
        "xml:lang" : "en",
        "type" : "literal",
        "value" : "Fugitive Telemetry"
      },
      "publicationDate" : {
        "datatype" : "http://www.w3.org/2001/XMLSchema#dateTime",
        "type" : "literal",
        "value" : "2021-04-27T00:00:00Z"
      }
    }, {
      "work" : {
        "type" : "uri",
        "value" : "http://www.wikidata.org/entity/Q100540380"
      },
      "authorLabel" : {
        "xml:lang" : "en",
        "type" : "literal",
        "value" : "Martha Wells"
      },
      "title" : {
        "xml:lang" : "en",
        "type" : "literal",
        "value" : "Network Effect"
      },
      "publicationDate" : {
        "datatype" : "http://www.w3.org/2001/XMLSchema#dateTime",
        "type" : "literal",
        "value" : "2020-05-05T00:00:00Z"
      }
    } ]
  }
}`
	expectedXml := `<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">
  <id>https://github.com/jwoudenberg/book-alert</id>
  <title>Book Alert</title>
  <link rel="self" href="https://github.com/jwoudenberg/book-alert"></link>
  <updated>2021-04-27T00:00:00Z</updated>
  <entry>
    <title>Fugitive Telemetry</title>
    <link rel="alternate" href="http://www.wikidata.org/entity/Q107123470"></link>
    <id>http://www.wikidata.org/entity/Q107123470</id>
    <updated>2021-04-27T00:00:00Z</updated>
    <author>
      <name>Martha Wells</name>
    </author>
  </entry>
  <entry>
    <title>Network Effect</title>
    <link rel="alternate" href="http://www.wikidata.org/entity/Q100540380"></link>
    <id>http://www.wikidata.org/entity/Q100540380</id>
    <updated>2020-05-05T00:00:00Z</updated>
    <author>
      <name>Martha Wells</name>
    </author>
  </entry>
</feed>`
	actualXml, err := rawSparqlToRawFeed([]byte(inputJson))
	if err != nil {
		t.Errorf("Conversion failed: %s", err)
	}
	if string(actualXml) != expectedXml {
		t.Errorf("Actual XML output did not match expected.\n%s", actualXml)
	}

}
