package rest

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"reflect"
	"testing"
)

type MockHttp struct {
	response *http.Response
	err      error
}

func (m *MockHttp) Do(request *http.Request) (*http.Response, error) {
	return m.response, m.err
}

type MockReadCloserError struct {
}

func (r *MockReadCloserError) Read(p []byte) (n int, err error) {
	return 0, errors.New("read error")
}

func (r *MockReadCloserError) Close() error {
	return nil
}

func TestHttp_Request(t *testing.T) {
	type fields struct {
		Client Doer
	}
	type args struct {
		ctx        context.Context
		method     string
		requestURL string
		body       []byte
		header     http.Header
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Response
		wantErr bool
	}{
		{
			name: "context is nil expect error",
			fields: fields{
				Client: &MockHttp{},
			},
			args: args{
				ctx:        nil,
				method:     "",
				requestURL: "",
				body:       nil,
				header:     nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "net call return error",
			fields: fields{
				Client: &MockHttp{
					err: errors.New("network error"),
				},
			},
			args: args{
				ctx:        context.Background(),
				method:     "",
				requestURL: "",
				body:       nil,
				header:     nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "empty method should work as GET",
			fields: fields{
				Client: &MockHttp{
					response: &http.Response{
						StatusCode:    200,
						Body:          io.NopCloser(bytes.NewBufferString("")),
						ContentLength: 0,
						Header:        Header(WithContent(ApplicationJSON)),
					},
				},
			},
			args: args{
				ctx:        context.Background(),
				method:     "",
				requestURL: "",
				body:       nil,
				header:     nil,
			},
			want: &Response{
				Body:       []byte(""),
				StatusCode: 200,
				Header: http.Header{
					Content: []string{string(ApplicationJSON)},
				},
			},
			wantErr: false,
		},
		{
			name: "response Body error",
			fields: fields{
				Client: &MockHttp{
					response: &http.Response{
						Body: &MockReadCloserError{},
					},
				},
			},
			args: args{
				ctx:        context.Background(),
				method:     "",
				requestURL: "",
				body:       nil,
				header:     nil,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Http{
				Client: tt.fields.Client,
			}
			got, err := h.Request(tt.args.ctx, tt.args.method, tt.args.requestURL, tt.args.body, tt.args.header)
			if (err != nil) != tt.wantErr {
				t.Errorf("Request() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Request() got = %v, want %v", got, tt.want)
			}
		})
	}
}
