package main

import (
	"fmt"

	"github.com/Kirizu-Official/windows-camera-go/camera/v1"
	"github.com/Kirizu-Official/windows-camera-go/windows/guid"
)

func main() {
	err := camera.Init()
	if err != nil {
		panic(err)
	}

	device, err := camera.OpenDevice("\\\\?\\xiaomi#virtualcamera#1&2463b2b7&0&04#{e5323777-f976-4f5b-9b55-b94699c46e44}\\{06ae74cc-87a4-40e8-9b79-961950bfaecf}")
	if err != nil {
		panic(err)
	}
	capture, err := device.StartCapture(&camera.CaptureFormats{
		DescriptorIndex:    0,
		MediaTypeIndex:     2,
		IsCompressedFormat: false,
		MajorType:          &guid.MajorTypeVideo,
		SubType:            &guid.SubTypeMediaSubTypeNV12,
		Width:              1440,
		Height:             1080,
		Fps:                30,
	})
	if err != nil {
		panic(err)
	}

	config, err := capture.GetConfig()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Config: %+v\n", config)
	device.CloseDevice()
	camera.Shutdown()

}
