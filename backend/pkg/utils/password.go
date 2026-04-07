package utils

import (
	"crypto/rand"
	"errors"
	"math/big"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

// 错误定义
var (
	ErrOldPasswordIncorrect = errors.New("原密码错误")
	ErrPasswordComplexity   = errors.New("密码复杂度不足")
)

// 密码加密
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// 验证密码
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// 生成随机密码
func GenerateRandomPassword(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+"
	b := make([]byte, length)
	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		b[i] = charset[n.Int64()]
	}
	return string(b), nil
}

// ValidatePasswordComplexity 验证密码复杂度
func ValidatePasswordComplexity(password string) bool {
	// 至少8个字符，包含大小写字母、数字和特殊字符
	if len(password) < 8 {
		return false
	}
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+]`).MatchString(password)
	return hasLower && hasUpper && hasDigit && hasSpecial
}
