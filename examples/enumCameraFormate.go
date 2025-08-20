package main

import (
	"fmt"

	"github.com/Kirizu-Official/windows-camera-go/camera/v1"
)

func main() {
	err := camera.Init()
	if err != nil {
		panic(err)
	}
	// symbolLink can be found in the output of the `EnumDeviceSources` function
	device, err := camera.OpenDevice("\\\\?\\usb#vid_2c7f&pid_2910&mi_00#8&2412c02a&0&0000#{e5323777-f976-4f5b-9b55-b94699c46e44}\\global")
	if err != nil {
		panic(err)
	}
	formats, err := device.EnumerateCaptureFormats()
	if err != nil {
		panic(err)
	}

	for _, format := range formats {
		fmt.Printf("Index(desc-mt): %d,%d", format.DescriptorIndex, format.MediaTypeIndex)
		fmt.Printf("\tMajor Type: %s", format.MajorType)
		fmt.Printf("\tSub Type:%s", format.SubType)
		fmt.Printf("\tFrame: %dx%d", format.Width, format.Height)
		fmt.Printf("\tisComposed: %d", format.IsCompressedFormat)
		fmt.Printf("\tFPS:%.2f\n", format.Fps)
		fmt.Println("-----------------------------")
	}

	device.CloseDevice()
	camera.Shutdown()
}
