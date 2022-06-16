package main

import (
	"image"
	"image/color"
	imggen "../image_generator"
	"log"
)

const (
	ImageFolder string = "../telegram_bot/img/"
	BackgroundImage string = "../telegram_bot/img/background.png"
)

type ImgItem struct {
	ImgURL			string
	Title 			string
	Subtitle		string
	Additional		string
}

func GenerateRatingImage(items []ImgItem) {
	background, err := imggen.LoadImage(BackgroundImage)
	if err != nil {
		log.Println("GenerateRatingImageError:", err)
		return
	}

	img := background.(*image.NRGBA)

	for i, item := range items {
		pasteImgRaw, err := imggen.DownloadImage(item.ImgURL)
		pasteImg := imggen.Resize(pasteImgRaw, 300, 300)
		if err != nil {
			log.Println("GenerateRatingImageError:", err)
			return
		}

		x := 60
		y := 60 + 360 * i
		img = imggen.Paste(img, pasteImg, x, y)
		
		if item.Subtitle == "" {
			if len(item.Title) > 14 {
				item.Title = item.Title[0: 13] + "..."
			}
			if len(item.Additional) > 23 {
				item.Additional = item.Additional[0: 22] + "..."
			}
			imggen.WriteText(img, x + 330, y + 25, 64, item.Title, color.White)
			imggen.WriteText(img, x + 330, y + 180, 50, item.Additional, color.Opaque)
		} else {
			if len(item.Title) > 15 {
				item.Title = item.Title[0:14] + "..."
			}
			imggen.WriteText(img, x + 330, y - 5, 54, item.Title, color.White)
			if len(item.Subtitle) > 20 {
				item.Subtitle = item.Subtitle[0:19] + "..."
			}
			imggen.WriteText(img, x + 330, y + 120, 42, item.Subtitle, color.Opaque)
			if len(item.Additional) > 20 {
				item.Additional = item.Additional[0:19] + "..."
			}
			imggen.WriteText(img, x + 330, y + 215, 42, item.Additional, color.Opaque)
		}
	}
	
	log.Println("Saving the image")
	err = imggen.SaveImage(img, ImageFolder + "stat.png")
	if err != nil {
		log.Println("GenerateRatingImageError:", err)
	}
}

func InitImgGen() {
	imggen.InitImageGenerator()
}