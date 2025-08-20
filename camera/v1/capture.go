package camera

import (
	"github.com/Kirizu-Official/windows-camera-go/utils"
	"github.com/Kirizu-Official/windows-camera-go/windows/consts"
	"github.com/Kirizu-Official/windows-camera-go/windows/guid"
	"github.com/Kirizu-Official/windows-camera-go/windows/mf"
)

type CaptureConfig struct {
	enableAsync bool
	callBack    uintptr
	MediaType   *CaptureFormats
}
type CaptureAsync struct {
	Device   *Device
	Config   *CaptureConfig
	callback *SourceReaderCallback
}
type CaptureSync struct {
	Device *Device
	Config *CaptureConfig
}

func (d *Device) StartCaptureAsync(mediaType *CaptureFormats) (*CaptureAsync, error) {
	if mediaType == nil {
		return nil, utils.ErrorParameterInvalid
	}
	a := &CaptureAsync{
		Device: d,
	}
	a.callback = newSourceReaderCallback(a)
	config := &CaptureConfig{
		enableAsync: true,
		callBack:    a.callback.ToUintptr(),
		MediaType:   mediaType,
	}
	a.Config = config
	err := d.startCapture(config)
	if err != nil {
		return nil, err
	}

	return a, nil
}
func (d *Device) StartCapture(mediaType *CaptureFormats) (*CaptureSync, error) {
	if mediaType == nil {
		return nil, utils.ErrorParameterInvalid
	}
	config := &CaptureConfig{
		enableAsync: false,
		MediaType:   mediaType,
	}
	err := d.startCapture(config)
	if err != nil {
		return nil, err
	}
	return &CaptureSync{
		Device: d,
		Config: config,
	}, nil
}

func (d *Device) startCapture(config *CaptureConfig) error {
	if d.pDeviceSource == nil {
		return utils.ErrorDeviceNotOpen
	}
	if d.pDevPresent == nil {
		return utils.ErrorDeviceNotOpen
	}
	if d.pDeviceSource == nil {
		return utils.ErrorDeviceNotOpen
	}
	if config.enableAsync && config.callBack == 0 {
		return utils.ErrorAsyncNeedCallBack
	}
	d.config = config
	var err error
	d.pMediaType, err = d.getSelectedMediaType(config.MediaType)
	if err != nil {
		return err
	}
	var sourceReaderAttributes *mf.IMFAttributes
	err = mfPlat.MFCreateAttributes(&sourceReaderAttributes, 2)
	if err != nil {
		return err
	}

	if config.enableAsync {
		err = sourceReaderAttributes.SetUnknown(guid.MF_SOURCE_READER_ASYNC_CALLBACK, config.callBack)
		if err != nil {
			return err
		}
	}

	err = mfReadWrite.MFCreateSourceReaderFromMediaSource(d.pDeviceSource, sourceReaderAttributes, &d.pSourceReader)
	if err != nil {
		return err
	}
	//sourceReaderAttributes.Release()
	err = d.pSourceReader.SetCurrentMediaType(0xFFFFFFFC, d.pMediaType)
	if err != nil {
		return err
	}
	return nil
}

//func (d *Device) StartCapture() error {
//	fmt.Println(d.pDeviceSource.S, "start capture")
//	if d.config.EnableAsync {
//		err := d.pSourceReader.ReadSampleASync(consts.MF_SOURCE_READER_FIRST_VIDEO_STREAM, 0)
//		if err != nil {
//			return err
//		}
//	} else {
//		var pdwActualStreamIndex, pdwStreamFlags uint32
//		var pllTimestamp int64
//		var ppSample *mf.IMFSample
//
//		err := d.pSourceReader.ReadSample(consts.MF_SOURCE_READER_FIRST_VIDEO_STREAM, 0, &pdwActualStreamIndex, &pdwStreamFlags, &pllTimestamp, &ppSample)
//		if err != nil {
//			return err
//		}
//	}
//
//	return nil
//}

func (d *CaptureAsync) GetCallBackInterface() *SourceReaderCallback {
	return d.callback
}

type SampleData struct {
	PdwActualStreamIndex uint32
	PdwStreamFlags       consts.MF_SOURCE_READER_FLAG
	PllTimestamp         int64
	PpSample             *mf.IMFSample
}

func (s *SampleData) Release() error {
	if s.PpSample != nil {
		err := s.PpSample.Release()
		s.PpSample = nil
		return err
	}
	return nil

}

