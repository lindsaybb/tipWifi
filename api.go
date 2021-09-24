package uClig

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

// GetRequest structures the HTTP "GET" call to the supplied endpoint.
func GetRequest(oAuth *OAuth2, endpoint, uri string) (*http.Response, error) {
	url := fmt.Sprintf("https://%s/api/v1/%s", endpoint, uri)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", oAuth.AccessToken))
	if Debug {
		fmt.Printf("|-| GET %s |-|\n", url)
	}
	return http.DefaultClient.Do(req)
}

// PostRequest structures the HTTP "POST" call to the supplied endpoint.
func PostRequest(oAuth *OAuth2, endpoint, uri string, data []byte) (*http.Response, error) {
	var raw io.Reader
	if data != nil {
		raw = bytes.NewReader(data)
	}
	url := fmt.Sprintf("https://%s/api/v1/%s", endpoint, uri)
	req, err := http.NewRequest("POST", url, raw)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	// if Posting to request OAuth2 token, skip the Bearer declaration
	if oAuth != nil {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", oAuth.AccessToken))
	}
	if Debug {
		fmt.Printf("|-| POST %s |-|\n", url)
	}
	return http.DefaultClient.Do(req)
}

// PutRequest structures the HTTP "PUT" call to the supplied endpoint.
func PutRequest(oAuth *OAuth2, endpoint, uri string, data []byte) (*http.Response, error) {
	var raw io.Reader
	if data != nil {
		raw = bytes.NewReader(data)
	}
	url := fmt.Sprintf("https://%s/api/v1/%s", endpoint, uri)
	req, err := http.NewRequest("PUT", url, raw)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", oAuth.AccessToken))
	if Debug {
		fmt.Printf("|-| PUT %s |-|\n", url)
	}
	return http.DefaultClient.Do(req)
}

// DeleteRequest structures the HTTP "DELETE" call to the supplied endpoint.
func DeleteRequest(oAuth *OAuth2, endpoint, uri string, data []byte) (*http.Response, error) {
	var raw io.Reader
	if data != nil {
		raw = bytes.NewReader(data)
	}
	// will this error if passed no bytes? Same as nil or NO?
	url := fmt.Sprintf("https://%s/api/v1/%s", endpoint, uri)
	req, err := http.NewRequest("DELETE", url, raw)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", oAuth.AccessToken))
	if Debug {
		fmt.Printf("|-| DELETE %s |-|\n", url)
	}
	return http.DefaultClient.Do(req)
}
