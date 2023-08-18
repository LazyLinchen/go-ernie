package go_ernie

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	erinie "github.com/LazyLinchen/go-ernie/internal"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"strings"
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
	// set access token
	// TODO 不应该在这里直接设置，而是需要在后台配置并且定时刷新，并且集群共用token
	err := c.setAccessToken()
	if err != nil {
		return nil, errors.New("setAccessToken error: " + err.Error())
	}
	client[config.AK] = c
	return client[config.AK], nil
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
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")

	resp, err := client.config.HTTPClient.Do(req)
	if err != nil {
		return new(streamReader[T]), err
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		return new(streamReader[T]), client.handleErrorResp(resp)
	}

	// 百度的傻逼接口，报错不直接返回 400，而是返回 200，然后返回一个 json，里面有 error 字段，非SSE格式
	// 需要取出来前面6个byte做判断是不是SSE，如果不是SSE那就是报错了，直接返回报错内容
	pr, err := newPeekingReader(resp.Body, 6)
	if err != nil {
		return new(streamReader[T]), err
	}
	buf := make([]byte, 6)
	pr.Read(buf)

	// SSE response starts with "data: "
	// 取出来之后还得还回去，不然就不是完整的请求体了，后面的流读取会有问题
	resp.Body = io.NopCloser(io.MultiReader(bytes.NewReader(buf), resp.Body))

	// 如果不是 SSE 格式，那就是报错了，直接给报错内容返回
	if !bytes.HasPrefix(buf, []byte("data: ")) {
		var respError APIError
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return new(streamReader[T]), errors.New("Unexpected SSE response" + err.Error())
		}
		err = json.Unmarshal(b, &respError)
		if err != nil {
			return new(streamReader[T]), errors.New("Unexpected SSE response" + err.Error())
		}
		return new(streamReader[T]), &respError
	}

	return &streamReader[T]{
		emptyMessagesLimit: client.config.EmptyMessagesLimit,
		reader:             bufio.NewReader(resp.Body),
		response:           resp,
		errAccumulator:     erinie.NewErrorAccumulator(),
		unmarshaler:        &erinie.JSONUnmarshaler{},
	}, nil
}
