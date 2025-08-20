package mf

import (
	"syscall"
	"unsafe"

	"github.com/Kirizu-Official/windows-camera-go/windows/consts"
)

type MfReadWrite struct {
	mfreadwriteDll                      *syscall.DLL
	mfCreateSourceReaderFromMediaSource *syscall.Proc
}

func CreateNewMfReadWrite() (*MfReadWrite, error) {
	loadRes, err := syscall.LoadDLL("Mfreadwrite.dll")
	if err != nil {
		return nil, err
	}
	newMfReadWrite := &MfReadWrite{
		mfreadwriteDll: loadRes,
	}

	proc, err := newMfReadWrite.mfreadwriteDll.FindProc("MFCreateSourceReaderFromMediaSource")
	if err == nil {
		newMfReadWrite.mfCreateSourceReaderFromMediaSource = proc
	}

	return newMfReadWrite, nil
}

func (m *MfReadWrite) MFCreateSourceReaderFromMediaSource(pMediaSource *IMFMediaSource, pAttributes *IMFAttributes, ppSourceReader **IMFSourceReader) error {
	code, _, _ := m.mfCreateSourceReaderFromMediaSource.Call(
		uintptr(unsafe.Pointer(pMediaSource)),
		uintptr(unsafe.Pointer(pAttributes)),
		uintptr(unsafe.Pointer(ppSourceReader)),
	)

	if code != consts.S_OK {
		return &consts.HResultError{Code: code}
	}
	return nil
}
