package image

import (
	"bytes"
	gosha256 "crypto/sha256"
	"encoding/hex"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
)

type Image struct {
	Hash    string
	Size    int
	Width   int
	Height  int
	entropy *float64
	trans   *bool
	Format  string
	Data    []byte
	Label   string
	img     image.Image
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

func Load(b []byte) (*Image, error) {
	img, form, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	h, err := sha256(b)
	if err != nil {
		return nil, err
	}

	return &Image{
		Hash:   h,
		Size:   len(b),
		Width:  img.Bounds().Dx(),
		Height: img.Bounds().Dy(),
		Format: form,
		Data:   b,
		img:    img,
	}, nil
}

func (img *Image) Trans() bool {
	if img.trans == nil {
		// relies on the fact that jpeg don't support transparency
		if img.Format == "jpeg" {
			*img.trans = false
		}

		for x := img.img.Bounds().Min.X; x < img.img.Bounds().Max.X; x++ {
			for y := img.img.Bounds().Min.Y; y < img.img.Bounds().Max.Y; y++ {
				_, _, _, alpha := img.img.At(x, y).RGBA()
				if alpha < 0xffff {
					*img.trans = true
				}
			}
		}

		*img.trans = false
	}

	return *img.trans
}

func (img *Image) Entropy() float64 {
	if img.entropy != nil {
		return *img.entropy
	}

	hist := make([]int, 256)
	total := 0

	for y := img.img.Bounds().Min.Y; y < img.img.Bounds().Max.Y; y++ {
		for x := img.img.Bounds().Min.X; x < img.img.Bounds().Max.X; x++ {
			gray := color.GrayModel.Convert(img.img.At(x, y)).(color.Gray)
			hist[gray.Y]++
			total++
		}
	}

	entropy := 0.0
	for _, count := range hist {
		if count > 0 {
			prob := float64(count) / float64(total)
			entropy -= prob * math.Log2(prob)
		}
	}

	img.entropy = &entropy
	return entropy
}
