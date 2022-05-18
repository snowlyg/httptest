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

type paramFunc func(req *httpexpect.Request) *httpexpect.Request

//NewWithJsonParamFunc return req.WithJSON
func NewWithJsonParamFunc(query map[string]interface{}) paramFunc {
	return func(req *httpexpect.Request) *httpexpect.Request {
		return req.WithJSON(query)
	}
}

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

func Instance(t *testing.T, url string, handler http.Handler) *Client {
	httpTestClient = &Client{
		t: t,
		expect: httpexpect.WithConfig(httpexpect.Config{
			BaseURL: url,
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
		}),
	}
	return httpTestClient
}

func (c *Client) Login(url string, res Responses, paramFuncs ...paramFunc) (uint, error) {
	var id uint
	if len(paramFuncs) == 0 {
		paramFuncs = append(paramFuncs, LoginFunc)
	}
	if res == nil {
		res = LoginResponse
	}
	c.POST(url, res, paramFuncs...)
	token := res.GetString("data.AccessToken")
	fmt.Printf("access_token is '%s'\n", token)
	if token == "" {
		return id, fmt.Errorf("access_token is empty")
	}
	id = res.GetUint("data.user.id")
	c.expect = c.expect.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", str.Join("Bearer ", token))
	})
	return id, nil
}

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

// POST
func (c *Client) POST(url string, res Responses, paramFuncs ...paramFunc) {
	req := c.expect.POST(url)
	if len(paramFuncs) > 0 {
		for _, f := range paramFuncs {
			req = f(req)
		}
	}
	obj := req.Expect().Status(http.StatusOK).JSON().Object()
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
	obj := req.Expect().Status(http.StatusOK).JSON().Object()
	res.Test(obj)
}

// UPLOAD 上传文件
func (c *Client) UPLOAD(url string, res Responses, paramFuncs ...paramFunc) {
	req := c.expect.POST(url)
	if len(paramFuncs) > 0 {
		for _, f := range paramFuncs {
			req = f(req)
		}
	}
	obj := req.Expect().Status(http.StatusOK).JSON().Object()
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
	obj := req.Expect().Status(http.StatusOK).JSON().Object()
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
	return req.Expect().Status(http.StatusOK).ContentType("application/octet-stream").Body().NotEmpty().Raw()
}

// DELETE
func (c *Client) DELETE(url string, res Responses, paramFuncs ...paramFunc) {
	req := c.expect.DELETE(url)
	if len(paramFuncs) > 0 {
		for _, f := range paramFuncs {
			req = f(req)
		}
	}
	obj := req.Expect().Status(http.StatusOK).JSON().Object()
	res.Test(obj)
}
