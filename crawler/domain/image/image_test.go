package image

import (
	"os"
	"testing"
)

func loadTestData(key string) (*Image, error) {
	b, err := os.ReadFile(key)
	if err != nil {
		return nil, err
	}
	img, err := Load(b)
	if err != nil {
		return nil, err
	}
	return img, nil
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
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			img, err := loadTestData(test.input)
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
