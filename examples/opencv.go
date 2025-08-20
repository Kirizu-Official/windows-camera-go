package main

import (
	"github.com/Kirizu-Official/windows-camera-go/camera/v1"
	"github.com/Kirizu-Official/windows-camera-go/windows/guid"
	"gocv.io/x/gocv"
)

func main() {
	err := camera.Init()
	if err != nil {
		panic(err)
	}

	device, err := camera.OpenDevice("\\\\?\\usb#vid_2c7f&pid_2910&mi_00#8&2412c02a&0&0000#{e5323777-f976-4f5b-9b55-b94699c46e44}\\global")
	if err != nil {
		panic(err)
	}
	// mediaType can be obtained from EnumerateCaptureFormats, and only supports these formats.
	capture, err := device.StartCapture(&camera.CaptureFormats{
		DescriptorIndex:    0,
		MediaTypeIndex:     0,
		IsCompressedFormat: false,
		MajorType:          &guid.MajorTypeVideo,
		SubType:            &guid.SubTypeMediaSubTypeNV12,
		Width:              1920,
		Height:             1080,
		Fps:                30,
	})
	if err != nil {
		panic(err)
	}

	frame, err := capture.GetFrame()
	if err != nil {
		panic(err)
	}
	buffer, err := capture.Device.ParseSampleToBuffer(frame.PpSample)
	if err != nil {
		panic(err)
	}

	matHeight := int(capture.Config.MediaType.Height) + int(capture.Config.MediaType.Height)/2
	matWidth := int(capture.Config.MediaType.Width)
	mat, err := gocv.NewMatFromBytes(matHeight, matWidth, gocv.MatTypeCV8UC1, buffer.Buffer[:buffer.Length])
	if err != nil {
		panic(err)
	}
	defer mat.Close()

	resultMat := gocv.NewMat()
	defer resultMat.Close()
	err = gocv.CvtColor(mat, &resultMat, gocv.ColorYUVToBGRNV12)
	if err != nil {
		panic(err)
	}

	window := gocv.NewWindow("Camera Feed")
	defer window.Close()
	window.IMShow(resultMat)
	window.WaitKey(0)
	window.Close()

	buffer.Release()
	frame.Release()

	device.CloseDevice()
	camera.Shutdown()

}
