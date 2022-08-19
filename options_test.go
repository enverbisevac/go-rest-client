package rest

import (
	"net/http"
	"testing"
)

func TestWithBody(t *testing.T) {
	s := struct {
		Title string
	}{
		Title: "Some title",
	}
	base := requestData{}

	got := WithBody(s)
	got(&base)

	equals(t, s, base.Body)
}

func TestWithHeaders(t *testing.T) {
	base := requestData{}

	input := http.Header{
		Content: []string{string(ApplicationJSON)},
	}

	got := WithHeaders(input)
	got(&base)

	equals(t, input, base.Header)
}
