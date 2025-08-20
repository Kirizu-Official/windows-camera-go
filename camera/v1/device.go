package camera

import (
	"github.com/Kirizu-Official/windows-camera-go/utils"
	"github.com/Kirizu-Official/windows-camera-go/windows/consts"
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
	config        *CaptureConfig
	errorCallBack func(err error)
}

func OpenDevice(symbolLink string) (*Device, error) {

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
			device := &Device{
				Name:          info.Name,
				SymbolLink:    info.SymbolLink,
				pDeviceActive: source,
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

type ConfigParameterStatus struct {
	Current  int32 `json:"current" yaml:"current"`
	Min      int32 `json:"min" yaml:"min"`
	Max      int32 `json:"max" yaml:"max"`
	Default  int32 `json:"default" yaml:"default"`
	Flags    int32 `json:"flags" yaml:"flags"`       // Flags indicate the control
	Stepping int32 `json:"stepping" yaml:"stepping"` // Setting indicates the current setting of the control
	Found    bool  `json:"found" yaml:"found"`       // Found indicates if the control was found
}
type ConfigParameter struct {
	Pan                   ConfigParameterStatus `json:"pan" yaml:"pan"`
	Tilt                  ConfigParameterStatus `json:"tilt" yaml:"tilt"`
	Roll                  ConfigParameterStatus `json:"roll" yaml:"roll"`
	Zoom                  ConfigParameterStatus `json:"zoom" yaml:"zoom"`
	Exposure              ConfigParameterStatus `json:"exposure" yaml:"exposure"`
	Iris                  ConfigParameterStatus `json:"iris" yaml:"iris"`
	Focus                 ConfigParameterStatus `json:"focus" yaml:"focus"`
	Brightness            ConfigParameterStatus `json:"brightness" yaml:"brightness"`
	Contrast              ConfigParameterStatus `json:"contrast" yaml:"contrast"`
	Hue                   ConfigParameterStatus `json:"hue" yaml:"hue"`
	Saturation            ConfigParameterStatus `json:"saturation" yaml:"saturation"`
	Sharpness             ConfigParameterStatus `json:"sharpness" yaml:"sharpness"`
	Gamma                 ConfigParameterStatus `json:"gamma" yaml:"gamma"`
	ColorEnable           ConfigParameterStatus `json:"color_enable" yaml:"color_enable"`
	WhiteBalance          ConfigParameterStatus `json:"white_balance" yaml:"white_balance"`
	BacklightCompensation ConfigParameterStatus `json:"backlight_compensation" yaml:"backlight_compensation"`
	Gain                  ConfigParameterStatus `json:"gain" yaml:"gain"`
}

func (d *Device) getConfig() (*ConfigParameter, error) {
	cameraCtl, procAmp, err := d.GetControl()
	if err != nil {
		return nil, err
	}
	pan, err := cameraCtl.GetRange(consts.CameraControl_Pan)
	if err != nil {
		return nil, err
	}

	tilt, err := cameraCtl.GetRange(consts.CameraControl_Tilt)
	if err != nil {
		return nil, err
	}
	roll, err := cameraCtl.GetRange(consts.CameraControl_Roll)
	if err != nil {
		return nil, err
	}
	zoom, err := cameraCtl.GetRange(consts.CameraControl_Zoom)
	if err != nil {
		return nil, err
	}
	exposure, err := cameraCtl.GetRange(consts.CameraControl_Exposure)
	if err != nil {
		return nil, err
	}
	iris, err := cameraCtl.GetRange(consts.CameraControl_Iris)
	if err != nil {
		return nil, err
	}
	focus, err := cameraCtl.GetRange(consts.CameraControl_Focus)
	if err != nil {
		return nil, err
	}
	brightness, err := procAmp.GetRange(consts.VideoProcAmp_Brightness)
	if err != nil {
		return nil, err
	}
	contrast, err := procAmp.GetRange(consts.VideoProcAmp_Contrast)
	if err != nil {
		return nil, err
	}
	hue, err := procAmp.GetRange(consts.VideoProcAmp_Hue)
	if err != nil {
		return nil, err
	}
	saturation, err := procAmp.GetRange(consts.VideoProcAmp_Saturation)
	if err != nil {
		return nil, err
	}
	sharpness, err := procAmp.GetRange(consts.VideoProcAmp_Sharpness)
	if err != nil {
		return nil, err
	}
	gamma, err := procAmp.GetRange(consts.VideoProcAmp_Gamma)
	if err != nil {
		return nil, err
	}
	colorEnable, err := procAmp.GetRange(consts.VideoProcAmp_ColorEnable)
	if err != nil {
		return nil, err
	}
	whiteBalance, err := procAmp.GetRange(consts.VideoProcAmp_WhiteBalance)
	if err != nil {
		return nil, err
	}
	backlightCompensation, err := procAmp.GetRange(consts.VideoProcAmp_BacklightCompensation)
	if err != nil {
		return nil, err
	}
	gain, err := procAmp.GetRange(consts.VideoProcAmp_Gain)
	if err != nil {
		return nil, err
	}
	panCurrent, _, err := cameraCtl.Get(consts.CameraControl_Pan)
	if err != nil {
		return nil, err
	}
	tiltCurrent, _, err := cameraCtl.Get(consts.CameraControl_Tilt)
	if err != nil {
		return nil, err
	}
	rollCurrent, _, err := cameraCtl.Get(consts.CameraControl_Roll)
	if err != nil {
		return nil, err
	}
	zoomCurrent, _, err := cameraCtl.Get(consts.CameraControl_Zoom)
	if err != nil {
		return nil, err
	}
	exposureCurrent, _, err := cameraCtl.Get(consts.CameraControl_Exposure)
	if err != nil {
		return nil, err
	}
	irisCurrent, _, err := cameraCtl.Get(consts.CameraControl_Iris)
	if err != nil {
		return nil, err
	}
	focusCurrent, _, err := cameraCtl.Get(consts.CameraControl_Focus)
	if err != nil {
		return nil, err
	}
	brightnessCurrent, _, err := procAmp.Get(consts.VideoProcAmp_Brightness)
	if err != nil {
		return nil, err
	}
	contrastCurrent, _, err := procAmp.Get(consts.VideoProcAmp_Contrast)
	if err != nil {
		return nil, err
	}
	hueCurrent, _, err := procAmp.Get(consts.VideoProcAmp_Hue)
	if err != nil {
		return nil, err
	}
	saturationCurrent, _, err := procAmp.Get(consts.VideoProcAmp_Saturation)
	if err != nil {
		return nil, err
	}
	sharpnessCurrent, _, err := procAmp.Get(consts.VideoProcAmp_Sharpness)
	if err != nil {
		return nil, err
	}
	gammaCurrent, _, err := procAmp.Get(consts.VideoProcAmp_Gamma)
	if err != nil {
		return nil, err
	}
	colorEnableCurrent, _, err := procAmp.Get(consts.VideoProcAmp_ColorEnable)
	if err != nil {
		return nil, err
	}
	whiteBalanceCurrent, _, err := procAmp.Get(consts.VideoProcAmp_WhiteBalance)
	if err != nil {
		return nil, err
	}
	backlightCompensationCurrent, _, err := procAmp.Get(consts.VideoProcAmp_BacklightCompensation)
	if err != nil {
		return nil, err
	}
	gainCurrent, _, err := procAmp.Get(consts.VideoProcAmp_Gain)
	if err != nil {
		return nil, err
	}
	panStatus := ConfigParameterStatus{
		Current:  panCurrent,
		Min:      pan.Min,
		Max:      pan.Max,
		Default:  pan.DefaultValue,
		Flags:    int32(pan.CapsFlags),
		Stepping: pan.SteppingDelta,
		Found:    pan.Found,
	}
	tiltStatus := ConfigParameterStatus{
		Current:  tiltCurrent,
		Min:      tilt.Min,
		Max:      tilt.Max,
		Default:  tilt.DefaultValue,
		Flags:    int32(tilt.CapsFlags),
		Stepping: tilt.SteppingDelta,
		Found:    tilt.Found,
	}
	rollStatus := ConfigParameterStatus{
		Current:  rollCurrent,
		Min:      roll.Min,
		Max:      roll.Max,
		Default:  roll.DefaultValue,
		Flags:    int32(roll.CapsFlags),
		Stepping: roll.SteppingDelta,
		Found:    roll.Found,
	}
	zoomStatus := ConfigParameterStatus{
		Current:  zoomCurrent,
		Min:      zoom.Min,
		Max:      zoom.Max,
		Default:  zoom.DefaultValue,
		Flags:    int32(zoom.CapsFlags),
		Stepping: zoom.SteppingDelta,
		Found:    zoom.Found,
	}
	exposureStatus := ConfigParameterStatus{
		Current:  exposureCurrent,
		Min:      exposure.Min,
		Max:      exposure.Max,
		Default:  exposure.DefaultValue,
		Flags:    int32(exposure.CapsFlags),
		Stepping: exposure.SteppingDelta,
		Found:    exposure.Found,
	}
	irisStatus := ConfigParameterStatus{
		Current:  irisCurrent,
		Min:      iris.Min,
		Max:      iris.Max,
		Default:  iris.DefaultValue,
		Flags:    int32(iris.CapsFlags),
		Stepping: iris.SteppingDelta,
		Found:    iris.Found,
	}
	focusStatus := ConfigParameterStatus{
		Current:  focusCurrent,
		Min:      focus.Min,
		Max:      focus.Max,
		Default:  focus.DefaultValue,
		Flags:    int32(focus.CapsFlags),
		Stepping: focus.SteppingDelta,
		Found:    focus.Found,
	}
	brightnessStatus := ConfigParameterStatus{
		Current:  brightnessCurrent,
		Min:      brightness.Min,
		Max:      brightness.Max,
		Default:  brightness.DefaultValue,
		Flags:    int32(brightness.CapsFlags),
		Stepping: brightness.SteppingDelta,
		Found:    brightness.Found,
	}
	contrastStatus := ConfigParameterStatus{
		Current:  contrastCurrent,
		Min:      contrast.Min,
		Max:      contrast.Max,
		Default:  contrast.DefaultValue,
		Flags:    int32(contrast.CapsFlags),
		Stepping: contrast.SteppingDelta,
		Found:    contrast.Found,
	}
	hueStatus := ConfigParameterStatus{
		Current:  hueCurrent,
		Min:      hue.Min,
		Max:      hue.Max,
		Default:  hue.DefaultValue,
		Flags:    int32(hue.CapsFlags),
		Stepping: hue.SteppingDelta,
		Found:    hue.Found,
	}
	saturationStatus := ConfigParameterStatus{
		Current:  saturationCurrent,
		Min:      saturation.Min,
		Max:      saturation.Max,
		Default:  saturation.DefaultValue,
		Flags:    int32(saturation.CapsFlags),
		Stepping: saturation.SteppingDelta,
		Found:    saturation.Found,
	}
	sharpnessStatus := ConfigParameterStatus{
		Current:  sharpnessCurrent,
		Min:      sharpness.Min,
		Max:      sharpness.Max,
		Default:  sharpness.DefaultValue,
		Flags:    int32(sharpness.CapsFlags),
		Stepping: sharpness.SteppingDelta,
		Found:    sharpness.Found,
	}
	gammaStatus := ConfigParameterStatus{
		Current:  gammaCurrent,
		Min:      gamma.Min,
		Max:      gamma.Max,
		Default:  gamma.DefaultValue,
		Flags:    int32(gamma.CapsFlags),
		Stepping: gamma.SteppingDelta,
		Found:    gamma.Found,
	}
	colorEnableStatus := ConfigParameterStatus{
		Current:  colorEnableCurrent,
		Min:      colorEnable.Min,
		Max:      colorEnable.Max,
		Default:  colorEnable.DefaultValue,
		Flags:    int32(colorEnable.CapsFlags),
		Stepping: colorEnable.SteppingDelta,
		Found:    colorEnable.Found,
	}
	whiteBalanceStatus := ConfigParameterStatus{
		Current:  whiteBalanceCurrent,
		Min:      whiteBalance.Min,
		Max:      whiteBalance.Max,
		Default:  whiteBalance.DefaultValue,
		Flags:    int32(whiteBalance.CapsFlags),
		Stepping: whiteBalance.SteppingDelta,
		Found:    whiteBalance.Found,
	}
	backlightCompensationStatus := ConfigParameterStatus{
		Current:  backlightCompensationCurrent,
		Min:      backlightCompensation.Min,
		Max:      backlightCompensation.Max,
		Default:  backlightCompensation.DefaultValue,
		Flags:    int32(backlightCompensation.CapsFlags),
		Stepping: backlightCompensation.SteppingDelta,
		Found:    backlightCompensation.Found,
	}
	gainStatus := ConfigParameterStatus{
		Current:  gainCurrent,
		Min:      gain.Min,
		Max:      gain.Max,
		Default:  gain.DefaultValue,
		Flags:    int32(gain.CapsFlags),
		Stepping: gain.SteppingDelta,
		Found:    gain.Found,
	}

	return &ConfigParameter{
		Pan:                   panStatus,
		Tilt:                  tiltStatus,
		Roll:                  rollStatus,
		Zoom:                  zoomStatus,
		Exposure:              exposureStatus,
		Iris:                  irisStatus,
		Focus:                 focusStatus,
		Brightness:            brightnessStatus,
		Contrast:              contrastStatus,
		Hue:                   hueStatus,
		Saturation:            saturationStatus,
		Sharpness:             sharpnessStatus,
		Gamma:                 gammaStatus,
		ColorEnable:           colorEnableStatus,
		WhiteBalance:          whiteBalanceStatus,
		BacklightCompensation: backlightCompensationStatus,
		Gain:                  gainStatus,
	}, nil

}

func (d *Device) CloseDevice() {
	if d.pSourceReader != nil {
		d.pSourceReader.Release()
	}
	if d.pMediaType != nil {
		d.pMediaType.Release()
	}

	if d.pDevPresent != nil {
		d.pDevPresent.Release()
	}
	if d.pDeviceSource != nil {
		d.pDeviceSource.Release()
	}
	if d.pDeviceActive != nil {
		d.pDeviceActive.Release()
	}

}
