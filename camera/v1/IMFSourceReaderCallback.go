package camera

import (
	"syscall"
	"unsafe"

	"github.com/Kirizu-Official/windows-camera-go/windows/consts"
	"github.com/Kirizu-Official/windows-camera-go/windows/guid"
	"github.com/Kirizu-Official/windows-camera-go/windows/mf"
)

var IID_IMFSourceReaderCallback = &guid.GUID{
	Data1: 0xdeec8d99,
	Data2: 0xfa1d,
	Data3: 0x4d82,
	Data4: [8]byte{0x84, 0xc2, 0x2c, 0x89, 0x69, 0x94, 0x48, 0x67},
}

type IMFSourceReaderCallbackVtbl struct {
	// IUnknown 方法
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr

	// IMFSourceReaderCallback 方法
	OnReadSample uintptr
	OnFlush      uintptr
	OnEvent      uintptr
}

type IMFSourceReaderCallback struct {
	Vtbl *IMFSourceReaderCallbackVtbl
}

type SourceReaderCallback struct {
	vtbl     IMFSourceReaderCallbackVtbl
	intf     IMFSourceReaderCallback
	refCount int32
	this     *CaptureAsync
	// 回调函数
	OnReadSampleFunc func(this *CaptureAsync, dwStreamIndex uint32, dwStreamFlags uint32, llTimestamp int64, pSample *mf.IMFSample)
	OnFlushFunc      func(this *CaptureAsync, dwStreamIndex uint32)
	OnEventFunc      func(this *CaptureAsync, dwStreamIndex uint32, pEvent *mf.IMFMediaEvent)
}

func newSourceReaderCallback(this *CaptureAsync) *SourceReaderCallback {
	callback := &SourceReaderCallback{
		this:     this,
		refCount: 1,
	}

	callback.vtbl = IMFSourceReaderCallbackVtbl{
		QueryInterface: syscall.NewCallback(callback.queryInterface),
		AddRef:         syscall.NewCallback(callback.addRef),
		Release:        syscall.NewCallback(callback.release),
		OnReadSample:   syscall.NewCallback(callback.onReadSample),
		OnFlush:        syscall.NewCallback(callback.onFlush),
		OnEvent:        syscall.NewCallback(callback.onEvent),
	}

	callback.intf = IMFSourceReaderCallback{
		Vtbl: &callback.vtbl,
	}

	return callback
}

func (c *SourceReaderCallback) GetInterface() *IMFSourceReaderCallback {
	return &c.intf
}

// SetOnReadSample set the callback function that will be called when a sample is read from the camera when use async.
//
// this: *camera.CaptureAsync is the async capture instance which from.
//
// dwStreamIndex: The zero-based index of the stream that delivered the sample.
//
// dwStreamFlags: A bitwise OR of zero or more flags from the consts.MF_SOURCE_READER_FLAG enumeration.
//
// llTimestamp: The time stamp of the sample, or the time of the stream event indicated in dwStreamFlags. The time is given in 100-nanosecond units.
//
// pSample: DO NOT RELEASE IT in callback! the sample that contains the buffer data, it will never be nil
//
// For more information, see https://learn.microsoft.com/en-us/windows/win32/api/mfreadwrite/nf-mfreadwrite-imfsourcereadercallback-onreadsample
func (c *SourceReaderCallback) SetOnReadSample(fn func(this *CaptureAsync, dwStreamIndex uint32, dwStreamFlags uint32, llTimestamp int64, pSample *mf.IMFSample)) {
	c.OnReadSampleFunc = fn
}

func (c *SourceReaderCallback) SetOnFlush(fn func(this *CaptureAsync, dwStreamIndex uint32)) {
	c.OnFlushFunc = fn
}

func (c *SourceReaderCallback) SetOnEvent(fn func(this *CaptureAsync, dwStreamIndex uint32, pEvent *mf.IMFMediaEvent)) {
	c.OnEventFunc = fn
}

// IUnknown 方法实现
func (c *SourceReaderCallback) queryInterface(this uintptr, riidIn *syscall.GUID, ppvObject *uintptr) uintptr {
	if riidIn == nil || ppvObject == nil {
		return 0x80004003
	}
	riid := &guid.GUID{
		Data1: riidIn.Data1,
		Data2: riidIn.Data2,
		Data3: riidIn.Data3,
		Data4: riidIn.Data4,
	}

	*ppvObject = 0

	// 检查请求的接口
	if riid.IsMatch(IID_IMFSourceReaderCallback) {
		*ppvObject = this
		c.addRef(this)
		return 0
	}

	// IUnknown GUID
	iidUnknown := guid.GUID{
		Data1: 0x00000000,
		Data2: 0x0000,
		Data3: 0x0000,
		Data4: [8]byte{0xc0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x46},
	}

	if *riid == iidUnknown {
		*ppvObject = this
		c.addRef(this)
		return 0 // S_OK
	}

	return 0x80004002
}

func (c *SourceReaderCallback) addRef(this uintptr) uintptr {
	c.refCount++
	return uintptr(c.refCount)
}

func (c *SourceReaderCallback) release(this uintptr) uintptr {
	c.refCount--
	return uintptr(c.refCount)
}

// IMFSourceReaderCallback 方法实现
func (c *SourceReaderCallback) onReadSample(this uintptr, hrStatus uint32, dwStreamIndex uint32, dwStreamFlags uint32, llTimestamp int64, pSample uintptr) uintptr {
	HResult := consts.HResultError{Code: consts.HRESULT(hrStatus)}
	if HResult.Code != consts.S_OK {
		//panic(HResult.Code)
		c.this.Device.SetError(HResult)
		return uintptr(hrStatus)
	}

	if pSample == 0 {
		err := c.this.GetNextFrame()
		if err != nil {
			c.this.Device.SetError(err)
			//panic(err)
			return consts.S_OK
		}
		return uintptr(hrStatus)
	}

	if c.OnReadSampleFunc != nil {
		sample := (*mf.IMFSample)(unsafe.Pointer(pSample))
		if sample.S == nil {
			return consts.S_OK
		}

		c.OnReadSampleFunc(c.this, dwStreamIndex, dwStreamFlags, llTimestamp, sample)
	}

	return consts.S_OK
}

func (c *SourceReaderCallback) onFlush(this uintptr, dwStreamIndex uint32) uintptr {
	if c.OnFlushFunc != nil {
		c.OnFlushFunc(c.this, dwStreamIndex)
	}
	return consts.S_OK
}

func (c *SourceReaderCallback) onEvent(this uintptr, dwStreamIndex uint32, pEvent uintptr) uintptr {

	if c.OnEventFunc != nil && pEvent != 0 {
		event := (*mf.IMFMediaEvent)(unsafe.Pointer(pEvent))
		c.OnEventFunc(c.this, dwStreamIndex, event)
		_ = event.Release()
	}
	return consts.S_OK
}

func (c *SourceReaderCallback) ToUintptr() uintptr {
	return uintptr(unsafe.Pointer(c.GetInterface()))
}
