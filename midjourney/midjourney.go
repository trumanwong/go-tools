package midjourney

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/trumanwong/go-tools/helper"
)

// GetPromptAndParametersRequest is a struct that represents the request for the GetPromptAndParameters function.
// It contains two fields: Content and DisableParams.
// Content is a string that represents the content to be processed.
// DisableParams is a map that represents the parameters to be disabled.
type GetPromptAndParametersRequest struct {
	Content       string
	DisableParams []string
}

// GetPromptAndParametersResponse is a struct that represents the response from the GetPromptAndParameters function.
// It contains two fields: Prompt and Parameters.
// Prompt is a string that represents the prompt extracted from the content.
// Parameters is a map that represents the parameters extracted from the content.
type GetPromptAndParametersResponse struct {
	Prompt     string
	Parameters map[string]string
}

// GetPromptAndParameters is a function that takes a GetPromptAndParametersRequest pointer as input and returns a GetPromptAndParametersResponse pointer and an error.
// It processes the content from the request, extracts the prompt and parameters, and returns them in the response.
// If the content is empty, it returns an error.
// If the content does not contain any parameters, it returns the prompt with nil parameters.
// It also validates the parameters based on predefined rules. If a parameter does not meet the rules, it returns an error.
func GetPromptAndParameters(req *GetPromptAndParametersRequest) (*GetPromptAndParametersResponse, error) {
	var prompt string
	var parameters = make(map[string]string)
	content := strings.TrimSpace(req.Content)
	if content == "" {
		return nil, errors.New("content is empty")
	}

	arr := strings.Split(content, "--")
	prompt = arr[0]
	if len(arr) <= 1 {
		return &GetPromptAndParametersResponse{
			Prompt:     prompt,
			Parameters: nil,
		}, nil
	}

	// 获取参数
	for i := 1; i < len(arr); i++ {
		paramValue := strings.Split(arr[i], " ")
		if len(paramValue) <= 1 {
			paramValue = append(paramValue, "")
		}
		param := strings.TrimSpace(strings.ToLower(paramValue[0]))
		if req.DisableParams != nil && helper.InArray(param, req.DisableParams) {
			continue
		}
		val := strings.TrimSpace(strings.Join(paramValue[1:], " "))
		switch param {
		case "aspect", "ar":
			parameters["aspect"] = val
		case "chaos", "c":
			temp, err := strconv.ParseInt(val, 10, 64)
			if err != nil || temp < 0 || temp > 100 {
				return nil, errors.New("chaos参数值范围必须在0~100之间")
			}
			parameters["chaos"] = val
		case "iw":
			temp, err := strconv.ParseFloat(val, 64)
			if err != nil || temp < 0 || temp > 2 {
				return nil, errors.New("iw参数值范围必须在0~2之间")
			}
			parameters["iw"] = val
		case "quality", "q":
			temp, err := strconv.ParseFloat(val, 64)
			if err != nil || temp < 0 || temp > 2 {
				return nil, errors.New("quality值范围必须在0~2之间")
			}
			parameters["quality"] = val
		case "repeat", "r":
			temp, err := strconv.ParseInt(val, 10, 64)
			if err != nil || temp < 1 || temp > 40 {
				return nil, fmt.Errorf("%s参数值范围必须在0~40之间", param)
			}
			parameters["repeat"] = val
		case "seed":
			temp, err := strconv.ParseInt(val, 10, 64)
			if err != nil || temp < 0 || temp > 4294967295 {
				return nil, fmt.Errorf("%s参数值范围必须在0~4294967295之间", param)
			}
			parameters["seed"] = val
		case "stop":
			temp, err := strconv.ParseInt(val, 10, 64)
			if err != nil || temp < 10 || temp > 100 {
				return nil, fmt.Errorf("%s参数值范围必须在10~100之间", param)
			}
			parameters["stop"] = val
		case "stylize", "s":
			temp, err := strconv.ParseInt(val, 10, 64)
			if err != nil || temp < 0 || temp > 1000 {
				return nil, fmt.Errorf("%s参数值范围必须在0~1000之间", param)
			}
			parameters["stylize"] = val
		case "weird", "w":
			temp, err := strconv.ParseInt(val, 10, 64)
			if err != nil || temp < 0 || temp > 3000 {
				return nil, fmt.Errorf("%s参数值范围必须在0~3000之间", param)
			}
			parameters["weird"] = val
		case "niji":
			if val != "4" && val != "5" && val != "6" {
				return nil, errors.New("niji参数值范围必须是4、5或6")
			}
			parameters["niji"] = val
		case "version", "v":
			temp, err := strconv.ParseFloat(val, 10)
			tempVal := int(temp * 10)
			// "1", "2", "3", "4", "5.0", "5.1", "5.2", "6", "6.1"
			if err != nil || (temp != 40 && temp != 50 && tempVal != 51 && tempVal != 52 && tempVal != 60 && tempVal != 61) {
				return nil, fmt.Errorf("%s参数值必须是4, 5, 5.1, 5.2，6，6.1", param)
			}

			parameters["version"] = val
		case "cw":
			temp, err := strconv.ParseInt(val, 10, 64)
			if err != nil || temp < 0 || temp > 100 {
				return nil, fmt.Errorf("%s参数值范围必须在0~100之间", param)
			}
			parameters["cw"] = val
		case "sref":
			links := strings.Split(val, " ")
			for _, link := range links {
				if link == "" {
					return nil, fmt.Errorf("%s参数值不能为空", param)
				}
				u, err := url.Parse(link)
				if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
					return nil, fmt.Errorf("%s参数值必须是一个有效的URL", param)
				}
			}
			parameters[param] = val
		case "cref":
			if val == "" {
				return nil, fmt.Errorf("%s参数值不能为空", param)
			}
			u, err := url.Parse(val)
			if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
				return nil, fmt.Errorf("%s参数值必须是一个有效的URL", param)
			}
			parameters[param] = val
		case "p", "personalize":
			parameters["p"] = val
		case "no", "style":
			parameters[param] = val
		default:
			if helper.InArray(param, []string{"tile", "relax", "fast", "turbo"}) {
				parameters[param] = ""
			}
		}
	}
	return &GetPromptAndParametersResponse{
		Prompt:     prompt,
		Parameters: parameters,
	}, nil
}
