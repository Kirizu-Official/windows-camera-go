package mf

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/Kirizu-Official/windows-camera-go/windows/consts"
	"github.com/Kirizu-Official/windows-camera-go/windows/guid"
)

type IMFActivate struct {
	S *IMFActivate_
}
type IMFActivate_ struct {
	IMFAttributes_
	ActivateObject uintptr
	ShutdownObject uintptr
	DetachObject   uintptr
}

func (i *IMFActivate) ActivateObject(pSource **IMFMediaSource) error {

	r1, _, _ := syscall.SyscallN(i.S.ActivateObject, uintptr(unsafe.Pointer(i)), uintptr(unsafe.Pointer(&guid.MF_MEDIASOURCE_SERVICE)), uintptr(unsafe.Pointer(pSource)))
	if r1 != consts.S_OK {
		return consts.HResultError{Code: r1}
	}
	return nil
}

func (i *IMFActivate) GetString(guidKey *guid.GUID, pName uintptr, maxReceiveLen uint32) (uint32, error) {
	var readLen uint32
	r1, _, _ := syscall.SyscallN(i.S.GetString, uintptr(unsafe.Pointer(i)), uintptr(unsafe.Pointer(guidKey)), pName, uintptr(maxReceiveLen), uintptr(unsafe.Pointer(&readLen)))
	if r1 != consts.S_OK {
		return 0, consts.HResultError{Code: r1}
	}

	return readLen, nil
}

func (i *IMFActivate) GetStrLength(guidKey *guid.GUID) (uint32, error) {
	var length uint32
	r1, _, _ := syscall.SyscallN(i.S.GetStringLength, uintptr(unsafe.Pointer(i)), uintptr(unsafe.Pointer(guidKey)), uintptr(unsafe.Pointer(&length)))
	if r1 != consts.S_OK {
		fmt.Printf("GetStrLength failed: %d\n", r1)
		return 0, consts.HResultError{Code: r1}
	}
	return length, nil
}

func (i *IMFActivate) Release() error {
	r1, _, _ := syscall.SyscallN(i.S.Release, uintptr(unsafe.Pointer(i)))
	if r1 != consts.S_OK {
		return consts.HResultError{Code: r1}
	}

	return nil
}
