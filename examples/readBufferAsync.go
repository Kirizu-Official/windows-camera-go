package main

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/Kirizu-Official/windows-camera-go/camera/v1"
	"github.com/Kirizu-Official/windows-camera-go/windows/guid"
	"github.com/Kirizu-Official/windows-camera-go/windows/mf"
)

var frames int

func main() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	fmt.Printf("Allocated memory at Start: %v bytes\n", memStats.Alloc)

	err := camera.Init()
	if err != nil {
		panic(err)
	}

	device, err := camera.OpenDevice(context.Background(), "\\\\?\\usb#vid_2c7f&pid_2910&mi_00#8&2412c02a&0&0000#{e5323777-f976-4f5b-9b55-b94699c46e44}\\global")
	if err != nil {
		panic(err)
	}
	// mediaType can be obtained from EnumerateCaptureFormats, and only supports these formats.
	asyncCapture, err := device.StartCaptureAsync(&camera.CaptureFormats{
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
	asyncCapture.GetCallBackInterface().SetOnReadSample(OnSample)
	runtime.ReadMemStats(&memStats)
	fmt.Printf("Allocated memory after init: %v bytes\n", memStats.Alloc)

	err = asyncCapture.GetNextFrame()
	if err != nil {
		panic(err)
	}

	// get 60s frame
	time.Sleep(time.Second * 60)
	runtime.ReadMemStats(&memStats)
	device.CloseDevice()
	camera.Shutdown()

	fmt.Printf("Allocated memory after run: %v bytes\n", memStats.Alloc)

}

func OnSample(this *camera.CaptureAsync, dwStreamIndex uint32, dwStreamFlags uint32, llTimestamp int64, pSample *mf.IMFSample) {
	// opt, if you want to get the buffer, you can use ParseSampleToBuffer
	buffer, err := this.Device.ParseSampleToBuffer(pSample)
	if err != nil {
		panic(err)
	}

	//you can process buffer.Buffer[:buffer.Length]
	//fmt.println(hex.EncodeToString(buffer.Buffer[:buffer.Length]))

	//You should release the buffer after use
	buffer.Release()

	frames++

	//to get next frame, you should call GetNextFrame
	err = this.GetNextFrame()
	if err != nil {
		panic(err)
	}

}
