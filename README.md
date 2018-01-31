httpwrapper
----

[![GoDoc][1]][2] [![License: MIT][3]][4] [![Release][5]][6] [![Build Status][7]][8] [![Codecov Coverage][11]][12] [![Go Report Card][13]][14] [![Downloads][15]][16]

[1]: https://godoc.org/github.com/evalphobia/httpwrapper?status.svg
[2]: https://godoc.org/github.com/evalphobia/httpwrapper
[3]: https://img.shields.io/badge/License-MIT-blue.svg
[4]: LICENSE.md
[5]: https://img.shields.io/github/release/evalphobia/httpwrapper.svg
[6]: https://github.com/evalphobia/httpwrapper/releases/latest
[7]: https://travis-ci.org/evalphobia/httpwrapper.svg?branch=master
[8]: https://travis-ci.org/evalphobia/httpwrapper
[9]: https://coveralls.io/repos/evalphobia/httpwrapper/badge.svg?branch=master&service=github
[10]: https://coveralls.io/github/evalphobia/httpwrapper?branch=master
[11]: https://codecov.io/github/evalphobia/httpwrapper/coverage.svg?branch=master
[12]: https://codecov.io/github/evalphobia/httpwrapper?branch=master
[13]: https://goreportcard.com/badge/github.com/evalphobia/httpwrapper
[14]: https://goreportcard.com/report/github.com/evalphobia/httpwrapper
[15]: https://img.shields.io/github/downloads/evalphobia/httpwrapper/total.svg?maxAge=1800
[16]: https://github.com/evalphobia/httpwrapper/releases
[17]: https://img.shields.io/github/stars/evalphobia/httpwrapper.svg
[18]: https://github.com/evalphobia/httpwrapper/stargazers

HTTP client wrapper using [h2non/gentleman](https://github.com/h2non/gentleman)

# Quick Usage

```go
package main

import (
	"time"

	// automatically set debug flag (output HTTP request/response to stdout)
	//  _ "github.com/evalphobia/httpwrapper/debug"
	"github.com/evalphobia/httpwrapper/request"
)

func main() {
	// Example #1
	// GET request
	resp, err := request.GET("http://example.com", request.Option{
		Headers: map[string]string{"x-example-token": "abcdefg"},
		Timeout: 10 * time.Second,
		Retry:   true,
		// output HTTP request/response
		Debug: false,
		// for basic auth
		User:  "MyName",
		Pass:  "secret",
		Query: "name=MyName&pass=secret",
		// Query: map[string]string{"name": "MyName", "pass": "secret"},
	})
	if err != nil {
		panic(err)
	}
	// if status code is not 2xx.
	if err := resp.HasStatusCodeError(); err != nil {
		panic(err)
	}

	// for JSON response
	user := User{}
	err = resp.JSON(&user)
	if err != nil {
		panic(err)
	}

	// for XML response
	err = resp.XML(&user, nil)
	if err != nil {
		panic(err)
	}

	// Example #2
	// POST request and set response result.
	user = User{}
	err = request.CallWithResult(request.Option{
		URL:         "http://example.com",
		Method:      request.MethodPOST,
		PayloadType: request.PayloadTypeJSON,
		Payload:     map[string]interface{}{"param1": 100, "param2": "abc"},
	}, &user)
	if err != nil {
		panic(err)
	}

	// ...
}

// struct for HTTP response data
type User struct {
	LastName  string `json:"last_name" xml:"last_name"`
	FirstName string `json:"first_name" xml:"first_name"`
	Age       int    `json:"age" xml:"age"`
}
```

## Option

|Name|Description|Example|
|:--|:--|:--|
| URL | Request URL | `http://exmaple.com` |
| Method | HTTP method name | `GET`, `POST`, `PUT`, `DELETE`. request package has constants (`request.MethodGET` and others) |
| Headers | HTTP headers | `map[string]string{"X-HTTP-SOMETHING": "1"}` |
| Timeout | HTTP timeout| `30 * time.Second` |
| Retry | retry flag, if set `true` and timeout has occured, then it tries another 2 requests | |
| Debug | if set `true`, then it outputs HTTP request/response to stdout. | |
| Plugins | [gentleman's plugins](https://github.com/h2non/gentleman/tree/master/plugin) | |
| UserAgent | User Agent | `Mozilla/5.0 (iPhone; CPU iPhone OS 11_0 like Mac OS X) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Mobile/15A372 Safari/604.1` |
| User | Basic auth username | |
| Pass | Basic auth password | |
| Query | HTTP Query String in `string`, `map[string]string` or `map[string]interface{}` | `name=aaa&pass=bbb&param=ccc` |
| PayloadType | Define type of `Payload` | `JSON`, `XML`, `BODY`, `FORM`. request package has variables (`request.PayloadTypeJSON` and others) |
| Payload | Request body parameters | |
