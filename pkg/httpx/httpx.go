package httpx

import (
	"bytes"
	"context"
	"encoding/json"

	"io"
	"net/http"

	"time"
)

// req, err := http.NewRequestWithContext(ctx, "POST", "https://api.flutterwave.com/v3/virtual-account-numbers", bytes.NewBuffer(flwPayload))
// if err != nil {
// 	return nil, err
// }
// req.Header.Set("Authorization", "Bearer "+os.Getenv("FLW_SECRET_KEY"))
// req.Header.Set("Content-Type", "application/json")

// client := &http.Client{}
// resp, err := client.Do(req)
// if err != nil {
// 	return nil, err
// }
// defer resp.Body.Close()

// respbody, err := io.ReadAll(resp.Body)
// if err != nil {
// 	return nil, fmt.Errorf("failed to read response body: %v", err)
// }

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

