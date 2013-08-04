package main

import (
	"fmt"
	"net/url"
)

var (
	serviceURL = "http://www.rsk.is/fyrirtaekjaskra/leit"
)

func ReadOrGetSSID(ssid string) (content []byte, err error) {
	filename := fmt.Sprintf("./cache/fskra-%s.html", ssid)
	url := fmt.Sprintf("%s/kennitala/%s", serviceURL, ssid)
	content, err = ReadOrGetURL(filename, url)
	return
}

func ReadOrGetSearch(street string) (content []byte, err error) {
	filename := fmt.Sprintf("./cache/fskra-search-%s.html", street)
	u, err := url.Parse(serviceURL)
	if err != nil {
		return
	}
	q := u.Query()
	q.Set("heimili", street)
	u.RawQuery = q.Encode()
	content, err = ReadOrGetURL(filename, u.String())
	return
}
