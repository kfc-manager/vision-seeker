package data

import (
	"encoding/json"
	gourl "net/url"

	"github.com/kfc-manager/vision-seeker/crawler/adapter/bucket"
	"github.com/kfc-manager/vision-seeker/crawler/adapter/cache"
	"github.com/kfc-manager/vision-seeker/crawler/adapter/database"
	"github.com/kfc-manager/vision-seeker/crawler/adapter/queue"
	"github.com/kfc-manager/vision-seeker/crawler/domain"
	"github.com/kfc-manager/vision-seeker/crawler/domain/image"
)

type Service interface {
	StoreImage(img *image.Image, label string) error
	Visit(url *gourl.URL, alt string) error
	Next() (*gourl.URL, string, error)
}

type service struct {
	db       database.Database
	cache    cache.Cache
	bucket   bucket.Bucket
	urlQueue queue.Queue
	imgQueue queue.Queue
}

func New(
	db database.Database,
	c cache.Cache,
	b bucket.Bucket,
	url queue.Queue,
	img queue.Queue,
) *service {
	return &service{
		db:       db,
		cache:    c,
		bucket:   b,
		urlQueue: url,
		imgQueue: img,
	}
}

func (s *service) StoreImage(img *image.Image, label string) error {
	imgHash, err := domain.Sha256(img.Data)
	if err != nil {
		return err
	}
	ok, err := s.db.InsertImage(imgHash, img)
	if err != nil {
		return err
	}
	if ok {
		err := s.bucket.Put(imgHash, img.Data)
		if err != nil {
			return err
		}
		s.imgQueue.Push([]byte(imgHash))
	}

	lblHash, err := domain.Sha256([]byte(label))
	if err != nil {
		return err
	}
	_, err = s.db.InsertLabel(lblHash, label)
	if err != nil {
		return err
	}

	_, err = s.db.InsertMapping(imgHash, lblHash)
	if err != nil {
		return err
	}

	return nil
}

type message struct {
	Url string `json:"url"`
	Alt string `json:"alt"`
}

func (s *service) Visit(url *gourl.URL, alt string) error {
	hash, err := domain.Sha256([]byte(url.Host + url.Path + url.RawQuery))
	if err != nil {
		return err
	}
	exist, err := s.cache.Exist(hash)
	if err != nil {
		return err
	}
	if exist {
		return nil
	}

	ok, err := s.db.InsertUrl(hash)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}

	b, err := json.Marshal(&message{Url: url.String(), Alt: alt})
	if err != nil {
		return err
	}

	return s.urlQueue.Push(b)
}

func (s *service) Next() (*gourl.URL, string, error) {
	b, err := s.urlQueue.Pull()
	if err != nil {
		return nil, "", err
	}

	msg := &message{}
	if err := json.Unmarshal(b, msg); err != nil {
		return nil, "", err
	}

	url, err := gourl.Parse(msg.Url)
	if err != nil {
		return nil, "", err
	}

	return url, msg.Alt, nil
}
