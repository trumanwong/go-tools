package helper

import (
	"crypto/sha1"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/trumanwong/go-tools/crawler"
	"io"
	"math"
	"math/big"
	"math/rand"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
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

// Response is a function that sends a JSON response to the client.
//
// The function uses the gin package to send a JSON response with a given HTTP status code, message, and data.
// The response is a JSON object with two properties: "message" and "data".
// The "message" property is a string that represents the message to be sent to the client.
// The "data" property is an interface{} that represents the data to be sent to the client.
//
// Parameters:
// ctx: a pointer to a gin.Context that represents the context of the request.
// data: an interface{} that represents the data to be sent to the client.
// code: an int that represents the HTTP status code of the response.
// message: a string that represents the message to be sent to the client.
//
// Returns:
// The function does not return a value.
func Response(ctx *gin.Context, data interface{}, code int, message string) {
	ctx.JSON(code, gin.H{
		"message": message,
		"data":    data,
	})
}

// PathExists is a function that checks if a given path exists in the file system.
//
// The function uses the os package to get the file or directory information of the given path.
// If there is an error in getting the information, the function checks if the error is because the file or directory exists.
// If the file or directory exists, the function returns true; otherwise, it returns false.
//
// If there is no error in getting the information, the function returns true, indicating that the path exists.
//
// Parameters:
// path: a string representing the path to be checked.
//
// Returns:
// A boolean value indicating whether the given path exists in the file system.
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

// ConvertDownloadCount is a function that converts a download count to a string representation.
// The function supports download counts up to 100 million and above.
//
// The function checks if the download count is greater than or equal to 100 million.
// If it is, the function converts the download count to a float, divides it by 100 million, and formats it as a string with two decimal places followed by "亿次下载".
//
// If the download count is less than 100 million but greater than or equal to 10,000, the function divides it by 10,000 and formats it as a string followed by "万次下载".
//
// If the download count is less than 10,000, the function formats it as a string followed by "次下载".
//
// Parameters:
// downloadCount: a uint64 representing the download count.
//
// Returns:
// A string representing the download count in a more readable format.
func ConvertDownloadCount(downloadCount uint64) string {
	if downloadCount >= 100000000 {
		return fmt.Sprintf("%.2f亿次下载", float64(downloadCount)/100000000)
	} else if downloadCount >= 10000 {
		return fmt.Sprintf("%d万次下载", downloadCount/10000)
	}
	return fmt.Sprintf("%d次下载", downloadCount)
}

// IP2Long is a function that converts an IP address to a big integer.
// The function supports both IPv4 and IPv6 addresses.
//
// The function uses the net package to parse the input string into an IP address.
// If the input string is not a valid IP address, the function returns nil.
//
// The function then checks if the IP address is an IPv6 address by looking for a colon in the input string.
// If the IP address is an IPv6 address, the function converts it to a 16-byte representation and sets it to a big integer.
// If the IP address is an IPv4 address, the function converts it to a 4-byte representation and sets it to a big integer.
//
// Parameters:
// ipAddress: a string representing an IP address.
//
// Returns:
// A pointer to a big integer representing the IP address.
// If the input string is not a valid IP address, the function returns nil.
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

// DownloadFile is a function that downloads a file from a given URL and saves it to a specified path.
// It uses the http package to send a GET request to the URL and receive the response.
// If there is an error in sending the request or receiving the response, the function returns the error.
//
// The function then gets the directory path of the save path.
// If the directory does not exist, the function creates it.
// If there is an error in creating the directory, the function returns the error.
//
// The function then creates a new file at the save path.
// If there is an error in creating the file, the function returns the error.
//
// The function then writes the body of the response to the file.
// If there is an error in writing to the file, the function returns the error.
//
// Parameters:
// request: a pointer to a crawler.Request representing the request to be sent.
// savePath: a string representing the path where the file is to be saved.
// checkContentLength: a boolean indicating whether to check if the downloaded file size matches the content length.
//
// Returns:
// The size of the downloaded file and an error if there was a problem in downloading or saving the file.
func DownloadFile(request *crawler.Request, savePath string, checkContentLength bool) (int64, error) {
	// Send a GET request to the URL.
	resp, err := crawler.Send(request)
	if err != nil {
		return 0, errors.New("failed to send request: " + err.Error())
	}
	// Ensure the response body is closed after the function returns.
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("status code not 200, status: %d", resp.StatusCode)
	}

	// Write the body of the response to the file.
	size, err := SaveFile(resp.Body, savePath)
	if err != nil {
		return 0, errors.New("failed to write to file: " + err.Error())
	}
	// If checkContentLength is true, compare the downloaded file size with the content length.
	if checkContentLength && size != resp.ContentLength {
		return size, fmt.Errorf("downloaded file size [%d] does not match content length [%d]", size, resp.ContentLength)
	}
	return size, nil
}

