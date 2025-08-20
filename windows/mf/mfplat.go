package mf

import (
	"syscall"
	"unsafe"

	"github.com/Kirizu-Official/windows-camera-go/utils"
	"github.com/Kirizu-Official/windows-camera-go/windows/consts"
)

type MfApi struct {
	mfPlatDll          *syscall.DLL
	mfCreateAttributes *syscall.Proc
	mFStartup          *syscall.Proc
	mfShutdown         *syscall.Proc
}

func CreateNewMfPlat() (*MfApi, error) {
	loadRes, err := syscall.LoadDLL("Mfplat.dll")
	if err != nil {
		return nil, err
	}
	newMfApi := &MfApi{
		mfPlatDll: loadRes,
	}

	proc, err := newMfApi.mfPlatDll.FindProc("MFCreateAttributes")
	if err == nil {
		newMfApi.mfCreateAttributes = proc
	}

	proc, err = newMfApi.mfPlatDll.FindProc("MFStartup")
	if err == nil {
		newMfApi.mFStartup = proc
	}

	proc, err = newMfApi.mfPlatDll.FindProc("MFShutdown")
	if err == nil {
		newMfApi.mfShutdown = proc
	}

	return newMfApi, nil
}

func (m *MfApi) MFCreateAttributes(ppMFAttributes **IMFAttributes, cInitialSize uint32) error {

	code, _, err := m.mfCreateAttributes.Call(uintptr(unsafe.Pointer(ppMFAttributes)), uintptr(cInitialSize))
	err = utils.CheckError(err)
	if err != nil {
		return err
	}

	if code != consts.S_OK {
		return &consts.HResultError{Code: code}
	}
	return nil
}

func (m *MfApi) MFStartup() error {
	code, _, err := m.mFStartup.Call(uintptr(consts.MF_VERSION), uintptr(consts.MFSTARTUP_NOSOCKET))
	err = utils.CheckError(err)
	if err != nil {
		return err
	}

	if code != consts.S_OK {
		return &consts.HResultError{Code: code}
	}
	return nil
}
func (m *MfApi) MFShutdown() error {
	code, _, err := m.mfShutdown.Call()
	err = utils.CheckError(err)
	if err != nil {
		return err
	}

	if code != consts.S_OK {
		return &consts.HResultError{Code: code}
	}
	return nil
}
