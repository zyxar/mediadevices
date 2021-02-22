package frame

import (
	"fmt"
	"image"
)

var decoderMap = map[Format]decoderFunc{
	FormatI420: decodeI420,
	FormatI444: decodeI444,
	FormatNV21: decodeNV21,
	FormatNV12: decodeNV12,
	FormatYUY2: decodeYUY2,
	FormatUYVY: decodeUYVY,
	FormatYV12: decodeYV12,
	FormatMJPG: decodeMJPG,
}

func NewDecoder(f Format) (Decoder, error) {
	decoder, ok := decoderMap[f]

	if !ok {
		return nil, fmt.Errorf("%s is not supported", f)
	}

	return decoder, nil
}

type Decoder interface {
	Decode(frame []byte, width, height int) (image.Image, func(), error)
}

// DecoderFunc is a proxy type for Decoder
type decoderFunc func(frame []byte, width, height int) (image.Image, func(), error)

func (f decoderFunc) Decode(frame []byte, width, height int) (image.Image, func(), error) {
	return f(frame, width, height)
}
