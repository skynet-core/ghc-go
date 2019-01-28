package ghc

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
)

// Options ...
type Options struct {
	Header http.Header
}

// Client ....
type Client struct {
	u          *url.URL
	opts       *Options
	httpClient *http.Client
}

// New ...
func New(apiURL string, opts *Options) (*Client, error) {
	var (
		err error
		url *url.URL
		cl  *http.Client
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
	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewReader(data))
	if err != nil {
		return nil, errors.New("execute: new request: " + err.Error())
	}

	if c.opts != nil && c.opts.Header != nil {
		for k, v := range c.opts.Header {
			for _, h := range v {
				req.Header.Add(k, h)
			}
		}
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.New("execute: request: " + err.Error())
	}

	data, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New("execute: read response: " + err.Error())
	}

	defer res.Body.Close()

	resp := &Response{}

	decoder := json.NewDecoder(bytes.NewReader(data))

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

func initHTTPClient(opts *Options) (*http.Client, error) {
	if opts != nil {
		return nil, errors.New("isn't supported yet")
	}
	return http.DefaultClient, nil
}
