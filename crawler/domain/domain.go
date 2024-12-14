package domain

import (
	"crypto/sha256"
	"encoding/hex"
)

func Sha256(b []byte) (string, error) {
	hash := sha256.New()
	_, err := hash.Write(b)
	if err != nil {
		return "", err
	}
	s := hash.Sum(nil)
	return hex.EncodeToString(s), nil
}
