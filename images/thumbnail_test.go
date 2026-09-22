package images_test

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"testing"

	"github.com/gonstruct/providers/images"
)

func TestMakeThumbnailShrinksAWidePNG(t *testing.T) {
	source := image.NewRGBA(image.Rect(0, 0, 640, 400))

	for y := range 400 {
		for x := range 640 {
			source.Set(x, y, color.RGBA{R: 200, G: 40, B: 40, A: 255})
		}
	}

	encoded := new(bytes.Buffer)
	if err := png.Encode(encoded, source); err != nil {
		t.Fatal(err)
	}

	thumbnail, err := images.MakeThumbnail(encoded, 320)
	if err != nil || thumbnail.Extension != ".png" || thumbnail.Width != 320 || thumbnail.Height != 200 {
		t.Fatalf("expected a 320 by 200 png, got %+v and %v", thumbnail, err)
	}

	decoded, err := png.Decode(bytes.NewReader(thumbnail.Bytes))
	if err != nil {
		t.Fatal(err)
	}

	if r, _, _, _ := decoded.At(100, 100).RGBA(); r>>8 != 200 {
		t.Fatalf("expected the red to survive, got %d", r>>8)
	}
}

func TestMakeThumbnailShrinksATallJPEGByHeight(t *testing.T) {
	tall := new(bytes.Buffer)
	if err := jpeg.Encode(tall, image.NewRGBA(image.Rect(0, 0, 300, 900)), nil); err != nil {
		t.Fatal(err)
	}

	thumbnail, err := images.MakeThumbnail(tall, 320)
	if err != nil || thumbnail.Extension != ".jpg" || thumbnail.Width != 106 || thumbnail.Height != 320 {
		t.Fatalf("expected a 106 by 320 jpeg, got %+v and %v", thumbnail, err)
	}
}

func TestMakeThumbnailKeepsASmallImageItsSize(t *testing.T) {
	small := new(bytes.Buffer)
	if err := png.Encode(small, image.NewRGBA(image.Rect(0, 0, 12, 8))); err != nil {
		t.Fatal(err)
	}

	thumbnail, err := images.MakeThumbnail(small, 320)
	if err != nil || thumbnail.Width != 12 || thumbnail.Height != 8 {
		t.Fatalf("expected 12 by 8, got %+v and %v", thumbnail, err)
	}
}

func TestMakeThumbnailRefusesWhatIsNotAnImage(t *testing.T) {
	if _, err := images.MakeThumbnail(bytes.NewReader([]byte("png-one")), 320); !errors.Is(err, images.ErrNotAnImage) {
		t.Fatalf("expected ErrNotAnImage, got %v", err)
	}
}
