package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Kirizu-Official/windows-camera-go/camera/v1"
	"github.com/Kirizu-Official/windows-camera-go/windows/guid"
)

func main() {
	s := NewStream()
	go func() {
		err := http.ListenAndServe(":8080", s)
		if err != nil {
			panic(err)
		}
	}()
	fmt.Println("HTTP MJPEG Stream started at http://localhost:8080")
	err := camera.Init()
	if err != nil {
		panic(err)
	}
	fmt.Println("Camera initialized successfully")
	device, err := camera.OpenDevice(context.Background(), "\\\\?\\usb#vid_2c7f&pid_2910&mi_00#8&2412c02a&0&0000#{e5323777-f976-4f5b-9b55-b94699c46e44}\\global")
	if err != nil {
		panic(err)
	}
	fmt.Println("Camera opened successfully")
	// mjpeg server need mjpeg format which from enumerateCaptureFormats
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
	fps := 30
	fmt.Println("MJPEG Stream started at 1920x1080 with", fps, "FPS :", "http://localhost:8080")
	c := make(chan os.Signal)
	// 监听信号
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	for {
		select {
		case <-time.After(time.Second / time.Duration(fps)):
			frame, err := capture.GetFrame()
			if err != nil {
				panic(err)
			}
			buffer, err := device.ParseSampleToBuffer(frame.PpSample)
			if err != nil {
				panic(err)
			}
			s.UpdateJPEG(buffer.Buffer[:buffer.Length])
			buffer.Release()
			frame.Release()
		case <-c:
			device.CloseDevice()
			camera.Shutdown()
			os.Exit(0)
		}
	}

}

/**
The source code below is from: https://github.com/hybridgroup/mjpeg/blob/master/stream.go
*/

// Stream represents a single video feed.
type Stream struct {
	ctx           context.Context
	start         time.Time
	m             map[chan []byte]bool
	frame         []byte
	lock          sync.Mutex
	FrameInterval time.Duration
}

// NewStream initializes and returns a new Stream.
func NewStream() *Stream {
	return &Stream{
		m:             make(map[chan []byte]bool),
		frame:         make([]byte, len(headerf)),
		FrameInterval: 50 * time.Millisecond,
		ctx:           context.Background(),
	}
}

// NewStreamWithContext initializes and returns a new Stream.
func NewStreamWithContext(ctx context.Context) *Stream {
	return &Stream{
		m:             make(map[chan []byte]bool),
		frame:         make([]byte, len(headerf)),
		FrameInterval: 50 * time.Millisecond,
		ctx:           ctx,
	}
}

const boundaryWord = "MJPEGBOUNDARY"
const headerf = "\r\n" +
	"--" + boundaryWord + "\r\n" +
	"Content-Type: image/jpeg\r\n" +
	"Content-Length: %d\r\n" +
	"X-Timestamp: %d.%d\r\n" +
	"\r\n"

// ServeHTTP responds to HTTP requests with the MJPEG stream, implementing the http.Handler interface.
func (s *Stream) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	slog.Info("Stream:", r.RemoteAddr, "connected")
	w.Header().Add("Content-Type", "multipart/x-mixed-replace;boundary="+boundaryWord)

	c := make(chan []byte)
	s.lock.Lock()
	s.m[c] = true
	s.lock.Unlock()
	s.start = time.Now()

FRAMELOOP:
	for {
		time.Sleep(s.FrameInterval)
		select {
		case <-s.ctx.Done():
			break FRAMELOOP
		case b := <-c:
			_, err := w.Write(b)
			if err != nil {
				slog.Debug("Stream:", r.RemoteAddr, err.Error())
				break FRAMELOOP
			}
		}
	}

	s.lock.Lock()
	delete(s.m, c)
	s.lock.Unlock()

	slog.Info("Stream:", r.RemoteAddr, "disconnected")
}

// UpdateJPEG pushes a new JPEG frame onto the clients.
func (s *Stream) UpdateJPEG(jpeg []byte) {
	if len(jpeg) == 0 {
		return
	}
	elapsed := time.Since(s.start)
	s.updateFrame(jpeg, elapsed)

	s.lock.Lock()
	defer s.lock.Unlock()

	for c := range s.m {
		// Select to skip streams which are sleeping to drop frames.
		// This might need more thought.
		select {
		case c <- s.frame:
		case <-s.ctx.Done():
			return
		default:
		}
	}
}

func (s *Stream) updateFrame(jpeg []byte, elapsed time.Duration) {
	header := s.frameHeader(jpeg, elapsed)
	if len(s.frame) < len(jpeg)+len(header) {
		s.frame = make([]byte, (len(jpeg)+len(header))*2)
	}

	copy(s.frame, header)
	copy(s.frame[len(header):], jpeg)
}

func (s *Stream) frameHeader(jpeg []byte, elapsed time.Duration) string {
	sec := int64(elapsed.Seconds())
	usec := int64(elapsed.Microseconds() % 1e6)
	return fmt.Sprintf(headerf, len(jpeg), sec, usec)
}
