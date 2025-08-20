package mf

import (
	"github.com/Kirizu-Official/windows-camera-go/windows/consts"
	"syscall"
	"unsafe"
)

type IMFStreamDescriptor struct {
	S *IMFStreamDescriptor_
}
type IMFStreamDescriptor_ struct {
	IMFAttributes_
	GetStreamIdentifier uintptr
	GetMediaTypeHandler uintptr
}

func (s *IMFStreamDescriptor) GetMediaTypeHandler(ppHandler **IMFMediaTypeHandler) error {
	r1, _, _ := syscall.SyscallN(s.S.GetMediaTypeHandler, uintptr(unsafe.Pointer(s)), uintptr(unsafe.Pointer(ppHandler)))
	if r1 != consts.S_OK {
		return consts.HResultError{Code: r1}
	}
	return nil

}

func (s *IMFStreamDescriptor) Release() error {
	r1, _, _ := syscall.SyscallN(s.S.Release, uintptr(unsafe.Pointer(s)))
	if r1 != consts.S_OK {
		return consts.HResultError{Code: r1}
	}

	return nil
}
