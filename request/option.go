package request

import (
	"fmt"
	"strings"
	"time"

	"gopkg.in/h2non/gentleman.v1/plugin"
)

// Option is wrapper struct of http option
type Option struct {
	URL       string
	Method    Method
	Headers   map[string]string
	Timeout   time.Duration
	Retry     bool
	Debug     bool
	Plugins   []plugin.Plugin
	UserAgent string

	// Basic Auth
	User string
	Pass string

	// Query Parameter
	Query interface{}

	// Body Parameter
	Payload interface{}
	// PayloadType is used for body payload type
	PayloadType PayloadType
}

func (o Option) hasHeaders() bool {
	return len(o.Headers) != 0
}

func (o Option) hasTimeout() bool {
	return o.Timeout > 0
}

func (o Option) hasUserAgent() bool {
	return o.UserAgent != ""
}

func (o Option) hasBasicAuth() bool {
	return o.User != ""
}

func (o Option) hasPayload() bool {
	return o.Payload != nil
}

func (o Option) hasQuery() bool {
	return o.Query != nil
}

func (o Option) queryToMap() map[string]string {
	switch v := o.Query.(type) {
	case map[string]string:
		return v
	case map[string]interface{}:
		m := make(map[string]string)
		for key, val := range v {
			m[key] = fmt.Sprint(val)
		}
		return m
	case string:
		m := make(map[string]string)
		for _, kv := range strings.Split(v, "&") {
			values := strings.Split(kv, "=")
			if len(values) == 2 {
				m[values[0]] = values[1]
			}
		}
		return m
	default:
		return nil
	}
}

// Method is HTTP Method.
type Method string

// HTTP method list.
const (
	MethodGET    Method = "GET"
	MethodPOST   Method = "POST"
	MethodPUT    Method = "PUT"
	MethodDELETE Method = "DELETE"
)

func (m Method) String() string {
	return string(m)
}

func (m Method) isEmpty() bool {
	return string(m) == ""
}

func (m Method) isGET() bool {
	return string(m) == "GET"
}

func (m Method) isPOST() bool {
	return string(m) == "POST"
}

func (m Method) isPUT() bool {
	return string(m) == "PUT"
}

func (m Method) isDELETE() bool {
	return string(m) == "DELETE"
}

// PayloadType is payload type for POST
type PayloadType string

// POST Payload type variables
var (
	PayloadTypeBODY PayloadType = "BODY"
	PayloadTypeJSON PayloadType = "JSON"
	PayloadTypeXML  PayloadType = "XML"
	PayloadTypeFORM PayloadType = "FORM"
)

func (p PayloadType) isBody() bool {
	return p == PayloadTypeBODY || string(p) == ""
}

func (p PayloadType) isJSON() bool {
	return p == PayloadTypeJSON
}

func (p PayloadType) isXML() bool {
	return p == PayloadTypeXML
}

func (p PayloadType) isForm() bool {
	return p == PayloadTypeFORM
}
