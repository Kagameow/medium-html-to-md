package main

import (
	"fmt"
	md "github.com/JohannesKaufmann/html-to-markdown"
	"log"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	f, err := os.Open("example.html")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		log.Fatal(err)
	}

	article := doc.Find(".ch")

	properLinks := md.Rule{
		Filter: []string{"a"},
		Replacement: func(content string, selec *goquery.Selection, opt *md.Options) *string {
			// This converter can't handle footnotes from medium, so we skip them
			href, exists := selec.Attr("href")
			if exists && !strings.Contains(href, "data:text/html;") {
				return nil
			}

			return md.String(content)
		},
	}

	converter := md.NewConverter("", true, nil)
	converter.AddRules(properLinks)
	markdown := converter.Convert(article)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create("output.md")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = file.WriteString(markdown)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Markdown content successfully saved to 'output.md'")
}
