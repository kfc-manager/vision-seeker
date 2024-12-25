package storage

import (
	"log"

	"github.com/kfc-manager/vision-seeker/server/adapter/bucket"
	"github.com/kfc-manager/vision-seeker/server/adapter/cache"
	"github.com/kfc-manager/vision-seeker/server/adapter/database"
	"github.com/kfc-manager/vision-seeker/server/adapter/queue"
	"github.com/kfc-manager/vision-seeker/server/domain/image"
	"github.com/kfc-manager/vision-seeker/server/service"
)

type Service interface {
	Seen(hash string) error
	StoreImage(img *image.Image) error
}

type storage struct {
	db     database.Database
	bucket bucket.Bucket
	cache  cache.Cache
	clip   queue.Queue
	dino   queue.Queue
}

func New(
	db database.Database,
	b bucket.Bucket,
	c cache.Cache,
	clip queue.Queue,
	dino queue.Queue,
) *storage {
	return &storage{
		db:     db,
		bucket: b,
		cache:  c,
		clip:   clip,
		dino:   dino,
	}
}

func (s *storage) Seen(hash string) error {
	seen, err := s.cache.Exist(hash)
	if err != nil {
		log.Println(err.Error())
		return &service.Error{
			Msg:    "Internal Server",
			Status: 500,
		}
	}

	if seen {
		return &service.Error{
			Msg:    "image already processed",
			Status: 409,
		}
	}

	err = s.cache.Set(hash)
	if err != nil {
		log.Println(err.Error())
		return &service.Error{
			Msg:    "Internal Server",
			Status: 500,
		}
	}

	return nil
}

func (s *storage) StoreImage(img *image.Image) error {
	err := s.bucket.Put(img.Hash, img.Data)
	if err != nil {
		log.Println(err.Error())
		return &service.Error{
			Msg:    "Internal Server",
			Status: 500,
		}
	}

	ok, err := s.db.InsertImage(img)
	if err != nil {
		log.Println(err.Error())
		return &service.Error{
			Msg:    "Internal Server",
			Status: 500,
		}
	}
	if !ok {
		return &service.Error{
			Msg:    "image already processed",
			Status: 409,
		}
	}

	err = s.clip.Push([]byte(img.Hash))
	if err != nil {
		log.Println(err.Error())
	}
	err = s.dino.Push([]byte(img.Hash))
	if err != nil {
		log.Println(err.Error())
	}

	return nil
}
