package file

import (
	"bytes"
	"image"
	"image/png"
	"io"
	"net/http"
	"path/filepath"
	"strings"
)

type File struct {
	Body io.ReadSeeker
	Name string
}

func (file File) Extension() string {
	return strings.ToLower(filepath.Ext(file.Name))
}

func FromReader(name string, body io.ReadSeeker) File {
	return File{
		Body: body,
		Name: name,
	}
}

func FromBytes(name string, data []byte) File {
	return FromReader(name, bytes.NewReader(data))
}

func FromRequest(request *http.Request, key string) (File, error) {
	file, fileHeader, err := request.FormFile(key)
	if err != nil {
		return File{}, err
	}
	defer file.Close()

	return FromReader(fileHeader.Filename, file), nil
}

func FromImage(name string, image image.Image) (File, error) {
	buffer := new(bytes.Buffer)
	err := png.Encode(buffer, image)
	if err != nil {
		return File{}, err
	}

	return FromBytes(name, buffer.Bytes()), nil
}
