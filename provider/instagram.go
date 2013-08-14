package provider

import (
	"fmt"

	"github.com/gedex/go-instagram/instagram"
)

type Instagram struct {
	client *instagram.Client
}

func NewInstagram() *Instagram {
	c := instagram.NewClient(nil)
	return &Instagram{
		client: c,
	}
}

func (p *Instagram) Configure(c map[string]string) {
	if accessToken, ok := c["access_token"]; ok {
		p.client.AccessToken = accessToken
	}
}

func (p *Instagram) Request(tag string, n uint) (ProviderResponse, error) {
	numOfRequestedImages := int(n)
	listToDownloads := make(ProviderResponse, numOfRequestedImages)
	cursorIndex := 0
	doneFilling := false // true if cursorIndex >= numOfRequestedImages
	maxId := ""

	tags, _, err := p.client.Tags.Search(tag)
	if err != nil {
		return nil, err
	}

	for _, t := range tags {

		// Iterate until the end of media pagination or error occured
		for {
			media, pageMedia, err := p.recentMedia(t.Name, maxId)
			if err != nil {
				break
			}

			for _, m := range media {
				if cursorIndex >= numOfRequestedImages {
					doneFilling = true
					break
				}
				if m.Type != "image" {
					continue
				}
				item := &ProviderItem{
					Filename: fmt.Sprintf("%s.jpg", m.ID),
					Link:     m.Images.StandardResolution.URL,
				}
				listToDownloads[cursorIndex] = item
				cursorIndex++
			}

			maxId = pageMedia.NextMaxID
			if maxId == "" || doneFilling {
				break
			}
		}

		if doneFilling {
			break
		}
	}

	return listToDownloads, nil
}

func (p *Instagram) recentMedia(tag string, maxId string) (media []instagram.Media, page *instagram.ResponsePagination, err error) {
	opt := new(instagram.Parameters)
	if maxId != "" {
		opt.MaxID = maxId
	}

	media, page, err = p.client.Tags.RecentMedia(tag, opt)
	return
}
