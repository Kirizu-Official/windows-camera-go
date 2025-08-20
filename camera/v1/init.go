package camera

import (
	"fmt"
	"syscall"

	"github.com/Kirizu-Official/windows-camera-go/windows/consts"
	"github.com/Kirizu-Official/windows-camera-go/windows/mf"
)

var mfApi *mf.MF
var mfPlat *mf.MfApi
var mfReadWrite *mf.MfReadWrite

func Init() error {
	var err error
	mfApi, err = mf.CreateNewMF()
	if err != nil {
		return err
	}
	mfPlat, err = mf.CreateNewMfPlat()
	if err != nil {
		return err
	}
	mfReadWrite, err = mf.CreateNewMfReadWrite()
	if err != nil {
		return err
	}
	err = initWindowsCOM()
	if err != nil {
		return err
	}

	return nil
}

var ole32dll *syscall.DLL

func initWindowsCOM() error {
	var err error
	ole32dll, err = syscall.LoadDLL("Ole32.dll")
	if err != nil {
		return fmt.Errorf("failed to load Ole32.dll: %w", err)
	}
	proc, err := ole32dll.FindProc("CoInitializeEx")
	if err != nil {
		return err
	}
	r, _, _ := proc.Call(0, uintptr(consts.COINIT_MULTITHREADED))
	if r != consts.S_OK {
		return fmt.Errorf("CoInitialize failed: %d", r)
	}
	err = mfPlat.MFStartup()
	if err != nil {
		return err
	}

	return nil
}

func Shutdown() {

	proc, err := ole32dll.FindProc("CoUninitialize")
	if err == nil && proc != nil {
		r, _, _ := proc.Call(0)
		if r != consts.S_OK {
			fmt.Printf("[WARN] GoCamera: CoInitialize failed: %d\n", r)
		}
	} else {
		fmt.Println("[WARN] GoCamera: CoUninitialize failed.")
	}

	err = mfPlat.MFShutdown()
	if err != nil {
		fmt.Println("[WARN] GoCamera: MFShutdown failed.")
	}
}