func (d *CaptureSync) GetFrame() (*SampleData, error) {
	var pdwActualStreamIndex, pdwStreamFlags uint32
	var pllTimestamp int64
	var ppSample *mf.IMFSample
RetryGetSample:
	err := d.Device.pSourceReader.ReadSample(consts.MF_SOURCE_READER_FIRST_VIDEO_STREAM, 0, &pdwActualStreamIndex, &pdwStreamFlags, &pllTimestamp, &ppSample)
	if err != nil {
		return nil, err
	}
	if consts.MF_SOURCE_READER_FLAG(pdwStreamFlags) == consts.MF_SOURCE_READERF_STREAMTICK {
		goto RetryGetSample
	}
	if ppSample == nil {
		return &SampleData{
			PdwActualStreamIndex: pdwActualStreamIndex,
			PdwStreamFlags:       consts.MF_SOURCE_READER_FLAG(pdwStreamFlags),
			PllTimestamp:         pllTimestamp,
			PpSample:             nil,
		}, utils.ErrorInternalError
	} else {
		return &SampleData{
			PdwActualStreamIndex: pdwActualStreamIndex,
			PdwStreamFlags:       consts.MF_SOURCE_READER_FLAG(pdwStreamFlags),
			PllTimestamp:         pllTimestamp,
			PpSample:             ppSample,
		}, nil
	}

}
func (d *CaptureAsync) GetNextFrame() error {
	err := d.Device.pSourceReader.ReadSample(consts.MF_SOURCE_READER_FIRST_VIDEO_STREAM, 0, nil, nil, nil, nil)
	if err != nil {
		return err
	}
	return nil
}

type CaptureBuffer struct {
	Buffer         []byte
	Length         uint32
	Total          uint32
	internalBuffer *mf.IMFMediaBuffer
	interSample    *mf.IMFSample
}

func (c *CaptureBuffer) Release() {
	if c.internalBuffer != nil {
		_ = c.internalBuffer.Unlock()
		_ = c.internalBuffer.Release()
		c.internalBuffer = nil
	}

	c.Buffer = nil
	c.Length = 0
	c.Total = 0
}

func (d *Device) ParseSampleToBuffer(sample *mf.IMFSample) (buffer *CaptureBuffer, err error) {
	if sample == nil {
		return nil, utils.ErrorParameterInvalid
	}
	cBuffer, err := sample.ConvertToContiguousBuffer()
	if err != nil {
		return nil, err
	}

	goBuffer, totalLen, currentLen, err := cBuffer.Lock()
	if err != nil {
		return nil, err
	}

	return &CaptureBuffer{
		Buffer:         goBuffer[:totalLen],
		Length:         currentLen,
		Total:          totalLen,
		internalBuffer: cBuffer,
		interSample:    sample,
	}, nil
}

func (d *CaptureAsync) GetConfig() (*ConfigParameter, error) {
	return d.Device.getConfig()
}

func (d *CaptureSync) GetConfig() (*ConfigParameter, error) {
	return d.Device.getConfig()
}

// SetConfigWithCheck will call getConfig to diff the current config and the new config, then set the new config if they are different.
func (d *CaptureSync) SetConfigWithCheck(newConfig *ConfigParameter) error {
	return d.Device.setConfigWithCheck(newConfig)
}
func (d *CaptureAsync) SetConfigWithCheck(newConfig *ConfigParameter) error {
	return d.Device.setConfigWithCheck(newConfig)
}

