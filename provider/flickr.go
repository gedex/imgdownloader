package provider

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	flickrBaseURL = "http://api.flickr.com/services/rest"
)

type Flickr struct {
	baseURL *url.URL
	config  map[string]string
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

func (p *Flickr) Request(tag string, n uint) (map[string]string, error) {
	q := p.baseURL.Query()

	q.Set("method", "flickr.photos.search")
	q.Set("format", "json")
	q.Set("nojsoncallback", "1")
	q.Set("tags", tag)
	if apiKey, ok := p.config["api_key"]; ok {
		q.Set("api_key", apiKey)
	}

	p.baseURL.RawQuery = q.Encode()
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

	listToDownloads := make(map[string]string)
	for _, fp := range fr.Photos.Photo {
		filename := fmt.Sprintf("%s.jpg", fp.ID)
		listToDownloads[filename] = fp.getURL()
	}
	return listToDownloads, nil
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
