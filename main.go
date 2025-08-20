package main

import "C"
import (
	"context"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/Kirizu-Official/windows-camera-go/camera/v1"
	"github.com/Kirizu-Official/windows-camera-go/windows/guid"
	"github.com/Kirizu-Official/windows-camera-go/windows/mf"
)

var las uintptr
var s sync.Mutex
var frames int

func main() {

	err := camera.Init()
	if err != nil {
		panic(err)
	}
	enumDevice, err := camera.EnumDevice()
	if err != nil {
		panic(err)
	}
	fmt.Println("Device Count:", len(enumDevice))
	for i, device := range enumDevice {
		fmt.Println("Device Index:", i)
		fmt.Println("Device Name:", device.Name)
		fmt.Println("Device Symbolic Link:", device.SymbolLink)
	}

	openDevice, err := camera.OpenDevice(context.Background(), "\\\\?\\usb#vid_0c45&pid_6368&mi_00#8&1cc66dda&0&0000#{e5323777-f976-4f5b-9b55-b94699c46e44}\\global")
	if err != nil {
		panic(err)
	}
	formats, err := openDevice.EnumerateCaptureFormats()
	if err != nil {
		panic(err)
	}
	fmt.Println("Capture Formats Count:", len(formats))
	for i, format := range formats {
		fmt.Println("Format Index:", i)

		fmt.Printf("%+v\n", format)

		//mediaType.GetGUID(guid.MF_MT_SUBTYPE, &subType)
		//
		//mediaType.GetMajorType(&majorType)
		//var subType guid.GUID
		//mediaType.GetGUID(guid.MF_MT_SUBTYPE, &subType)
	}

	// 创建回调实例

	cpt, err := openDevice.StartCaptureAsync(&camera.CaptureFormats{
		DescriptorIndex:    0,
		MediaTypeIndex:     1,
		IsCompressedFormat: false,
		MajorType:          &guid.MajorTypeVideo,
		SubType:            &guid.SubTypeMediaSubTypeMJPG,
		Width:              1280,
		Height:             720,
		Fps:                30,
	})
	if err != nil {
		panic(err)
	}
	cpt.GetCallBackInterface().SetOnReadSample(SampleCallBack)
	//cpt.GetCallBackInterface().SetOnReadSample(SampleCallBack)
	//err = cpt.GetNextFrame()
	//if err != nil {
	//	panic(err)
	//}
	//frame, err := cpt.GetFrame()
	//if err != nil {
	//	panic(err)
	//}
	//buffer, err := openDevice.ParseSampleToBuffer(frame.PpSample)
	err = cpt.GetNextFrame()
	if err != nil {
		panic(err)
	}
	//fmt.Println("Buffer:", buffer.Buffer[:buffer.Length])
	//fmt.Println("Buffer Length:", buffer.Length)
	//
	//fmt.Println(frame)
	fmt.Println("call ok")
	//openDevice.GetFrame()
	for true {
		time.Sleep(time.Second)
		cpt.GetNextFrame()
	}
	return

}

func SampleCallBack(this *camera.CaptureAsync, dwStreamIndex uint32, dwStreamFlags uint32, llTimestamp int64, sample *mf.IMFSample) {
	s.Lock()
	buffer, err := this.Device.ParseSampleToBuffer(sample)
	if err != nil {
		panic(err)
	}

	fmt.Println(hex.EncodeToString(buffer.Buffer[:buffer.Length]))
	fmt.Println(buffer.Length, buffer.Total)

	// 处理读取到的样本数据
	println("OnReadSample called")
	println("Stream Index:", dwStreamIndex)
	println("Stream Flags:", dwStreamFlags)
	println("Timestamp:", llTimestamp)
	buffer.Release()
	//sample.Release()
	fmt.Println("free done!")
	s.Unlock()
}
