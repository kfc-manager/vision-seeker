package image

import (
	"bytes"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math"

	"github.com/kfc-manager/vision-seeker/domain"
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
)

type Image struct {
	Width  int
	Height int
	format string
	Raw    []byte
	img    image.Image
}

func LoadImage(b []byte) (*Image, error) {
	img, form, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	return &Image{
		Width:  img.Bounds().Dx(),
		Height: img.Bounds().Dy(),
		format: form,
		Raw:    b,
		img:    img,
	}, nil
}

func (img *Image) Hash() (string, error) {
	h, err := domain.Sha256(img.Raw)
	if err != nil {
		return "", err
	}

	return h, nil
}

func (img *Image) trans() bool {
	if img.format == "jpeg" {
		return false
	}

	for x := img.img.Bounds().Min.X; x < img.img.Bounds().Max.X; x++ {
		for y := img.img.Bounds().Min.Y; y < img.img.Bounds().Max.Y; y++ {
			_, _, _, alpha := img.img.At(x, y).RGBA()
			if alpha != 0 {
				return true
			}
		}
	}

	return false
}

func (img *Image) entropy() float64 {
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

	return entropy
}

func (img *Image) Valid(width, height int, entropy float64, trans bool) bool {
	if img.img.Bounds().Dx() < width || img.img.Bounds().Dy() < height {
		return false
	}

	if !trans {
		if img.trans() {
			return false
		}
	}

	if img.entropy() < entropy {
		return false
	}

	return true
}
