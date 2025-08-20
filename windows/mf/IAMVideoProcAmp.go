package mf

import (
	"syscall"
	"unsafe"

	"github.com/Kirizu-Official/windows-camera-go/windows/consts"
)

type IAMVideoProcAmp struct {
	S *IAMVideoProcAmp_
}
type IAMVideoProcAmp_ struct {
	Unknown_
	GetRange uintptr
	Set      uintptr
	Get      uintptr
}
type IAMVideoProcAmpRange struct {
	Min, Max, SteppingDelta, DefaultValue int32
	CapsFlags                             consts.VideoProcAmpFlags
	Found                                 bool
}

func (i *IAMVideoProcAmp) GetRange(property consts.VideoProcAmpProperty) (*IAMVideoProcAmpRange, error) {
	var pMin, pMax, pSteppingDelta, pDefault, pCapsFlags int32

	r1, _, _ := syscall.SyscallN(i.S.GetRange, uintptr(unsafe.Pointer(i)), uintptr(property),
		uintptr(unsafe.Pointer(&pMin)), uintptr(unsafe.Pointer(&pMax)),
		uintptr(unsafe.Pointer(&pSteppingDelta)), uintptr(unsafe.Pointer(&pDefault)), uintptr(unsafe.Pointer(&pCapsFlags)))
	if r1 != consts.S_OK && r1 != consts.ERROR_NOT_FOUND {
		return nil, consts.HResultError{Code: r1}
	}

	return &IAMVideoProcAmpRange{
		Min:           pMin,
		Max:           pMax,
		SteppingDelta: pSteppingDelta,
		DefaultValue:  pDefault,
		CapsFlags:     consts.VideoProcAmpFlags(pCapsFlags),
		Found:         r1 == consts.S_OK,
	}, nil

}
func (i *IAMVideoProcAmp) Set(Property consts.VideoProcAmpProperty, lValue int32, Flags consts.VideoProcAmpFlags) error {
	r1, _, _ := syscall.SyscallN(i.S.Set, uintptr(unsafe.Pointer(i)),
		uintptr(Property), uintptr(lValue), uintptr(Flags))
	if r1 != consts.S_OK && r1 != consts.ERROR_NOT_FOUND {
		return consts.HResultError{Code: r1}
	}

	return nil
}

func (i *IAMVideoProcAmp) Get(Property consts.VideoProcAmpProperty) (value int32, flags consts.VideoProcAmpFlags, err error) {
	var lValue int32
	var lFlags consts.VideoProcAmpFlags
	r1, _, _ := syscall.SyscallN(i.S.Get, uintptr(unsafe.Pointer(i)),
		uintptr(Property), uintptr(unsafe.Pointer(&lValue)), uintptr(unsafe.Pointer(&lFlags)))
	if r1 != consts.S_OK && r1 != consts.ERROR_NOT_FOUND {
		return 0, 0, consts.HResultError{Code: r1}
	}
	return lValue, lFlags, nil
}
func (i *IAMVideoProcAmp) Release() error {
	r1, _, _ := syscall.SyscallN(i.S.Release, uintptr(unsafe.Pointer(i)))
	if r1 != consts.S_OK {
		return consts.HResultError{Code: r1}
	}

	return nil
}
