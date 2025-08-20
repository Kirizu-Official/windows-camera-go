package camera

import (
	"fmt"

	"github.com/Kirizu-Official/windows-camera-go/utils"
	"github.com/Kirizu-Official/windows-camera-go/windows/guid"
	"github.com/Kirizu-Official/windows-camera-go/windows/mf"
)

type DeviceInfo struct {
	Name       string `json:"name" yaml:"name"`
	SymbolLink string `json:"symbol_link" yaml:"symbol_link"`
}
type CaptureFormats struct {
	DescriptorIndex    uint32
	MediaTypeIndex     uint32
	IsCompressedFormat bool
	MajorType          *guid.GUID
	SubType            *guid.GUID
	Width              uint32
	Height             uint32
	Fps                float64
}

func EnumDevice() ([]*DeviceInfo, error) {
	list, err := callEnumDeviceSources()
	if err != nil {
		return nil, err
	}
	var deviceList []*DeviceInfo
	var info *DeviceInfo
	for i := uint32(0); i < uint32(len(list)); i++ {
		info, err = callGetBaseInfo(list[i])
		if err != nil {
			continue
		}
		deviceList = append(deviceList, info)
		_ = list[i].Release()
	}
	return deviceList, nil
}

func (d *Device) EnumerateCaptureFormats() ([]*CaptureFormats, error) {
	decCount, err := d.pDevPresent.GetStreamDescriptorCount()
	if err != nil {
		return nil, err
	}
	var result []*CaptureFormats
	var typeCount uint32

	for index := 0; index < int(decCount); index++ {
		descriptor, _, err := d.getStreamDescriptorByIndex(uint32(index))
		if err != nil {
			fmt.Println("[WARN]GetStreamDescriptorByIndex: " + err.Error())
			continue
		}
		var MediaTypeHandler *mf.IMFMediaTypeHandler
		err = descriptor.GetMediaTypeHandler(&MediaTypeHandler)
		if err != nil {
			fmt.Println("[WARN]GetMediaTypeHandler: " + err.Error())
			descriptor.Release()
			continue
		}

		typeCount, err = MediaTypeHandler.GetMediaTypeCount()
		for n := 0; n < int(typeCount); n++ {
			var mediaType *mf.IMFMediaType
			err = MediaTypeHandler.GetMediaTypeByIndex(uint32(n), &mediaType)
			if err != nil {
				fmt.Println("[WARN]GetMediaTypeByIndex: " + err.Error())
				mediaType.Release()
				continue
			}
			info, err := d.getMediaTypeInfo(mediaType)
			if err != nil {
				fmt.Println("[WARN]GetMediaTypeInfo: " + err.Error())
				mediaType.Release()
				continue
			}
			info.DescriptorIndex = uint32(index)
			info.MediaTypeIndex = uint32(n)
			result = append(result, info)
			mediaType.Release()
		}
		MediaTypeHandler.Release()
		descriptor.Release()

	}
	return result, nil
}
func (d *Device) getSelectedMediaType(format *CaptureFormats) (*mf.IMFMediaType, error) {

	descriptor, _, err := d.getStreamDescriptorByIndex(format.DescriptorIndex)
	if err != nil {
		return nil, err
	}

	var MediaTypeHandler *mf.IMFMediaTypeHandler
	err = descriptor.GetMediaTypeHandler(&MediaTypeHandler)
	defer descriptor.Release()
	if err != nil {
		return nil, err
	}

	var mediaType *mf.IMFMediaType
	err = MediaTypeHandler.GetMediaTypeByIndex(format.MediaTypeIndex, &mediaType)
	if err != nil {
		return nil, err
	}

	info, err := d.getMediaTypeInfo(mediaType)
	if err != nil {
		mediaType.Release()
		return nil, err
	}

	if info.IsCompressedFormat != format.IsCompressedFormat || !(*guid.GUID)(info.MajorType).IsMatch(format.MajorType) || !(*guid.GUID)(info.SubType).IsMatch(format.SubType) || info.Height != format.Height || info.Width != format.Width || info.Fps != format.Fps {
		mediaType.Release()
		return nil, utils.ErrorFormatNotMatched
	}

	return mediaType, nil
}

func (d *Device) getStreamDescriptorByIndex(index uint32) (*mf.IMFStreamDescriptor, bool, error) {
	var descriptor *mf.IMFStreamDescriptor
	var isSelected bool
	err := d.pDevPresent.GetStreamDescriptorByIndex(index, &isSelected, &descriptor)
	return descriptor, isSelected, err
}
func (d *Device) getMediaTypeInfo(mediaType *mf.IMFMediaType) (*CaptureFormats, error) {
	isCompressed, err := mediaType.IsCompressedFormat()
	if err != nil {
		return nil, err
	}
	var majorType guid.GUID
	err = mediaType.GetMajorType(&majorType)
	if err != nil {
		return nil, err
	}

	var subType guid.GUID
	err = mediaType.GetGUID(guid.MF_MT_SUBTYPE, &subType)
	if err != nil {
		return nil, err
	}

	width, height, err := mediaType.GetFrameSize()
	if err != nil {
		return nil, err
	}
	fps, err := mediaType.GetFrameRate()
	if err != nil {
		return nil, err
	}
	return &CaptureFormats{
		IsCompressedFormat: isCompressed,
		MajorType:          &majorType,
		SubType:            &subType,
		Width:              width,
		Height:             height,
		Fps:                fps,
	}, nil
}
