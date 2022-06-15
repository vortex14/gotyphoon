// Package main shows an example for transferring jpeg stream over HTTP
package main

import (
	"bytes"
	"fmt"
	"github.com/kbinani/screenshot"
	"github.com/vortex14/gotyphoon/extensions/agents/capturer"
	utils2 "github.com/vortex14/gotyphoon/utils"
	//"github.com/disintegration/imaging"
	"github.com/vortex14/gotyphoon/utils/img"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"log"
	"math"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
)

// Define common colors for convenience
var (
	blue   = color.RGBA{0, 0, 255, 255}
	red    = color.RGBA{255, 0, 0, 255}
	green  = color.RGBA{0, 255, 0, 255}
	yellow = color.RGBA{255, 255, 0, 255}
)

// Boundary will separate frames in M-JPEG animation transfer
const boundary = "abcd4321"

func main() {
	println("Start ")
	// Static files (such as html file) are served from /static folder
	http.Handle("/", http.FileServer(http.Dir("./static")))
	// Handle retrieval of a single jpeg image
	http.HandleFunc("/picture", getPicture)
	// Handle simple animation request
	http.HandleFunc("/animation", getAnimation)
	// Handle sine wave animation request
	http.HandleFunc("/wave", getSinewaves)

	http.HandleFunc("/screencast", getScreenCast)
	// Start a server on the port
	port := "8080"
	log.Fatal(http.ListenAndServe(":"+port, nil))

}

// getJPEG creates a single color image with given dimensions and color.
// Returns the image as a slice of jpeg bytes
func getJPEG(w int, h int, color color.RGBA) []byte {
	im := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.Draw(im, im.Bounds(), &image.Uniform{color}, image.ZP, draw.Src)
	var buff bytes.Buffer
	jpeg.Encode(&buff, im, nil)
	return buff.Bytes()
}

// getPicture sends jpeg image bytes over http as well as content description
// for browser to able to render the image properly
func getPicture(w http.ResponseWriter, r *http.Request) {
	imgBytes := getJPEG(200, 200, blue)
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(imgBytes)))
	w.Write(imgBytes)
}

// getAnimation creates sample images and sends them one after the other to client
func getAnimation(w http.ResponseWriter, r *http.Request) {
	const (
		// Size of images
		size = 200
		// Delay between frames in miliseconds
		delay = 500 * time.Millisecond
	)
	// To send buffered data to client right away
	f, ok := w.(http.Flusher)
	if !ok {
		log.Println("HTTP buffer flushing is not implemented")
	}
	// Sample images
	imgRed := getJPEG(size, size, red)
	imgYellow := getJPEG(size, size, yellow)
	imgGreen := getJPEG(size, size, green)
	// Set headers and content to send as a response
	w.Header().Set("Content-Type", "multipart/x-mixed-replace; boundary="+boundary)
	w.Write([]byte("\r\n--" + boundary + "\r\n"))
	w.Write([]byte("Content-Type: image/jpeg\r\nContent-Length: " + strconv.Itoa(len(imgRed)) + "\r\n\r\n"))
	w.Write(imgRed)
	w.Write([]byte("\r\n--" + boundary + "\r\n"))
	// Otherwise buffer will be flushed after handler exits or buffer maxsize is full
	f.Flush()
	// Delay

	time.Sleep(delay)
	w.Write([]byte("Content-Type: image/jpeg\r\nContent-Length: " + strconv.Itoa(len(imgYellow)) + "\r\n\r\n"))
	w.Write(imgYellow)
	w.Write([]byte("\r\n--" + boundary + "\r\n"))
	f.Flush()
	time.Sleep(delay)
	w.Write([]byte("Content-Type: image/jpeg\r\nContent-Length: " + strconv.Itoa(len(imgGreen)) + "\r\n\r\n"))
	w.Write(imgGreen)
	w.Write([]byte("\r\n--" + boundary + "\r\n"))
}



func getScreenCast(w http.ResponseWriter, r *http.Request)  {
	const (
		// Size of an image frame
		width  = 600
		height = 400
		// Number of frames in animation
		nframes = 33
		// Delay between frames in miliseconds
		delay = 150 * time.Millisecond
	)

	f, ok := w.(http.Flusher)
	if !ok {
		log.Println("HTTP buffer flushing is not implemented")
	}
	w.Header().Set("Content-Type", "multipart/x-mixed-replace; boundary="+boundary)
	// Start animation frame by frame

	agent := &capturer.Capturer{
		Delay:        0.1,
		Quality:      35,
		IsFullScreen: true,
	}


	for imgBytes := range agent.Capture().Output {
		println(">>>>> ", len(imgBytes))
		if agent.Count == 0 {
			w.Write([]byte("\r\n--" + boundary + "\r\n"))
		}
		w.Write([]byte("Content-Type: image/jpeg\r\nContent-Length: " + strconv.Itoa(len(imgBytes)) + "\r\n\r\n"))
		w.Write(imgBytes)
		w.Write([]byte("\r\n--" + boundary + "\r\n"))
		f.Flush()
		println(agent.Count)
	}

	agent.Await()


}


