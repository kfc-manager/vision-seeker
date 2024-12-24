package domain

import (
	"crypto/rand"
	gosha256 "crypto/sha256"
	"encoding/hex"
	"errors"
	"math/big"
	"time"
	"unicode"

	"github.com/google/uuid"
)

type Image struct {
	Hash      string
	CreatedAt time.Time
}

type Dataset struct {
	Id        uuid.UUID
	Name      string
	CreatedAt time.Time
}

type Mapping struct {
	Id        uuid.UUID
	Hash      string
	CreatedAt time.Time
}

type Label struct {
	Mapping *Mapping
	Hash    string
	Label   string
}

func sha256(b []byte) (string, error) {
	hash := gosha256.New()
	_, err := hash.Write(b)
	if err != nil {
		return "", err
	}
	s := hash.Sum(nil)
	return hex.EncodeToString(s), nil
}

func LoadImgHash(b []byte) (*Image, error) {
	h, err := sha256(b)
	if err != nil {
		return nil, err
	}
	return &Image{Hash: h, CreatedAt: time.Now()}, nil
}

func valid(name string) error {
	if len(name) < 3 {
		return errors.New("dataset name must be at least 3 characters long")
	}
	if len(name) > 32 {
		return errors.New("dataset name can't be longer than 32 characters")
	}
	for _, char := range name {
		if unicode.IsLetter(char) {
			continue
		}
		if unicode.IsDigit(char) {
			continue
		}
		if string(char) == "_" || string(char) == "-" {
			continue
		}
		return errors.New(
			"dataset name can only contain letters, numbers, hyphens or underscores",
		)
	}
	return nil
}

func CreateDataset(name string) (*Dataset, error) {
	if err := valid(name); err != nil {
		return nil, err
	}
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	return &Dataset{
		Id:        id,
		Name:      name,
		CreatedAt: time.Now(),
	}, nil
}

func (m *Mapping) Label(label string) (*Label, error) {
	h, err := sha256([]byte(label))
	if err != nil {
		return nil, err
	}
	return &Label{Mapping: m, Hash: h, Label: label}, nil
}

func RandomStr() (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, 16)
	for i := range result {
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[index.Int64()]
	}
	return string(result), nil
}
