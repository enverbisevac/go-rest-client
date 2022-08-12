package rest

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type article struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

func TestGet(t *testing.T) {
	title := "golang generics"

	// Start a local HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		equals(t, req.URL.String(), "/resource/article/1")
		// Send response to be tested
		rw.Header().Set(Content, ApplicationJSON.String())
		art := article{
			Title: title,
			Body:  "",
		}
		bytes, err := json.Marshal(art)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, _ = rw.Write(bytes)
	}))
	// Close the server when test finishes
	defer server.Close()

	art, err := Get[article](context.Background(), server.URL+"/resource/article/1", nil, Header())
	ok(t, err)
	equals(t, title, art.Title)
}

func TestGetWrongUrl(t *testing.T) {
	// Start a local HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Send response to be tested
		rw.WriteHeader(http.StatusNotFound)
	}))
	// Close the server when test finishes
	defer server.Close()

	_, err := Get[article](context.Background(), server.URL+"/resource/article/1", nil, Header())
	var httpError Error
	ok := errors.As(err, &httpError)
	assert(t, ok && httpError.StatusCode == http.StatusNotFound, "Expected "+err.Error())
}

func TestDelete(t *testing.T) {
	// Start a local HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		equals(t, req.URL.String(), "/resource/article/1")
		// Send response to be tested
		rw.WriteHeader(http.StatusNoContent)
	}))
	// Close the server when test finishes
	defer server.Close()

	err := Delete(context.Background(), server.URL+"/resource/article/1", Header())
	ok(t, err)
}

func TestDeleteWrongUrl(t *testing.T) {
	// Start a local HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Send response to be tested
		rw.WriteHeader(http.StatusNotFound)
	}))
	// Close the server when test finishes
	defer server.Close()

	err := Delete(context.Background(), server.URL+"/resource/article/1", Header())
	var httpError Error
	ok := errors.As(err, &httpError)
	assert(t, ok && httpError.StatusCode == http.StatusNotFound, "Expected "+err.Error())
}

func Test_marshal(t *testing.T) {
	type args struct {
		object  any
		content ContentType
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "test json marshaller - happy path",
			args: args{
				object:  []string{"enver", "bisevac"},
				content: ApplicationJSON,
			},
			want: []byte("[\"enver\",\"bisevac\"]"),
		},
		{
			name: "test xml marshaller - happy path",
			args: args{
				object:  []string{"enver", "bisevac"},
				content: ApplicationXML,
			},
			want: []byte("<string>enver</string><string>bisevac</string>"),
		},
		{
			name: "wrong content type",
			args: args{
				object:  []string{"enver", "bisevac"},
				content: "custom content type",
			},
			want:    []byte{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := marshal(tt.args.object, tt.args.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("marshal() got = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}

func Test_unmarshal(t *testing.T) {
	var s []string
	type args struct {
		data    []byte
		object  any
		content ContentType
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test json unmarshaller - happy path",
			args: args{
				data:    []byte("[\"enver\",\"bisevac\"]"),
				object:  &s,
				content: ApplicationJSON,
			},
		},
		{
			name: "test xml unmarshaller - happy path",
			args: args{
				data:    []byte("<string>enver</string><string>bisevac</string>"),
				object:  &s,
				content: ApplicationXML,
			},
		},
		{
			name: "wrong content type for unmarshaller",
			args: args{
				data:    []byte("<string>enver</string><string>bisevac</string>"),
				object:  &s,
				content: "some content type",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := unmarshal(tt.args.data, tt.args.object, tt.args.content)
			if (err != nil) != tt.wantErr {
				ok(t, err)
			}
		})
	}
}
