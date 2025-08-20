package mf

import (
	"github.com/Kirizu-Official/windows-camera-go/windows/consts"
	"syscall"
	"unsafe"
)

type IMFExtendedCameraController struct {
	S *IMFExtendedCameraController_
}
type IMFExtendedCameraController_ struct {
	Unknown_
	GetExtendedCameraControl uintptr
}

func (i *IMFExtendedCameraController) GetExtendedCameraControl() error {
	var control *IMFExtendedCameraControl
	r1, _, _ := syscall.SyscallN(i.S.GetExtendedCameraControl, uintptr(unsafe.Pointer(i)), uintptr(0), uintptr(14), uintptr(unsafe.Pointer(&control)))
	if r1 != consts.S_OK {
		return consts.HResultError{Code: r1}
	}
	return nil

}
