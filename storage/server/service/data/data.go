package data

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/kfc-manager/vision-seeker/storage/server/adapter/bucket"
	"github.com/kfc-manager/vision-seeker/storage/server/adapter/cache"
	"github.com/kfc-manager/vision-seeker/storage/server/adapter/database"
	"github.com/kfc-manager/vision-seeker/storage/server/adapter/queue"
	"github.com/kfc-manager/vision-seeker/storage/server/domain"
)

type Service interface {
	StoreImg(hash string, b []byte) error
	StoreDataset(name string) (string, error)
	StoreMetadata(id, hash, label string) error
}

type service struct {
	db     database.Database
	bucket bucket.Bucket
	cache  cache.Cache
	queue  queue.Queue
}

func New(
	db database.Database,
	b bucket.Bucket,
	c cache.Cache,
	q queue.Queue,
) *service {
	return &service{db: db, bucket: b, cache: c, queue: q}
}

func (s *service) StoreImg(hash string, b []byte) error {
	img, err := domain.LoadImgHash(b)
	if err != nil {
		log.Printf("hash loading error: %s", err.Error())
		return err
	}

	if img.Hash != hash {
		return fmt.Errorf("checksum mismatch")
	}

	exist, err := s.cache.Exist(img.Hash)
	if err != nil {
		log.Printf("cache error: %s", err.Error())
	}
	if exist {
		return fmt.Errorf("image with hash '%s' already exists", img.Hash)
	}

	err = s.bucket.Put(img.Hash, b)
	if err != nil {
		err = fmt.Errorf("bucket error: %s", err.Error())
		log.Printf(err.Error())
		return err
	}
	ok, err := s.db.InsertImgHash(img)
	if err != nil {
		err = fmt.Errorf("database error: %s", err.Error())
		log.Printf(err.Error())
		return err
	}
	if !ok {
		return fmt.Errorf("image with hash '%s' already exists", img.Hash)
	}
	err = s.queue.Push([]byte(img.Hash))
	if err != nil {
		err = fmt.Errorf("queue error: %s", err.Error())
		log.Printf(err.Error())
		return err
	}

	return nil
}

func (s *service) StoreDataset(name string) (string, error) {
	ds, err := domain.CreateDataset(name)
	if err != nil {
		return "", err
	}

	ok, err := s.db.InsertDataset(ds)
	if err != nil {
		err = fmt.Errorf("database error: %s", err.Error())
		log.Printf(err.Error())
		return "", err
	}
	if !ok {
		return "", fmt.Errorf("dataset with name '%s' already exists", name)
	}

	return ds.Id.String(), nil
}

func (s *service) StoreMetadata(id, hash, label string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid uuid '%s'", id)
	}
	mapping := &domain.Mapping{
		Id:        uid,
		Hash:      hash,
		CreatedAt: time.Now(),
	}
	ok, err := s.db.InsertImgMap(mapping)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf(
			"mapping from dataset '%s' to image '%s' already exists",
			id, hash,
		)
	}

	if len(label) < 1 {
		return nil
	}

	l, err := mapping.Label(label)
	if err != nil {
		log.Printf("hash loading error: %s", err.Error())
		return err
	}

	ok, err = s.db.InsertLabel(l)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf(
			"label with hash '%s' already exists",
			l.Hash,
		)
	}

	ok, err = s.db.InsertLabelMap(l)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf(
			"label mapping with hash '%s' to dataset '%s' and image '%s' already exists",
			l.Hash, l.Mapping.Id.String(), l.Mapping.Hash,
		)
	}

	return nil
}
