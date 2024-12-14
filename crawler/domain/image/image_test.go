package image

import (
	"os"
	"testing"
)

func loadTestData(key string) ([]byte, error) {
	b, err := os.ReadFile(key)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func TestLoad(t *testing.T) {
	var tests = []struct {
		name  string
		input string
		want  *Image
	}{
		{
			"transparent png",
			"../../test/trans.png",
			&Image{
				Width:  16,
				Height: 16,
				Format: "png",
			},
		},
		{
			"non transparent png",
			"../../test/non-trans.png",
			&Image{
				Width:  137,
				Height: 91,
				Format: "png",
			},
		},
		{
			"jpeg",
			"../../test/jpeg.jpeg",
			&Image{
				Width:  241,
				Height: 180,
				Format: "jpeg",
			},
		},
		{
			"non transparent webp",
			"../../test/non-trans.webp",
			&Image{
				Width:  16,
				Height: 12,
				Format: "webp",
			},
		},
		{
			"transparent webp",
			"../../test/trans.webp",
			&Image{
				Width:  16,
				Height: 16,
				Format: "webp",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			b, err := loadTestData(test.input)
			if err != nil {
				t.Error(err.Error())
				return
			}

			got, err := Load(b)
			if err != nil {
				t.Error(err.Error())
				return
			}

			if got.Width != test.want.Width {
				t.Errorf("got width: %d, want: %d", got.Width, test.want.Width)
			}
			if got.Height != test.want.Height {
				t.Errorf("got height: %d, want: %d", got.Height, test.want.Height)
			}
			if got.Format != test.want.Format {
				t.Errorf("got format: %s, want: %s", got.Format, test.want.Format)
			}
		})
	}
}

func TestTrans(t *testing.T) {
	var tests = []struct {
		name  string
		input string
		want  bool
	}{
		{
			"transparent png",
			"../../test/trans.png",
			true,
		},
		{
			"non transparent png",
			"../../test/non-trans.png",
			false,
		},
		{
			"jpeg",
			"../../test/jpeg.jpeg",
			false,
		},
		{
			"non transparent webp",
			"../../test/non-trans.webp",
			false,
		},
		{
			"transparent webp",
			"../../test/trans.webp",
			true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			b, err := loadTestData(test.input)
			if err != nil {
				t.Error(err.Error())
				return
			}

			img, err := Load(b)
			if err != nil {
				t.Error(err.Error())
				return
			}

			got := img.trans()
			if got != test.want {
				t.Errorf("got: %t, want: %t", got, test.want)
			}
		})
	}
}