// InArray is a function that checks if a given element (needle) is present in a given collection (haystack).
// The function supports collections of type slice, array, and map.
//
// The function uses the reflect package to get the value of the haystack and its kind.
// If the kind of the haystack is either a slice or an array, the function iterates over the elements of the haystack.
// For each element, it checks if the element is deeply equal to the needle.
// If it finds a match, it returns true.
//
// If the kind of the haystack is a map, the function iterates over the keys of the map.
// For each key, it checks if the value associated with the key is deeply equal to the needle.
// If it finds a match, it returns true.
//
// If the kind of the haystack is neither a slice, an array, nor a map, the function panics with a message.
//
// If the function does not find a match after checking all elements or values, it returns false.
//
// Parameters:
// needle: an element to be searched in the haystack.
// haystack: a collection where the needle is to be searched.
//
// Returns:
// A boolean value indicating whether the needle is present in the haystack.
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

// CheckPort is a function that checks if a given port on a given IP address is open.
// It uses the net.DialTimeout function from the net package to attempt to establish a connection to the specified IP address and port within the specified timeout duration.
// If the connection is successful, the function closes the connection and returns nil, indicating that the port is open.
// If the connection is not successful, the function returns the error returned by net.DialTimeout, indicating that the port is not open or that there was a problem in establishing the connection.
//
// Parameters:
// ip: a string representing the IP address to check.
// port: a string representing the port to check.
// timeout: a time.Duration representing the maximum amount of time to wait for the connection to be established.
//
// Returns:
// An error if the connection could not be established within the specified timeout duration or if there was a problem in establishing the connection; otherwise, nil.
func CheckPort(ip, port string, timeout time.Duration) error {
	conn, err := net.DialTimeout("tcp", ip+":"+port, timeout)
	if err != nil {
		return err
	}
	defer conn.Close()
	return nil
}

// CheckHttp is a function that checks if a given URL is accessible via HTTP.
// It uses the http.Head function from the net/http package to send a HEAD request to the URL.
// A HEAD request is similar to a GET request, but it only requests the headers and not the body of the response.
// This makes the function efficient for checking if a URL is accessible without downloading the entire content.
//
// If the HEAD request is successful, the function returns nil, indicating that the URL is accessible.
// If the HEAD request is not successful, the function returns the error returned by http.Head, indicating that the URL is not accessible or that there was a problem in sending the request.
//
// Parameters:
// link: a string representing the URL to check.
//
// Returns:
// An error if the HEAD request could not be sent or if the URL is not accessible; otherwise, nil.
func CheckHttp(link string, timeout time.Duration) (*http.Response, error) {
	return crawler.Send(&crawler.Request{
		Url:     link,
		Method:  http.MethodHead,
		Timeout: timeout,
	})
}

// GenerateShortUrl GenerateShortUrlKey is a function that generates a short URL key from a given URL.
func GenerateShortUrl(shortLinkPrefix string, link string) (string, error) {
	if shortLinkPrefix == "" {
		return "", fmt.Errorf("short link prefix is empty")
	}
	hasher := sha1.New()
	hasher.Write([]byte(link))
	sha := hasher.Sum(nil)
	shortUrl := fmt.Sprintf("%x", sha)
	return strings.TrimLeft(shortLinkPrefix, "/") + "/" + shortUrl[:8], nil
}

func PaginateData(list interface{}, total int64, page, pageSize int) map[string]interface{} {
	if page <= 0 {
		page = 1
	}

	switch {
	case pageSize > 100:
		pageSize = 100
	case pageSize <= 0:
		pageSize = 10
	}
	res := map[string]interface{}{
		"list":  list,
		"total": total,
	}
	res["current_page"] = page
	res["first_page"] = 1
	res["page_size"] = pageSize
	res["last_page"] = int64(math.Ceil(float64(total) / float64(pageSize)))
	return res
}

// ShuffleArray is a function that shuffles the elements of a given array.
// The function uses the math/rand package to generate a new random source based on the current time.
// It then uses this random source to shuffle the elements of the array.
//
// The function uses a type parameter T, which means it can be used with arrays of any type.
// The function takes an array of type T as a parameter and returns an array of type T.
//
// The function uses a closure in the call to the Shuffle method.
// This closure swaps the elements at the provided indices.
//
// Parameters:
// arr: an array of type T representing the array to be shuffled.
//
// Returns:
// An array of type T representing the shuffled array.
func ShuffleArray[T any](arr []T) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(arr), func(i, j int) {
		arr[i], arr[j] = arr[j], arr[i]
	})
}

// SaveFile is a function that saves the content read from an io.Reader to a specified file path.
// The function first checks if the directory of the save path exists, if not, it creates it.
// Then, it creates a new file at the save path.
// After that, it writes the content read from the io.Reader to the file.
// If there is an error at any step, the function returns the error and the number of bytes written so far.
//
// Parameters:
// reader: an io.Reader from which the content is read.
// savePath: a string representing the path where the content is to be saved.
//
// Returns:
// The number of bytes written to the file and an error if there was a problem in creating the directory, creating the file, or writing to the file.
func SaveFile(reader io.Reader, savePath string) (int64, error) {
	// Get the directory path of the save path.
	dirPath := filepath.Dir(savePath)
	// If the directory does not exist, create it.
	if !PathExists(dirPath) {
		err := os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			return 0, errors.New("failed to create directory: " + err.Error())
		}
	}
	// Create a new file at the save path.
	out, err := os.Create(savePath)
	if err != nil {
		return 0, errors.New("failed to create file: " + err.Error())
	}
	// Ensure the file is closed after the function returns.
	defer out.Close()
	// Write the body of the response to the file.
	size, err := io.Copy(out, reader)
	if err != nil {
		return 0, errors.New("failed to write to file: " + err.Error())
	}
	return size, nil
}

