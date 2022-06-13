package httptest

import (
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/snowlyg/helper/str"
)

var (
	httpTestClient *Client
)

// paramFunc
type paramFunc func(req *httpexpect.Request) *httpexpect.Request

//NewWithJsonParamFunc return req.WithJSON
func NewWithJsonParamFunc(query map[string]interface{}) paramFunc {
	return func(req *httpexpect.Request) *httpexpect.Request {
		return req.WithJSON(query)
	}
}

// NewWithQueryObjectParamFunc query for get method
func NewWithQueryObjectParamFunc(query map[string]interface{}) paramFunc {
	return func(req *httpexpect.Request) *httpexpect.Request {
		return req.WithQueryObject(query)
	}
}

//NewWithFileParamFunc return req.WithFile
func NewWithFileParamFunc(fs []File) paramFunc {
	return func(req *httpexpect.Request) *httpexpect.Request {
		if len(fs) == 0 {
			return req
		}
		req = req.WithMultipart()
		for _, f := range fs {
			req = req.WithFile(f.Key, f.Path, f.Reader)
		}
		return req
	}
}

//NewResponsesWithLength return Responses with length value for data key
func NewResponsesWithLength(status int, message string, data []Responses, length int) Responses {
	return Responses{
		{Key: "status", Value: status},
		{Key: "message", Value: message},
		{Key: "data", Value: data, Length: length},
	}
}

//NewResponsesWithHttpStatus return Responses with http response status
func NewResponsesWithHttpStatus(status int, message string, data []Responses, httpStatus int) Responses {
	return Responses{
		{Key: "http_status", Value: httpStatus},
		{Key: "status", Value: status},
		{Key: "message", Value: message},
		{Key: "data", Value: data},
	}
}

//NewResponses return Responses
func NewResponses(status int, message string, data ...Responses) Responses {
	if status != http.StatusOK {
		return Responses{
			{Key: "status", Value: status},
			{Key: "message", Value: message},
		}
	}
	if len(data) == 0 {
		return Responses{
			{Key: "status", Value: status},
			{Key: "message", Value: message},
		}
	}
	if len(data) == 1 {
		return Responses{
			{Key: "status", Value: status},
			{Key: "message", Value: message},
			{Key: "data", Value: data[0]},
		}
	}
	return Responses{
		{Key: "status", Value: status},
		{Key: "message", Value: message},
		{Key: "data", Value: data},
	}
}

type Client struct {
	t      *testing.T
	expect *httpexpect.Expect
}

// Instance return test client instance
func Instance(t *testing.T, handler http.Handler, url ...string) *Client {
	config := httpexpect.Config{
		Client: &http.Client{
			Transport: httpexpect.NewBinder(handler),
			Jar:       httpexpect.NewJar(),
		},
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
			httpexpect.NewCurlPrinter(t),
			httpexpect.NewCompactPrinter(t),
		},
	}
	if len(url) == 1 && url[0] != "" {
		config.BaseURL = url[0]
	}
	httpTestClient = &Client{
		t:      t,
		expect: httpexpect.WithConfig(config),
	}
	return httpTestClient
}

// Login for http login
func (c *Client) Login(url, tokenIndex string, res Responses, paramFuncs ...paramFunc) error {
	if len(paramFuncs) == 0 {
		paramFuncs = append(paramFuncs, LoginFunc)
	}
	c.POST(url, res, paramFuncs...)
	token := res.GetString("data.accessToken")
	if tokenIndex != "" {
		token = res.GetString(tokenIndex)
	}
	fmt.Printf("access_token is '%s'\n", token)
	if token == "" {
		return fmt.Errorf("access_token is empty")
	}
	c.expect = c.expect.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", str.Join("Bearer ", token))
	})
	return nil
}

// Logout for http logout
func (c *Client) Logout(url string, res Responses) {
	if res == nil {
		res = LogoutResponse
	}
	c.GET(url, res)
}

type File struct {
	Key    string
	Path   string
	Reader io.Reader
}

// checkStatus check what's http response stauts want
func checkStatus(res Responses) int {
	if len(res) == 0 {
		return http.StatusOK
	}
	if res[0].Key != "http_status" {
		return http.StatusOK
	}

	return res[0].Value.(int)
}

// POST
func (c *Client) POST(url string, res Responses, paramFuncs ...paramFunc) {
	req := c.expect.POST(url)
	if len(paramFuncs) > 0 {
		for _, f := range paramFuncs {
			req = f(req)
		}
	}
	obj := req.Expect().Status(checkStatus(res)).JSON().Object()
	res.Test(obj)
}

// PUT
func (c *Client) PUT(url string, res Responses, paramFuncs ...paramFunc) {
	req := c.expect.PUT(url)
	if len(paramFuncs) > 0 {
		for _, f := range paramFuncs {
			req = f(req)
		}
	}
	obj := req.Expect().Status(checkStatus(res)).JSON().Object()
	res.Test(obj)
}

// UPLOAD
func (c *Client) UPLOAD(url string, res Responses, paramFuncs ...paramFunc) {
	req := c.expect.POST(url)
	if len(paramFuncs) > 0 {
		for _, f := range paramFuncs {
			req = f(req)
		}
	}
	obj := req.Expect().Status(checkStatus(res)).JSON().Object()
	res.Test(obj)
}

// GET
func (c *Client) GET(url string, res Responses, paramFuncs ...paramFunc) {
	req := c.expect.GET(url)
	if len(paramFuncs) > 0 {
		for _, f := range paramFuncs {
			req = f(req)
		}
	}
	obj := req.Expect().Status(checkStatus(res)).JSON().Object()
	res.Test(obj)
}

// DOWNLOAD
func (c *Client) DOWNLOAD(url string, res Responses, paramFuncs ...paramFunc) string {
	req := c.expect.GET(url)
	if len(paramFuncs) > 0 {
		for _, f := range paramFuncs {
			req = f(req)
		}
	}
	return req.Expect().Status(checkStatus(res)).ContentType("application/octet-stream").Body().NotEmpty().Raw()
}

// DELETE
func (c *Client) DELETE(url string, res Responses, paramFuncs ...paramFunc) {
	req := c.expect.DELETE(url)
	if len(paramFuncs) > 0 {
		for _, f := range paramFuncs {
			req = f(req)
		}
	}
	obj := req.Expect().Status(checkStatus(res)).JSON().Object()
	res.Test(obj)
}
