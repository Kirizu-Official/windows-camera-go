package mf

import (
	"syscall"
	"unsafe"

	"github.com/Kirizu-Official/windows-camera-go/windows/consts"
)

type IMFSample struct {
	S *IMFSample_
}
type IMFSample_ struct {
	IMFAttributes_
	GetSampleFlags            uintptr
	SetSampleFlags            uintptr
	GetSampleTime             uintptr
	SetSampleTime             uintptr
	GetSampleDuration         uintptr
	SetSampleDuration         uintptr
	GetBufferCount            uintptr
	GetBufferByIndex          uintptr
	ConvertToContiguousBuffer uintptr
	AddBuffer                 uintptr
	RemoveBufferByIndex       uintptr
	RemoveAllBuffers          uintptr
	GetTotalLength            uintptr
	CopyToBuffer              uintptr
}

func (s *IMFSample) ConvertToContiguousBuffer() (*IMFMediaBuffer, error) {

	var buffer *IMFMediaBuffer
	r1, _, _ := syscall.SyscallN(s.S.ConvertToContiguousBuffer, uintptr(unsafe.Pointer(s)), uintptr(unsafe.Pointer(&buffer)))
	if r1 != consts.S_OK {
		return nil, consts.HResultError{Code: r1}
	}

	return buffer, nil
}
func (s *IMFSample) GetTotalLength() error {
	var length uint32
	r1, _, _ := syscall.SyscallN(s.S.GetTotalLength, uintptr(unsafe.Pointer(s)), uintptr(unsafe.Pointer(&length)))
	if r1 != consts.S_OK {
		return consts.HResultError{Code: r1}
	}
	return nil

}
func (s *IMFSample) GetBufferCount() error {
	var count uint32
	r1, _, _ := syscall.SyscallN(s.S.GetBufferCount, uintptr(unsafe.Pointer(s)), uintptr(unsafe.Pointer(&count)))
	if r1 != consts.S_OK {
		return consts.HResultError{Code: r1}
	}
	return nil

}

func (s *IMFSample) Release() error {
	if s == nil {
		return nil
	}
	r1, _, _ := syscall.SyscallN(s.S.Release, uintptr(unsafe.Pointer(s)))
	if r1 != consts.S_OK {
		return consts.HResultError{Code: r1}
	}

	return nil
}
