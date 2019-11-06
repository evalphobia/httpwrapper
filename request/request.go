package request

import (
	"fmt"
	"io"

	retry "gopkg.in/h2non/gentleman-retry.v2"
	"gopkg.in/h2non/gentleman.v2"
	"gopkg.in/h2non/gentleman.v2/plugins/auth"
	"gopkg.in/h2non/gentleman.v2/plugins/body"
	"gopkg.in/h2non/gentleman.v2/plugins/bodytype"
	"gopkg.in/h2non/gentleman.v2/plugins/headers"
	"gopkg.in/h2non/gentleman.v2/plugins/multipart"
	"gopkg.in/h2non/gentleman.v2/plugins/query"
	"gopkg.in/h2non/gentleman.v2/plugins/timeout"
)

// GET sends GET request with option.
func GET(url string, opt Option) (*Response, error) {
	opt.URL = url
	opt.Method = MethodGET
	return Call(opt)
}

// POST sends POST request with option.
func POST(url string, opt Option) (*Response, error) {
	opt.URL = url
	opt.Method = MethodPOST
	return Call(opt)
}

// PUT sends PUT request with option.
func PUT(url string, opt Option) (*Response, error) {
	opt.URL = url
	opt.Method = MethodPUT
	return Call(opt)
}

// DELETE sends DELETE request with option.
func DELETE(url string, opt Option) (*Response, error) {
	opt.URL = url
	opt.Method = MethodDELETE
	return Call(opt)
}

// Call sneds HTTP request by given option.
func Call(opt Option) (*Response, error) {
	cli := gentleman.New()
	cli.URL(opt.URL)

	req := cli.Request()
	if !opt.Method.isEmpty() {
		req.Method(opt.Method.String())
	}

	// Set plugins
	for _, p := range opt.Plugins {
		req.Use(p)
	}

	// Set User-Agent
	if opt.hasUserAgent() {
		req.Use(headers.Set("User-Agent", opt.UserAgent))
	}

	// Set basic auth
	if opt.hasBasicAuth() {
		req.Use(auth.Basic(opt.User, opt.Pass))
	}
	// Set custom headers
	if opt.hasHeaders() {
		req.Use(headers.SetMap(opt.Headers))
	}
	// Set Auth Bearer header
	if opt.hasAuthBearer() {
		req.Use(auth.Bearer(opt.Bearer))
	}
	// Set timeout
	if opt.hasTimeout() {
		req.Use(timeout.Request(opt.Timeout))
	}
	// Set retry (3 times)
	if opt.Retry {
		req.Use(retry.New(retry.ConstantBackoff))
	}

	// Set Query String
	if opt.hasQuery() {
		req.Use(query.SetMap(opt.queryToMap()))
	}

	// Set parameter
	if opt.hasPayload() {
		payload := opt.Payload
		switch {
		case opt.PayloadType.isJSON():
			req.Use(body.JSON(payload))
		case opt.PayloadType.isXML():
			req.Use(body.XML(payload))
		case opt.PayloadType.isForm():
			req.Use(headers.Set("Content-Type", "application/x-www-form-urlencoded"))
			req.Use(body.String(parseParam(payload)))
		case opt.PayloadType.isData():
			req.Use(multipart.Fields(convertToMultipartData(tagName, payload)))
		case opt.PayloadType.isStream():
			if v, ok := payload.(io.Reader); ok {
				req.Use(body.Reader(v))
			}
		default:
			req.Use(body.String(fmt.Sprint(payload)))
		}
	}

	// Set Custom content-type
	if opt.hasContentType() {
		req.Use(bodytype.Set(opt.ContentType))
	}

	// show debug request
	if opt.Debug || debug {
		req.Use(debugRequest())
	}

	resp, err := req.Send()
	if opt.Debug || debug {
		showDebugResponse(resp, err)
	}
	if err != nil {
		return nil, err
	}
	return &Response{resp}, nil
}

// CallWithResult sends request and set result.
func CallWithResult(opt Option, result interface{}) error {
	resp, respErr := Call(opt)

	var unmarshalErr error
	switch {
	case opt.PayloadType.isXML():
		unmarshalErr = resp.XML(result, nil)
	default:
		unmarshalErr = resp.JSON(result)
	}

	if respErr != nil {
		return respErr
	}
	if err := resp.HasStatusCodeError(); err != nil {
		return err
	}
	return unmarshalErr
}
