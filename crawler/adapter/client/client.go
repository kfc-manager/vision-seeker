package client

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type ResType string

const (
	Image  ResType = "image"
	Html   ResType = "html"
	Unkown ResType = "unkown"
)

type response struct {
	Type ResType
	Body []byte
}

type Client interface {
	Get(url string) (*response, error)
}

type client struct {
	client *http.Client
}

func New() *client {
	return &client{client: &http.Client{
		Timeout: 10 * time.Second,
	}}
}

func (c *client) Get(url string) (*response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode > 299 || res.StatusCode < 200 {
		return nil, fmt.Errorf("response status: '%d', with body: '%s'", res.StatusCode, string(b))
	}

	t := Unkown
	for _, h := range res.Header["Content-Type"] {
		if strings.Contains(h, string(Html)) {
			t = Html
			break
		}
		if strings.Contains(h, string(Image)) {
			t = Image
			break
		}
	}

	return &response{Body: b, Type: t}, nil
}
