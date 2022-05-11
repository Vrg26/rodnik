package image_service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

type client struct {
	url        string
	httpClient *http.Client
}

func NewClient(url string, httpClient *http.Client) *client {
	return &client{url, httpClient}
}

func (c *client) GetURL() string {
	return c.url
}

func (c *client) Upload(ctx context.Context, image []byte) (*http.Response, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fileWriter, err := writer.CreateFormFile("file", "test.png")

	if err != nil {
		return nil, err
	}
	r := bytes.NewReader(image)
	_, err = io.Copy(fileWriter, r)
	if err != nil {
		return nil, err
	}
	writer.Close()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/upload", c.url), bytes.NewReader(body.Bytes()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return c.httpClient.Do(req)
}
