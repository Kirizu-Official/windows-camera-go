package mf

import (
	"github.com/Kirizu-Official/windows-camera-go/windows/consts"
	"syscall"
	"unsafe"
)

type IMFMediaTypeHandler struct {
	S *IMFMediaTypeHandler_
}
type IMFMediaTypeHandler_ struct {
	Unknown_
	IsMediaTypeSupported uintptr
	GetMediaTypeCount    uintptr
	GetMediaTypeByIndex  uintptr
	SetCurrentMediaType  uintptr
	GetCurrentMediaType  uintptr
	GetMajorType         uintptr
}

func (i *IMFMediaTypeHandler) GetMediaTypeByIndex(dwIndex uint32, ppType **IMFMediaType) error {
	r1, _, _ := syscall.SyscallN(i.S.GetMediaTypeByIndex, uintptr(unsafe.Pointer(i)), uintptr(dwIndex), uintptr(unsafe.Pointer(ppType)))
	if r1 != consts.S_OK {
		return consts.HResultError{Code: r1}
	}
	return nil
}
func (i *IMFMediaTypeHandler) GetMediaTypeCount() (uint32, error) {
	var count uint32
	r1, _, _ := syscall.SyscallN(i.S.GetMediaTypeCount, uintptr(unsafe.Pointer(i)), uintptr(unsafe.Pointer(&count)))
	if r1 != consts.S_OK {
		return 0, consts.HResultError{Code: r1}
	}
	return count, nil
}

func (i *IMFMediaTypeHandler) Release() error {
	r1, _, _ := syscall.SyscallN(i.S.Release, uintptr(unsafe.Pointer(i)))
	if r1 != consts.S_OK {
		return consts.HResultError{Code: r1}
	}

	return nil
}
