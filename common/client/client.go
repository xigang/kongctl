package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/net/context/ctxhttp"

	"kongctl/utils/types"
)

var GatewayClient *Client

const (
	MAX_IDLE_CONNS          = 100
	MAX_IDLE_CONNS_PER_HOST = 100
	IDLE_CONN_TIMEOUT       = 90
)

type Client struct {
	scheme            string
	host              string
	client            *http.Client
	customHTTPHeaders map[string]string
}

func NewHTTPClient(host string, headers map[string]string) (*Client, error) {
	url, err := ParseHostURL(host)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:        MAX_IDLE_CONNS,
			MaxIdleConnsPerHost: MAX_IDLE_CONNS_PER_HOST,
			IdleConnTimeout:     IDLE_CONN_TIMEOUT * time.Second,
		},
	}

	return &Client{
		scheme:            url.Scheme,
		host:              url.Host,
		client:            client,
		customHTTPHeaders: headers,
	}, nil
}

func ParseHostURL(host string) (*url.URL, error) {
	protoAddrParts := strings.SplitN(host, "://", 2)
	if len(protoAddrParts) == 1 {
		return nil, fmt.Errorf("unable to parse host `%s`", host)
	}

	return &url.URL{
		Scheme: protoAddrParts[0],
		Host:   protoAddrParts[1],
	}, nil
}

func (cli *Client) Close() error {
	if t, ok := cli.client.Transport.(*http.Transport); ok {
		t.CloseIdleConnections()
	}
	return nil
}

type ServerResponse struct {
	Body       io.ReadCloser
	Header     http.Header
	StatusCode int
	ReqURL     *url.URL
}

type headers map[string][]string

func (cli *Client) Head(ctx context.Context, path string, query url.Values, headers map[string][]string) (ServerResponse, error) {
	return cli.sendRequest(ctx, "HEAD", path, query, nil, headers)
}

func (cli *Client) Get(ctx context.Context, path string, query url.Values, headers map[string][]string) (ServerResponse, error) {
	return cli.sendRequest(ctx, "GET", path, query, nil, headers)
}

func (cli *Client) Post(ctx context.Context, path string, query url.Values, obj interface{}, headers map[string][]string) (ServerResponse, error) {
	body, headers, err := encodeBody(obj, headers)
	if err != nil {
		return ServerResponse{}, err
	}
	return cli.sendRequest(ctx, "POST", path, query, body, headers)
}

func (cli *Client) PostRaw(ctx context.Context, path string, query url.Values, body io.Reader, headers map[string][]string) (ServerResponse, error) {
	return cli.sendRequest(ctx, "POST", path, query, body, headers)
}

func (cli *Client) PATCH(ctx context.Context, path string, query url.Values, obj interface{}, headers map[string][]string) (ServerResponse, error) {
	body, headers, err := encodeBody(obj, headers)
	if err != nil {
		return ServerResponse{}, err
	}
	return cli.sendRequest(ctx, "PATCH", path, query, body, headers)
}

func (cli *Client) Put(ctx context.Context, path string, query url.Values, obj interface{}, headers map[string][]string) (ServerResponse, error) {
	body, headers, err := encodeBody(obj, headers)
	if err != nil {
		return ServerResponse{}, err
	}
	return cli.sendRequest(ctx, "PUT", path, query, body, headers)
}

func (cli *Client) PutRaw(ctx context.Context, path string, query url.Values, body io.Reader, headers map[string][]string) (ServerResponse, error) {
	return cli.sendRequest(ctx, "PUT", path, query, body, headers)
}

func (cli *Client) Delete(ctx context.Context, path string, query url.Values, headers map[string][]string) (ServerResponse, error) {
	return cli.sendRequest(ctx, "DELETE", path, query, nil, headers)
}

func (cli *Client) sendRequest(ctx context.Context, method, path string, query url.Values, body io.Reader, headers headers) (ServerResponse, error) {
	req, err := cli.buildRequest(method, cli.getAPIPath(path, query), body, headers)
	if err != nil {
		return ServerResponse{}, err
	}

	resp, err := cli.doRequest(ctx, req)
	if err != nil {
		return resp, err
	}

	return resp, cli.checkResponseErr(resp)
}

