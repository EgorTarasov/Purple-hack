package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const url = "https://www.cbr.ru/news/eventandpress/?page=%d&IsEng=false&type=100&dateFrom=&dateTo=&Tid=&vol=&phrase="

type News struct {
	Id      int      `json:"id"`
	Content []string `json:"content"`
}

func getPage(client *http.Client, i int, wg *sync.WaitGroup) {
	defer wg.Done()
	req, err := http.NewRequest("GET", fmt.Sprintf(url, i), nil)
	if err != nil {
		log.Printf("request creat failed: %v\n", err)
		return
	}
	req.Header.Set("authority", "www.cbr.ru")
	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("accept-language", "en-US,en;q=0.9")
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("dnt", "1")
	req.Header.Set("pragma", "no-cache")
	req.Header.Set("sec-ch-ua", `"Not(A:Brand";v="24", "Chromium";v="122"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"macOS"`)
	req.Header.Set("sec-fetch-dest", "document")
	req.Header.Set("sec-fetch-mode", "navigate")
	req.Header.Set("sec-fetch-site", "none")
	req.Header.Set("sec-fetch-user", "?1")
	req.Header.Set("upgrade-insecure-requests", "1")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("parsing page error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	raw := make([]any, 0)
	if err = json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		log.Printf("unmarshal error: %v\n", err)
		return
	}
	news := News{
		Id:      i,
		Content: make([]string, 0),
	}
	for _, obj := range raw {
		data, ok := obj.(map[string]any)
		if !ok {
			log.Printf("not json obj: %v\n", err)
			continue
		}
		news.Content = append(news.Content, data["doc_htm"].(string))
	}

	for _, page := range news.Content {
		go func() {
			time.Sleep(time.Millisecond * 1500)
			req, err := http.NewRequest("GET", fmt.Sprintf("https://www.cbr.ru/press/event/?id=%s", page), nil)
			if err != nil {
				log.Printf("failed to send request to page: %v\n", err)
				return
			}
			req.Header.Set("authority", "www.cbr.ru")
			req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
			req.Header.Set("accept-language", "en-US,en;q=0.9")
			req.Header.Set("cache-control", "no-cache")
			req.Header.Set("cookie", "__ddg1_=a6FN2I7jVv1ffLNbkguc; _ym_uid=1705079485223578977; _ym_d=1709985139; _ym_isad=1; _ym_visorc=b; accept=1; ASPNET_SessionID=4xrryweghcad5pf3xk3gwjna")
			req.Header.Set("dnt", "1")
			req.Header.Set("pragma", "no-cache")
			req.Header.Set("sec-ch-ua", `"Not(A:Brand";v="24", "Chromium";v="122"`)
			req.Header.Set("sec-ch-ua-mobile", "?0")
			req.Header.Set("sec-ch-ua-platform", `"macOS"`)
			req.Header.Set("sec-fetch-dest", "document")
			req.Header.Set("sec-fetch-mode", "navigate")
			req.Header.Set("sec-fetch-site", "none")
			req.Header.Set("sec-fetch-user", "?1")
			req.Header.Set("upgrade-insecure-requests", "1")
			req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36")
			resp, err := client.Do(req)
			if err != nil {
				log.Printf("failed to parse page: %v\n", err)
				return
			}

			data, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Printf("failed to read page body: %v\n", err)
				return
			}
			resp.Body.Close()

			if err = os.WriteFile(fmt.Sprintf("news-parser/pages/%s.html", page), data, fs.ModePerm); err != nil {
				log.Printf("failed to store page file: %v\n", err)
				return
			}
		}()
	}
}

func main() {
	//wg := sync.WaitGroup{}
	//client := &http.Client{
	//	Transport: &http.Transport{
	//		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	//	},
	//}
	//for i := 1; i < 100; i++ {
	//	fmt.Printf("parsing page with id %d\n", i)
	//	if i%15 == 0 {
	//		time.Sleep(time.Second * 3)
	//	}
	//	wg.Add(1)
	//	go getPage(client, i, &wg)
	//}
	//wg.Wait()
	data := getData()

	ctx := context.Background()
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://LasuriaRobert:HilbertSpace@larek.tech:9500"))
	if err != nil {
		log.Fatal(err)
	}
	toMongo := make([]interface{}, 0)
	for _, document := range data {
		toMongo = append(toMongo, document)
	}
	if err = mongoClient.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	fmt.Println(len(data))
	//_, err = mongoClient.Database("cbr").Collection("news").InsertMany(ctx, toMongo)
	//if err != nil {
	//	log.Fatal(err)
	//}
}
