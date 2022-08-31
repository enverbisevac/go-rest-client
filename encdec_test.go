package rest

import (
	"encoding/json"
	"errors"
	"reflect"
	"runtime"
	"testing"
)

func TestDecoderRegistry_Clone(t *testing.T) {
	r := DecoderRegistry{
		ApplicationJSON: json.Unmarshal,
	}

	c := r.Clone()

	funcName1 := runtime.FuncForPC(reflect.ValueOf(r[ApplicationJSON]).Pointer()).Name()
	funcName2 := runtime.FuncForPC(reflect.ValueOf(c[ApplicationJSON]).Pointer()).Name()

	equals(t, funcName1, funcName2)
}

func TestDecoderRegistry_Decode(t *testing.T) {
	type args struct {
		data    []byte
		object  any
		content Parser
	}
	tests := []struct {
		name    string
		r       DecoderRegistry
		args    args
		wantErr bool
		want    any
	}{
		{
			name: "nil object or content returns nil",
			r: DecoderRegistry{
				ApplicationJSON: json.Unmarshal,
			},
			args: args{
				data:    nil,
				object:  nil,
				content: nil,
			},
		},
		{
			name: "nil data should return error",
			r: DecoderRegistry{
				ApplicationJSON: json.Unmarshal,
			},
			args: args{
				data:    nil,
				object:  &mockArticle{},
				content: ApplicationJSON,
			},
			wantErr: true,
			want:    &mockArticle{},
		},
		{
			name: "unmarshaler not found",
			r: DecoderRegistry{
				ApplicationJSON: json.Unmarshal,
			},
			args: args{
				data:    nil,
				object:  &mockArticle{},
				content: ApplicationXML,
			},
			wantErr: true,
			want:    &mockArticle{},
		},
		{
			name: "decoder error for uncompleted json data",
			r: DecoderRegistry{
				ApplicationJSON: func(data []byte, v any) error {
					return errors.New("decoder error")
				},
			},
			args: args{
				data:    nil,
				object:  &mockArticle{},
				content: ApplicationJSON,
			},
			wantErr: true,
			want:    &mockArticle{},
		},
		{
			name: "happy path",
			r: DecoderRegistry{
				ApplicationJSON: json.Unmarshal,
			},
			args: args{
				data:    []byte("{\"title\": \"some title\", \"content\": \"some content\"}"),
				object:  &mockArticle{},
				content: ApplicationJSON,
			},
			want: &mockArticle{
				Title:   "some title",
				Content: "some content",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.Decode(tt.args.data, tt.args.object, tt.args.content); (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
			}

			equals(t, tt.want, tt.args.object)
		})
	}
}

func TestDecoderRegistry_Set(t *testing.T) {
	r := DecoderRegistry{}

	r.Set(ApplicationJSON, json.Unmarshal)

	funcName1 := runtime.FuncForPC(reflect.ValueOf(json.Unmarshal).Pointer()).Name()
	funcName2 := runtime.FuncForPC(reflect.ValueOf(r[ApplicationJSON]).Pointer()).Name()

	equals(t, funcName1, funcName2)
}

func TestDecoder_Decode(t *testing.T) {
	type fields struct {
		Registry interface {
			Decode(data []byte, object any, content Parser) error
			Set(contentType ContentType, f UnmarshallFunc)
			Clone() DecoderRegistry
		}
	}
	type args struct {
		data           []byte
		val            *mockArticle
		contentType    ContentType
		unmarshallFunc UnmarshallFunc
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
			d := Decoder[mockArticle]{
				Registry: tt.fields.Registry,
			}
			if err := d.Decode(tt.args.data, tt.args.val, tt.args.contentType, tt.args.unmarshallFunc); (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEncoderRegistry_Clone(t *testing.T) {
	r := EncoderRegistry{
		ApplicationJSON: json.Marshal,
	}

	c := r.Clone()

	funcName1 := runtime.FuncForPC(reflect.ValueOf(r[ApplicationJSON]).Pointer()).Name()
	funcName2 := runtime.FuncForPC(reflect.ValueOf(c[ApplicationJSON]).Pointer()).Name()

	equals(t, funcName1, funcName2)
}

func TestEncoderRegistry_Encode(t *testing.T) {
	type args struct {
		object  any
		content Parser
	}
	tests := []struct {
		name    string
		r       EncoderRegistry
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.Encode(tt.args.object, tt.args.content)
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

func TestEncoderRegistry_Set(t *testing.T) {
	r := EncoderRegistry{}

	r.Set(ApplicationJSON, json.Marshal)

	funcName1 := runtime.FuncForPC(reflect.ValueOf(json.Marshal).Pointer()).Name()
	funcName2 := runtime.FuncForPC(reflect.ValueOf(r[ApplicationJSON]).Pointer()).Name()

	equals(t, funcName1, funcName2)
}

func TestEncoder_Encode(t *testing.T) {
	type fields struct {
		Registry interface {
			Encode(object any, content Parser) ([]byte, error)
			Set(contentType ContentType, f MarshallFunc)
			Clone() EncoderRegistry
		}
	}
	type args struct {
		value       any
		contentType ContentType
		marshalFunc MarshallFunc
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Encoder{
				Registry: tt.fields.Registry,
			}
			got, err := e.Encode(tt.args.value, tt.args.contentType, tt.args.marshalFunc)
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
