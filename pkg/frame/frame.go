package frame

type Format string

const ( // ref: https://chromium.googlesource.com/libyuv/libyuv/+show/master/docs/formats.md
	// 10 Primary YUV formats: 5 planar, 2 biplanar, 2 packed.
	FormatI420 Format = "I420"
	FormatI422 Format = "I422"
	FormatI444 Format = "I444"
	FormatI400 Format = "I400"
	FormatNV21 Format = "NV21"
	FormatNV12 Format = "NV12"
	FormatYUY2 Format = "YUY2"
	FormatUYVY Format = "UYVY"
	FormatH010 Format = "H010"
	FormatU010 Format = "U010"
	// 11 Primary RGB formats: 4 32 bpp, 2 24 bpp, 3 16 bpp, 1 10 bpc
	FormatARGB Format = "ARGB"
	FormatBGRA Format = "BGRA"
	FormatABGR Format = "ABGR"
	FormatAR30 Format = "AR30"
	FormatAB30 Format = "AB30"
	Format24BG Format = "24BG"
	FormatRAW  Format = "RAW"
	FormatRGBA Format = "RGBA"
	FormatRGBP Format = "RGBP"
	FormatRGBO Format = "RGBO"
	FormatR444 Format = "R444"
	// 1 Primary Compressed YUV format.
	FormatMJPG Format = "MJPG"
	// 11 Auxiliary YUV variations: 3 with U and V planes are swapped, 1 Alias.
	FormatYV12 Format = "YV12"
	FormatYV16 Format = "YV16"
	FormatYV24 Format = "YV24"
	FormatYU12 Format = "YU12"
	FormatJ420 Format = "J420"
	FormatJ400 Format = "J400"
	FormatH420 Format = "H420"
	FormatH422 Format = "H422"
	FormatU420 Format = "U420"
	FormatU422 Format = "U422"
	FormatU444 Format = "U444"
	// 14 Auxiliary aliases.  CanonicalFourCC() maps these to canonical fourcc.
	FormatIYUV  = FormatI420 // Alias for I420.
	FormatYU16  = FormatI422 // Alias for I422.
	FormatYU24  = FormatI444 // Alias for I444.
	FormatYUYV  = FormatYUY2 // Alias for YUY2.
	FormatYUVS  = FormatYUY2 // Alias for YUY2 on Mac.
	FormatHDYC  = FormatUYVY // Alias for UYVY.
	Format2VUY  = FormatUYVY // Alias for UYVY on Mac.
	FormatJPEG  = FormatMJPG // Alias for MJPG.
	FormatDMB1  = FormatMJPG // Alias for MJPG on Mac.
	FormatMJPEG = FormatMJPG // Alias for MJPG.
	FormatRGB3  = FormatRAW  // Alias for RAW.
	FormatBGR3  = Format24BG // Alias for 24BG.
	FormatCM32  = FormatBGRA // Alias for BGRA kCVPixelFormatType_32ARGB
	FormatCM24  = FormatRAW  // Alias for RAW kCVPixelFormatType_24RGB
	FormatL555  = FormatRGBO // Alias for RGBO.
	FormatL565  = FormatRGBP // Alias for RGBP.
	Format5551  = FormatRGBO // Alias for RGBO.
	//FormatBA81 = FormatBGGR // Alias for BGGR.
)
