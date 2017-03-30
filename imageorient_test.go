package imageorient

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"testing"
)

var testFiles = []struct {
	path        string
	orientation int
}{
	{"testdata/orientation_0.jpg", 0},
	{"testdata/orientation_1.jpg", 1},
	{"testdata/orientation_2.jpg", 2},
	{"testdata/orientation_3.jpg", 3},
	{"testdata/orientation_4.jpg", 4},
	{"testdata/orientation_5.jpg", 5},
	{"testdata/orientation_6.jpg", 6},
	{"testdata/orientation_7.jpg", 7},
	{"testdata/orientation_8.jpg", 8},
}

var (
	testWidth  = 50
	testHeight = 70
)

func TestReadOrientation(t *testing.T) {
	for _, tf := range testFiles {
		f, err := os.Open(tf.path)
		if err != nil {
			t.Fatalf("os.Open(%q): %v", tf.path, err)
		}

		o := readOrientation(f)
		if o != tf.orientation {
			t.Fatalf("expected orientation=%d but got %d (%s)", tf.orientation, o, tf.path)
		}
	}
}

func TestDecodeConfig(t *testing.T) {
	for _, tf := range testFiles {
		f, err := os.Open(tf.path)
		if err != nil {
			t.Fatalf("Open(%q): %v", tf.path, err)
		}

		cfg, format, err := DecodeConfig(f)
		if err != nil {
			t.Fatalf("Decode(%q): %v", tf.path, err)
		}

		if cfg.Width != testWidth {
			t.Fatalf("bad decoded width (%q): want %d got %d", tf.path, testWidth, cfg.Width)
		}

		if cfg.Height != testHeight {
			t.Fatalf("bad decoded width (%q): want %d got %d", tf.path, testHeight, cfg.Height)
		}

		if format != "jpeg" {
			t.Fatalf("bad decoded format (%q): %s", tf.path, format)
		}
	}
}

func TestDecode(t *testing.T) {
	f, err := os.Open("testdata/orientation_0.jpg")
	if err != nil {
		t.Fatalf("os.Open(%q): %v", "testdata/orientation_0.jpg", err)
	}
	orig, _, err := image.Decode(f)
	if err != nil {
		t.Fatalf("image.Decode(%q): %v", "testdata/orientation_0.jpg", err)
	}

	for _, tf := range testFiles {
		f, err := os.Open(tf.path)
		if err != nil {
			t.Fatalf("Open(%q): %v", tf.path, err)
		}

		img, format, err := Decode(f)
		if err != nil {
			t.Fatalf("Decode(%q): %v", tf.path, err)
		}

		if err := compareImages(img, orig); err != nil {
			t.Fatalf("decoded image differs from orig (%q): %v", tf.path, err)
		}

		if img.Bounds().Dx() != testWidth {
			t.Fatalf("bad decoded width (%q): want %d got %d", tf.path, testWidth, img.Bounds().Dx())
		}

		if img.Bounds().Dy() != testHeight {
			t.Fatalf("bad decoded width (%q): want %d got %d", tf.path, testHeight, img.Bounds().Dy())
		}

		if format != "jpeg" {
			t.Fatalf("bad decoded format (%q): %s", tf.path, format)
		}
	}
}

func compareImages(m1, m2 image.Image) error {
	if m1.Bounds() != m2.Bounds() {
		return fmt.Errorf("bounds not equal: %v vs %v", m1.Bounds(), m2.Bounds())
	}
	for x := m1.Bounds().Min.X; x < m1.Bounds().Max.X; x++ {
		for y := m1.Bounds().Min.Y; y < m1.Bounds().Max.Y; y++ {
			c1 := color.GrayModel.Convert(m1.At(x, y)).(color.Gray)
			c2 := color.GrayModel.Convert(m2.At(x, y)).(color.Gray)
			d := int(c1.Y) - int(c2.Y)
			if d < -5 || d > 5 {
				return fmt.Errorf("pixels at (%d, %d) not equal: %v vs %v", x, y, c1, c2)
			}
		}
	}
	return nil
}
