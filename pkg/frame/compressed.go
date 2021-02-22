package frame

import (
	"bytes"
	"image"
	"image/jpeg"
)

func decodeMJPG(frame []byte, width, height int) (image.Image, func(), error) {
	img, err := jpeg.Decode(bytes.NewReader(frame))
	return img, func() {}, err
}
