package storage

import (
	s3 "github.com/gonstruct/providers/adapters/storage/amazon_s3"
	"github.com/gonstruct/providers/entities/file"
	"github.com/gonstruct/providers/storage"
)

func Example() {
	// Set up S3 storage
	storage.Adapt(s3.Adapter{
		AccessKeyID:     "your-access-key",
		SecretAccessKey: "your-secret-key",
		Region:          "us-east-1",
		Bucket:          "my-bucket",
	})

	// Upload a file
	f := file.FromBytes("report.pdf", []byte("content"))
	obj, _ := storage.PutFile("uploads", f)

	println(obj.Path) // "uploads/report.pdf"
}
