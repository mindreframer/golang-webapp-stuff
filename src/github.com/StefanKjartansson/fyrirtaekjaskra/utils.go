package fyrirtaekjaskra

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

func ReadOrGetURL(filename string, url string) (io.Reader, error) {

	if _, err := os.Stat(filename); err == nil {
		fi, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		return fi, nil
	}

	logger.Debugf("Fetching: %s\n", url)

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, err
	}

	err = ioutil.WriteFile(filename, content, 0644)

	return bytes.NewReader(content), nil

}
