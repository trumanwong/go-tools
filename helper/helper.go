package helper

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"math/big"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
)

// CheckIdCard is a function that checks if a given string is a valid Chinese ID card number.
// It supports both 15-digit and 18-digit ID card numbers, and the last digit of 18-digit ID card numbers can be 'X' or 'x'.
//
// The function uses regular expressions to match the input string with the pattern of a valid ID card number.
// If the input string does not match the pattern, the function returns false.
//
// If the input string matches the pattern, the function then checks the last digit of the ID card number.
// For 18-digit ID card numbers, the last digit is a check digit that is calculated based on the first 17 digits.
// The function calculates the check digit and compares it with the actual last digit of the input string.
// If they match, the function returns true; otherwise, it returns false.
//
// Parameters:
// idCardStr: a string representing a Chinese ID card number.
//
// Returns:
// A boolean value indicating whether the input string is a valid Chinese ID card number.
func CheckIdCard(idCardStr string) bool {
	//18位身份证 ^(\d{17})([0-9]|X)$
	// 匹配规则
	// (^\d{15}$) 15位身份证
	// (^\d{18}$) 18位身份证
	// (^\d{17}(\d|X|x)$) 18位身份证 最后一位为X的用户
	regRuler := "(^\\d{15}$)|(^\\d{18}$)|(^\\d{17}(\\d|X|x)$)"

	// 正则调用规则
	reg := regexp.MustCompile(regRuler)

	// 返回 MatchString 是否匹配
	if !reg.MatchString(idCardStr) {
		return false
	}

	idCardBytes := []byte(idCardStr)

	// 通过前17位计算最后一位
	array := make([]int, 17)

	// 强制类型转换，将[]byte转换成[]int ,变化过程
	// []byte -> byte -> string -> int
	// 将通过range 将[]byte转换成单个byte,再用强制类型转换string()，将byte转换成string
	// 再通过strconv.Atoi()将string 转换成int 类型
	for index, value := range idCardBytes[0:17] {
		array[index], _ = strconv.Atoi(string(value))
	}

	var weight = [...]int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
	var res int
	for i := 0; i < 17; i++ {
		res += array[i] * weight[i]
	}

	lastByte := res % 11
	a18 := [11]byte{'1', '0', 'X', '9', '8', '7', '6', '5', '4', '3', '2'}
	if a18[lastByte] == idCardBytes[len(idCardBytes)-1] {
		return true
	}
	return false
}

func Response(ctx *gin.Context, data interface{}, code int, message string) {
	ctx.JSON(code, gin.H{
		"message": message,
		"data":    data,
	})
}

// PathExists 判断路径是否存在
func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func ConvertDownloadCount(downloadCount uint64) string {
	if downloadCount >= 100000000 {
		return fmt.Sprintf("%.2f亿次下载", float64(downloadCount)/100000000)
	} else if downloadCount >= 10000 {
		return fmt.Sprintf("%d万次下载", downloadCount/10000)
	}
	return fmt.Sprintf("%d次下载", downloadCount)
}

func IP2Long(ipAddress string) *big.Int {
	ip := net.ParseIP(ipAddress)
	if ip == nil {
		return nil
	}
	isIPv6 := false
	for i := 0; i < len(ipAddress); i++ {
		switch ipAddress[i] {
		case '.':
			break
		case ':':
			isIPv6 = true
			break
		}
	}
	ipInt := big.NewInt(0)
	if isIPv6 {
		return ipInt.SetBytes(ip.To16())
	}
	return ipInt.SetBytes(ip.To4())
}

// DownloadFile 下载文件
// url: 下载地址
// savePath: 保存路径
func DownloadFile(url, savePath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	dirPath := filepath.Dir(savePath)
	// 判断目录是否存在
	if !PathExists(dirPath) {
		err = os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			return err
		}
	}
	out, err := os.Create(savePath)
	if err != nil {
		return err
	}
	defer out.Close()
	// 写入文件
	_, err = io.Copy(out, resp.Body)
	return err
}

// InArray 判断元素是否在数组中
func InArray(needle interface{}, haystack interface{}) bool {
	val := reflect.ValueOf(haystack)
	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			if reflect.DeepEqual(needle, val.Index(i).Interface()) {
				return true
			}
		}
	case reflect.Map:
		for _, k := range val.MapKeys() {
			if reflect.DeepEqual(needle, val.MapIndex(k).Interface()) {
				return true
			}
		}
	default:
		panic("haystack: haystack type muset be slice, array or map")
	}

	return false
}
