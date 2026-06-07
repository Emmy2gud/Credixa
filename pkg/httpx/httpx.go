package httpx

import (
	"bytes"
	"context"
	"encoding/json"

	"io"
	"net/http"

	"time"
)


type Client struct {
	baseURL    string
	httpClient *http.Client
	headers    map[string]string
	timeout    time.Duration
}

//a constructor function to create a new httpx client instance 

func New(baseURL string, headers map[string]string) *Client {
	return &Client{
		baseURL:    baseURL,
		headers:    headers,
		httpClient: &http.Client{},
		timeout: 30*time.Second,
	}
}

func (c *Client) DoRequest(ctx context.Context, method string, endpoint string, payload interface{}, reqest interface{}) ([]byte, error) {

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		method,
		c.baseURL+endpoint,
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, err
	}

	for k, v := range c.headers {
		req.Header.Set(k, v)
	}



	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respbody, err := io.ReadAll(resp.Body)
	return respbody, err
}

