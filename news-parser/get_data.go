package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Document struct {
	Url   string `json:"url" bson:"url"`
	Title string `json:"title" bson:"title"`
	Date  string `json:"date" bson:"date"`
	Body  string `json:"body" bson:"body"`
}

func getData() []Document {
	//args := []string{
	//	"-nopgbrk",   // Don't insert page breaks (form feed characters) between pages.
	//	"news-parser/%s.pdf", // The input file.
	//	"-",          // Send the output to stdout.
	//}
	//cmd := exec.CommandContext(context.Background(), "pdftotext", args...)
	//
	//var buf bytes.Buffer
	//cmd.Stdout = &buf
	//
	//if err := cmd.Run(); err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//
	//fmt.Println(buf.String())

	pages, err := os.ReadDir("pages")
	if err != nil {
		log.Fatal(err)
	}

	storage := make([]Document, 0)
	set := map[string]struct{}{}

	for _, page := range pages {
		f, err := os.Open("pages/" + page.Name())
		if err != nil {
			log.Fatal(err)
		}
		doc, err := goquery.NewDocumentFromReader(f)
		if err != nil {
			log.Fatal(err)
		}

		document := Document{
			Url:   fmt.Sprintf("https://cbr.ru/press/event/?id=%s", page.Name()),
			Date:  doc.Find(".news-info-line_date").First().Text(),
			Title: doc.Find("title").Text(),
			Body:  strings.ReplaceAll(strings.Split(strings.ReplaceAll(doc.Find(".landing-text").AfterSelection(doc.Find(".landing-text").First()).BeforeSelection(doc.Find(".landing-text").Last()).Text(), " ", " "), "\nФото на превью:")[0], "\n", " "),
		}
		if _, ok := set[document.Body]; ok {
			continue
		}
		set[document.Body] = struct{}{}
		if document.Title != "500 | Банк России" {
			storage = append(storage, document)
		}
	}
	data, _ := json.Marshal(storage)
	if err = os.WriteFile("data.json", data, fs.ModePerm); err != nil {
		log.Fatal(err)
	}
	return storage
}
