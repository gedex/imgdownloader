// Copyright 2013 Akeda Bagus <admin@gedex.web.id>. All rights reversed.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.

// Package provider provides images link that can be dowloaded, typically
// via REST API of image-sharing-sites such as Flickr and Instagram.
package provider

import (
	"fmt"
)

const (
	defaultUserAgent = "imgdownloader"
)

var (
	providers map[string]Provider
)

// A type that satisfies provider.Provider can be requested
// to return list of image links tagged with particular string.
type Provider interface {
	// Set configuration of this provider.
	Configure(map[string]string)
	// Request n links (to image) tagged with tag.
	Request(tag string, n uint) (map[string]string, error)
}

func init() {
	providers = make(map[string]Provider)
	providers["flickr"] = NewFlickr()
	providers["instagram"] = NewInstagram()
}

// Get gets provider of a given provider string.
func Get(provider string) (p Provider, err error) {
	p, err = getProvider(provider)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func getProvider(name string) (p Provider, err error) {
	var ok bool
	if _, ok = providers[name]; !ok {
		return nil, fmt.Errorf("undefined provider %s", name)
	}
	if p, ok = providers[name].(Provider); !ok {
		return nil, fmt.Errorf("provider %s doesn't implement Provider", name)
	}
	return
}