func GetSSLExpireDate(domain string) (*time.Time, error) {
	conn, err := tls.Dial("tcp", domain+":443", nil)
	if err != nil {
		return nil, fmt.Errorf("server doesn't support SSL certificate err:: %v", err)
	}
	defer conn.Close()
	err = conn.VerifyHostname(domain)
	if err != nil {
		return nil, fmt.Errorf("hostname doesn't match with certificate: : %v", err)
	}
	certs := conn.ConnectionState().PeerCertificates
	if len(certs) == 0 {
		return nil, fmt.Errorf("no certificate found")
	}
	return &certs[0].NotAfter, nil
}

// FormatByte converts a byte value to a human-readable string representation in bytes (B) or kilobytes (KB).
//
// Parameters:
// - data: The byte value to be formatted.
// - dividend: The threshold value to determine if the byte value should be converted to kilobytes.
//
// Returns:
// - A string representing the formatted byte value.
func FormatByte(data, dividend float64) string {
	if data >= dividend {
		return FormatKB(data/dividend, dividend)
	}
	return fmt.Sprintf("%.2fB", data+0.0000000001)
}

// FormatKB converts a byte value to a human-readable string representation in kilobytes (KB) or megabytes (MB).
//
// Parameters:
// - data: The byte value to be formatted.
// - dividend: The threshold value to determine if the byte value should be converted to megabytes.
//
// Returns:
// - A string representing the formatted kilobyte value.
func FormatKB(data, dividend float64) string {
	if data >= dividend {
		return FormatMB(data/dividend, dividend)
	}
	return fmt.Sprintf("%.2fKB", data+0.0000000001)
}

// FormatMB converts a byte value to a human-readable string representation in megabytes (MB) or gigabytes (GB).
//
// Parameters:
// - data: The byte value to be formatted.
// - dividend: The threshold value to determine if the byte value should be converted to gigabytes.
//
// Returns:
// - A string representing the formatted megabyte value.
func FormatMB(data, dividend float64) string {
	if data >= dividend {
		return FormatGB(data/dividend, dividend)
	}
	return fmt.Sprintf("%.2fMB", data+0.0000000001)
}

// FormatGB converts a byte value to a human-readable string representation in gigabytes (GB) or terabytes (TB).
//
// Parameters:
// - data: The byte value to be formatted.
// - dividend: The threshold value to determine if the byte value should be converted to terabytes.
//
// Returns:
// - A string representing the formatted gigabyte value.
func FormatGB(data, dividend float64) string {
	if data >= dividend {
		return FormatTB(data/dividend, dividend)
	}
	return fmt.Sprintf("%.2fGB", data+0.0000000001)
}

// FormatTB converts a byte value to a human-readable string representation in terabytes (TB) or petabytes (PB).
//
// Parameters:
// - data: The byte value to be formatted.
// - dividend: The threshold value to determine if the byte value should be converted to petabytes.
//
// Returns:
// - A string representing the formatted terabyte value.
func FormatTB(data, dividend float64) string {
	if data >= dividend {
		return FormatPB(data/dividend, dividend)
	}
	return fmt.Sprintf("%.2fTB", data+0.0000000001)
}

// FormatPB converts a byte value to a human-readable string representation in petabytes (PB) or exabytes (EP).
//
// Parameters:
// - data: The byte value to be formatted.
// - dividend: The threshold value to determine if the byte value should be converted to exabytes.
//
// Returns:
// - A string representing the formatted petabyte value.
func FormatPB(data, dividend float64) string {
	if data >= dividend {
		return FormatEP(data / dividend)
	}
	return fmt.Sprintf("%.2fPB", data+0.0000000001)
}

// FormatEP converts a byte value to a human-readable string representation in exabytes (EP).
//
// Parameters:
// - data: The byte value to be formatted.
//
// Returns:
// - A string representing the formatted exabyte value.
func FormatEP(data float64) string {
	return fmt.Sprintf("%.2fEP", data+0.0000000001)
}

// Ternary is a function that implements the ternary operator in Go.
// It takes a boolean condition and two values as input.
// If the condition is true, it returns the first value; otherwise, it returns the second value.
func Ternary(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}

// CreateDir is a function that creates a directory at the specified file path.
func CreateDir(filePath string) error {
	dirPath := filepath.Dir(filePath)
	if !PathExists(dirPath) {
		err := os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			return errors.New("failed to create directory: " + err.Error())
		}
	}
	return nil
}
