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
	device, err := camera.EnumDevice()
	if err != nil {
		panic(err)
	}

	fmt.Println("Total Devices Found:", len(device))
	fmt.Println("======================================================================================================================================")
	for _, dev := range device {
		fmt.Println("Device Name:", dev.Name)
		fmt.Println("Symbol Link:", dev.SymbolLink)
		fmt.Println("======================================================================================================================================")
	}
	camera.Shutdown()
}
