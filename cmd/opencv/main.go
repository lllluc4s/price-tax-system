package main

import "gocv.io/x/gocv"

func main(){
	webcam, err := gocv.VideoCaptureDevice(0)
	if err != nil {
		println("Error opening webcam")
		return
	}
	defer webcam.Close()

	window := gocv.NewWindow("Detector")
	defer window.Close()

	classifier := gocv.NewCascadeClassifier()
	defer classifier.Close()

	classifier.Load("haarcascade_frontalface_default.xml")

	for{
		img := gocv.NewMat()

		if ok := webcam.Read(&img); !ok {
			println("Cannot read device")
			return
		}
		if img.Empty() {
			continue
		}

		myface := classifier.DetectMultiScale(img)
		for _, r := range myface {
			gocv.Rectangle(&img, r, gocv.RGB(0, 255, 0), 3)
		}

		window.IMShow(img)
		if window.WaitKey(1) >= 0 {
			break
		}
	}
}