package sendgrid

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type Client struct {
	client  *http.Client
	baseURL *url.URL

	// options
	opts Options

	common service

	Contact *ContactService
	List    *ListService
	Mail    *MailService
}

// Options can be used to create a customized client
type Options struct {
	apiKey    string
	apiSecret string

	debug bool
}

type service struct {
	client *Client
	opts   Options
}

const (
	apiURL = "https://api.sendgrid.com/v3/"
)

type Option func(*Options) error

func WithApp(apiKey, apiSecret string) Option {
	return func(o *Options) error {
		o.apiKey = apiKey
		o.apiSecret = apiSecret
		return nil
	}
}

func WithDebug(v bool) Option {
	return func(o *Options) error {
		o.debug = v
		return nil
	}
}

func NewClient(httpClient *http.Client, options ...Option) (*Client, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	c := &Client{client: httpClient}
	for _, opt := range options {
		if err := opt(&c.opts); err != nil {
			return nil, err
		}
	}

	c.baseURL, _ = url.Parse(apiURL)
	c.common.client = c
	c.common.opts = c.opts

	c.Contact = (*ContactService)(&c.common)
	c.List = (*ListService)(&c.common)
	c.Mail = (*MailService)(&c.common)

	return c, nil
}

func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	if !strings.HasSuffix(c.baseURL.Path, "/") {
		return nil, fmt.Errorf("BaseURL must have a trailing slash, but %q does not", c.baseURL)
	}
	u, err := c.baseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	if c.opts.debug {
		fmt.Printf("[sendgrid] %s %s\n", method, u.String())
	}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	defer func() {
		if c.opts.debug {
			b, _ := httputil.DumpRequest(req, true)
			fmt.Printf("[sendgrid] %s", string(b))
		}
	}()
	if err != nil {
		return nil, err
	}

	// add authorization header
	req.Header.Set("Authorization", "BEARER "+c.opts.apiSecret)

	if body != nil {
		req.Header.Add("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
	}

	return req, nil
}

// Do sends an API request and returns the API response. The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred. If v implements the io.Writer
// interface, the raw response body will be written to v, without attempting to
// first decode it.
//
// The provided ctx must be non-nil. If it is canceled or times out,
// ctx.Err() will be returned.
// TODO: Rate limiting
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	req = req.WithContext(ctx)

	resp, err := c.client.Do(req)
	if err != nil {
		// If we got an error, and the context has been canceled,
		// the context's error is probably more useful.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		return nil, err
	}

	defer func() {
		// Drain up to 512 bytes and close the body to let the Transport reuse the connection
		io.CopyN(io.Discard, resp.Body, 512)
		resp.Body.Close()
	}()

	err = CheckResponse(resp)
	if err != nil {
		return resp, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			if c.opts.debug {
				b, _ := httputil.DumpResponse(resp, true)
				fmt.Printf("[sendgrid]: %s\n", string(b))
			}
			err = json.NewDecoder(resp.Body).Decode(v)
			if err == io.EOF {
				err = nil // ignore EOF errors caused by empty response body
			}
		}
	}

	return resp, err
}

type ErrorResponse struct {
	Response *http.Response // HTTP response that caused this error
	Errors   []*Error       `json:"errors,omitempty"`
}

type Error struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	ErrorID string `json:"error_id"`
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("[%v] %v: %d | (1/%d) %s",
		r.Response.Request.Method,
		sanitizeURL(r.Response.Request.URL),
		r.Response.StatusCode,
		len(r.Errors),
		r.Errors[0].Message,
	)
}

// sanitizeURL redacts the client_secret parameter from the URL which may be
// exposed to the user.
func sanitizeURL(uri *url.URL) *url.URL {
	if uri == nil {
		return nil
	}
	params := uri.Query()
	if len(params.Get("client_secret")) > 0 {
		params.Set("client_secret", "REDACTED")
		uri.RawQuery = params.Encode()
	}
	return uri
}

// CheckResponse checks the API response for errors, and returns them if
// present. A response is considered an error if it has a status code outside
// the 200 range or equal to 202 Accepted.
// API error responses are expected to have response
// body, and a JSON response body that maps to ErrorResponse.
//
// The error type will be *RateLimitError for rate limit exceeded errors,
// *AcceptedError for 202 Accepted status codes,
// and *TwoFactorAuthError for two-factor authentication errors.
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}
	errorResponse := &ErrorResponse{Response: r}
	data, err := io.ReadAll(r.Body)
	if err == nil && data != nil {
		json.Unmarshal(data, errorResponse)
	}
	return errorResponse
}
