package utils

import (
	"io/ioutil"
	"net/http"
)



func ReadBody(resp *http.Response) ([]byte, error) {
	var body []byte = nil
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		if bodyBytes, err := ioutil.ReadAll(resp.Body); err != nil {
			return nil, err
		} else {
			body = bodyBytes
		}
	}
	return body, nil
}
