package bucket

import (
	"errors"
	"os"
	"path/filepath"
)

type Bucket interface {
	Get(key string) ([]byte, error)
	Put(key string, body []byte) error
	Delete(key string) error
}

type bucket struct {
	path string
}

func New(path string) (*bucket, error) {
	if len(path) < 1 {
		return nil, errors.New("path can't be empty")
	}
	if path[len(path)-1:] == "/" {
		path = path[:len(path)-1]
	}
	return &bucket{path: path}, nil
}

func (b *bucket) Get(key string) ([]byte, error) {
	body, err := os.ReadFile(b.path + "/" + key)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (b *bucket) Put(key string, body []byte) error {
	dir := filepath.Dir(key)
	if err := os.MkdirAll(b.path+"/"+dir, 0777); err != nil {
		return err
	}
	return os.WriteFile(b.path+"/"+key, body, 0777)
}

func (b *bucket) Delete(key string) error {
	return os.RemoveAll(b.path + "/" + key)
}
