package validate

import "testing"

func TestNormalizeURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		in      string
		want    string
		wantErr error
	}{
		{name: "empty", in: "", wantErr: ErrEmptyURL},
		{name: "whitespace", in: "   ", wantErr: ErrEmptyURL},
		{name: "relative", in: "/path", wantErr: ErrInvalidURL},
		{name: "no scheme", in: "example.com", wantErr: ErrInvalidURL},
		{name: "ftp", in: "ftp://example.com", wantErr: ErrBadScheme},
		{name: "https ok", in: "https://Example.com/Path", want: "https://Example.com/Path"},
		{name: "http ok", in: "http://example.com", want: "http://example.com"},
		{name: "strip fragment", in: "https://example.com/a#frag", want: "https://example.com/a"},
		{name: "default https port", in: "https://example.com:443/x", want: "https://example.com/x"},
		{name: "default http port", in: "http://example.com:80/", want: "http://example.com/"},
		{name: "trim space", in: "  https://example.com  ", want: "https://example.com"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := NormalizeURL(tt.in)
			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Fatalf("NormalizeURL(%q) error = %v, want %v", tt.in, err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("NormalizeURL(%q) unexpected error: %v", tt.in, err)
			}
			if got != tt.want {
				t.Fatalf("NormalizeURL(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
