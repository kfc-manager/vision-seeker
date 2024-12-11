package data

import (
	gourl "net/url"

	"github.com/kfc-manager/vision-seeker/adapter/bucket"
	"github.com/kfc-manager/vision-seeker/adapter/cache"
	"github.com/kfc-manager/vision-seeker/adapter/database"
	"github.com/kfc-manager/vision-seeker/adapter/queue"
	"github.com/kfc-manager/vision-seeker/domain"
	"github.com/kfc-manager/vision-seeker/domain/image"
)

type Service interface {
	StoreImage(img *image.Image) error
	Url() (string, error)
	SetUrl(url *gourl.URL) error
}

type service struct {
	db     database.Database
	cache  cache.Cache
	bucket bucket.Bucket
	queue  queue.Queue
}

func New(
	db database.Database,
	c cache.Cache,
	b bucket.Bucket,
	q queue.Queue,
) *service {
	return &service{db: db, cache: c, bucket: b, queue: q}
}

func (s *service) StoreImage(img *image.Image) error {
	h, err := img.Hash()
	if err != nil {
		// TODO log
		return err
	}

	err = s.bucket.Put(h, img.Raw)
	if err != nil {
		// TODO log
		return err
	}

	return nil
}

func (s *service) Url() (string, error) {
	return s.queue.Pull()
}

func (s *service) SetUrl(url *gourl.URL) error {
	hash, err := domain.Sha256([]byte(url.Host + url.Path + url.RawQuery))
	if err != nil {
		// TODO log
		return err
	}

	exist, err := s.cache.Exists(hash)
	if err != nil {
		// TODO log
	}
	if exist {
		return nil
	}

	ok, err := s.db.InsertUrl(hash)
	if err != nil {
		// TODO log
		return err
	}
	// url already exist in database table
	if !ok {
		return nil
	}

	err = s.cache.Set(hash)
	if err != nil {
		// TODO log
	}

	err = s.queue.Push(url.String())
	if err != nil {
		return err
	}

	return nil
}
