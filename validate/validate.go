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
