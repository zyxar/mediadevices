package frame

import (
	"fmt"
	"image"
)

func decodeARGB(frame []byte, width, height int) (image.Image, func(), error) {
	size := 4 * width * height
	if size > len(frame) {
		return nil, func() {}, fmt.Errorf("frame length (%d) less than expected (%d)", len(frame), size)
	}
	r := image.Rect(0, 0, width, height)
	return &image.RGBA{
		Pix:    frame[:size:size],
		Stride: 4 * r.Dx(),
		Rect:   r,
	}, func() {}, nil
}

func decodeBGRA(frame []byte, width, height int) (image.Image, func(), error) {
	size := 4 * width * height
	if size > len(frame) {
		return nil, func() {}, fmt.Errorf("frame length (%d) less than expected (%d)", len(frame), size)
	}
	r := image.Rect(0, 0, width, height)
	for i := 0; i < size; i += 4 {
		frame[i], frame[i+1], frame[i+2], frame[i+3] = frame[i+1], frame[i+2], frame[i+3], frame[i]
	}
	return &image.RGBA{
		Pix:    frame[:size:size],
		Stride: 4 * r.Dx(),
		Rect:   r,
	}, func() {}, nil
}
