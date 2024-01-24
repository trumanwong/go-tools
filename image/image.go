package image

import (
	"bufio"
	"image"
	"image/png"
	"io"
	"os"
)

func ReSave(reader io.Reader, savePath string) error {
	img, _, err := image.Decode(reader)
	if err != nil {
		return err
	}
	saveFile, err := os.Create(savePath)
	if err != nil {
		return err
	}
	defer saveFile.Close()
	b := bufio.NewWriter(saveFile)
	if err := png.Encode(b, img); err != nil {
		return err
	}
	return b.Flush()
}
