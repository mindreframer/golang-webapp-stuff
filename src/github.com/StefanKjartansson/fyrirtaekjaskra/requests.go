package fyrirtaekjaskra

import (
	"fmt"
	"io"
	"net/url"
)

var (
	serviceURL = "http://www.rsk.is/fyrirtaekjaskra/leit"
)

func ReadOrGetSSID(ssid string) (io.Reader, error) {
	filename := fmt.Sprintf("./cache/fskra-%s.html", ssid)
	url := fmt.Sprintf("%s/kennitala/%s", serviceURL, ssid)
	return ReadOrGetURL(filename, url)
}

func ReadOrGetSearch(street string) (io.Reader, error) {
	filename := fmt.Sprintf("./cache/fskra-search-%s.html", street)
	u, err := url.Parse(serviceURL)
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("heimili", street)
	u.RawQuery = q.Encode()
	return ReadOrGetURL(filename, u.String())
}
