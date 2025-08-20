package consts

import "strconv"

//S_OK	Operation successful	0x00000000
//E_ABORT	Operation aborted	0x80004004
//E_ACCESSDENIED	General access denied error	0x80070005
//E_FAIL	Unspecified failure	0x80004005
//E_HANDLE	Handle that is not valid	0x80070006
//E_INVALIDARG	One or more arguments are not valid	0x80070057
//E_NOINTERFACE	No such interface supported	0x80004002
//E_NOTIMPL	Not implemented	0x80004001
//E_OUTOFMEMORY	Failed to allocate necessary memory	0x8007000E
//E_POINTER	Pointer that is not valid	0x80004003
//E_UNEXPECTED	Unexpected failure	0x8000FFFF

type HRESULT = uintptr

const (
	S_OK            HRESULT = 0x00000000
	E_ABORT         HRESULT = 0x80004004
	E_ACCESSDENIED  HRESULT = 0x80070005
	E_FAIL          HRESULT = 0x80004005
	E_HANDLE        HRESULT = 0x80070006
	E_INVALIDARG    HRESULT = 0x80070057
	E_NOINTERFACE   HRESULT = 0x80004002
	E_NOTIMPL       HRESULT = 0x80004001
	E_OUTOFMEMORY   HRESULT = 0x8007000E
	E_POINTER       HRESULT = 0x80004003
	E_UNEXPECTED    HRESULT = 0x8000FFFF
	ERROR_NOT_FOUND HRESULT = 0x80070490 // #define ERROR_NOT_FOUND 1168L
)

type HResultError struct {
	Code HRESULT
}

func (e HResultError) Error() string {
	return "HRESULT Error: " + strconv.Itoa(int(e.Code))
}
func (e HResultError) String() string {
	return "HRESULT Error: " + strconv.Itoa(int(e.Code))
}
