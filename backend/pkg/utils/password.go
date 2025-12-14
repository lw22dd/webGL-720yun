package utils

import (
	"errors"
	"math/rand"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// 错误定义
var (
	ErrOldPasswordIncorrect = errors.New("原密码错误")
	ErrPasswordComplexity   = errors.New("密码复杂度不足")
)

// HashPassword 密码加密
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPassword 验证密码
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateRandomPassword 生成随机密码
func GenerateRandomPassword(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+"
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	return string(b)
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
