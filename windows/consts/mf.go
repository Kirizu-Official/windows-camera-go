package consts

const MF_SOURCE_READER_FIRST_VIDEO_STREAM = uint32(0xFFFFFFFC)
const MF_SOURCE_READER_FIRST_AUDIO_STREAM = uint32(0xFFFFFFFD)
const MF_SOURCE_READER_ALL_STREAMS = uint32(0xFFFFFFFF)

type CameraControlFlags int32

const (
	KSPROPERTY_CAMERACONTROL_FLAGS_AUTO     = CameraControlFlags(0X0001)
	KSPROPERTY_CAMERACONTROL_FLAGS_MANUAL   = CameraControlFlags(0X0002)
	KSPROPERTY_CAMERACONTROL_FLAGS_ABSOLUTE = CameraControlFlags(0X0000)
	KSPROPERTY_CAMERACONTROL_FLAGS_RELATIVE = CameraControlFlags(0X0010)
)

type CameraControlProperty int32

const (
	CameraControl_Pan CameraControlProperty = iota
	CameraControl_Tilt
	CameraControl_Roll
	CameraControl_Zoom
	CameraControl_Exposure
	CameraControl_Iris
	CameraControl_Focus
)

type VideoProcAmpFlags int32

const (
	KSPROPERTY_VIDEOPROCAMP_FLAGS_AUTO   VideoProcAmpFlags = 0X0001
	KSPROPERTY_VIDEOPROCAMP_FLAGS_MANUAL VideoProcAmpFlags = 0X0002
)

type VideoProcAmpProperty int32

const (
	VideoProcAmp_Brightness VideoProcAmpProperty = iota
	VideoProcAmp_Contrast
	VideoProcAmp_Hue
	VideoProcAmp_Saturation
	VideoProcAmp_Sharpness
	VideoProcAmp_Gamma
	VideoProcAmp_ColorEnable
	VideoProcAmp_WhiteBalance
	VideoProcAmp_BacklightCompensation
	VideoProcAmp_Gain
)

type MFMEDIASOURCE_CHARACTERISTICS uint32

const (
	MFMEDIASOURCE_IS_LIVE                    = 0x1
	MFMEDIASOURCE_CAN_SEEK                   = 0x2
	MFMEDIASOURCE_CAN_PAUSE                  = 0x4
	MFMEDIASOURCE_HAS_SLOW_SEEK              = 0x8
	MFMEDIASOURCE_HAS_MULTIPLE_PRESENTATIONS = 0x10
	MFMEDIASOURCE_CAN_SKIPFORWARD            = 0x20
	MFMEDIASOURCE_CAN_SKIPBACKWARD           = 0x40
	MFMEDIASOURCE_DOES_NOT_USE_NETWORK       = 0x80
)

const MF_SDK_VERSION = 0x0002
const MF_API_VERSION = 0x0070
const MF_VERSION uint64 = (MF_SDK_VERSION << 16) | MF_API_VERSION
const MFSTARTUP_NOSOCKET uint32 = 0x1
const MFSTARTUP_LITE uint32 = MFSTARTUP_NOSOCKET
const MFSTARTUP_FULL uint32 = 0

type COINIT uint32

const (
	COINIT_APARTMENTTHREADED COINIT = 0x2
	COINIT_MULTITHREADED     COINIT = 0x0
	COINIT_DISABLE_OLE1DDE   COINIT = 0x4
	COINIT_SPEED_OVER_MEMORY COINIT = 0x8
)

