package image

import (
	"bytes"
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
	Size    int
	Width   int
	Height  int
	entropy *float64
	Format  string
	Data    []byte
	img     image.Image
}

func Load(b []byte) (*Image, error) {
	img, form, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	return &Image{
		Size:   len(b),
		Width:  img.Bounds().Dx(),
		Height: img.Bounds().Dy(),
		Format: form,
		Data:   b,
		img:    img,
	}, nil
}

func (img *Image) trans() bool {
	if img.Format == "jpeg" {
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

func (img *Image) Valid(width, height int, entropy float64, trans bool) bool {
	if img.img.Bounds().Dx() < width || img.img.Bounds().Dy() < height {
		return false
	}

	if !trans {
		if img.trans() {
			return false
		}
	}

	if img.Entropy() < entropy {
		return false
	}

	return true
}
