// Package imageorient provides image decoding functions similar to standard library's
// image.Decode and image.DecodeConfig with the addition that they also handle the
// EXIF orientation tag (if present).
//
// See also: http://www.daveperrett.com/articles/2012/07/28/exif-orientation-handling-is-a-ghetto/
//
package imageorient

import (
	"bytes"
	"image"
	"io"

	"github.com/disintegration/gift"
	"github.com/rwcarlsen/goexif/exif"
)

// maxBufLen is the maximum size of a buffer that should be enough to read
// the EXIF metadata. According to the EXIF sepcs, it is located inside the
// APP1 block that goes right after the start of image (SOI).
const maxBufLen = 1 << 20

// Decode decodes an image and changes its orientation
// according to the EXIF orientation tag (if present).
func Decode(r io.Reader) (image.Image, string, error) {
	orientation, r := getOrientation(r)

	img, format, err := image.Decode(r)
	if err != nil {
		return img, format, err
	}

	return fixOrientation(img, orientation), format, nil
}

// DecodeConfig decodes the color model and dimensions of an image
// with the respect to the EXIF orientation tag (if present).
func DecodeConfig(r io.Reader) (image.Config, string, error) {
	orientation, r := getOrientation(r)

	cfg, format, err := image.DecodeConfig(r)
	if err != nil {
		return cfg, format, err
	}

	if orientation >= 5 && orientation <= 8 {
		cfg.Width, cfg.Height = cfg.Height, cfg.Width
	}

	return cfg, format, nil
}

// getOrientation returns the EXIF orientation tag from the given image
// and a new io.Reader with the same state as the original reader r.
func getOrientation(r io.Reader) (int, io.Reader) {
	buf := new(bytes.Buffer)
	tr := io.TeeReader(io.LimitReader(r, maxBufLen), buf)
	orientation := readOrientation(tr)
	return orientation, io.MultiReader(buf, r)
}

// readOrientation reads the EXIF orientation tag from the given image.
// It returns 0 if the orientation tag is not found or invalid.
func readOrientation(r io.Reader) int {
	x, err := exif.Decode(r)
	if err != nil {
		return 0
	}

	tag, err := x.Get(exif.Orientation)
	if err != nil {
		return 0
	}

	orientation, err := tag.Int(0)
	if err != nil {
		return 0
	}

	if orientation < 1 || orientation > 8 {
		return 0
	}

	return orientation
}

// Filters needed to fix the given image orientation.
var filters = map[int]gift.Filter{
	2: gift.FlipHorizontal(),
	3: gift.Rotate180(),
	4: gift.FlipVertical(),
	5: gift.Transpose(),
	6: gift.Rotate270(),
	7: gift.Transverse(),
	8: gift.Rotate90(),
}

// fixOrientation changes the image orientation based on the EXIF orientation tag value.
func fixOrientation(img image.Image, orientation int) image.Image {
	filter, ok := filters[orientation]
	if !ok {
		return img
	}
	g := gift.New(filter)
	newImg := image.NewRGBA(g.Bounds(img.Bounds()))
	g.Draw(newImg, img)
	return newImg
}
