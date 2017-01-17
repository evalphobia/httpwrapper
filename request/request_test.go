package request

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGET(t *testing.T) {
	assert := assert.New(t)
	ts, record := createTestHandler()
	defer ts.Close()

	tests := []struct {
		query    interface{}
		expected string
	}{
		{map[string]string{"key": "value", "user_id": "100"}, "key=value&user_id=100"},
		{map[string]interface{}{"key": "value", "user_id": "100"}, "key=value&user_id=100"},
		{"key=value&user_id=100", "key=value&user_id=100"},
		{nil, ""},
		{"", ""},
	}

	for _, tt := range tests {
		target := fmt.Sprintf("%+v", tt)

		opt := Option{
			Query: tt.query,
		}

		resp, err := GET(ts.URL, opt)
		assert.NoError(err, target)
		assert.NotNil(resp, target)
		assert.True(resp.Ok, target)
		assert.Equal("GET", record["method"], target)
		assert.Empty(record["body"], target)

		assert.Equal(tt.expected, record["query"], target)
	}
}

func TestPOST(t *testing.T) {
	assert := assert.New(t)
	ts, record := createTestHandler()
	defer ts.Close()

	tests := []struct {
		typ      PayloadType
		payload  interface{}
		expected string
	}{
		{PayloadTypeJSON, `{"key": "value", "user_id": "100"}`, `{"key": "value", "user_id": "100"}`},
		{PayloadTypeJSON, map[string]interface{}{"key": "value", "user_id": "100"}, "{\"key\":\"value\",\"user_id\":\"100\"}\n"},
		{PayloadTypeBODY, `{"key": "value", "user_id": "100"}`, `{"key": "value", "user_id": "100"}`},
		{PayloadTypeBODY, nil, ``},
	}

	for _, tt := range tests {
		target := fmt.Sprintf("%+v", tt)

		opt := Option{
			Payload:     tt.payload,
			PayloadType: tt.typ,
		}

		resp, err := POST(ts.URL, opt)
		assert.NoError(err, target)
		assert.NotNil(resp, target)
		assert.True(resp.Ok, target)
		assert.Equal("POST", record["method"], target)
		assert.Empty(record["query"], target)

		assert.Equal(tt.expected, record["body"], target)
	}
}

func TestPUT(t *testing.T) {
	assert := assert.New(t)
	ts, record := createTestHandler()
	defer ts.Close()

	tests := []struct {
		typ      PayloadType
		payload  interface{}
		expected string
	}{
		{PayloadTypeJSON, `{"key": "value", "user_id": "100"}`, `{"key": "value", "user_id": "100"}`},
		{PayloadTypeJSON, map[string]interface{}{"key": "value", "user_id": "100"}, "{\"key\":\"value\",\"user_id\":\"100\"}\n"},
		{PayloadTypeBODY, `{"key": "value", "user_id": "100"}`, `{"key": "value", "user_id": "100"}`},
		{PayloadTypeBODY, nil, ``},
	}

	for _, tt := range tests {
		target := fmt.Sprintf("%+v", tt)

		opt := Option{
			Payload:     tt.payload,
			PayloadType: tt.typ,
		}

		resp, err := PUT(ts.URL, opt)
		assert.NoError(err, target)
		assert.NotNil(resp, target)
		assert.True(resp.Ok, target)
		assert.Equal("PUT", record["method"], target)
		assert.Empty(record["query"], target)

		assert.Equal(tt.expected, record["body"], target)
	}
}

func TestDELETE(t *testing.T) {
	assert := assert.New(t)
	ts, record := createTestHandler()
	defer ts.Close()

	tests := []struct {
		query    interface{}
		expected string
	}{
		{map[string]string{"key": "value", "user_id": "100"}, "key=value&user_id=100"},
		{map[string]interface{}{"key": "value", "user_id": "100"}, "key=value&user_id=100"},
		{"key=value&user_id=100", "key=value&user_id=100"},
		{nil, ""},
		{"", ""},
	}

	for _, tt := range tests {
		target := fmt.Sprintf("%+v", tt)

		opt := Option{
			Query: tt.query,
		}

		resp, err := DELETE(ts.URL, opt)
		assert.NoError(err, target)
		assert.NotNil(resp, target)
		assert.True(resp.Ok, target)
		assert.Equal("DELETE", record["method"], target)
		assert.Empty(record["body"], target)

		assert.Equal(tt.expected, record["query"], target)
	}
}

func TestCall(t *testing.T) {
	assert := assert.New(t)
	ts, record := createTestHandler()
	defer ts.Close()

	tests := []struct {
		method Method
	}{
		{MethodGET},
		{MethodPOST},
		{MethodPUT},
		{MethodDELETE},
	}

	for _, tt := range tests {
		target := fmt.Sprintf("%+v", tt)

		opt := Option{
			URL:    ts.URL,
			Method: tt.method,
		}

		resp, err := Call(opt)
		assert.NoError(err, target)
		assert.NotNil(resp, target)
		assert.True(resp.Ok, target)

		assert.Equal(tt.method.String(), record["method"], target)
		assert.Empty(record["query"], target)
		assert.Empty(record["body"], target)
	}
}

func createTestHandler() (*httptest.Server, map[string]string) {
	record := make(map[string]string)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)

		record["body"] = string(body)
		record["query"] = r.URL.Query().Encode()
		record["method"] = r.Method
		record["content_type"] = r.Header.Get("Content-Type")

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintln(w, `{"error": true}`)
		} else {
			w.WriteHeader(200)
			fmt.Fprintln(w, `{"error": false}`)
		}
		return
	}))
	return ts, record
}
