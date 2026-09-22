// Package images is what an application does to pictures: small copies
// for lists, made once where the file is stored.
package images

import (
	"bytes"
	"errors"
	"image"
	_ "image/gif" // a format a thumbnail can be read from
	"image/jpeg"
	"image/png"
	"io"

	"golang.org/x/image/draw"
)

// ErrNotAnImage is a source that no known image format decodes.
var ErrNotAnImage = errors.New("images: not an image")

// jpegQuality is what a thumbnail is encoded at: small, and no artefact
// shows at thumbnail size.
const jpegQuality = 82

// Thumbnail is a small copy of an image.
type Thumbnail struct {
	Bytes  []byte
	Width  int
	Height int
	// Extension says what the bytes are: ".png" when the source was png,
	// so transparency survives, ".jpg" otherwise.
	Extension string
}

// MakeThumbnail scales an image down so its longest side is at most
// `longest` pixels, keeping its proportions. An image already small enough
// comes back re-encoded at its own size.
func MakeThumbnail(source io.Reader, longest int) (*Thumbnail, error) {
	decoded, format, err := image.Decode(source)
	if err != nil {
		return nil, errors.Join(ErrNotAnImage, err)
	}

	bounds := decoded.Bounds()

	width, height := bounds.Dx(), bounds.Dy()
	if width > longest || height > longest {
		if width >= height {
			width, height = longest, max(1, height*longest/bounds.Dx())
		} else {
			width, height = max(1, width*longest/bounds.Dy()), longest
		}
	}

	scaled := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.CatmullRom.Scale(scaled, scaled.Bounds(), decoded, bounds, draw.Over, nil)

	buffer := new(bytes.Buffer)
	if format == "png" {
		if err := png.Encode(buffer, scaled); err != nil {
			return nil, err
		}

		return &Thumbnail{Bytes: buffer.Bytes(), Width: width, Height: height, Extension: ".png"}, nil
	}

	if err := jpeg.Encode(buffer, scaled, &jpeg.Options{Quality: jpegQuality}); err != nil {
		return nil, err
	}

	return &Thumbnail{Bytes: buffer.Bytes(), Width: width, Height: height, Extension: ".jpg"}, nil
}
