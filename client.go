package go_ernie

import (
	"context"
	"encoding/json"
	"fmt"
	erinie "github.com/LazyLinchen/go-ernie/internal"
	"github.com/pkg/errors"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

var client = map[string]*Client{}

type Client struct {
	config         ClientConfig
	accessToken    string
	requestBuilder erinie.RequestBuilder
}

func NewClient(ak, sk string) (*Client, error) {
	config := DefaultConfig(ak, sk)
	return NewClientWithConfig(config)
}

func NewClientWithConfig(config ClientConfig) (*Client, error) {
	if config.AK == "" || config.SK == "" {
		return nil, errors.New("config is error")
	}
	if c, ok := client[config.AK]; ok {
		return c, nil
	}
	c := &Client{
		config:         config,
		requestBuilder: erinie.NewRequestBuilder(),
	}
	err := c.setAccessToken()
	if err != nil {
		return nil, errors.New("setAccessToken error: " + err.Error())
	}
	refreshAccessToken(c)
	client[config.AK] = c
	return client[config.AK], nil
}

func refreshAccessToken(c *Client) {
	go func() {
		for {
			time.Sleep(118 * time.Minute)
			err := c.setAccessToken()
			if err != nil {
				log.Println("refreshAccessToken error: " + err.Error())
			}
		}
	}()
}

func (c *Client) GetAccessToken() string {
	return c.accessToken
}

func (c *Client) setAccessToken() error {
	uri := c.config.BaseURL + "/oauth/2.0/token?grant_type=client_credentials&client_id=" + c.config.AK + "&client_secret=" + c.config.SK
	req, err := c.newRequest(context.Background(), "POST", uri, withBody(strings.NewReader(``)))
	if err != nil {
		return errors.New("newRequest error: " + err.Error())
	}
	var resp = struct {
		AccessToken  string `json:"access_token"`
		ExpiresIn    int64  `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
		Error        string `json:"error"`
	}{}
	err = c.sendRequest(req, &resp)
	if err != nil {
		return errors.New("sendRequest error: " + err.Error())
	}
	if resp.Error != "" {
		return errors.New("resp.Error: " + resp.Error)
	}
	c.accessToken = resp.AccessToken
	return nil
}

type requestOptions struct {
	body   any
	header http.Header
}

type requestOption func(*requestOptions)

func withBody(body any) requestOption {
	return func(args *requestOptions) {
		args.body = body
	}
}

func withHeader(header http.Header) requestOption {
	return func(args *requestOptions) {
		args.header = header
	}
}

func (c *Client) newRequest(ctx context.Context, method, url string, setters ...requestOption) (*http.Request, error) {
	args := &requestOptions{
		body:   nil,
		header: make(http.Header),
	}
	for _, setter := range setters {
		setter(args)
	}
	req, err := c.requestBuilder.Build(ctx, method, url, args.body, args.header)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (c *Client) sendRequest(req *http.Request, v any) error {
	req.Header.Set("Accept", "application/json; charset=utf-8")
	contentType := req.Header.Get("Content-Type")
	if contentType == "" {
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
	}
	res, err := c.config.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		return c.handleErrorResp(res)
	}

	return decodeResponse(res.Body, v)
}

func (c *Client) sendRequestRaw(req *http.Request) (body io.ReadCloser, err error) {
	resp, err := c.config.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		err = c.handleErrorResp(resp)
		return
	}
	return resp.Body, nil
}

func (c *Client) fullURL(suffix string, args ...any) string {
	return fmt.Sprintf("%s%s%s?access_token=%s", c.config.BaseURL, c.config.AiApiURL, suffix, c.accessToken)
}

func decodeResponse(body io.Reader, v any) error {
	if v == nil {
		return nil
	}
	if result, ok := v.(*string); ok {
		return decodeString(body, result)
	}
	return json.NewDecoder(body).Decode(v)
}

func decodeString(body io.Reader, output *string) error {
	b, err := io.ReadAll(body)
	if err != nil {
		return err
	}
	*output = string(b)
	return nil
}

func (c *Client) handleErrorResp(res *http.Response) error {
	return errors.New("HTTP status code: " + res.Status)
}

func sendRequestStream[T streamable](client *Client, req *http.Request) (*streamReader[T], error) {
	stream := &streamReader[T]{
		emptyMessagesLimit: client.config.EmptyMessagesLimit,
		isFinished:         false,
		response:           nil,
		scanner:            nil,
		errAccumulator:     erinie.NewErrorAccumulator(),
		unmarshaler:        &erinie.JSONUnmarshaler{},
	}
	resp, err := client.config.HTTPClient.Do(req)
	if err != nil {
		stream.isFinished = true
		return stream, err
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		stream.isFinished = true
		err = client.handleErrorResp(resp)
		return stream, err
	}
	stream.response = resp
	return stream, nil
}
