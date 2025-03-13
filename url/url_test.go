package url_test

import (
	"net/url"
	"reflect"
	"testing"

	urlpkg "github.com/truewebber/gopkg/url"
)

func TestNormalizeWithOptions(t *testing.T) {
	t.Parallel()

	type args struct {
		urlToNormalize string
		options        []urlpkg.AllowOption
	}

	tests := []struct {
		name    string
		args    args
		want    *url.URL
		wantErr bool
	}{
		// cases from https://en.wikipedia.org/wiki/URI_normalization
		{
			name: "Converting the scheme and host to lowercase",
			args: args{
				urlToNormalize: "HTTP://Example.COM/Foo",
				options:        nil,
			},
			want: &url.URL{
				Scheme: "http",
				Host:   "example.com",
				Path:   "/Foo",
			},
			wantErr: false,
		},
		{
			name: "Removing dot-segments",
			args: args{
				urlToNormalize: "http://example.com/foo/./bar/baz/../qux",
				options:        nil,
			},
			want: &url.URL{
				Scheme: "http",
				Host:   "example.com",
				Path:   "/foo/bar/qux",
			},
			wantErr: false,
		},
		{
			name: "Removing the default port",
			args: args{
				urlToNormalize: "http://example.com:80",
				options:        nil,
			},
			want: &url.URL{
				Scheme: "http",
				Host:   "example.com",
			},
			wantErr: false,
		},
		{
			name: "Removing the default port",
			args: args{
				urlToNormalize: "http://example.com:443",
				options:        nil,
			},
			want: &url.URL{
				Scheme: "http",
				Host:   "example.com",
			},
			wantErr: false,
		},
		{
			name: "Preserving non-default port",
			args: args{
				urlToNormalize: "http://example.com:81",
				options:        nil,
			},
			want: &url.URL{
				Scheme: "http",
				Host:   "example.com:81",
			},
			wantErr: false,
		},
		{
			name: `Removing a trailing "/" from empty path`,
			args: args{
				urlToNormalize: "http://example.com/",
				options:        nil,
			},
			want: &url.URL{
				Scheme: "http",
				Host:   "example.com",
			},
			wantErr: false,
		},
		{
			name: `Removing a trailing "/" from non-empty path`,
			args: args{
				urlToNormalize: "http://example.com/foo/",
				options:        nil,
			},
			want: &url.URL{
				Scheme: "http",
				Host:   "example.com",
				Path:   "/foo",
			},
			wantErr: false,
		},
		{
			name: "Sorting (and encoding) keys of query parameters",
			args: args{
				urlToNormalize: "http://example.com/?cars[]=Saab&auto[]=Audi",
				options:        nil,
			},
			want: &url.URL{
				Scheme:   "http",
				Host:     "example.com",
				Path:     "",
				RawQuery: "auto%5B%5D=Audi&cars%5B%5D=Saab",
			},
			wantErr: false,
		},
		{
			name: "Sorting (and encoding) values of query parameters",
			args: args{
				urlToNormalize: "http://example.com/?cars[]=Saab&cars[]=Audi",
				options:        nil,
			},
			want: &url.URL{
				Scheme:   "http",
				Host:     "example.com",
				Path:     "",
				RawQuery: "cars%5B%5D=Audi&cars%5B%5D=Saab",
			},
			wantErr: false,
		},
		{
			name: "Removing duplicate slashes",
			args: args{
				urlToNormalize: "http://example.com/foo//bar.html",
				options:        nil,
			},
			want: &url.URL{
				Scheme:   "http",
				Host:     "example.com",
				Path:     "/foo/bar.html",
				RawQuery: "",
			},
			wantErr: false,
		},
		{
			name: "Removing directory index",
			args: args{
				urlToNormalize: "http://example.com/a/index.html",
				options:        nil,
			},
			want: &url.URL{
				Scheme:   "http",
				Host:     "example.com",
				Path:     "/a",
				RawQuery: "",
			},
			wantErr: false,
		},

		// not implemented:
		// 1. "Converting percent-encoded triplets to uppercase"
		// 2. "Decoding percent-encoded triplets of unreserved characters"

		// other cases
		{
			name: "Allow http",
			args: args{
				urlToNormalize: "http://example.com",
				options:        nil,
			},
			want: &url.URL{
				Scheme:   "http",
				Host:     "example.com",
				Path:     "",
				RawQuery: "",
			},
			wantErr: false,
		},
		{
			name: "Allow https",
			args: args{
				urlToNormalize: "https://example.com",
				options:        nil,
			},
			want: &url.URL{
				Scheme:   "https",
				Host:     "example.com",
				Path:     "",
				RawQuery: "",
			},
			wantErr: false,
		},
		{
			name: "Disallow other schemes",
			args: args{
				urlToNormalize: "file://example.com/foo",
				options:        nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Disallow userinfo",
			args: args{
				urlToNormalize: "https://user@semrush.com:443",
				options:        nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Preserve fragment",
			args: args{
				urlToNormalize: "https://semrush.com/eyeon#fragment",
				options:        nil,
			},
			want: &url.URL{
				Scheme:   "https",
				Host:     "semrush.com",
				Path:     "/eyeon",
				RawQuery: "",
				Fragment: "fragment",
			},
			wantErr: false,
		},
		{
			name: "Remove multiple trailing shashes",
			args: args{
				urlToNormalize: "https://semrush.com//////",
				options:        nil,
			},
			want: &url.URL{
				Scheme: "https",
				Host:   "semrush.com",
				Path:   "",
			},
			wantErr: false,
		},
		{
			name: "Punycode url",
			args: args{
				urlToNormalize: "https://почта.рф/курьеры",
				options:        nil,
			},
			want: &url.URL{
				Scheme:  "https",
				Host:    "xn--80a1acny.xn--p1ai",
				Path:    "/курьеры",
				RawPath: "/курьеры",
			},
			wantErr: false,
		},
		{
			name: "Clean path",
			args: args{
				urlToNormalize: "https://semrush.com/eyeon/folder1/folder2/../folder3",
				options:        nil,
			},
			want: &url.URL{
				Scheme: "https",
				Host:   "semrush.com",
				Path:   "/eyeon/folder1/folder3",
			},
			wantErr: false,
		},
		{
			name: "Invalid url",
			args: args{
				urlToNormalize: "semrush dot com",
				options:        nil,
			},
			want: nil,

			wantErr: true,
		},
		{
			name: "Host without schema",
			args: args{
				urlToNormalize: "semrush.com",
				options:        nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Invalid host with schema",
			args: args{
				urlToNormalize: "https://semrush",
				options:        nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Invalid host with schema and path",
			args: args{
				urlToNormalize: "https://semrush/eyeon",
				options:        nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Punycode host with custom port",
			args: args{
				urlToNormalize: "https://почта.рф:8080",
				options:        nil,
			},
			want: &url.URL{
				Scheme: "https",
				Host:   "xn--80a1acny.xn--p1ai:8080",
				Path:   "",
			},
			wantErr: false,
		},
		{
			name: "Punycode host with port",
			args: args{
				urlToNormalize: "https://почта.рф:443",
				options:        nil,
			},
			want: &url.URL{
				Scheme: "https",
				Host:   "xn--80a1acny.xn--p1ai",
				Path:   "",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := urlpkg.NormalizeWithOptions(tt.args.urlToNormalize, tt.args.options...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NormalizeWithOptions() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NormalizeWithOptions() got = %#v, want %#v", got, tt.want)
			}
		})
	}
}
