// Package avfoundation provides AVFoundation binding for Go
package avfoundation

// #cgo CFLAGS: -x objective-c
// #cgo LDFLAGS: -framework AVFoundation -framework Foundation -framework CoreMedia -framework CoreVideo
// #include "AVFoundationBind/AVFoundationBind.h"
// #include "AVFoundationBind/AVFoundationBind.m"
// extern void onData(void*, void*, int);
// void onDataBridge(void *userData, void *buf, int len) {
// 	onData(userData, buf, len);
// }
import "C"
import (
	"fmt"
	"io"
	"unsafe"

	"github.com/pion/mediadevices/pkg/frame"
	"github.com/pion/mediadevices/pkg/prop"
)

type MediaType C.AVBindMediaType

const (
	Video = MediaType(C.AVBindMediaTypeVideo)
	Audio = MediaType(C.AVBindMediaTypeAudio)
)

// Device represents a metadata that later can be used to retrieve back the
// underlying device given by AVFoundation
type Device struct {
	// UID is a unique identifier for a device
	UID     string
	cDevice C.AVBindDevice
}

type pixelFormatSpec struct {
	fmt frame.Format
	fcc C.FourCharCode
}

// Useful mapping reference from ffmpeg:
// https://github.com/FFmpeg/FFmpeg/blob/c810a9502cebe32e1dd08ee3d0d17053dde44aa9/libavdevice/avfoundation.m#L53-L80
var pixelFormats = []pixelFormatSpec{
	{frame.FormatUYVY, C.kCVPixelFormatType_422YpCbCr8},
	{frame.FormatI420, C.kCVPixelFormatType_420YpCbCr8Planar},
	{frame.FormatNV12, C.kCVPixelFormatType_420YpCbCr8BiPlanarVideoRange},
	{frame.FormatNV21, C.kCVPixelFormatType_420YpCbCr8BiPlanarFullRange},
	{frame.FormatYUY2, C.kCVPixelFormatType_422YpCbCr8_yuvs},
	//C.kCVPixelFormatType_1Monochrome,
	//C.kCVPixelFormatType_16BE555,
	//C.kCVPixelFormatType_16LE555,
	//C.kCVPixelFormatType_16BE565,
	//C.kCVPixelFormatType_16LE565,
	//C.kCVPixelFormatType_24RGB,
	//C.kCVPixelFormatType_24BGR,
	//C.kCVPixelFormatType_32ARGB,
	//C.kCVPixelFormatType_32BGRA,
	//C.kCVPixelFormatType_32ABGR,
	//C.kCVPixelFormatType_32RGBA,
	//C.kCVPixelFormatType_48RGB,
	//C.kCVPixelFormatType_4444YpCbCrA8R,
	//C.kCVPixelFormatType_4444AYpCbCr16,
	//C.kCVPixelFormatType_444YpCbCr8,
	//C.kCVPixelFormatType_422YpCbCr16,
	//C.kCVPixelFormatType_422YpCbCr10,
	//C.kCVPixelFormatType_444YpCbCr10,
}

func pixelFormatToFourCharCode(f frame.Format) (C.FourCharCode, bool) {
	for _, spec := range pixelFormats {
		if spec.fmt == f {
			return spec.fcc, true
		}
	}
	return 0, false
}

func fourCharCodeToPixelFormat(c C.FourCharCode) (frame.Format, bool) {
	for _, spec := range pixelFormats {
		if spec.fcc == c {
			return spec.fmt, true
		}
	}
	return "", false
}

// Devices uses AVFoundation to query a list of devices based on the media type
func Devices(mediaType MediaType) ([]Device, error) {
	var cDevicesPtr C.PAVBindDevice
	var cDevicesLen C.int

	status := C.AVBindDevices(C.AVBindMediaType(mediaType), &cDevicesPtr, &cDevicesLen)
	if status != nil {
		return nil, fmt.Errorf("%s", C.GoString(status))
	}

	// https://github.com/golang/go/wiki/cgo#turning-c-arrays-into-go-slices
	cDevices := (*[1 << 28]C.AVBindDevice)(unsafe.Pointer(cDevicesPtr))[:cDevicesLen:cDevicesLen]
	devices := make([]Device, cDevicesLen)

	for i := range devices {
		devices[i].UID = C.GoString(&cDevices[i].uid[0])
		devices[i].cDevice = cDevices[i]
	}

	return devices, nil
}

