package rest

import (
	"errors"
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
	base := RequestData[any]{}

	got := WithBody(s)
	got(&base)

	equals(t, s, base.Body)
}

func TestWithHeaders(t *testing.T) {
	base := RequestData[any]{}

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

func TestRequestData_Decode(t *testing.T) {
	type fields struct {
		Body           any
		Header         http.Header
		encoder        Encoder
		decoder        Decoder
		MarshalFunc    MarshallFunc
		UnmarshallFunc UnmarshallFunc
	}
	type args struct {
		data []byte
		v    *mockArticle
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RequestData[mockArticle]{
				Body:           tt.fields.Body,
				Header:         tt.fields.Header,
				encoder:        tt.fields.encoder,
				decoder:        tt.fields.decoder,
				MarshalFunc:    tt.fields.MarshalFunc,
				UnmarshallFunc: tt.fields.UnmarshallFunc,
			}
			if err := r.Decode(tt.args.data, tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRequestData_Encode(t *testing.T) {
	type fields struct {
		Body        any
		Header      http.Header
		encoder     Encoder
		MarshalFunc MarshallFunc
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "happy path",
			fields: fields{
				Body: mockArticle{
					Title:   "some title",
					Content: "some content",
				},
				Header: nil,
				encoder: &mockEncoder{
					content: []byte("{\"title\":\"some title\",\"content\":\"some content\"}"),
				},
				MarshalFunc: nil,
			},
			want: []byte("{\"title\":\"some title\",\"content\":\"some content\"}"),
		},
		{
			name: "encoder error",
			fields: fields{
				Body: mockArticle{
					Title:   "some title",
					Content: "some content",
				},
				Header: nil,
				encoder: &mockEncoder{
					err: errors.New("encoder error"),
				},
				MarshalFunc: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RequestData[mockArticle]{
				Body:        tt.fields.Body,
				Header:      tt.fields.Header,
				encoder:     tt.fields.encoder,
				MarshalFunc: tt.fields.MarshalFunc,
			}
			got, err := r.Encode()
			if (err != nil) != tt.wantErr {
				t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encode() got = %v, want %v", got, tt.want)
			}
		})
	}
}
