/*
http test
通过使用假的的request/response
通过起服务器
*/
package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// panic handler
func errPanic(writer http.ResponseWriter, request *http.Request) error {
	panic(123)
}

// 定义一个user error 结构
type testingUserError struct {
	message string
}

func (t testingUserError) Error() string {
	return t.Message()
}

// 实现userError interface
func (t testingUserError) Message() string {
	return t.message
}

// 定义system error 404
func errNotFound(w http.ResponseWriter, r *http.Request) error {
	return os.ErrNotExist
}

// 定义system error 403
func errNoPermission(w http.ResponseWriter, r *http.Request) error {
	return os.ErrPermission
}

// 定义system error 500
func errUnKnown(w http.ResponseWriter, r *http.Request) error {
	return errors.New("unknown error")
}

// no error
func noError(writer http.ResponseWriter, request *http.Request) error {
	fmt.Fprintln(writer, "no error")
	return nil
}

// userError handler
func errUserError(writer http.ResponseWriter, request *http.Request) error {
	return testingUserError{"this is an user error"}
}

// 定义测试结构
var tests = []struct {
	h       appHandler // 一个包里所有类型 即使不在一个文件中也能直接用
	code    int
	message string
}{
	{errPanic, 500, "Internal Server Error"},
	{errUserError, 400, "this is an user error"},
	{errNotFound, 404, "Not Found"},
	{errNoPermission, 403, "Forbidden"},
	{errUnKnown, 500, "Internal Server Error"},
	{errUnKnown, 500, "Internal Server Error"},
	{noError, 200, "no error"},
}

// test wrapper func
func TestErrWrapper(t *testing.T) {
	// 遍历结构体
	for _, tt := range tests {
		// 包装器返回一个 httpHandler
		f := errWrapper(tt.h)
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "http://www.baidu.com", nil)
		// 执行httpHandler
		f(response, request)
		// 从NewRecorder response 获取http response 否则类型不对
		resp := response.Result()

		verifyResponse(resp, tt.code, tt.message, t)
	}

}

// test wrapper server
func TestErrWrapperInServer(t *testing.T) {
	for _, tt := range tests {
		f := errWrapper(tt.h)
		// new an http server
		server := httptest.NewServer(http.HandlerFunc(f))
		// send a get request to the server
		// 返回一个http response指针
		resp, _ := http.Get(server.URL)

		verifyResponse(resp, tt.code, tt.message, t)
	}
}

// verify 校验响应
func verifyResponse(resp *http.Response, expectedCode int, expectedMsg string, t *testing.T) {
	//从内存中获取http response Body 为[]byte类型
	b, _ := ioutil.ReadAll(resp.Body)
	// 获取http response body&code
	responseBody := strings.Trim(string(b), "\n")
	responseCode := resp.StatusCode

	if responseCode != expectedCode || responseBody != expectedMsg {
		t.Errorf("expected (%d,%s); got (%d,%s)", expectedCode, expectedMsg, responseCode, responseBody)
	}
}
