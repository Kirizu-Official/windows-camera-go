package camera

import (
	"context"

	"github.com/Kirizu-Official/windows-camera-go/utils"
	"github.com/Kirizu-Official/windows-camera-go/windows/mf"
)

type Device struct {
	Name          string
	SymbolLink    string
	pDeviceActive *mf.IMFActivate
	pDeviceSource *mf.IMFMediaSource
	pDevPresent   *mf.IMFPresentationDescriptor
	pSourceReader *mf.IMFSourceReader
	pMediaType    *mf.IMFMediaType
	ctx           context.Context
	ctxCancel     context.CancelFunc
	config        *CaptureConfig
	errorCallBack func(err error)
}

func OpenDevice(ctx context.Context, symbolLink string) (*Device, error) {
	if ctx == nil {
		return nil, utils.ErrorContextNil
	}
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	sources, err := callEnumDeviceSources()
	if err != nil {
		return nil, err
	}

	var info *DeviceInfo
	for _, source := range sources {
		info, err = callGetBaseInfo(source)
		if err != nil {
			continue
		}
		if info.SymbolLink == symbolLink {
			ctxUse, ctxCancel := context.WithCancel(ctx)
			device := &Device{
				Name:          info.Name,
				SymbolLink:    info.SymbolLink,
				pDeviceActive: source,
				ctx:           ctxUse,
				ctxCancel:     ctxCancel,
			}
			err = device.init()
			if err != nil {
				return nil, err
			}
			return device, nil
		} else {
			_ = source.Release()
		}
	}
	return nil, utils.ErrorDeviceNotFound
}

func (d *Device) init() error {
	if d.pDeviceActive == nil {
		return utils.ErrorDeviceNotAllowedCall
	}

	err := d.pDeviceActive.ActivateObject(&d.pDeviceSource)
	if err != nil {
		return err
	}

	err = d.pDeviceSource.CreatePresentationDescriptor(&d.pDevPresent)
	if err != nil {
		return err
	}

	return nil
}

func (d *Device) internalErrorSet(err error) {
	if d.errorCallBack != nil {
		d.errorCallBack(err)
	}
}

func (d *Device) GetControl() (*mf.IAMCameraControl, *mf.IAMVideoProcAmp, error) {
	return d.pSourceReader.GetServiceForStream()
}

func (d *Device) CloseDevice() {
	d.ctxCancel()
	d.pSourceReader.Release()
	d.pMediaType.Release()
	d.pDevPresent.Release()
	d.pDeviceSource.Release()
	d.pDeviceActive.Release()
}
