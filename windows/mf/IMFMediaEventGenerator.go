package mf

type IMFMediaEventGenerator struct {
	S *IMFMediaEventGenerator_
}
type IMFMediaEventGenerator_ struct {
	Unknown_
	GetEvent      uintptr
	BeginGetEvent uintptr
	EndGetEvent   uintptr
	QueueEvent    uintptr
}
