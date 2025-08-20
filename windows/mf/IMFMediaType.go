package mf

import (
	"syscall"
	"unsafe"

	"github.com/Kirizu-Official/windows-camera-go/windows/consts"
	"github.com/Kirizu-Official/windows-camera-go/windows/guid"
)

type IMFMediaType struct {
	S *IMFMediaType_
}
type IMFMediaType_ struct {
	IMFAttributes_
	GetMajorType       uintptr
	IsCompressedFormat uintptr
	IsEqual            uintptr
	GetRepresentation  uintptr
	FreeRepresentation uintptr
}

func (i *IMFMediaType) GetMajorType(pguidMajorType *guid.GUID) error {
	r1, _, _ := syscall.SyscallN(i.S.GetMajorType, uintptr(unsafe.Pointer(i)), uintptr(unsafe.Pointer(pguidMajorType)))
	if r1 != consts.S_OK {
		return consts.HResultError{Code: r1}
	}
	return nil

}

func (i *IMFMediaType) GetGUID(guidKey, pguidValue *guid.GUID) error {
	r1, _, _ := syscall.SyscallN(i.S.GetGUID, uintptr(unsafe.Pointer(i)), uintptr(unsafe.Pointer(guidKey)), uintptr(unsafe.Pointer(pguidValue)))
	if r1 != consts.S_OK {
		return consts.HResultError{Code: r1}
	}
	return nil
}

func (i *IMFMediaType) IsCompressedFormat() (bool, error) {
	var isCompressed bool
	r1, _, _ := syscall.SyscallN(i.S.IsCompressedFormat, uintptr(unsafe.Pointer(i)), uintptr(unsafe.Pointer(&isCompressed)))
	if r1 != consts.S_OK {
		return false, consts.HResultError{Code: r1}
	}
	return isCompressed, nil
}

func (i *IMFMediaType) GetFrameSize() (uint32, uint32, error) {
	var data uint64

	r1, _, _ := syscall.SyscallN(i.S.GetUINT64, uintptr(unsafe.Pointer(i)), uintptr(unsafe.Pointer(guid.MF_MT_FRAME_SIZE)), uintptr(unsafe.Pointer(&data)))
	if r1 != consts.S_OK {
		return 0, 0, consts.HResultError{Code: r1}
	}

	//上 32 位包含宽度，下 32 位包含高度。
	width := uint32(data >> 32)
	height := uint32(data & 0xFFFFFFFF)

	return width, height, nil

}

func (i *IMFMediaType) GetFrameRate() (float64, error) {
	var data uint64

	r1, _, _ := syscall.SyscallN(i.S.GetUINT64, uintptr(unsafe.Pointer(i)), uintptr(unsafe.Pointer(guid.MF_MT_FRAME_RATE)), uintptr(unsafe.Pointer(&data)))
	if r1 != consts.S_OK {
		return 0, consts.HResultError{Code: r1}
	}

	//上 32 位包含帧率的分子，下 32 位包含帧率的分母。
	numerator := float64(data >> 32)
	denominator := float64(data & 0xFFFFFFFF)

	if denominator == 0 {
		return 0, consts.HResultError{Code: consts.E_INVALIDARG}
	}

	return numerator / denominator, nil
}
func (i *IMFMediaType) Release() error {
	r1, _, _ := syscall.SyscallN(i.S.Release, uintptr(unsafe.Pointer(i)))
	if r1 != consts.S_OK {
		return consts.HResultError{Code: r1}
	}
	return nil
}
