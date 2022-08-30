package rest

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestDecodeRegistry_Clone(t *testing.T) {
	tests := []struct {
		name string
		r    DecoderRegistry
		want DecoderRegistry
	}{
		{
			name: "happy path",
			r: DecoderRegistry{
				ApplicationJSON: json.Unmarshal,
			},
			want: DecoderRegistry{
				ApplicationJSON: json.Unmarshal,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.Clone(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Clone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecodeRegistry_Decode(t *testing.T) {
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
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.Decode(tt.args.data, tt.args.object, tt.args.content); (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDecodeRegistry_Set(t *testing.T) {
	type args struct {
		contentType ContentType
		f           UnmarshallFunc
	}
	tests := []struct {
		name string
		r    DecoderRegistry
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.r.Set(tt.args.contentType, tt.args.f)
		})
	}
}

func TestEncodeRegistry_Clone(t *testing.T) {
	tests := []struct {
		name string
		r    EncoderRegistry
		want EncoderRegistry
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.Clone(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Clone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncodeRegistry_Encode(t *testing.T) {
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

func TestEncodeRegistry_Set(t *testing.T) {
	type args struct {
		contentType ContentType
		f           MarshallFunc
	}
	tests := []struct {
		name string
		r    EncoderRegistry
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.r.Set(tt.args.contentType, tt.args.f)
		})
	}
}
