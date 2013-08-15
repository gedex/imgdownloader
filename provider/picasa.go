package provider

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

const (
	picasaFeedURL = "https://picasaweb.google.com/data/feed/api/all"
)

type Picasa struct {
	baseURL *url.URL
	config  map[string]string
}

func NewPicasa() *Picasa {
	baseURL, _ := url.Parse(picasaFeedURL)
	return &Picasa{
		baseURL: baseURL,
	}
}

func (p *Picasa) Configure(c map[string]string) {
	p.config = c
}

func (p *Picasa) Request(tag string, n uint) (ProviderResponse, error) {
	numOfRequestedImages := int(n)
	listToDownloads := make(ProviderResponse, numOfRequestedImages)
	cursorIndex := 0

	p.prepareURL(tag, numOfRequestedImages)
	u := p.baseURL.String()

	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%v returned %v", u, resp.StatusCode)
	}

	result := new(PicasaResponse)

	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		return nil, err
	}

	h := md5.New()

	for _, entry := range result.Feed.Entry {
		if entry.MediaGroup == nil {
			continue
		}
		if len(entry.MediaGroup.Content) == 0 {
			continue
		}

		io.WriteString(h, entry.ID.String)

		item := &ProviderItem{
			Filename: fmt.Sprintf("%x.jpg", h.Sum(nil)),
			Link:     entry.MediaGroup.Content[0].URL,
		}
		listToDownloads[cursorIndex] = item

		cursorIndex += 1
		h.Reset()
	}

	return listToDownloads, nil
}

func (p *Picasa) prepareURL(tag string, maxResults int) {
	q := p.baseURL.Query()

	q.Set("q", tag)
	q.Set("max-results", strconv.Itoa(maxResults))
	q.Set("fields", "entry(id,media:group)")
	q.Set("alt", "json")

	p.baseURL.RawQuery = q.Encode()
}

type PicasaResponse struct {
	Feed *PicasaFeed `json:"feed,omitempty"`
}

type PicasaFeed struct {
	Entry []PicasaEntry `json:"entry,omitempty"`
}

type PicasaEntry struct {
	ID         *PicasaMediaID    `json:"id,omitempty"`
	MediaGroup *PicasaMediaGroup `json:"media$group,omitempty"`
}

type PicasaMediaID struct {
	String string `json:"$t"`
}

type PicasaMediaGroup struct {
	Content []PicasaMediaContent `json:"media$content,omitempty"`
}

type PicasaMediaContent struct {
	URL string `json:"url"`
}
