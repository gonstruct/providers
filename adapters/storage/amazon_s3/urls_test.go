package amazon_s3_test

import (
	"testing"

	"github.com/gonstruct/providers/adapters/storage/amazon_s3"
)

func TestURL(t *testing.T) {
	cases := []struct {
		name    string
		adapter amazon_s3.Adapter
		want    string
	}{
		{
			name: "base url wins over the endpoint",
			adapter: amazon_s3.Adapter{
				Bucket:   "files",
				Endpoint: "https://account.r2.cloudflarestorage.com",
				BaseURL:  "https://cdn.example.com/",
			},
			want: "https://cdn.example.com/a/b.png",
		},
		{
			name:    "endpoint without a base url",
			adapter: amazon_s3.Adapter{Bucket: "files", Endpoint: "http://minio:9000"},
			want:    "http://minio:9000/files/a/b.png",
		},
		{
			name:    "aws without either",
			adapter: amazon_s3.Adapter{Bucket: "files", Region: "eu-west-1"},
			want:    "https://files.s3.eu-west-1.amazonaws.com/a/b.png",
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			if got := test.adapter.URL("a/b.png"); got != test.want {
				t.Errorf("URL() = %q, want %q", got, test.want)
			}
		})
	}
}
