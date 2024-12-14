package crawler

import (
	gourl "net/url"

	"github.com/kfc-manager/vision-seeker/crawler/adapter/client"
	"github.com/kfc-manager/vision-seeker/crawler/domain/html"
	"github.com/kfc-manager/vision-seeker/crawler/domain/image"
	"github.com/kfc-manager/vision-seeker/crawler/service/data"
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
	url, alt, err := s.data.Next()
	if err != nil {
		return
	}
	defer s.Crawl()

	res, err := s.client.Get(url.String())
	if err != nil {
		return
	}

	if res.Type == client.Image {
		img, err := image.Load(res.Body)
		if err != nil {
			return
		}
		if !img.Valid(300, 300, 3.0, false) {
			return
		}
		_ = s.data.StoreImage(img, alt)
	}

	if res.Type == client.Html {
		doc, err := html.Parse(res.Body)
		if err != nil {
			return
		}

		for _, img := range doc.Images() {
			src := img.Attribute("src")
			if len(src) < 1 {
				continue
			}
			imgUrl, err := gourl.Parse(src)
			if err != nil {
				continue
			}
			if len(imgUrl.Scheme) < 1 {
				imgUrl.Scheme = url.Scheme
			}
			if len(imgUrl.Host) < 1 {
				imgUrl.Host = url.Host
			}
			_ = s.data.Visit(imgUrl, img.Attribute("alt"))
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
			_ = s.data.Visit(link, "")
		}
	}

	s.Crawl()
}
