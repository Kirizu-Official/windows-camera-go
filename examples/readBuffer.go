package main

import (
	"fmt"
	"runtime"

	"github.com/Kirizu-Official/windows-camera-go/camera/v1"
	"github.com/Kirizu-Official/windows-camera-go/windows/guid"
)

func main() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	fmt.Printf("Allocated memory at Start: %v bytes\n", memStats.Alloc)

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
		MediaTypeIndex:     1,
		IsCompressedFormat: false,
		MajorType:          &guid.MajorTypeVideo,
		SubType:            &guid.SubTypeMediaSubTypeMJPG,
		Width:              1920,
		Height:             1080,
		Fps:                30,
	})
	if err != nil {
		panic(err)
	}
	for i := 0; i < 1800; i++ {
		frame, err := capture.GetFrame()
		if err != nil {
			panic(err)
		}
		buffer, err := capture.Device.ParseSampleToBuffer(frame.PpSample)
		if err != nil {
			panic(err)
		}
		// you can process buffer data by buffer.Buffer[:buffer.Length]
		//fmt.Println(hex.EncodeToString(buffer.Buffer[:buffer.Length]), buffer.Length)

		buffer.Release()
		frame.Release()

	}
	device.CloseDevice()
	camera.Shutdown()
	runtime.ReadMemStats(&memStats)
	fmt.Printf("Allocated memory after run: %v bytes\n", memStats.Alloc)
}
