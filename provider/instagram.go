package provider

import (
	"fmt"
	"net/url"
)

const (
	instagramBaseURL = "https://instagram.com"
)

type Instagram struct {
	baseURL *url.URL
	config  map[string]string
}

func NewInstagram() *Instagram {
	baseURL, _ := url.Parse(instagramBaseURL)
	return &Instagram{
		baseURL: baseURL,
	}
}

func (p *Instagram) Configure(c map[string]string) {
	fmt.Printf("Flickr.Configure(%v)\n", c)
	p.config = c
}

func (p *Instagram) Request(tag string, n uint) (map[string]string, error) {
	fmt.Printf("Instagram.Request(%s, %d)\n", tag, n)
	return nil, nil
}
