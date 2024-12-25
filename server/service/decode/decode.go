package decode

import (
	"log"

	"github.com/kfc-manager/vision-seeker/server/domain/image"
	"github.com/kfc-manager/vision-seeker/server/service"
)

type Service interface {
	DecodeImage(hash, label string, b []byte) (*image.Image, error)
}

type decode struct{}

func New() *decode {
	return &decode{}
}

func (s *decode) DecodeImage(hash, label string, b []byte) (*image.Image, error) {
	img, err := image.Load(b)
	if err != nil {
		log.Println(err.Error())
		return nil, &service.Error{Msg: "Internal Server", Status: 500}
	}

	if img.Hash != hash {
		return nil, &service.Error{Msg: "checksum mismatch", Status: 400}
	}

	img.Label = label
	return img, nil
}
