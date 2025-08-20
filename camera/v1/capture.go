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
	mediaType   *CaptureFormats
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
		mediaType:   mediaType,
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
		mediaType:   mediaType,
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
	d.pMediaType, err = d.getSelectedMediaType(config.mediaType)
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