func (d *Device) setConfigWithCheck(newConfig *ConfigParameter) error {
	if newConfig == nil {
		return utils.ErrorParameterInvalid
	}
	currentConfig, err := d.getConfig()
	if err != nil {
		return err
	}
	cameraCtl, procImp, err := d.GetControl()
	if err != nil {
		return err
	}
	defer cameraCtl.Release()
	defer procImp.Release()

	if newConfig.Pan != currentConfig.Pan {
		err = cameraCtl.Set(consts.CameraControl_Pan, newConfig.Pan.Current, consts.KSPROPERTY_CAMERACONTROL_FLAGS_MANUAL)
		if err != nil {
			return err
		}
	}
	if newConfig.Tilt != currentConfig.Tilt {
		err = cameraCtl.Set(consts.CameraControl_Tilt, newConfig.Tilt.Current, consts.KSPROPERTY_CAMERACONTROL_FLAGS_MANUAL)
		if err != nil {
			return err
		}
	}
	if newConfig.Roll != currentConfig.Roll {
		err = cameraCtl.Set(consts.CameraControl_Roll, newConfig.Roll.Current, consts.KSPROPERTY_CAMERACONTROL_FLAGS_MANUAL)
		if err != nil {
			return err
		}
	}
	if newConfig.Zoom != currentConfig.Zoom {
		err = cameraCtl.Set(consts.CameraControl_Zoom, newConfig.Zoom.Current, consts.KSPROPERTY_CAMERACONTROL_FLAGS_MANUAL)
		if err != nil {
			return err
		}
	}
	if newConfig.Exposure != currentConfig.Exposure {
		err = cameraCtl.Set(consts.CameraControl_Exposure, newConfig.Exposure.Current, consts.KSPROPERTY_CAMERACONTROL_FLAGS_MANUAL)
		if err != nil {
			return err
		}
	}
	if newConfig.Iris != currentConfig.Iris {
		err = cameraCtl.Set(consts.CameraControl_Iris, newConfig.Iris.Current, consts.KSPROPERTY_CAMERACONTROL_FLAGS_MANUAL)
		if err != nil {
			return err
		}
	}
	if newConfig.Focus != currentConfig.Focus {
		err = cameraCtl.Set(consts.CameraControl_Focus, newConfig.Focus.Current, consts.KSPROPERTY_CAMERACONTROL_FLAGS_MANUAL)
		if err != nil {
			return err
		}
	}

	if newConfig.Brightness != currentConfig.Brightness {
		err = procImp.Set(consts.VideoProcAmp_Brightness, newConfig.Brightness.Current, consts.KSPROPERTY_VIDEOPROCAMP_FLAGS_MANUAL)
		if err != nil {
			return err
		}
	}
	if newConfig.Contrast != currentConfig.Contrast {
		err = procImp.Set(consts.VideoProcAmp_Contrast, newConfig.Contrast.Current, consts.KSPROPERTY_VIDEOPROCAMP_FLAGS_MANUAL)
		if err != nil {
			return err
		}
	}
	if newConfig.Hue != currentConfig.Hue {
		err = procImp.Set(consts.VideoProcAmp_Hue, newConfig.Hue.Current, consts.KSPROPERTY_VIDEOPROCAMP_FLAGS_MANUAL)
		if err != nil {
			return err
		}
	}
	if newConfig.Saturation != currentConfig.Saturation {
		err = procImp.Set(consts.VideoProcAmp_Saturation, newConfig.Saturation.Current, consts.KSPROPERTY_VIDEOPROCAMP_FLAGS_MANUAL)
		if err != nil {
			return err
		}
	}
	if newConfig.Sharpness != currentConfig.Sharpness {
		err = procImp.Set(consts.VideoProcAmp_Sharpness, newConfig.Sharpness.Current, consts.KSPROPERTY_VIDEOPROCAMP_FLAGS_MANUAL)
		if err != nil {
			return err
		}
	}
	if newConfig.Gamma != currentConfig.Gamma {
		err = procImp.Set(consts.VideoProcAmp_Gamma, newConfig.Gamma.Current, consts.KSPROPERTY_VIDEOPROCAMP_FLAGS_MANUAL)
		if err != nil {
			return err
		}
	}
	if newConfig.ColorEnable != currentConfig.ColorEnable {
		err = procImp.Set(consts.VideoProcAmp_ColorEnable, newConfig.ColorEnable.Current, consts.KSPROPERTY_VIDEOPROCAMP_FLAGS_MANUAL)
		if err != nil {
			return err
		}
	}
	if newConfig.WhiteBalance != currentConfig.WhiteBalance {
		err = procImp.Set(consts.VideoProcAmp_WhiteBalance, newConfig.WhiteBalance.Current, consts.KSPROPERTY_VIDEOPROCAMP_FLAGS_MANUAL)
		if err != nil {
			return err
		}
	}
	if newConfig.BacklightCompensation != currentConfig.BacklightCompensation {
		err = procImp.Set(consts.VideoProcAmp_BacklightCompensation, newConfig.BacklightCompensation.Current, consts.KSPROPERTY_VIDEOPROCAMP_FLAGS_MANUAL)
		if err != nil {
			return err
		}
	}
	if newConfig.Gain != currentConfig.Gain {
		err = procImp.Set(consts.VideoProcAmp_Gain, newConfig.Gain.Current, consts.KSPROPERTY_VIDEOPROCAMP_FLAGS_MANUAL)
		if err != nil {
			return err
		}
	}
	return nil
}