//
//var (
//	MAJOR_TYPE_Video MajorType = &guid.MajorTypeVideo
//	MAJOR_TYPE_Audio MajorType = &guid.MajorTypeAudio
//)
//
//var (
//	// The most frequently used
//
//	SubTypeMediaSubTypeNV12 SubType = &guid.SubTypeMediaSubTypeNV12
//	SubTypeMediaSubTypeMJPG SubType = &guid.SubTypeMediaSubTypeMJPG
//	SubTypeMediaSubTypeYUY2 SubType = &guid.SubTypeMediaSubTypeYUY2
//
//	// Other formats
//
//	SubTypeMediaSubTypeYV12             SubType = &guid.SubTypeMediaSubTypeYV12
//	SubTypeMediaSubTypeYUYV             SubType = &guid.SubTypeMediaSubTypeYUYV
//	SubTypeMediaSubTypeIYUV             SubType = &guid.SubTypeMediaSubTypeIYUV
//	SubTypeMediaSubTypeYVYU             SubType = &guid.SubTypeMediaSubTypeYVYU
//	SubTypeMediaSubTypeUYVY             SubType = &guid.SubTypeMediaSubTypeUYVY
//	SubTypeMediaSubTypeRGB24            SubType = &guid.SubTypeMediaSubTypeRGB24
//	SubTypeMediaSubTypeRGB32            SubType = &guid.SubTypeMediaSubTypeRGB32
//	SubTypeMediaSubTypeRGB555           SubType = &guid.SubTypeMediaSubTypeRGB555
//	SubTypeMediaSubTypeRGB565           SubType = &guid.SubTypeMediaSubTypeRGB565
//	SubTypeMediaSubTypeRGB8             SubType = &guid.SubTypeMediaSubTypeRGB8
//	SubTypeMediaSubTypeRGB4             SubType = &guid.SubTypeMediaSubTypeRGB4
//	SubTypeMediaSubTypeRGB1             SubType = &guid.SubTypeMediaSubTypeRGB1
//	SubTypeMediaSubTypeY211             SubType = &guid.SubTypeMediaSubTypeY211
//	SubTypeMediaSubTypeY41P             SubType = &guid.SubTypeMediaSubTypeY41P
//	SubTypeMediaSubTypeY411             SubType = &guid.SubTypeMediaSubTypeY411
//	SubTypeMediaSubTypeYVU9             SubType = &guid.SubTypeMediaSubTypeYVU9
//	SubTypeMediaSubTypeCLJR             SubType = &guid.SubTypeMediaSubTypeCLJR
//	SubTypeMediaSubTypeIF09             SubType = &guid.SubTypeMediaSubTypeIF09
//	SubTypeMediaSubTypeCPLA             SubType = &guid.SubTypeMediaSubTypeCPLA
//	SubTypeMediaSubTypeTVMJ             SubType = &guid.SubTypeMediaSubTypeTVMJ
//	SubTypeMediaSubTypeWAKE             SubType = &guid.SubTypeMediaSubTypeWAKE
//	SubTypeMediaSubTypeCFCC             SubType = &guid.SubTypeMediaSubTypeCFCC
//	SubTypeMediaSubTypeIJPG             SubType = &guid.SubTypeMediaSubTypeIJPG
//	SubTypeMediaSubTypePlum             SubType = &guid.SubTypeMediaSubTypePlum
//	SubTypeMediaSubTypeDVCS             SubType = &guid.SubTypeMediaSubTypeDVCS
//	SubTypeMediaSubTypeH264             SubType = &guid.SubTypeMediaSubTypeH264
//	SubTypeMediaSubTypeDVSD             SubType = &guid.SubTypeMediaSubTypeDVSD
//	SubTypeMediaSubTypeMDVF             SubType = &guid.SubTypeMediaSubTypeMDVF
//	SubTypeMediaSubTypeRGB32D3DDX7RT    SubType = &guid.SubTypeMediaSubTypeRGB32D3DDX7RT
//	SubTypeMediaSubTypeARGB1555         SubType = &guid.SubTypeMediaSubTypeARGB1555
//	SubTypeMediaSubTypeARGB4444         SubType = &guid.SubTypeMediaSubTypeARGB4444
//	SubTypeMediaSubTypeARGB32           SubType = &guid.SubTypeMediaSubTypeARGB32
//	SubTypeMediaSubTypeA2R10G10B10      SubType = &guid.SubTypeMediaSubTypeA2R10G10B10
//	SubTypeMediaSubTypeA2B10G10R10      SubType = &guid.SubTypeMediaSubTypeA2B10G10R10
//	SubTypeMediaSubTypeAYUV             SubType = &guid.SubTypeMediaSubTypeAYUV
//	SubTypeMediaSubTypeAI44             SubType = &guid.SubTypeMediaSubTypeAI44
//	SubTypeMediaSubTypeIA44             SubType = &guid.SubTypeMediaSubTypeIA44
//	SubTypeMediaSubTypeCLPL             SubType = &guid.SubTypeMediaSubTypeCLPL
//	SubTypeMediaSubTypeRGB16D3DDX7RT    SubType = &guid.SubTypeMediaSubTypeRGB16D3DDX7RT
//	SubTypeMediaSubTypeARGB32D3DDX7RT   SubType = &guid.SubTypeMediaSubTypeARGB32D3DDX7RT
//	SubTypeMediaSubTypeARGB4444D3DDX7RT SubType = &guid.SubTypeMediaSubTypeARGB4444D3DDX7RT
//	SubTypeMediaSubTypeARGB1555D3DDX7RT SubType = &guid.SubTypeMediaSubTypeARGB1555D3DDX7RT
//	SubTypeMediaSubTypeRGB32D3DDX9RT    SubType = &guid.SubTypeMediaSubTypeRGB32D3DDX9RT
//	SubTypeMediaSubTypeRGB16D3DDX9RT    SubType = &guid.SubTypeMediaSubTypeRGB16D3DDX9RT
//	SubTypeMediaSubTypeARGB32D3DDX9RT   SubType = &guid.SubTypeMediaSubTypeARGB32D3DDX9RT
//	SubTypeMediaSubTypeARGB4444D3DDX9RT SubType = &guid.SubTypeMediaSubTypeARGB4444D3DDX9RT
//	SubTypeMediaSubTypeARGB1555D3DDX9RT SubType = &guid.SubTypeMediaSubTypeARGB1555D3DDX9RT
//)

type MF_SOURCE_READER_FLAG uint32

const (
	MF_SOURCE_READERF_ERROR                   MF_SOURCE_READER_FLAG = 0x1
	MF_SOURCE_READERF_ENDOFSTREAM             MF_SOURCE_READER_FLAG = 0x2
	MF_SOURCE_READERF_NEWSTREAM               MF_SOURCE_READER_FLAG = 0x4
	MF_SOURCE_READERF_NATIVEMEDIATYPECHANGED  MF_SOURCE_READER_FLAG = 0x10
	MF_SOURCE_READERF_CURRENTMEDIATYPECHANGED MF_SOURCE_READER_FLAG = 0x20
	MF_SOURCE_READERF_STREAMTICK              MF_SOURCE_READER_FLAG = 0x100
	MF_SOURCE_READERF_ALLEFFECTSREMOVED       MF_SOURCE_READER_FLAG = 0x200
)
