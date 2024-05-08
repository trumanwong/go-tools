package aliyun

import (
	"encoding/json"
	"fmt"
	"github.com/trumanwong/go-tools/crawler"
	"net/http"
	"strings"
)

type ImageInfo struct {
	Compression                 Compression                 `json:"Compression"`
	DateTime                    DateTime                    `json:"DateTime"`
	ExifTag                     ExifTag                     `json:"ExifTag"`
	FileSize                    FileSize                    `json:"FileSize"`
	Format                      Format                      `json:"Format"`
	GPSLatitude                 GPSLatitude                 `json:"GPSLatitude"`
	GPSLatitudeRef              GPSLatitudeRef              `json:"GPSLatitudeRef"`
	GPSLongitude                GPSLongitude                `json:"GPSLongitude"`
	GPSLongitudeRef             GPSLongitudeRef             `json:"GPSLongitudeRef"`
	GPSMapDatum                 GPSMapDatum                 `json:"GPSMapDatum"`
	GPSTag                      GPSTag                      `json:"GPSTag"`
	GPSVersionID                GPSVersionID                `json:"GPSVersionID"`
	ImageHeight                 ImageHeight                 `json:"ImageHeight"`
	ImageWidth                  ImageWidth                  `json:"ImageWidth"`
	JPEGInterchangeFormat       JPEGInterchangeFormat       `json:"JPEGInterchangeFormat"`
	JPEGInterchangeFormatLength JPEGInterchangeFormatLength `json:"JPEGInterchangeFormatLength"`
	Orientation                 Orientation                 `json:"Orientation"`
	ResolutionUnit              ResolutionUnit              `json:"ResolutionUnit"`
	Software                    Software                    `json:"Software"`
	XResolution                 XResolution                 `json:"XResolution"`
	YResolution                 YResolution                 `json:"YResolution"`
}

type Compression struct {
	Value string `json:"value"`
}

type DateTime struct {
	Value string `json:"value"`
}

type ExifTag struct {
	Value string `json:"value"`
}

type FileSize struct {
	Value string `json:"value"`
}

type Format struct {
	Value string `json:"value"`
}

type GPSLatitude struct {
	Value string `json:"value"`
}

type GPSLatitudeRef struct {
	Value string `json:"value"`
}

type GPSLongitude struct {
	Value string `json:"value"`
}

type GPSLongitudeRef struct {
	Value string `json:"value"`
}

type GPSMapDatum struct {
	Value string `json:"value"`
}

type GPSTag struct {
	Value string `json:"value"`
}

type GPSVersionID struct {
	Value string `json:"value"`
}

type ImageHeight struct {
	Value string `json:"value"`
}

type ImageWidth struct {
	Value string `json:"value"`
}

type JPEGInterchangeFormat struct {
	Value string `json:"value"`
}

type JPEGInterchangeFormatLength struct {
	Value string `json:"value"`
}

type Orientation struct {
	Value string `json:"value"`
}

type ResolutionUnit struct {
	Value string `json:"value"`
}

type Software struct {
	Value string `json:"value"`
}

type XResolution struct {
	Value string `json:"value"`
}

type YResolution struct {
	Value string `json:"value"`
}

// GetOssImageInfo 获取oss图片信息(https://help.aliyun.com/zh/oss/user-guide/query-the-exif-data-of-an-image-4?spm=a2c4g.11186623.0.0.160d27595PqPne)
func GetOssImageInfo(link string) (*ImageInfo, error) {
	resp, err := crawler.Send(&crawler.Request{
		Url:    fmt.Sprintf("%s?x-oss-process=image/info", strings.Split(link, "?")[0]),
		Method: http.MethodGet,
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var imageInfo ImageInfo
	err = json.NewDecoder(resp.Body).Decode(&imageInfo)
	if err != nil {
		return nil, err
	}
	return &imageInfo, nil
}
