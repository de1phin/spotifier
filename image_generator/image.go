package image_generator

import (
	"image"
	"net/http"
	"io/ioutil"
	"image/color"
	"log"

	"github.com/disintegration/imaging"
	"github.com/golang/freetype"
)

var (
	c *freetype.Context
)

const (
	FontFilename string = "../image_generator/arial.ttf"
)

func LoadImage(path string) (image.Image, error) {
	return imaging.Open(path)
}

func SaveImage(img image.Image, path string) error {
	return imaging.Save(img, path)
}

func Paste(background, img image.Image, x, y int) *image.NRGBA {
	return imaging.Paste(background, img, image.Point{X: x, Y: y})
}

func Resize(img image.Image, width, height int) *image.NRGBA {
	return imaging.Resize(img, width, height, imaging.Lanczos)
}

func DownloadImage(url string) (image.Image, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return imaging.Decode(resp.Body)
}

func WriteText(img *image.NRGBA, x, y, font_size int, text string, clr color.Color) {
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(image.NewUniform(clr))
	c.SetFontSize(float64(font_size))

	pt := freetype.Pt(x, y + int(c.PointToFixed(float64(font_size)) >> 6))
	_, err := c.DrawString(text, pt)
	if err != nil {
		log.Println(err)
	}
}

func InitImageGenerator() {
	c = freetype.NewContext()

	fontBytes, err := ioutil.ReadFile(FontFilename)
	if err != nil {
		log.Println(err)
		return
	}

	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println(err)
		return
	}

	c.SetDPI(96)
	c.SetFont(font)
	c.SetFontSize(72)
}