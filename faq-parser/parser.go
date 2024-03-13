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
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/ilyakaznacheev/cleanenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const baseUrl = "https://www.cbr.ru"

var client = &http.Client{}

type QA struct {
	Question string `json:"question" bson:"question"`
	Answer   string `json:"answer" bson:"answer"`
	Url      string `json:"url" bson:"url"`
}

type QALinked struct {
	Topic    string `json:"topic" bson:"topic"`
	Question string `json:"question" bson:"question"`
	Answer   string `json:"answer" bson:"answer"`
	Url      string `json:"url" bson:"url"`
}

func newRequest(url string) *http.Request {
	req, err := http.NewRequest("GET", baseUrl+url, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("authority", "www.cbr.ru")
	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("accept-language", "en-US,en;q=0.9")
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("cookie", "__ddg1_=mXQm8yz7IKrWTbzxzSrU; _ym_uid=1705079485223578977; _ym_d=1709989337; _ym_isad=1; _ym_visorc=b; accept=1")
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
	return req
}

func getRubricTitles(faqUrl string) []string {
	resp, err := client.Do(newRequest(faqUrl))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	if err = os.WriteFile("faq.html", raw, fs.ModePerm); err != nil {
		panic(err)
	}

	f, err := os.Open("faq.html")
	if err != nil {
		panic(err)
	}

	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		panic(err)
	}

	rubricTitles := make([]string, 0)
	doc.Find(".rubric_title").Each(func(i int, selection *goquery.Selection) {
		href, ok := selection.Attr("href")
		if ok {
			rubricTitles = append(rubricTitles, href)
		}
	})
	return rubricTitles
}

func getQA(ctx context.Context, url string, wg *sync.WaitGroup, db *mongo.Database) {
	defer wg.Done()
	time.Sleep(time.Millisecond * 1500)

	nameSplitted := strings.Split(url, "/")
	name := nameSplitted[2]

	resp, err := client.Do(newRequest(url))
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}

	if err = os.WriteFile(fmt.Sprintf("./faq-pages/%s.html", name), raw, fs.ModePerm); err != nil {
		log.Println(err)
		return
	}

	f, err := os.Open(fmt.Sprintf("./faq-pages/%s.html", name))
	if err != nil {
		log.Println(err)
		return
	}

	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		log.Println(err)
		return
	}

	qa := make([]QA, 0)
	doc.Find(".dropdown_content").Each(func(i int, content *goquery.Selection) {
		content.Find(".dropdown .question").Each(func(j int, block *goquery.Selection) {
			question := block.Find(".question_title").Text()
			question = strings.TrimSpace(strings.ReplaceAll(question, " ", " "))

			answer := block.Find(".additional-text-block").Text()
			answer = strings.TrimSpace(strings.ReplaceAll(answer, " ", " "))

			cur := QA{
				Question: question,
				Answer:   answer,
				Url:      baseUrl + url,
			}
			qa = append(qa, cur)
		})
	})

	doc.Find(".dropdown_title-link").Each(func(i int, block *goquery.Selection) {
		block.Children().Each(func(i int, link *goquery.Selection) {
			link.Children().Each(func(i int, selection *goquery.Selection) {
				href, ok := selection.Attr("href")
				if ok {
					wg.Add(1)
					go getLinkedQA(ctx, href, wg, db)
				}
			})
		})
	})

	toMongo := make([]interface{}, len(qa))
	for idx, obj := range qa {
		toMongo[idx] = obj
	}
	if len(toMongo) > 0 {
		_, err = db.Collection("faq").InsertMany(ctx, toMongo)
		if err != nil {
			log.Printf("failed to insert %s into monog: %v", name, err)
			return
		}
	}

	data, err := json.Marshal(qa)
	if err != nil {
		log.Println(err)
		return
	}

	if err = os.WriteFile(fmt.Sprintf("./qa/%s.json", name), data, fs.ModePerm); err != nil {
		log.Println(err)
		return
	}
}

func getLinkedQA(ctx context.Context, url string, wg *sync.WaitGroup, db *mongo.Database) {
	defer wg.Done()
	time.Sleep(time.Millisecond * 2000)

	nameSplitted := strings.Split(url, "/")
	name := ""
	if len(nameSplitted) > 3 {
		name = nameSplitted[3]
	} else {
		name = nameSplitted[2]
	}

	if name == "" {
		name = nameSplitted[len(nameSplitted)-2]
	}

	resp, err := client.Do(newRequest(url))
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}

	if err = os.WriteFile(fmt.Sprintf("./faq-pages-linked/%s.html", name), raw, fs.ModePerm); err != nil {
		log.Println(err)
		return
	}

	f, err := os.Open(fmt.Sprintf("./faq-pages-linked/%s.html", name))
	if err != nil {
		log.Println(err)
		return
	}

	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		log.Println(err)
		return
	}

	qaLinked := make([]QALinked, 0)
	doc.Find(".title-container").Each(func(i int, topicS *goquery.Selection) {
		topic := strings.ReplaceAll(strings.TrimSpace(strings.ReplaceAll(topicS.Find(".block-title").Text(), " ", " ")), "\n", "")
		topicS.Find(".dropdown").Each(func(j int, blockS *goquery.Selection) {
			question := strings.TrimSpace(strings.ReplaceAll(blockS.Find(".question_title").First().Text(), " ", " "))
			answer := ""
			blockS.Find(".additional-text-block").First().Each(func(k int, answerS *goquery.Selection) {
				answer += answerS.Text()
			})
			answer = strings.TrimSpace(strings.ReplaceAll(answer, " ", " "))
			qaLinked = append(qaLinked, QALinked{
				Topic:    topic,
				Question: question,
				Answer:   answer,
				Url:      baseUrl + url,
			})
		})
	})

	toMongo := make([]interface{}, len(qaLinked))
	for idx, obj := range qaLinked {
		toMongo[idx] = obj
	}
	if len(toMongo) > 0 {
		_, err = db.Collection("faq-links").InsertMany(ctx, toMongo)
		if err != nil {
			log.Printf("failed to insert %s into monog: %v", name, err)
			return
		}
	}

	data, err := json.Marshal(qaLinked)
	if err != nil {
		log.Println(err)
		return
	}

	if err = os.WriteFile(fmt.Sprintf("./qa-linked/%s.json", name), data, fs.ModePerm); err != nil {
		log.Println(err)
		return
	}
}

type config struct {
	MongoConnStr string `env:"MONGO_CONN_STR"`
}

func main() {
	var wg sync.WaitGroup

	cfg := &config{}
	if err := cleanenv.ReadConfig("../.env", cfg); err != nil {
		panic(err)
	}

	ctx := context.Background()
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoConnStr))
	if err != nil {
		panic(err)
	}
	if err = mongoClient.Ping(ctx, nil); err != nil {
		panic(err)
	}
	db := mongoClient.Database("cbr")

	for i, title := range getRubricTitles("/faq") {
		wg.Add(1)
		if i%5 == 0 {
			time.Sleep(time.Second)
		}
		go getQA(ctx, title, &wg, db)
	}
	wg.Wait()
}
