package aliyun

import (
	"os"
	"testing"

	"github.com/alibabacloud-go/tea/tea"
)

func TestGreenClient_ImageModeration(t *testing.T) {
	client, err := NewGreenClient(os.Getenv("AccessKeyId"), os.Getenv("AccessKeySecret"), os.Getenv("Endpoint"))
	if err != nil {
		t.Error(err)
		return
	}
	_, err = client.ImageModeration(&ImageModerationRequest{
		OssBucketName: tea.String(os.Getenv("Bucket")),
		OssObjectName: tea.String(os.Getenv("OssObjectName")),
		OssRegionID:   tea.String(os.Getenv("OssRegion")),
		Labels: map[string]float32{
			"pornographic_adultContent":              60,
			"pornographic_adultContent_tii":          60,
			"sexual_suggestiveContent":               60,
			"sexual_partialNudity":                   60,
			"political_historicalNihility":           60,
			"political_historicalNihility_tii":       60,
			"political_politicalFigure":              60,
			"political_politicalFigure_name_tii":     60,
			"political_politicalFigure_metaphor_tii": 60,
			"political_prohibitedPerson":             60,
			"political_prohibitedPerson_tii":         60,
			"political_taintedCelebrity":             60,
			"political_taintedCelebrity_tii":         60,
			"political_flag":                         60,
			"political_map":                          60,
			"political_religion_tii":                 60,
			"contraband_gamble":                      60,
			"contraband_gamble_tii":                  60,
		},
	})
	if err != nil {
		t.Error(err)
	}
}
