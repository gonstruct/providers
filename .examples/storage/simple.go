package storage

import (
	s3 "github.com/gonstruct/providers/adapters/storage/amazon_s3"
	"github.com/gonstruct/providers/entities/file"
	"github.com/gonstruct/providers/storage"
)

func Example() {
	storage.Adapt(s3.Adapter{
		AccessKeyID:     "your-access-key-id",
		SecretAccessKey: "your-secret-access-key",
		Region:          "your-region",
		Bucket:          "your-bucket",
		Endpoint:        "https://s3.your-region.amazonaws.com",
		UsePathStyle:    true,
	})

	file := file.FromBytes("hello.txt", []byte("hello world"))
	object, _ := storage.PutFile("this/folder", file)

	println("MimeType:", object.MimeType)
}
