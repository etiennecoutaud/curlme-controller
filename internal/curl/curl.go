package curl

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// Curl type with httpClient
type Curl struct {
	netClient *http.Client
}

// New init new curl type with http client
func New() *Curl {
	return &Curl{
		netClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// CallingURL call url to retrieve value
func (c *Curl) CallingURL(url string) (string, error) {
	resp, err := c.netClient.Get(formatURL(url))
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", errors.New("status code: " + string(resp.StatusCode))
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// formatURL add http:// to url for http client
func formatURL(url string) string {
	if strings.Contains(url, "http://") ||
		strings.Contains(url, "https://") {
		return url
	}
	return "http://" + url
}

// SetClientHTTP setter for http client in case of test
func (c *Curl) SetClientHTTP(client *http.Client) {
	c.netClient = client
}
