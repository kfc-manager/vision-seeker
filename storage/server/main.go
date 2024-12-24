package main

import (
	"fmt"
	"os"

	"github.com/kfc-manager/vision-seeker/storage/server/adapter/bucket"
	"github.com/kfc-manager/vision-seeker/storage/server/adapter/cache"
	"github.com/kfc-manager/vision-seeker/storage/server/adapter/database"
	"github.com/kfc-manager/vision-seeker/storage/server/adapter/queue"
	"github.com/kfc-manager/vision-seeker/storage/server/adapter/server"
	"github.com/kfc-manager/vision-seeker/storage/server/domain"
	"github.com/kfc-manager/vision-seeker/storage/server/service/data"
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
	buck, err := bucket.New(
		envOrPanic("BUCKET_PATH"),
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
	que, err := queue.New(
		envOrPanic("QUEUE_HOST"),
		envOrPanic("QUEUE_PORT"),
		envOrPanic("QUEUE_NAME"),
	)
	if err != nil {
		panic(err)
	}

	token, err := domain.RandomStr()
	if err != nil {
		panic(err)
	}
	fmt.Println(token)
	s := server.New(80, token, data.New(db, buck, cach, que))
	if err := s.Listen(); err != nil {
		panic(err)
	}
}

func envOrPanic(key string) string {
	val := os.Getenv(key)
	if len(val) < 1 {
		panic(fmt.Errorf("missing env variable: %s", key))
	}
	return val
}