// ReadCloser is a wrapper around the data callback from AVFoundation. The data received from the
// the underlying callback can be retrieved by calling Read.
type ReadCloser struct {
	dataChan chan []byte
	id       handleID
	onClose  func()
}

func newReadCloser(onClose func()) *ReadCloser {
	var rc ReadCloser
	rc.dataChan = make(chan []byte, 1)
	rc.onClose = onClose
	rc.id = register(rc.dataCb)
	return &rc
}

func (rc *ReadCloser) dataCb(data []byte) {
	// TODO: add a policy for slow reader
	rc.dataChan <- data
}

// Read reads raw data, the format is determined by the media type and property:
//   - For video, each call will return a frame.
//   - For audio, each call will return a chunk which its size configured by Latency
func (rc *ReadCloser) Read() ([]byte, func(), error) {
	data, ok := <-rc.dataChan
	if !ok {
		return nil, func() {}, io.EOF
	}
	return data, func() {}, nil
}

// Close closes the capturing session, and no data will flow anymore
func (rc *ReadCloser) Close() {
	if rc.onClose != nil {
		rc.onClose()
	}
	close(rc.dataChan)
	unregister(rc.id)
}

// Session represents a capturing session.
type Session struct {
	device   Device
	cSession C.PAVBindSession
}

// NewSession creates a new capturing session
func NewSession(device Device) (*Session, error) {
	var session Session

	status := C.AVBindSessionInit(device.cDevice, &session.cSession)
	if status != nil {
		return nil, fmt.Errorf("%s", C.GoString(status))
	}

	session.device = device
	return &session, nil
}

// Close stops capturing session and frees up resources
func (session *Session) Close() error {
	if session.cSession == nil {
		return nil
	}

	status := C.AVBindSessionFree(&session.cSession)
	if status != nil {
		return fmt.Errorf("%s", C.GoString(status))
	}
	return nil
}

// Open start capturing session. As soon as it returns successfully, the data will start
// flowing. The raw data can be retrieved by using ReadCloser's Read method.
func (session *Session) Open(property prop.Media) (*ReadCloser, error) {
	fourcc, ok := pixelFormatToFourCharCode(property.FrameFormat)
	if !ok {
		return nil, fmt.Errorf("Unsupported frame format")
	}

	cProperty := C.AVBindMediaProperty{
		width:  C.int(property.Width),
		height: C.int(property.Height),
		fourcc: fourcc,
	}

	rc := newReadCloser(func() {
		C.AVBindSessionClose(session.cSession)
	})
	status := C.AVBindSessionOpen(
		session.cSession,
		cProperty,
		C.AVBindDataCallback(unsafe.Pointer(C.onDataBridge)),
		unsafe.Pointer(&rc.id),
	)
	if status != nil {
		return nil, fmt.Errorf("%s", C.GoString(status))
	}
	return rc, nil
}

// Properties queries a list of properties that device supports
func (session *Session) Properties() []prop.Media {
	var cPropertiesPtr C.PAVBindMediaProperty
	var cPropertiesLen C.int

	status := C.AVBindSessionProperties(session.cSession, &cPropertiesPtr, &cPropertiesLen)
	if status != nil {
		return nil
	}
	defer C.free(unsafe.Pointer(cPropertiesPtr))
	// https://github.com/golang/go/wiki/cgo#turning-c-arrays-into-go-slices
	cProperties := (*[1 << 28]C.AVBindMediaProperty)(unsafe.Pointer(cPropertiesPtr))[:cPropertiesLen:cPropertiesLen]
	var properties []prop.Media
	for _, cProperty := range cProperties {
		frameFormat, ok := fourCharCodeToPixelFormat(cProperty.fourcc)
		if ok {
			properties = append(properties, prop.Media{
				Video: prop.Video{
					Width:       int(cProperty.width),
					Height:      int(cProperty.height),
					FrameFormat: frameFormat,
				},
			})
		}
	}
	return properties
}
