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
			"jpeg",
			"../../test/non-trans.jpeg",
			&Image{
				Width:  16,
				Height: 12,
				Format: "jpeg",
			},
		},
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
				Width:  16,
				Height: 12,
				Format: "png",
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
			"transparent avif",
			"../../test/trans.avif",
			&Image{
				Width:  16,
				Height: 16,
				Format: "webp",
			},
		},
		{
			"non transparent avif",
			"../../test/non-trans.avif",
			&Image{
				Width:  16,
				Height: 12,
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
			"jpeg",
			"../../test/non-trans.jpeg",
			false,
		},
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
			"transparent webp",
			"../../test/trans.webp",
			true,
		},
		{
			"non transparent webp",
			"../../test/non-trans.webp",
			false,
		},
		{
			"transparent avif",
			"../../test/trans.avif",
			true,
		},
		{
			"non transparent avif",
			"../../test/non-trans.avif",
			false,
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
