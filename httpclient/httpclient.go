package httpclient

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type HttpInterface interface {
	Get(url string) ([]byte, error)
}

type HttpClient struct {
	client *http.Client
}

func NewHttpClient(timeout time.Duration) *HttpClient {
	return &HttpClient{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *HttpClient) Get(url string) ([]byte, error) {
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.New("not found")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non 200 status code: %v", resp.StatusCode)
	}

	return data, nil
}
