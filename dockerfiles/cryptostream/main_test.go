package main

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEncrypt(t *testing.T) {
	res := httptest.NewRecorder()
	Encrypt(res, queryGET())
	assert.Equal(t, "test", readBody(res))
}

func queryGETWithParameter(params string) *http.Request {
	req, _ := http.NewRequest("GET", "URL"+params, nil)
	return req
}

func queryGET() *http.Request {
	req, _ := http.NewRequest("GET", "URL", nil)
	return req
}

func readBody(res *httptest.ResponseRecorder) string {
	content, _ := ioutil.ReadAll(res.Body)
	return string(content)
}