package validate

import "regexp"

// CheckEmail 验证邮箱是否合法
func CheckEmail(email string) bool {
	exp := regexp.MustCompile(`\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`)
	return exp.MatchString(email)
}

// CheckPhone 验证手机号是否合法
func CheckPhone(phone string) bool {
	exp := regexp.MustCompile(`^1[3456789]\d{9}$`)
	return exp.MatchString(phone)
}

// CheckContainUrl 验证字符串是否包含URL
func CheckContainUrl(str string) bool {
	exp := regexp.MustCompile(`(https?://)?(www\.)?[a-zA-Z0-9]+\.[a-zA-Z]+(\.[a-zA-Z]+)?`)
	return exp.FindString(str) != ""
}
