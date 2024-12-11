package crawler

import (
	gourl "net/url"

	"github.com/kfc-manager/vision-seeker/adapter/client"
	"github.com/kfc-manager/vision-seeker/domain/html"
	"github.com/kfc-manager/vision-seeker/domain/image"
	"github.com/kfc-manager/vision-seeker/service/data"
)

type Service interface {
	Crawl()
}

type service struct {
	client client.Client
	data   data.Service
}

func New(c client.Client, d data.Service) *service {
	return &service{client: c, data: d}
}

func (s *service) Crawl() {
	loc, err := s.data.Url()
	if err != nil {
		return
	}
	defer s.Crawl()

	res, err := s.client.Get(loc)
	if err != nil {
		return
	}

	if res.Type == client.Image {
		img, err := image.LoadImage(res.Body)
		if err != nil {
			return
		}
		if !img.Valid(500, 500, 6.0, false) {
			return
		}
		_ = s.data.StoreImage(img)
	}

	url, err := gourl.Parse(loc)
	if err != nil {
		return
	}

	if res.Type == client.Html {
		doc, err := html.Parse(res.Body)
		if err != nil {
			return
		}

		for _, i := range doc.Images() {
			src := i.Attribute("src")
			if len(src) < 1 {
				continue
			}
			img, err := gourl.Parse(src)
			if err != nil {
				continue
			}
			if len(img.Scheme) < 1 {
				img.Scheme = url.Scheme
			}
			if len(img.Host) < 1 {
				img.Host = url.Host
			}
			_ = s.data.SetUrl(img)
		}

		for _, l := range doc.Links() {
			link, err := gourl.Parse(l)
			if err != nil {
				continue
			}
			if len(link.Scheme) < 1 {
				link.Scheme = url.Scheme
			}
			if len(link.Host) < 1 {
				link.Host = url.Host
			}
			_ = s.data.SetUrl(link)
		}
	}

	s.Crawl()
}
