package mf

type IMFExtendedCameraControl struct {
	S *IMFExtendedCameraControl_
}
type IMFExtendedCameraControl_ struct {
	Unknown_
	GetCapabilities uintptr
	SetFlags        uintptr
	GetFlags        uintptr
	LockPayload     uintptr
	UnlockPayload   uintptr
	CommitSettings  uintptr
}
