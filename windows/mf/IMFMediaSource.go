package mf

import (
	"github.com/Kirizu-Official/windows-camera-go/windows/consts"
	"syscall"
	"unsafe"
)

type IMFMediaSource struct {
	S *IMFMediaSource_
}
type IMFMediaSource_ struct {
	IMFMediaEventGenerator_
	GetCharacteristics           uintptr
	CreatePresentationDescriptor uintptr
	Start                        uintptr
	Stop                         uintptr
	Pause                        uintptr
	Shutdown                     uintptr
}
type MediaSourceCharacteristics struct {
	MfMediaSourceIsLive                   bool
	MfMediaSourceCanSeek                  bool
	MfMediaSourceCanPause                 bool
	MfMediaSourceHasSlowSeek              bool
	MfMediaSourceHasMultiplePresentations bool
	MfMediaSourceCanSkipForward           bool
	MfMediaSourceCanSkipBackward          bool
	MfMediaSourceDoesNotUseNetwork        bool
}

func (i *IMFMediaSource) CreatePresentationDescriptor(pPresentationDescriptor **IMFPresentationDescriptor) error {
	r1, _, _ := syscall.SyscallN(i.S.CreatePresentationDescriptor, uintptr(unsafe.Pointer(&i.S)), uintptr(unsafe.Pointer(pPresentationDescriptor)))
	if r1 != consts.S_OK {
		return consts.HResultError{Code: r1}
	}
	return nil
}
func (i *IMFMediaSource) GetCharacteristics() (*MediaSourceCharacteristics, error) {
	var characteristics consts.MFMEDIASOURCE_CHARACTERISTICS
	r1, _, _ := syscall.SyscallN(i.S.GetCharacteristics, uintptr(unsafe.Pointer(i)), uintptr(unsafe.Pointer(&characteristics)))
	if r1 != consts.S_OK {
		return nil, consts.HResultError{Code: r1}
	}
	return &MediaSourceCharacteristics{
		MfMediaSourceIsLive:                   characteristics&consts.MFMEDIASOURCE_IS_LIVE != 0,
		MfMediaSourceCanSeek:                  characteristics&consts.MFMEDIASOURCE_CAN_SEEK != 0,
		MfMediaSourceCanPause:                 characteristics&consts.MFMEDIASOURCE_CAN_PAUSE != 0,
		MfMediaSourceHasSlowSeek:              characteristics&consts.MFMEDIASOURCE_HAS_SLOW_SEEK != 0,
		MfMediaSourceHasMultiplePresentations: characteristics&consts.MFMEDIASOURCE_HAS_MULTIPLE_PRESENTATIONS != 0,
		MfMediaSourceCanSkipForward:           characteristics&consts.MFMEDIASOURCE_CAN_SKIPFORWARD != 0,
		MfMediaSourceCanSkipBackward:          characteristics&consts.MFMEDIASOURCE_CAN_SKIPBACKWARD != 0,
		MfMediaSourceDoesNotUseNetwork:        characteristics&consts.MFMEDIASOURCE_DOES_NOT_USE_NETWORK != 0,
	}, nil
}
func (i *IMFMediaSource) Pause() error {
	r1, _, _ := syscall.SyscallN(i.S.Pause, uintptr(unsafe.Pointer(i)))
	if r1 != consts.S_OK {
		return consts.HResultError{Code: r1}
	}
	return nil
}
func (i *IMFMediaSource) Shutdown() error {
	r1, _, _ := syscall.SyscallN(i.S.Shutdown, uintptr(unsafe.Pointer(i)))
	if r1 != consts.S_OK {
		return consts.HResultError{Code: r1}
	}
	return nil
}

func (i *IMFMediaSource) Stop() error {
	r1, _, _ := syscall.SyscallN(i.S.Stop, uintptr(unsafe.Pointer(i)))
	if r1 != consts.S_OK {
		return consts.HResultError{Code: r1}
	}
	return nil
}
func (i *IMFMediaSource) Release() error {
	r1, _, _ := syscall.SyscallN(i.S.Release, uintptr(unsafe.Pointer(i)))
	if r1 != consts.S_OK {
		return consts.HResultError{Code: r1}
	}

	return nil
}
