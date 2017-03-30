package imageorient_test

import (
	"image/jpeg"
	"log"
	"os"

	"github.com/disintegration/imageorient"
)

func ExampleDecode() {
	// Open the test image. This particular image have the EXIF
	// orientation tag set to 3 (rotated by 180 deg).
	f, err := os.Open("testdata/orientation_3.jpg")
	if err != nil {
		log.Fatalf("os.Open failed: %v", err)
	}

	// Decode the test image using the imageorient.Decode function
	// to handle the image orientation correctly.
	img, _, err := imageorient.Decode(f)
	if err != nil {
		log.Fatalf("imageorient.Decode failed: %v", err)
	}

	// Save the decoded image to a new file. If we used image.Decode
	// instead of imageorient.Decode on the previous step, the saved
	// image would appear rotated.
	f, err = os.Create("testdata/example_output.jpg")
	if err != nil {
		log.Fatalf("os.Create failed: %v", err)
	}
	err = jpeg.Encode(f, img, nil)
	if err != nil {
		log.Fatalf("jpeg.Encode failed: %v", err)
	}
}
