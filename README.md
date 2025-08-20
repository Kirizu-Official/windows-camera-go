# Windows-Camera-Go

A simple golang camera library on Windows, based on the Windows Media Foundation API, support Windows 7 and later.

***⚠ WARN: This library is still under development and may not be stable. Use at your own risk. Api may change in the future.***

## Features
- Capture video from usb camera
- Read source binary data from camera
- Camera control (zoom, focus, brightness, contrast, etc.)
- Support multiple cameras
- Not need CGO

## Installation

```bash
go get github.com/Kirizu-Official/windows-camera-go
```

## Usage

For more examples, please refer to the [examples](https://github.com/Kirizu-Official/windows-camera-go/tree/main/examples).

Here is a simple example of how to use the library to capture video from a USB camera:

```go
package main

import (
	"context"
	"fmt"

	"github.com/Kirizu-Official/windows-camera-go/camera/v1"
	"github.com/Kirizu-Official/windows-camera-go/windows/guid"
)

func main() {
	err := camera.Init()
	if err != nil {
		panic(err)
	}

	// device can be obtained from EnumerateDevices, or you can use the device path directly.
	device, err := camera.OpenDevice(context.Background(), "\\\\?\\usb#vid_2c7f&pid_2910&mi_00#8&2412c02a&0&0000#{e5323777-f976-4f5b-9b55-b94699c46e44}\\global")
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


		// you can process mjpeg buffer data by buffer.Buffer[:buffer.Length]
		//fmt.Println(hex.EncodeToString(buffer.Buffer[:buffer.Length]), buffer.Length)

		buffer.Release()
		frame.Release()

	}

	device.CloseDevice()
	camera.Shutdown()

}

```
