package ghc

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"reflect"

	"github.com/valyala/fasthttp"
)

// Options ...
type Options struct {
	Header http.Header
}

// Client ....
type Client struct {
	u          *url.URL
	opts       *Options
	httpClient *fasthttp.Client
}

// New ...
func New(apiURL string, opts *Options) (*Client, error) {
	var (
		err error
		url *url.URL
		cl  *fasthttp.Client
	)

	if opts != nil {
		if cl, err = initHTTPClient(opts); err != nil {
			return nil, errors.New("can't configure client:" + err.Error())
		}
	}

	u, err := url.Parse(apiURL)
	if err != nil {
		return nil, err
	}
	return &Client{u: u, opts: opts, httpClient: cl}, nil
}

// Execute ...
func (c *Client) Execute(req *Request) (*Response, error) {
	data, err := json.Marshal(req)
	if err != nil {
		return nil, errors.New("execute: marshal error: " + err.Error())
	}

	switch req.Type() {
	case reflect.TypeOf(Query("")), reflect.TypeOf(Mutation("")):
		return c.httpExecute(c.u, data)
	case reflect.TypeOf(Subscription("")):
		return c.wsExecute(c.u, data)
	default:
		return nil, errors.New("execute: unsupported request type: " + req.Type().String())
	}
}

func (c *Client) httpExecute(u *url.URL, data []byte) (*Response, error) {
	req := fasthttp.AcquireRequest()
	req.Header.SetMethod("POST")
	req.SetRequestURI(u.String())
	req.SetBody(data)

	if c.opts != nil && c.opts.Header != nil {
		for k, v := range c.opts.Header {
			for _, h := range v {
				req.Header.Add(k, h)
			}
		}
	}

	res := fasthttp.AcquireResponse()
	err := c.httpClient.Do(req, res)
	if err != nil {
		return nil, errors.New("execute: request: " + err.Error())
	}

	resp := &Response{}

	decoder := json.NewDecoder(bytes.NewReader(res.Body()))

	if err = decoder.Decode(resp); err != nil {
		return nil, errors.New("execute: decode response: " + err.Error())
	}

	if resp.Errors != nil {
		return resp, errors.New("execute: server returned error: see response.Errors for details")
	}

	return resp, nil
}

func (c *Client) wsExecute(u *url.URL, data []byte) (*Response, error) {
	panic("isn't supported yet")
}

func initHTTPClient(opts *Options) (*fasthttp.Client, error) {
	return &fasthttp.Client{}, nil
}
