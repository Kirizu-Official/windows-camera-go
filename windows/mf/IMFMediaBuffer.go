package mf

import "C"
import (
	"github.com/Kirizu-Official/windows-camera-go/windows/consts"
	"syscall"
	"unsafe"
)

type IMFMediaBuffer struct {
	S *IMFMediaBuffer_
}
type IMFMediaBuffer_ struct {
	Unknown_
	Lock             uintptr
	Unlock           uintptr
	GetCurrentLength uintptr
	SetCurrentLength uintptr
	GetMaxLength     uintptr
}

func (m *IMFMediaBuffer) Lock() (data []byte, maxCount uint32, currentCount uint32, err error) {

	var pcbMaxLength uint32
	var pcbCurrentLength uint32
	var ppbData uintptr

	r1, _, _ := syscall.SyscallN(m.S.Lock, uintptr(unsafe.Pointer(m)), uintptr(unsafe.Pointer(&ppbData)), uintptr(unsafe.Pointer(&pcbMaxLength)), uintptr(unsafe.Pointer(&pcbCurrentLength)))

	if r1 != consts.S_OK {
		return nil, 0, 0, consts.HResultError{Code: r1}
	}

	data = C.GoBytes(unsafe.Pointer(ppbData), C.int(pcbMaxLength))
	return data, pcbMaxLength, pcbCurrentLength, nil
}

func (m *IMFMediaBuffer) Unlock() error {
	r1, _, _ := syscall.SyscallN(m.S.Unlock, uintptr(unsafe.Pointer(m)))
	if r1 != consts.S_OK {
		return consts.HResultError{Code: r1}
	}
	return nil
}

func (m *IMFMediaBuffer) Release() error {
	r1, _, _ := syscall.SyscallN(m.S.Release, uintptr(unsafe.Pointer(m)))
	if r1 != consts.S_OK {
		return consts.HResultError{Code: r1}
	}

	return nil
}
