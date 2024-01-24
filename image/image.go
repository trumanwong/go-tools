package image

import (
	"bufio"
	"image"
	"image/png"
	"os"
)

func ReSave(imagePath, savePath string) error {
	f, err := os.Open(imagePath)
	if err != nil {
		return err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return err
	}
	saveFile, err := os.Create(savePath)
	if err != nil {
		return err
	}
	defer f.Close()
	b := bufio.NewWriter(saveFile)
	if err := png.Encode(b, img); err != nil {
		return err
	}
	return b.Flush()
}
