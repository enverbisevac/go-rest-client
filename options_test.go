package rest

import (
	"net/http"
	"reflect"
	"testing"
)

type mockArticle struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

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
	type args struct {
		f MarshallFunc
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithMarshallFunc(tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithMarshallFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithUnmarshalFunc(t *testing.T) {
	type args struct {
		f UnmarshallFunc
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithUnmarshalFunc(tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithUnmarshalFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}
