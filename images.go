package butils

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"math"
	"net/http"
	"os"
	"sync"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/math/fixed"
)

type PlotDataModel struct {
	Rect  image.Rectangle
	Label string
}

func drawRectangle(img draw.Image, color color.Color, x1, y1, x2, y2, thickness int, label string) {
	wg := new(sync.WaitGroup)
	wg.Add(thickness * (((x2 - x1) * 2) + ((y2 - y1 + 1) * 2)))

	for t := 0; t < thickness; t++ {
		for i := x1; i < x2; i++ {
			go func(i, t int) {
				defer wg.Done()
				img.Set(i, y1+t, color)
			}(i, t)
			go func(i, t int) {
				defer wg.Done()
				img.Set(i, y2-t, color)
			}(i, t)
		}

		for i := y1; i <= y2; i++ {
			go func(i, t int) {
				defer wg.Done()
				img.Set(x1+t, i, color)
			}(i, t)
			go func(i, t int) {
				defer wg.Done()
				img.Set(x2-t, i, color)
			}(i, t)
		}
	}

	// draw label
	if label != "" {
		f, _ := truetype.Parse(goregular.TTF)
		d := &font.Drawer{
			Dst: img,
			Src: image.NewUniform(color),
			Face: truetype.NewFace(f, &truetype.Options{
				Size: float64(thickness) * 8,
			}),
			Dot: fixed.Point26_6{X: fixed.Int26_6(x1 * 64), Y: fixed.Int26_6((y1 - thickness) * 64)},
		}
		d.DrawString(label)
	}

	wg.Wait()
}

func addRectangleToFace(img draw.Image, rect image.Rectangle, label string) draw.Image {
	// กำหนดสีที่ใช้วาด
	myColor := color.RGBA{255, 0, 0, 255}

	min := rect.Min
	max := rect.Max

	// กำหนดความหนาเส้น
	thickness := math.Min(float64(max.X-min.X), float64(max.Y-min.Y)) * 0.01
	if thickness < 1 {
		thickness = 1
	}

	drawRectangle(img, myColor, min.X, min.Y, max.X, max.Y, int(thickness), label)

	return img
}

func getImageFromFilePath(filePath string) (draw.Image, string, error) {

	// read file
	f, err := os.Open(filePath)
	if err != nil {
		return nil, "", err
	}
	defer f.Close()

	// convert as image.Image
	orig, typeImage, err := image.Decode(f)
	if err != nil {
		return nil, "", err
	}

	// สร้าง instance image สำหรับแก้ไขไฟล์ภาพ
	b := orig.Bounds()                                     // ดึงขนาด orig img
	img := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy())) // กำหนด instance ขนาดเท่ากับ orig img
	draw.Draw(img, img.Bounds(), orig, b.Min, draw.Src)    // copy orig image ใส่ใน instance ขนาดตามที่กำหนด

	return img, typeImage, err
}

func getImageFromUrl(url string) (draw.Image, string, error) {
	// Read image from url
	res, err := http.Get(url)
	if err != nil {
		return nil, "", err
	}
	defer res.Body.Close()

	// Convert file to byte
	data, _ := io.ReadAll(res.Body)

	// convert as image.Image
	orig, typeImage, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, "", err
	}

	// สร้าง instance image สำหรับแก้ไขไฟล์ภาพ
	b := orig.Bounds()                                     // ดึงขนาด orig img
	img := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy())) // กำหนด instance ขนาดเท่ากับ orig img
	draw.Draw(img, img.Bounds(), orig, b.Min, draw.Src)    // copy orig image ใส่ใน instance ขนาดตามที่กำหนด

	return img, typeImage, err
}

func PlotImageFromUrl(url string, plotData []PlotDataModel) (result []byte, err error) {
	// read file and convert it
	src, t, err := getImageFromUrl(url)
	if err != nil {
		return nil, err
	}

	var dst draw.Image
	for idx, p := range plotData {
		if idx == 0 {
			dst = addRectangleToFace(src, p.Rect, p.Label)
			continue
		}
		dst = addRectangleToFace(dst, p.Rect, p.Label)
	}

	buf := new(bytes.Buffer)
	switch t {
	case "jpeg":
		jpeg.Encode(buf, dst, nil)
	case "png":
		png.Encode(buf, dst)
	}

	return buf.Bytes(), nil
}

func PlotImageFromBytes(data []byte, plotData []PlotDataModel) (result []byte, err error) {
	// convert as image.Image
	orig, t, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	// สร้าง instance image สำหรับแก้ไขไฟล์ภาพ
	b := orig.Bounds()                                     // ดึงขนาด orig img
	img := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy())) // กำหนด instance ขนาดเท่ากับ orig img
	draw.Draw(img, img.Bounds(), orig, b.Min, draw.Src)    // copy orig image ใส่ใน instance ขนาดตามที่กำหนด

	var dst draw.Image
	for idx, p := range plotData {
		if idx == 0 {
			dst = addRectangleToFace(img, p.Rect, p.Label)
			continue
		}
		dst = addRectangleToFace(dst, p.Rect, p.Label)
	}

	buf := new(bytes.Buffer)
	switch t {
	case "jpeg":
		jpeg.Encode(buf, dst, nil)
	case "png":
		png.Encode(buf, dst)
	}

	return buf.Bytes(), nil
}

func PlotImageFromDir(filePath string, plotData []PlotDataModel) (result []byte, err error) {
	// convert as image.Image
	src, t, err := getImageFromFilePath(filePath)
	if err != nil {
		return nil, err
	}

	var dst draw.Image
	for idx, p := range plotData {
		if idx == 0 {
			dst = addRectangleToFace(src, p.Rect, p.Label)
			continue
		}
		dst = addRectangleToFace(dst, p.Rect, p.Label)
	}

	buf := new(bytes.Buffer)
	switch t {
	case "jpeg":
		jpeg.Encode(buf, dst, nil)
	case "png":
		png.Encode(buf, dst)
	}

	return buf.Bytes(), nil
}
