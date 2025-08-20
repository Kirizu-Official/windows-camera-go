package mf

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/Kirizu-Official/windows-camera-go/windows/consts"
)

type IAMCameraControl struct {
	S *IAMCameraControl_
}
type IAMCameraControl_ struct {
	Unknown_
	GetRange uintptr
	Set      uintptr
	Get      uintptr
}
type IAMCameraControlRange struct {
	Min, Max, SteppingDelta, DefaultValue int32
	CapsFlags                             consts.CameraControlFlags
	Found                                 bool
}

func (i *IAMCameraControl) GetRange(property consts.CameraControlProperty) (*IAMCameraControlRange, error) {
	fmt.Println(i.S, property)
	var pMin, pMax, pSteppingDelta, pDefault, pCapsFlags int32

	r1, _, _ := syscall.SyscallN(i.S.GetRange, uintptr(unsafe.Pointer(i)), uintptr(property),
		uintptr(unsafe.Pointer(&pMin)), uintptr(unsafe.Pointer(&pMax)),
		uintptr(unsafe.Pointer(&pSteppingDelta)), uintptr(unsafe.Pointer(&pDefault)), uintptr(unsafe.Pointer(&pCapsFlags)))
	if r1 != consts.S_OK && r1 != consts.ERROR_NOT_FOUND {
		return nil, consts.HResultError{Code: r1}
	}

	return &IAMCameraControlRange{
		Min:           pMin,
		Max:           pMax,
		SteppingDelta: pSteppingDelta,
		DefaultValue:  pDefault,
		CapsFlags:     consts.CameraControlFlags(pCapsFlags),
		Found:         r1 == consts.S_OK,
	}, nil

}
func (i *IAMCameraControl) Set(Property consts.CameraControlProperty, lValue int32, Flags consts.CameraControlFlags) error {
	r1, _, _ := syscall.SyscallN(i.S.Set, uintptr(unsafe.Pointer(i)),
		uintptr(Property), uintptr(lValue), uintptr(Flags))
	if r1 != consts.S_OK && r1 != consts.ERROR_NOT_FOUND {
		return consts.HResultError{Code: r1}
	}

	return nil
}

func (i *IAMCameraControl) Get(Property consts.CameraControlProperty) (value int32, flags consts.CameraControlFlags, err error) {
	var lValue int32
	var lFlags consts.CameraControlFlags
	r1, _, _ := syscall.SyscallN(i.S.Get, uintptr(unsafe.Pointer(i)),
		uintptr(Property), uintptr(unsafe.Pointer(&lValue)), uintptr(unsafe.Pointer(&lFlags)))
	if r1 != consts.S_OK && r1 != consts.ERROR_NOT_FOUND {
		return 0, 0, consts.HResultError{Code: r1}
	}
	return lValue, lFlags, nil
}
func (i *IAMCameraControl) Release() error {
	r1, _, _ := syscall.SyscallN(i.S.Release, uintptr(unsafe.Pointer(i)))
	if r1 != consts.S_OK {
		return consts.HResultError{Code: r1}
	}

	return nil
}
