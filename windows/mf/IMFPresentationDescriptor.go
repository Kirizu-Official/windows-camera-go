package mf

import (
	"github.com/Kirizu-Official/windows-camera-go/windows/consts"
	"syscall"
	"unsafe"
)

type IMFPresentationDescriptor struct {
	S *IMFPresentationDescriptor_
}
type IMFPresentationDescriptor_ struct {
	IMFAttributes_
	GetStreamDescriptorCount   uintptr
	GetStreamDescriptorByIndex uintptr
	SelectStream               uintptr
	DeselectStream             uintptr
	Clone                      uintptr
}

func (i *IMFPresentationDescriptor) GetStreamDescriptorCount() (uint32, error) {
	var count uint32
	r1, _, _ := syscall.SyscallN(i.S.GetStreamDescriptorCount, uintptr(unsafe.Pointer(i)), uintptr(unsafe.Pointer(&count)))
	if r1 != consts.S_OK {
		return 0, consts.HResultError{Code: r1}
	}
	return count, nil
}
func (i *IMFPresentationDescriptor) GetStreamDescriptorByIndex(dwIndex uint32, pfSelected *bool, ppDescriptor **IMFStreamDescriptor) error {
	r1, _, _ := syscall.SyscallN(i.S.GetStreamDescriptorByIndex, uintptr(unsafe.Pointer(i)), uintptr(dwIndex), uintptr(unsafe.Pointer(pfSelected)), uintptr(unsafe.Pointer(ppDescriptor)))
	if r1 != consts.S_OK {
		return consts.HResultError{Code: r1}
	}
	return nil
}
func (i *IMFPresentationDescriptor) Release() error {
	r1, _, _ := syscall.SyscallN(i.S.Release, uintptr(unsafe.Pointer(i)))
	if r1 != consts.S_OK {
		return consts.HResultError{Code: r1}
	}

	return nil
}
