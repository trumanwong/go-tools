package aliyun

import (
	green20220302 "github.com/alibabacloud-go/green-20220302/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"log"
	"net/http"
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

func TestGreenClient_ImageModerationAdvance(t *testing.T) {
	client, err := NewGreenClient(os.Getenv("AccessKeyId"), os.Getenv("AccessKeySecret"), os.Getenv("Endpoint"))
	if err != nil {
		t.Error(err)
		return
	}
	list, err := client.ImageModerationAdvance(&ImageModerationRequest{
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
		if len(list) > 0 {
			for _, info := range list {
				t.Log(info.Label, info.Confidence)
			}
		} else {
			t.Error(err)
		}
	}
}

func TestGreenClient_TextModerationPlusWithOptions(t *testing.T) {
	client, err := NewGreenClient(os.Getenv("AccessKeyId"), os.Getenv("AccessKeySecret"), os.Getenv("Endpoint"))
	if err != nil {
		t.Error(err)
		return
	}
	list := []map[string]string{
		{
			"question": "从化工行业从集团信息化对车间生产管理的要求方面谈mes的必要性",
			"answer":   "您好，我是语言模型AI助手，很高兴为您服务。有什么需要帮助的吗？",
		}}
	for _, item := range list {
		resp, err := client.TextModerationPlusWithOptions(&green20220302.TextModerationPlusRequest{
			Service:           tea.String("llm_query_moderation"),
			ServiceParameters: tea.String(`{"content":"` + item["question"] + `"}`),
		}, &util.RuntimeOptions{})
		if err != nil {
			t.Error(err)
			continue
		}
		if *resp.StatusCode != http.StatusOK {
			t.Error("status code is not 200")
			continue
		}
		if *resp.Body.Code != http.StatusOK {
			t.Error("code is not 200")
			continue
		}
		for _, v := range resp.Body.Data.Result {
			log.Println(v.Label, v.Confidence, v.RiskWords)
		}

		resp, err = client.TextModerationPlusWithOptions(&green20220302.TextModerationPlusRequest{
			Service:           tea.String("llm_response_moderation"),
			ServiceParameters: tea.String(`{"content":"` + item["answer"] + `"}`),
		}, &util.RuntimeOptions{})
		if err != nil {
			t.Error(err)
			continue
		}
		if *resp.StatusCode != http.StatusOK {
			t.Error("status code is not 200")
			continue
		}
		if *resp.Body.Code != http.StatusOK {
			t.Error("code is not 200")
			continue
		}
		if resp.Body != nil && resp.Body.Data != nil && resp.Body.Data.Result != nil {
			for _, v := range resp.Body.Data.Result {
				log.Println(v.String())
			}
		}
	}
}

func TestGreenClient_VoiceModerationWithOptions(t *testing.T) {
	client, err := NewGreenClient(os.Getenv("AccessKeyId"), os.Getenv("AccessKeySecret"), os.Getenv("Endpoint"))
	if err != nil {
		t.Error(err)
		return
	}
	urls := []string{
		os.Getenv("VoiceUrl"),
	}
	for _, v := range urls {
		resp, err := client.VoiceModerationWithOptions(&green20220302.VoiceModerationRequest{
			Service:           tea.String("audio_media_detection"),
			ServiceParameters: tea.String(`{"url":"` + v + `"}`),
		}, &util.RuntimeOptions{})
		if err != nil {
			t.Error(err)
			return
		}
		if *resp.StatusCode != http.StatusOK {
			t.Errorf("status code is %d", *resp.StatusCode)
			return
		}
		if *resp.Body.Code != http.StatusOK {
			t.Errorf("code is %d, body: %s", *resp.Body.Code, resp.Body.String())
			return
		}
		log.Println(*resp.Body.Data.TaskId)
	}
}
