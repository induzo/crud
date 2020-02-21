package rest

import (
	"net/url"
	"reflect"
	"testing"

	"github.com/induzo/crud"
)

func TestListModifiersFromURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want crud.ListModifiers
	}{
		{
			name: "different querystring params",
			url:  "http://test.com?premier=1",
			want: crud.ListModifiers{
				"premier": []string{"1"},
			},
		},
		{
			name: "same querystring param name",
			url:  "http://test.com?orderby=test%20ASC&orderby=pol%20DESC",
			want: crud.ListModifiers{
				"orderby": []string{"test ASC", "pol DESC"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, _ := url.Parse(tt.url)
			if got := ListModifiersFromURL(u); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListModifiersFromURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
