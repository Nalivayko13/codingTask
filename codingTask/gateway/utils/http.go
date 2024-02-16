package utils

import (
	"io"
	"net/http"
)

func HttpGetCallWithHeader(url string, header map[string]string) (*http.Response, error) {
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	for key, val := range header {
		req.Header.Set(key, val)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return resp, nil
}

func HttpGetCallWithParam(url string, queryParam, header map[string]string) (string, int, error) {
	client := &http.Client{}

	newUrl := url + "?" + "login=" + queryParam["login"]

	req, err := http.NewRequest(http.MethodGet, newUrl, nil)
	if err != nil {
		return "", 0, err
	}

	for key, val := range header {
		req.Header.Set(key, val)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, err
	}

	return string(data), resp.StatusCode, nil
}
