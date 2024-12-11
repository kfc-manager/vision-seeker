package main

import (
	"fmt"
	"os"
	"time"

	"github.com/kfc-manager/vision-seeker/adapter/bucket"
	"github.com/kfc-manager/vision-seeker/adapter/cache"
	"github.com/kfc-manager/vision-seeker/adapter/client"
	"github.com/kfc-manager/vision-seeker/adapter/database"
	"github.com/kfc-manager/vision-seeker/adapter/queue"
	"github.com/kfc-manager/vision-seeker/service/crawler"
	"github.com/kfc-manager/vision-seeker/service/data"
)

func main() {
	b, err := bucket.New(envOrPanic("BUCKET_PATH"))
	if err != nil {
		panic(err)
	}

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

	c, err := cache.New(
		envOrPanic("CACHE_HOST"),
		envOrPanic("CACHE_PORT"),
		envOrPanic("CACHE_PASS"),
	)
	if err != nil {
		panic(err)
	}

	q, err := queue.New(
		envOrPanic("QUEUE_HOST"),
		envOrPanic("QUEUE_PORT"),
		envOrPanic("QUEUE_NAME"),
	)
	if err != nil {
		panic(err)
	}

	err = q.Push(envOrPanic("START"))
	if err != nil {
		panic(err)
	}
	time.Sleep(5 * time.Second)
	crawler.New(client.New(), data.New(db, c, b, q)).Crawl()
}

func envOrPanic(key string) string {
	val := os.Getenv(key)
	if len(val) < 1 {
		panic(fmt.Errorf("missing env variable: %s", key))
	}
	return val
}
