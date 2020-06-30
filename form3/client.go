package form3

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

const (
	apiVersion    string = "v1"
	contentType   string = "application/vnd.Form3+json"
	defaultScheme string = "http"
	defaultHost   string = "localhost:8080"
)

// Client is a Form3 client. Create one by calling NewClient.
type Client struct {
	httpClient *http.Client
	infoLog    *log.Logger // info log for non-critical messages
	errorLog   *log.Logger // error log for critical messages
	scheme     string      // http or https
	host       string      // host
}

// NewClient creates a new client to work with the Form3 API.
func NewClient(options ...ClientOptionFunc) (*Client, error) {
	c := &Client{
		scheme:     defaultScheme,
		host:       defaultHost,
		httpClient: &http.Client{},
		infoLog:    log.New(os.Stderr, "[form3_info]", log.LstdFlags),
		errorLog:   log.New(os.Stderr, "[form3_error]", log.LstdFlags),
	}

	// Apply passed options (if any), overriding defaults
	for _, option := range options {
		if err := option(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

// ClientOptionFunc is a function that configures a Client.
type ClientOptionFunc func(*Client) error

// MakeRequestOptions must be passed into MakeRequest.
type MakeRequestOptions struct {
	Method string
	Path   string
	Params url.Values
	Body   interface{}
}

// SetScheme sets the HTTP scheme (http by default)
func SetScheme(scheme string) ClientOptionFunc {
	return func(c *Client) error {
		c.scheme = scheme
		return nil
	}
}

// SetHost sets the host (localhost:8080 by default)
func SetHost(host string) ClientOptionFunc {
	return func(c *Client) error {
		c.host = host
		return nil
	}
}

// MakeRequest makes a HTTP request to the Form3 API.
// It returns a *http.Response and an error (on failure.
func (c *Client) MakeRequest(ctx context.Context, opt MakeRequestOptions) (*http.Response, error) {
	u := url.URL{
		Scheme:   c.scheme,
		Host:     c.host,
		Path:     fmt.Sprintf("%s%s", apiVersion, opt.Path),
		RawQuery: opt.Params.Encode(),
	}

	var payload []byte
	if opt.Body != nil {
		payload, _ = json.Marshal(opt.Body)
	}

	request, _ := http.NewRequest(opt.Method, u.String(), bytes.NewBuffer(payload))
	request.Header.Add("Accept", contentType)
	request.Header.Add("Content-Type", contentType)

	response, err := c.httpClient.Do(request)
	if err != nil {
		return response, err
	}

	err = checkResponse(response)
	if err != nil {
		c.errorf("%s -> %s -> %s", opt.Method, u.String(), err.Error())
		return response, err
	}

	c.infof("%s -> %s -> %s", opt.Method, u.String(), response.Status)
	return response, nil
}

func checkResponse(res *http.Response) error {
	// 200-299 are valid status codes
	if res.StatusCode >= 200 && res.StatusCode <= 299 {
		return nil
	}

	return errors.New(res.Status)
}

// Decode decodes with json.Unmarshal from the Go standard library.
func (c *Client) Decode(response *http.Response, v interface{}) error {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, v)
}

// errorf logs to the error log.
func (c *Client) errorf(format string, args ...interface{}) {
	if c.errorLog != nil {
		c.errorLog.Printf(format, args...)
	}
}

// infof logs info messages.
func (c *Client) infof(format string, args ...interface{}) {
	if c.infoLog != nil {
		c.infoLog.Printf(format, args...)
	}
}

// Accounts returns a service to handle accounts
func (c *Client) Accounts() *AccountsService {
	return NewAccountsService(c)
}
