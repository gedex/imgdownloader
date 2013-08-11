package provider

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

const (
	flickrBaseURL = "http://api.flickr.com/services/rest"
	MAX_PER_PAGE  = 500
)

type Flickr struct {
	baseURL *url.URL
	config  map[string]string
	n       uint
}

func NewFlickr() *Flickr {
	baseURL, _ := url.Parse(flickrBaseURL)
	return &Flickr{
		baseURL: baseURL,
	}
}

func (p *Flickr) Configure(c map[string]string) {
	p.config = c
}

func (p *Flickr) Request(tag string, n uint) (ProviderResponse, error) {
	numOfRequestedImages := int(n)
	listToDownloads := make(ProviderResponse, numOfRequestedImages)
	cursorIndex := 0

	// The maximum allowed value of `per_page` parameter is 500
	perPage, page, remaining := MAX_PER_PAGE, 1, 0

	if numOfRequestedImages <= MAX_PER_PAGE {
		perPage = numOfRequestedImages
	} else {
		remaining = numOfRequestedImages - perPage
	}
	p.prepareURL(tag, perPage, page)

	for {
		fp, err := p.makeRequest()
		if err != nil {
			return nil, err
		}

		for _, photo := range fp.Photo {
			item := &ProviderItem{
				Filename: fmt.Sprintf("%s.jpg", photo.ID),
				Link:     photo.getURL(),
			}
			listToDownloads[cursorIndex] = item
			cursorIndex += 1
		}

		if remaining == 0 {
			break
		} else {
			page, remaining := page+1, remaining-perPage
			if remaining < 0 {
				perPage = remaining
				remaining = 0
			}
			p.prepareURL(tag, perPage, page)
		}
	}

	return listToDownloads, nil
}

func (p *Flickr) prepareURL(tag string, perPage, page int) {
	q := p.baseURL.Query()

	q.Set("method", "flickr.photos.search")
	q.Set("format", "json")
	q.Set("nojsoncallback", "1")
	q.Set("tags", tag)
	q.Set("per_page", strconv.Itoa(perPage))
	q.Set("page", strconv.Itoa(page))

	if apiKey, ok := p.config["api_key"]; ok {
		q.Set("api_key", apiKey)
	}

	p.baseURL.RawQuery = q.Encode()
}

func (p *Flickr) makeRequest() (*FlickrPhotos, error) {
	u := p.baseURL.String()

	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%v returned %v")
	}

	fr := new(FlickrResponse)

	err = json.NewDecoder(resp.Body).Decode(fr)
	if err != nil {
		return nil, err
	}
	if fr.Stat != "ok" {
		return nil, fmt.Errorf("flickr responded with %v, reason: %v", fr.Stat, fr.Message)
	}
	return fr.Photos, nil
}

type FlickrResponse struct {
	Photos  *FlickrPhotos `json:"photos,omitempty"`
	Stat    string        `json:"stat,omitempty"`
	Message string        `json:"message,omitempty"`
}

type FlickrPhotos struct {
	Page    int            `json:"page,omitempty"`
	Pages   int            `json:"pages,omitempty"`
	PerPage int            `json:"perpage,omitempty"`
	Total   string         `json:"total,omitempty"`
	Photo   []*FlickrPhoto `json:"photo,omitempty"`
}

type FlickrPhoto struct {
	ID     string `json:"id,omitempty"`
	Owner  string `json:"owner,omitempty"`
	Secret string `json:"secret,omitempty"`
	Server string `json:"server,omitempty"`
	Farm   int    `json:"farm,omitempty"`
	Title  string `json:"title,omitempty"`
}

func (fp *FlickrPhoto) getURL() string {
	sourceURL := fmt.Sprintf("http://farm%d.staticflickr.com/%s/%s_%s.jpg", fp.Farm, fp.Server, fp.ID, fp.Secret)
	return sourceURL
}
