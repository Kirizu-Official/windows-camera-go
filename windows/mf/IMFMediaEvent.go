package mf

import (
	"github.com/Kirizu-Official/windows-camera-go/windows/consts"
	"syscall"
	"unsafe"
)

type IMFMediaEvent struct {
	S *IMFMediaEvent_
}
type IMFMediaEvent_ struct {
	IMFAttributes
	GetType         uintptr
	GetExtendedType uintptr
	GetStatus       uintptr
	GetValue        uintptr
}

func (i *IMFMediaEvent) Release() error {
	r1, _, _ := syscall.SyscallN(i.S.IMFAttributes.S.Release, uintptr(unsafe.Pointer(i)))
	if r1 != consts.S_OK {
		return consts.HResultError{Code: r1}
	}

	return nil
}
