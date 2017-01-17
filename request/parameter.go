package request

import (
	"fmt"
	"net/url"
)

const tagName = "url"

// parseParam creates string of form data from interface value.
func parseParam(param interface{}) string {
	switch v := param.(type) {
	case string:
		return v
	case fmt.Stringer:
		return v.String()
	case url.Values:
		return v.Encode()
	case *url.Values:
		return v.Encode()
	default:
		value := convertToURLValue(tagName, param)
		return value.Encode()
	}
}

// convertToURLValue convert url.Values data from struct data.
func convertToURLValue(tagName string, param interface{}) url.Values {
	result := url.Values{}
	convertTo(tagName, param, &result, func(result interface{}, key string, value interface{}) {
		r := result.(*url.Values)
		r.Add(key, fmt.Sprint(value))
	})
	return result
}
