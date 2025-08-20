package mf

import (
	"syscall"
	"unsafe"

	"github.com/Kirizu-Official/windows-camera-go/windows/consts"
	"github.com/Kirizu-Official/windows-camera-go/windows/guid"
)

type IMFSourceReader struct {
	S *IMFSourceReader_
}
type IMFSourceReader_ struct {
	Unknown_
	GetStreamSelection       uintptr
	SetStreamSelection       uintptr
	GetNativeMediaType       uintptr
	GetCurrentMediaType      uintptr
	SetCurrentMediaType      uintptr
	SetCurrentPosition       uintptr
	ReadSample               uintptr
	Flush                    uintptr
	GetServiceForStream      uintptr
	GetPresentationAttribute uintptr
}

func (i *IMFSourceReader) SetCurrentMediaType(dwStreamIndex uint32, pType *IMFMediaType) error {
	r1, _, _ := syscall.SyscallN(i.S.SetCurrentMediaType, uintptr(unsafe.Pointer(i)), uintptr(dwStreamIndex), uintptr(unsafe.Pointer(nil)), uintptr(unsafe.Pointer(pType)))
	if r1 != consts.S_OK {
		return consts.HResultError{Code: r1}
	}
	return nil
}
func (i *IMFSourceReader) ReadSample(dwStreamIndex uint32, dwControlFlags uint32, pdwActualStreamIndex *uint32, pdwStreamFlags *uint32, pllTimestamp *int64, ppSample **IMFSample) error {
	r1, _, _ := syscall.SyscallN(i.S.ReadSample, uintptr(unsafe.Pointer(i)), uintptr(dwStreamIndex), uintptr(dwControlFlags),
		uintptr(unsafe.Pointer(pdwActualStreamIndex)), uintptr(unsafe.Pointer(pdwStreamFlags)), uintptr(unsafe.Pointer(pllTimestamp)), uintptr(unsafe.Pointer(ppSample)))
	if r1 != consts.S_OK {
		return consts.HResultError{Code: r1}
	}

	return nil
}

//func (i *IMFSourceReader) ReadSampleASync(dwStreamIndex uint32, dwControlFlags uint32) error {
//
//	fmt.Println("nil value", uintptr(unsafe.Pointer(nil)))
//	r1, _, _ := syscall.SyscallN(i.S.ReadSample, uintptr(unsafe.Pointer(i)), uintptr(dwStreamIndex), uintptr(dwControlFlags),
//		uintptr(unsafe.Pointer(nil)), uintptr(unsafe.Pointer(nil)), uintptr(unsafe.Pointer(nil)), uintptr(unsafe.Pointer(nil)))
//	if r1 != consts.S_OK {
//		return consts.HResultError{Code: r1}
//	}
//	return nil
//}

func (i *IMFSourceReader) GetServiceForStream() (*IAMCameraControl, *IAMVideoProcAmp, error) {
	var control *IAMCameraControl

	r1, _, _ := syscall.SyscallN(i.S.GetServiceForStream, uintptr(unsafe.Pointer(i)), uintptr(0xFFFFFFFF), uintptr(unsafe.Pointer(&guid.NULL)), uintptr(unsafe.Pointer(&guid.IID_IAMCameraControl)), uintptr(unsafe.Pointer(&control)))
	if r1 != consts.S_OK {
		return nil, nil, consts.HResultError{Code: r1}
	}

	var videoProc *IAMVideoProcAmp
	r1, _, _ = syscall.SyscallN(i.S.GetServiceForStream, uintptr(unsafe.Pointer(i)), uintptr(0xFFFFFFFF), uintptr(unsafe.Pointer(&guid.NULL)), uintptr(unsafe.Pointer(&guid.IID_IAMVideoProcAmp)), uintptr(unsafe.Pointer(&videoProc)))
	if r1 != consts.S_OK {
		//panic(r1)
		if control != nil {
			control.Release()
		}
		return nil, nil, consts.HResultError{Code: r1}
	}
	//err = service.GetExtendedCameraControl()
	//if err != nil {
	//	panic(err)
	//}
	return control, videoProc, nil
}

func (i *IMFSourceReader) Release() error {
	r1, _, _ := syscall.SyscallN(i.S.Release, uintptr(unsafe.Pointer(i)))
	if r1 != consts.S_OK {
		return consts.HResultError{Code: r1}
	}

	return nil
}
