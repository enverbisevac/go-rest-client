package rest

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"testing"
)

type mockEncoder struct {
	EncodeFunc func(value any, contentType ContentType, marshalFunc MarshallFunc) ([]byte, error)
}

func (m *mockEncoder) Encode(value any, contentType ContentType, marshalFunc MarshallFunc) ([]byte, error) {
	return m.EncodeFunc(value, contentType, marshalFunc)
}

type mockDecoder[T any] struct {
	DecodeFunc func(data []byte, val *T, contentType ContentType, unmarshallFunc UnmarshallFunc) error
}

func (m *mockDecoder[T]) Decode(data []byte, val *T, contentType ContentType, unmarshallFunc UnmarshallFunc) error {
	return m.DecodeFunc(data, val, contentType, unmarshallFunc)
}

type mockRequest struct {
	response *Response
	err      error
}

func (m *mockRequest) Request(ctx context.Context, method string, requestURL string, body []byte,
	header http.Header) (*Response, error) {

	return m.response, m.err
}

func Test_request(t *testing.T) {
	article := mockArticle{
		Title:   "some title",
		Content: "some content",
	}

	type args struct {
		ctx     context.Context
		cfg     *config[mockArticle]
		options []Option
	}
	tests := []struct {
		name    string
		args    args
		want    mockArticle
		wantErr bool
	}{
		{
			name: "options test",
			args: args{
				ctx: context.Background(),
				cfg: &config[mockArticle]{
					method:     http.MethodGet,
					requestURL: "",
					requester: &mockRequest{
						response: &Response{
							Body:       []byte("{}"),
							StatusCode: 200,
							Header:     Header(WithContent(ApplicationJSON)),
						},
						err: nil,
					},
					encoder: &mockEncoder{
						func(value any, contentType ContentType, marshalFunc MarshallFunc) ([]byte, error) {
							return []byte("{}"), nil
						},
					},
					decoder: &mockDecoder[mockArticle]{
						DecodeFunc: func(data []byte, val *mockArticle, contentType ContentType, unmarshallFunc UnmarshallFunc) error {

							return nil
						},
					},
				},
				options: []Option{WithHeaders(Header(WithContent(ApplicationJSON)))},
			},
			want:    mockArticle{},
			wantErr: false,
		},
		{
			name: "encode error",
			args: args{
				ctx: context.Background(),
				cfg: &config[mockArticle]{
					method:     http.MethodGet,
					requestURL: "",
					requester:  nil,
					encoder: &mockEncoder{
						EncodeFunc: func(value any, contentType ContentType, marshalFunc MarshallFunc) ([]byte, error) {
							return nil, errors.New("encoding error")
						},
					},
					decoder: nil,
				},
				options: nil,
			},
			want:    mockArticle{},
			wantErr: true,
		},
		{
			name: "request return error",
			args: args{
				ctx: context.Background(),
				cfg: &config[mockArticle]{
					method:     http.MethodGet,
					requestURL: "",
					requester: &mockRequest{
						err: errors.New("request failed"),
					},
					encoder: &mockEncoder{
						EncodeFunc: func(value any, contentType ContentType, marshalFunc MarshallFunc) ([]byte, error) {
							return nil, nil
						},
					},
					decoder: &mockDecoder[mockArticle]{},
				},
				options: nil,
			},
			want:    mockArticle{},
			wantErr: true,
		},
		{
			name: "bad request",
			args: args{
				ctx: context.Background(),
				cfg: &config[mockArticle]{
					method:     http.MethodGet,
					requestURL: "",
					requester: &mockRequest{
						response: &Response{
							StatusCode: http.StatusBadRequest,
							Header:     Header(WithContent(ApplicationJSON)),
						},
						err: nil,
					},
					encoder: &mockEncoder{
						EncodeFunc: func(value any, contentType ContentType, marshalFunc MarshallFunc) ([]byte, error) {
							return nil, nil
						},
					},
					decoder: &mockDecoder[mockArticle]{},
				},
				options: []Option{WithHeaders(Header(WithContent(ApplicationJSON)))},
			},
			want:    mockArticle{},
			wantErr: true,
		},
		{
			name: "decode error",
			args: args{
				ctx: context.Background(),
				cfg: &config[mockArticle]{
					method:     http.MethodGet,
					requestURL: "",
					requester: &mockRequest{
						response: &Response{
							Body:       []byte("{}"),
							StatusCode: http.StatusOK,
							Header:     Header(WithContent(ApplicationJSON)),
						},
					},
					encoder: &mockEncoder{
						EncodeFunc: func(value any, contentType ContentType, marshalFunc MarshallFunc) ([]byte, error) {
							return nil, nil
						},
					},
					decoder: &mockDecoder[mockArticle]{
						DecodeFunc: func(data []byte, val *mockArticle, contentType ContentType, unmarshallFunc UnmarshallFunc) error {
							return errors.New("decoding error")
						},
					},
				},
				options: nil,
			},
			want:    mockArticle{},
			wantErr: true,
		},
		{
			name: "happy path",
			args: args{
				ctx: context.Background(),
				cfg: &config[mockArticle]{
					method:     http.MethodGet,
					requestURL: "",
					requester: &mockRequest{
						response: &Response{
							Body:       []byte("{\"title\":\"some title\", \"content\":\"some content\"}"),
							StatusCode: http.StatusOK,
							Header:     Header(WithContent(ApplicationJSON)),
						},
						err: nil,
					},
					encoder: &mockEncoder{
						EncodeFunc: func(value any, contentType ContentType, marshalFunc MarshallFunc) ([]byte, error) {
							return nil, nil
						},
					},
					decoder: &mockDecoder[mockArticle]{
						DecodeFunc: func(data []byte, val *mockArticle, contentType ContentType, unmarshallFunc UnmarshallFunc) error {
							*val = article
							return nil
						},
					},
				},
			},
			want:    article,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := request[mockArticle](tt.args.ctx, tt.args.cfg, tt.args.options...)
			if (err != nil) != tt.wantErr {
				t.Errorf("request() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("request() got = %v, want %v", got, tt.want)
			}
		})
	}
}
