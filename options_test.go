package rest

import (
	"net/http"
	"reflect"
	"runtime"
	"testing"
)

func TestWithBody(t *testing.T) {
	s := struct {
		Title string
	}{
		Title: "Some title",
	}
	base := RequestOption{}

	got := WithBody(s)
	got(&base)

	equals(t, s, base.Body)
}

func TestWithHeaders(t *testing.T) {
	base := RequestOption{}

	input := http.Header{
		Content: []string{string(ApplicationJSON)},
	}

	got := WithHeaders(input)
	got(&base)

	equals(t, input, base.Header)
}

func TestWithMarshallFunc(t *testing.T) {
	base := RequestOption{}

	input := func(v any) ([]byte, error) {
		return nil, nil
	}

	got := WithMarshallFunc(input)
	got(&base)

	funcName1 := runtime.FuncForPC(reflect.ValueOf(input).Pointer()).Name()
	funcName2 := runtime.FuncForPC(reflect.ValueOf(base.MarshalFunc).Pointer()).Name()

	equals(t, funcName1, funcName2)
}

func TestWithUnmarshalFunc(t *testing.T) {
	base := RequestOption{}

	input := func(data []byte, v any) error {
		return nil
	}

	got := WithUnmarshalFunc(input)
	got(&base)

	funcName1 := runtime.FuncForPC(reflect.ValueOf(input).Pointer()).Name()
	funcName2 := runtime.FuncForPC(reflect.ValueOf(base.UnmarshallFunc).Pointer()).Name()

	equals(t, funcName1, funcName2)
}
