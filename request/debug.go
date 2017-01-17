package request

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http/httputil"

	"gopkg.in/h2non/gentleman.v1"
	"gopkg.in/h2non/gentleman.v1/context"
	"gopkg.in/h2non/gentleman.v1/plugin"
)

// debug is flag for debugging HTTP request/response
var debug bool

// DebugOn activates request debugging.
func DebugOn() {
	debug = true
}

// DebugOff deactivates request debugging.
func DebugOff() {
	debug = false
}

// debugRequest returns plugin to show HTTP request details.
func debugRequest() plugin.Plugin {
	p := plugin.New()
	p.SetHandler("before dial", func(ctx *context.Context, h context.Handler) {
		req := ctx.Request
		body, _ := ioutil.ReadAll(req.Body)
		req.Body = ioutil.NopCloser(bytes.NewReader(body))
		dump, _ := httputil.DumpRequest(req, false)

		fmt.Printf("---> [HTTP Request] %s[Request Body]\n%s\n", string(dump), body)
		h.Next(ctx)
	})
	return p
}

// showDebugResponse shows HTTP response details.
func showDebugResponse(resp *gentleman.Response, err error) {
	if err, ok := err.(net.Error); ok && err.Timeout() {
		fmt.Printf("<--- [HTTP Response ] timeout\n\n")
		return
	}
	if resp == nil {
		return
	}

	res := resp.RawResponse
	body, _ := ioutil.ReadAll(res.Body)
	res.Body = ioutil.NopCloser(bytes.NewReader(body))
	dump, _ := httputil.DumpResponse(res, false)

	fmt.Printf("<--- [HTTP Response] %s[Response Body]\n%s\n", string(dump), body)
}