func (cli *Client) checkResponseErr(serverResp ServerResponse) error {
	if serverResp.StatusCode >= 200 && serverResp.StatusCode < 400 {
		return nil
	}

	body, err := ioutil.ReadAll(serverResp.Body)
	if err != nil {
		return err
	}
	if len(body) == 0 {
		return fmt.Errorf("request returned %s for API route ", http.StatusText(serverResp.StatusCode))
	}

	var ct string
	if serverResp.Header != nil {
		ct = serverResp.Header.Get("Content-Type")
	}

	var errorMessage string
	if ct == "application/json" {
		var errorResponse types.ErrorResponse
		if err := json.Unmarshal(body, &errorResponse); err != nil {
			return fmt.Errorf("Error reading JSON: %v", err)
		}
		errorMessage = errorResponse.Message
	} else {
		errorMessage = string(body)
	}

	return fmt.Errorf("Error response from daemon: %s", strings.TrimSpace(errorMessage))
}

func (cli *Client) doRequest(ctx context.Context, req *http.Request) (ServerResponse, error) {
	serverResp := ServerResponse{StatusCode: -1, ReqURL: req.URL}

	resp, err := ctxhttp.Do(ctx, cli.client, req)
	if err != nil {
		if cli.scheme != "https" && strings.Contains(err.Error(), "malformed HTTP response") {
			return serverResp, fmt.Errorf("%v.\n* Are you trying to connect to a TLS-enabled daemon without TLS?", err)
		}

		if cli.scheme == "https" && strings.Contains(err.Error(), "bad certificate") {
			return serverResp, fmt.Errorf("The server probably has client authentication (--tlsverify) enabled. Please check your TLS client certification settings: %v", err)
		}

		switch err {
		case context.Canceled, context.DeadlineExceeded:
			return serverResp, err
		}

		return serverResp, errors.Wrap(err, "error during connect")
	}

	if resp != nil {
		serverResp.StatusCode = resp.StatusCode
		serverResp.Body = resp.Body
		serverResp.Header = resp.Header
	}
	return serverResp, nil
}

func (cli *Client) getAPIPath(path string, query url.Values) string {
	return (&url.URL{Scheme: cli.scheme, Host: cli.host, Path: path, RawQuery: query.Encode()}).String()
}

func (cli *Client) buildRequest(method, path string, body io.Reader, headers headers) (*http.Request, error) {
	expectedPayload := (method == "POST" || method == "PUT")
	if expectedPayload && body == nil {
		body = bytes.NewReader([]byte{})
	}

	req, err := http.NewRequest(method, path, body)
	if err != nil {
		return nil, err
	}

	req = cli.addHeaders(req, headers)

	req.URL.Host = cli.host
	req.URL.Scheme = cli.scheme

	if expectedPayload && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "text/plain")
	}

	if body != nil {
		switch v := body.(type) {
		case *bytes.Buffer:
			req.Header.Set("Content-Length", strconv.Itoa(v.Len()))
		case *bytes.Reader:
			req.Header.Set("Content-Length", strconv.Itoa(v.Len()))
		case *strings.Reader:
			req.Header.Set("Content-Length", strconv.Itoa(v.Len()))
		}
	}

	return req, nil

}

func encodeBody(obj interface{}, headers headers) (io.Reader, headers, error) {
	if obj == nil {
		return nil, headers, nil
	}

	body, err := encodeData(obj)
	if err != nil {
		return nil, headers, err
	}
	if headers == nil {
		headers = make(map[string][]string)
	}
	headers["Content-Type"] = []string{"application/json"}
	return body, headers, nil
}

func (cli *Client) addHeaders(req *http.Request, headers headers) *http.Request {
	for k, v := range cli.customHTTPHeaders {
		req.Header.Set(k, v)
	}

	if headers != nil {
		for k, v := range headers {
			req.Header[k] = v
		}
	}
	return req
}

func encodeData(data interface{}) (*bytes.Buffer, error) {
	params := bytes.NewBuffer(nil)
	if data != nil {
		if err := json.NewEncoder(params).Encode(data); err != nil {
			return nil, err
		}
	}
	return params, nil
}

func ensureReaderClosed(response ServerResponse) {
	if response.Body != nil {
		io.CopyN(ioutil.Discard, response.Body, 512)
		response.Body.Close()
	}
}
