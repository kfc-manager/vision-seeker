package main

import (
	"fmt"
	gourl "net/url"
	"os"
	"time"

	"github.com/kfc-manager/vision-seeker/crawler/adapter/bucket"
	"github.com/kfc-manager/vision-seeker/crawler/adapter/cache"
	"github.com/kfc-manager/vision-seeker/crawler/adapter/client"
	"github.com/kfc-manager/vision-seeker/crawler/adapter/database"
	"github.com/kfc-manager/vision-seeker/crawler/adapter/queue"
	"github.com/kfc-manager/vision-seeker/crawler/service/crawler"
	"github.com/kfc-manager/vision-seeker/crawler/service/data"
)

func main() {
	db, err := database.New(
		envOrPanic("DB_HOST"),
		envOrPanic("DB_PORT"),
		envOrPanic("DB_NAME"),
		envOrPanic("DB_USER"),
		envOrPanic("DB_PASS"),
	)
	if err != nil {
		panic(err)
	}

	cach, err := cache.New(
		envOrPanic("CACHE_HOST"),
		envOrPanic("CACHE_PORT"),
		envOrPanic("CACHE_PASS"),
	)
	if err != nil {
		panic(err)
	}

	buck, err := bucket.New(
		envOrPanic("BUCKET_PATH"),
	)
	if err != nil {
		panic(err)
	}

	que, err := queue.New(
		envOrPanic("QUEUE_HOST"),
		envOrPanic("QUEUE_PORT"),
		envOrPanic("QUEUE_NAME"),
		0,
	)
	if err != nil {
		panic(err)
	}

	time.Sleep(5 * time.Second)
	dataServ := data.New(db, cach, buck, que)
	url, err := gourl.Parse(envOrPanic("START"))
	if err != nil {
		panic(err)
	}
	err = dataServ.Visit(url, "")
	if err != nil {
		panic(err)
	}

	crawler.New(client.New(), dataServ).Crawl()
}

func envOrPanic(key string) string {
	val := os.Getenv(key)
	if len(val) < 1 {
		panic(fmt.Errorf("env variable '%s' missing", key))
	}
	return val
}