func getScreenCast4(w http.ResponseWriter, r *http.Request)  {
	const (
		// Size of an image frame
		width  = 600
		height = 400
		// Number of frames in animation
		nframes = 33
		// Delay between frames in miliseconds
		delay = 150 * time.Millisecond
	)

	f, ok := w.(http.Flusher)
	if !ok {
		log.Println("HTTP buffer flushing is not implemented")
	}
	w.Header().Set("Content-Type", "multipart/x-mixed-replace; boundary="+boundary)
	// Start animation frame by frame


	//path := utils2.GetCurrentDir()
	for {
		n := screenshot.NumActiveDisplays()
		s := 0
		for i := 0; i < n; i++ {
			for {
				bounds := screenshot.GetDisplayBounds(i)
				imgC, err := screenshot.CaptureRect(bounds)

				//imgC, err := screenshot.Capture(200,600, 500, 500)
				if err != nil {
					panic(err)
				}


				//dstImage800 := imaging.Resize(imgC, 1000, 0, imaging.Lanczos)

				// Encode image to JPEG and get it's representation in bytes
				var buff bytes.Buffer
				jpeg.Encode(&buff, imgC, &jpeg.Options{Quality: 35})
				imgBytes := buff.Bytes()

				println(">>>>>> ",len(imgBytes))

				if s == 0 {
					w.Write([]byte("\r\n--" + boundary + "\r\n"))
				}
				w.Write([]byte("Content-Type: image/jpeg\r\nContent-Length: " + strconv.Itoa(len(imgBytes)) + "\r\n\r\n"))
				w.Write(imgBytes)
				w.Write([]byte("\r\n--" + boundary + "\r\n"))
				// Otherwise buffer will be flushed after handler exits or buffer maxsize is full
				f.Flush()
				time.Sleep(time.Duration(1) * time.Second)

				s++
			}

		}
	}

}

func getScreenCast2(w http.ResponseWriter, r *http.Request)  {
	const (
		// Size of an image frame
		width  = 600
		height = 400
		// Number of frames in animation
		nframes = 33
		// Delay between frames in miliseconds
		delay = 150 * time.Millisecond
	)

	f, ok := w.(http.Flusher)
	if !ok {
		log.Println("HTTP buffer flushing is not implemented")
	}
	w.Header().Set("Content-Type", "multipart/x-mixed-replace; boundary="+boundary)
	// Start animation frame by frame
	for t := 0; t < nframes; t++ {
		// Create paletted image of the size and aplette
		path := utils2.GetCurrentDir()
		fileName := fmt.Sprintf("0_%d:1680x1050.jpeg", t)
		path = filepath.Join(path, "main", "mjpeg", fileName)
		println(path)
		img := img.Read(path)

		println(img)



		// Encode image to JPEG and get it's representation in bytes
		var buff bytes.Buffer
		jpeg.Encode(&buff, img, nil)
		imgBytes := buff.Bytes()
		// Stream image back to client
		// For the first frame we need to draw boundry at the beginning
		if t == 0 {
			w.Write([]byte("\r\n--" + boundary + "\r\n"))
		}
		w.Write([]byte("Content-Type: image/jpeg\r\nContent-Length: " + strconv.Itoa(len(imgBytes)) + "\r\n\r\n"))
		w.Write(imgBytes)
		w.Write([]byte("\r\n--" + boundary + "\r\n"))
		// Otherwise buffer will be flushed after handler exits or buffer maxsize is full
		f.Flush()
		time.Sleep(delay)
	}

}


// getSinewaves generates and streams animation of a sine wave
func getSinewaves(w http.ResponseWriter, r *http.Request) {
	var palette = []color.Color{color.White, blue}
	const (
		// First color in palette
		whiteIndex = 0
		blueIndex  = 1
	)
	const (
		// Size of an image frame
		width  = 600
		height = 400
		// Number of frames in animation
		nframes = 60000
		// Delay between frames in miliseconds
		delay = 50 * time.Millisecond
	)
	// To send buffered data to client right away
	f, ok := w.(http.Flusher)
	if !ok {
		log.Println("HTTP buffer flushing is not implemented")
	}
	w.Header().Set("Content-Type", "multipart/x-mixed-replace; boundary="+boundary)
	// Start animation frame by frame
	for t := 0; t < nframes; t++ {
		// Create paletted image of the size and aplette
		img := image.NewPaletted(image.Rect(0, 0, width, height), palette)
		// Animate sine wave and store it in our image
		// y = a*sin(b*x + c) + d
		for n := 0; n < width; n++ {
			// Draw from left to right
			x := float64(n)
			// Amplitude
			a := height / 3.0
			// Period is 2*pi / b
			b := 0.01
			// Phase shift
			c := float64(t) / 6.0
			// Vertical shift
			d := height / 2.0
			y := a*math.Sin(x*b+c) + d
			img.SetColorIndex(int(x), int(y), blueIndex)
		}
		// Encode image to JPEG and get it's representation in bytes
		var buff bytes.Buffer
		jpeg.Encode(&buff, img, nil)
		imgBytes := buff.Bytes()
		// Stream image back to client
		// For the first frame we need to draw boundry at the beginning
		if t == 0 {
			w.Write([]byte("\r\n--" + boundary + "\r\n"))
		}
		w.Write([]byte("Content-Type: image/jpeg\r\nContent-Length: " + strconv.Itoa(len(imgBytes)) + "\r\n\r\n"))
		w.Write(imgBytes)
		w.Write([]byte("\r\n--" + boundary + "\r\n"))
		// Otherwise buffer will be flushed after handler exits or buffer maxsize is full
		f.Flush()
		time.Sleep(delay)
	}
}
