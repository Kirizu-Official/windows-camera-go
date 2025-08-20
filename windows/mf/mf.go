package mf

import "C"
import (
	"syscall"
	"unsafe"

	"github.com/Kirizu-Official/windows-camera-go/utils"
	"github.com/Kirizu-Official/windows-camera-go/windows/consts"
)

type MF struct {
	mf *syscall.DLL

	mfEnumDeviceSources *syscall.Proc
}

func CreateNewMF() (*MF, error) {
	loadRes, err := syscall.LoadDLL("Mf.dll")
	if err != nil {
		return nil, err
	}
	newMfApi := &MF{
		mf: loadRes,
	}

	proc, err := newMfApi.mf.FindProc("MFEnumDeviceSources")
	if err == nil {
		newMfApi.mfEnumDeviceSources = proc
	}

	return newMfApi, nil
}
func (mf *MF) MFEnumDeviceSources(pAttributes *IMFAttributes, pcSourceActivate *uint32) ([]*IMFActivate, error) {

	var listPtr uintptr
	call, _, err := mf.mfEnumDeviceSources.Call(uintptr(unsafe.Pointer(pAttributes)), uintptr(unsafe.Pointer(&listPtr)), uintptr(unsafe.Pointer(pcSourceActivate)))
	err = utils.CheckError(err)
	if err != nil {
		return nil, err
	}
	if call != consts.S_OK {
		return nil, consts.HResultError{Code: call}
	}

	var pppSourceActivate []*IMFActivate
	for i := uint32(0); i < *pcSourceActivate; i++ {
		arrayIthElementOffset := uintptr(i) * unsafe.Sizeof(uintptr(0))
		arrayIthElementAddress := listPtr + arrayIthElementOffset
		arrayIthElementPointer := (*uintptr)(unsafe.Pointer(arrayIthElementAddress))
		pppSourceActivate = append(pppSourceActivate, (*IMFActivate)(unsafe.Pointer(*arrayIthElementPointer)))
	}
	return pppSourceActivate, nil

}
