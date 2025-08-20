package camera

import (
	"github.com/Kirizu-Official/windows-camera-go/windows/guid"
	"github.com/Kirizu-Official/windows-camera-go/windows/mf"
	"syscall"
	"unsafe"
)

func callEnumDeviceSources() ([]*mf.IMFActivate, error) {
	var att *mf.IMFAttributes
	err := mfPlat.MFCreateAttributes(&att, 1)
	if err != nil {
		return nil, err
	}

	err = att.SetGUID(&guid.MF_DEVSOURCE_ATTRIBUTE_SOURCE_TYPE, &guid.MF_DEVSOURCE_ATTRIBUTE_SOURCE_TYPE_VIDCAP_GUID)
	if err != nil {
		return nil, err
	}

	var list []*mf.IMFActivate
	var count uint32
	list, err = mfApi.MFEnumDeviceSources(att, &count)
	if err != nil {
		return nil, err
	}
	_ = att.Release()
	return list, nil
}

func callGetBaseInfo(dev *mf.IMFActivate) (*DeviceInfo, error) {
	var err error
	var deviceNameLen uint32
	deviceNameLen, err = dev.GetStrLength(&guid.MF_DEVSOURCE_ATTRIBUTE_FRIENDLY_NAME)
	//fmt.Println("len", deviceNameLen, err)
	if err != nil {
		return nil, err
	}

	deviceName := make([]uint16, deviceNameLen+1)

	//deviceName := make([]uint16, deviceNameLen+1)
	_, err = dev.GetString(&guid.MF_DEVSOURCE_ATTRIBUTE_FRIENDLY_NAME, uintptr(unsafe.Pointer(&deviceName[0])), deviceNameLen+1)
	if err != nil {
		return nil, err
	}
	//fmt.Println("get result", deviceName, err)
	//deviceNameStr := syscall.UTF16ToString(deviceName[:deviceNameLen])
	//fmt.Println("Device Name:", syscall.UTF16ToString(deviceName[:readLen]))

	var symbolLinkLen uint32
	symbolLinkLen, err = dev.GetStrLength(&guid.MF_DEVSOURCE_ATTRIBUTE_SOURCE_TYPE_VIDCAP_SYMBOLIC_LINK)
	if err != nil {
		return nil, err
	}
	symbolLink := make([]uint16, symbolLinkLen+1)
	_, err = dev.GetString(&guid.MF_DEVSOURCE_ATTRIBUTE_SOURCE_TYPE_VIDCAP_SYMBOLIC_LINK, uintptr(unsafe.Pointer(&symbolLink[0])), symbolLinkLen+1)
	if err != nil {
		return nil, err
	}
	return &DeviceInfo{
		Name:       syscall.UTF16ToString(deviceName[:deviceNameLen]),
		SymbolLink: syscall.UTF16ToString(symbolLink[:symbolLinkLen]),
	}, nil
}
