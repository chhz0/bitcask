package utils

import (
	"math/rand"
	"strings"
)

// RandomString 支持随机长度的随机字符串生成
func RandomString(minLen, maxLen int) string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length := rand.Intn(maxLen-minLen+1) + minLen // 随机长度区间
	var sb strings.Builder
	for i := 0; i < length; i++ {
		sb.WriteByte(chars[rand.Intn(len(chars))])
	}
	return sb.String()
}
