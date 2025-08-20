package mf

import (
	"syscall"
	"unsafe"

	"github.com/Kirizu-Official/windows-camera-go/windows/consts"
	"github.com/Kirizu-Official/windows-camera-go/windows/guid"
)

type IMFAttributes struct {
	S *IMFAttributes_
}
type Unknown_ struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr
}
type IMFAttributes_ struct {
	Unknown_
	GetItem            uintptr
	GetItemType        uintptr
	CompareItem        uintptr
	Compare            uintptr
	GetUINT32          uintptr
	GetUINT64          uintptr
	GetDouble          uintptr
	GetGUID            uintptr
	GetStringLength    uintptr
	GetString          uintptr
	GetAllocatedString uintptr
	GetBlobSize        uintptr
	GetBlob            uintptr
	GetAllocatedBlob   uintptr
	GetUnknown         uintptr
	SetItem            uintptr
	DeleteItem         uintptr
	DeleteAllItems     uintptr
	SetUINT32          uintptr
	SetUINT64          uintptr
	SetDouble          uintptr
	SetGUID            uintptr
	SetString          uintptr
	SetBlob            uintptr
	SetUnknown         uintptr
	LockStore          uintptr
	UnlockStore        uintptr
	GetCount           uintptr
	GetItemByIndex     uintptr
	CopyAllItems       uintptr
}

func (i *IMFAttributes) SetGUID(guidKey *guid.GUID, guidValue *guid.GUID) error {
	r1, _, _ := syscall.SyscallN(i.S.SetGUID, uintptr(unsafe.Pointer(i)), uintptr(unsafe.Pointer(guidKey)), uintptr(unsafe.Pointer(guidValue)))
	if r1 != consts.S_OK {
		return consts.HResultError{Code: r1}
	}

	return nil
}

func (i *IMFAttributes) Release() error {
	r1, _, _ := syscall.SyscallN(i.S.Release, uintptr(unsafe.Pointer(i)))
	if r1 != consts.S_OK {
		return consts.HResultError{Code: r1}
	}

	return nil
}
func (i *IMFAttributes) SetUnknown(guidKey *guid.GUID, val uintptr) error {

	r1, _, _ := syscall.SyscallN(i.S.SetUnknown, uintptr(unsafe.Pointer(i)), uintptr(unsafe.Pointer(guidKey)), val)
	if r1 != consts.S_OK {
		return consts.HResultError{Code: r1}
	}
	return nil

}
