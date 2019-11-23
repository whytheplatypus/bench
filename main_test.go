package main

import (
	"encoding/json"
	"net/http"
	"testing"
)

var jsonReq = `{
  "Method": "GET",
  "URL": {
    "Scheme": "https",
    "Host": "www.example.com",
    "Path": "/test",
    "RawQuery": "param=test"
  },
  "Header": {
    "Accept-Encoding": ["gzip", "deflate"]
  }
}`

func TestUnmarshal(t *testing.T) {
	req := &http.Request{}
	err := json.Unmarshal([]byte(jsonReq), req)
	if err != nil {
		t.Fatal(err)
	}
	if req.Method != "GET" {
		t.Error("Wrong method", req.Method)
	}
	if req.Header["Accept-Encoding"][0] != "gzip" {
		t.Error("Wrong headers", req.Header)
	}
	if req.URL.String() != "https://www.example.com/test?param=test" {
		t.Error("Wrong url", req.URL)
	}
}
