package rest

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"testing"
)

type mockEncoder struct {
	content []byte
	err     error
}

func (m *mockEncoder) Encode(object any, content Parser) ([]byte, error) {
	return m.content, m.err
}

func (m *mockEncoder) Set(contentType ContentType, marshallFunc MarshallFunc) {

}

func (m *mockEncoder) Clone() EncodeRegistry {
	return EncodeRegistry{}
}

type mockDecoder struct {
	object any
	err    error
}

func (m *mockDecoder) Decode(data []byte, object any, content Parser) error {
	if m.object != nil {
		value := reflect.ValueOf(m.object)
		reflect.ValueOf(object).Elem().Set(value)
	}
	return m.err
}

func (m *mockDecoder) Set(contentType ContentType, unmarshaler UnmarshallFunc) {

}

func (m *mockDecoder) Clone() DecodeRegistry {
	return DecodeRegistry{}
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
		cfg     *config
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
				cfg: &config{
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
					encoder: &mockEncoder{},
					decoder: &mockDecoder{},
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
				cfg: &config{
					method:     http.MethodGet,
					requestURL: "",
					requester:  nil,
					encoder: &mockEncoder{
						err: errors.New("encoding error"),
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
				cfg: &config{
					method:     http.MethodGet,
					requestURL: "",
					requester: &mockRequest{
						err: errors.New("request failed"),
					},
					encoder: &mockEncoder{},
					decoder: &mockDecoder{},
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
				cfg: &config{
					method:     http.MethodGet,
					requestURL: "",
					requester: &mockRequest{
						response: &Response{
							StatusCode: http.StatusBadRequest,
							Header:     Header(WithContent(ApplicationJSON)),
						},
						err: nil,
					},
					encoder: &mockEncoder{},
					decoder: &mockDecoder{},
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
				cfg: &config{
					method:     http.MethodGet,
					requestURL: "",
					requester: &mockRequest{
						response: &Response{
							Body:       []byte("{}"),
							StatusCode: http.StatusOK,
							Header:     Header(WithContent(ApplicationJSON)),
						},
					},
					encoder: &mockEncoder{},
					decoder: &mockDecoder{
						err: errors.New("decoding error"),
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
				cfg: &config{
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
					encoder: &mockEncoder{},
					decoder: &mockDecoder{
						object: article,
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
