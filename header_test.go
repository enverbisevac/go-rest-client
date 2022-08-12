package rest

import (
	"net/http"
	"reflect"
	"testing"
)

func TestWithHeader(t *testing.T) {
	base := http.Header{}
	input := http.Header{
		Content: []string{string(ApplicationJSON)},
	}

	got := WithHeader(input)
	got(base)

	equals(t, input, base)
}

func TestWithAuth(t *testing.T) {
	token := "some_token"
	base := http.Header{}
	exp := http.Header{
		Authorization: []string{string(Bearer) + " " + token},
	}

	got := WithAuth(Bearer, token)
	got(base)

	equals(t, exp, base)
}

func TestWithContent(t *testing.T) {
	base := http.Header{}
	exp := http.Header{
		Content: []string{string(ApplicationJSON)},
	}

	got := WithContent(ApplicationJSON)
	got(base)

	equals(t, exp, base)
}

func TestWith(t *testing.T) {
	base := http.Header{}
	exp := http.Header{
		"Custom": []string{"Value"},
	}

	got := With("Custom", "Value")
	got(base)

	equals(t, exp, base)
}

func TestHeader(t *testing.T) {
	type args struct {
		opts []Func
	}
	tests := []struct {
		name string
		args args
		want http.Header
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Header(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Header() = %v, want %v", got, tt.want)
			}
		})
	}
}
