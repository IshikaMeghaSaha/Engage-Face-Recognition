package main

import (
	"fmt"
	"html/template"
	"log"

	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	"github.com/hybridgroup/mjpeg"
	"gocv.io/x/gocv"
)

var (
	deviceID int
	err      error
	webcam   *gocv.VideoCapture
	stream   *mjpeg.Stream
)

func home(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("./page.html"))
	t.Execute(w, t)

	r.ParseForm()
	if r.Method == "POST" {

		go Stream(w, r)

	}

}

func Stream(w http.ResponseWriter, r *http.Request) {

	deviceID := 0

	// open webcam

	webcam, err = gocv.OpenVideoCapture(deviceID)
	if err != nil {
		fmt.Printf("Error opening capture device: %v\n", deviceID)
		return
	}

	// create the mjpeg stream
	stream = mjpeg.NewStream()

	// start capturing

	go mjpegCapture()

	//fmt.Println("Capturing. Point your browser to " + host)

	// start http server
	http.Handle("/stream", stream)
	time.Sleep(10 * time.Second)
	os.Exit(0)

}

func main() {
	host := "0.0.0.0:8000"

	http.HandleFunc("/", home)

	log.Fatal(http.ListenAndServe(host, nil))

}

func mjpegCapture() {
	img := gocv.NewMat()
	defer img.Close()
	defer webcam.Close()
	for {
		if ok := webcam.Read(&img); !ok {
			fmt.Printf("Device closed: %v\n", deviceID)
			return
		}
		if img.Empty() {
			continue
		}

		buf, _ := gocv.IMEncode(".jpg", img)
		gocv.IMWrite("../comp/img1.jpg", img)
		stream.UpdateJPEG(buf.GetBytes())
		buf.Close()
	}
}
