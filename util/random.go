package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "qwertyuiopasdfghjklzxcvbnm"

func init() {
	// 设置随机种子
	rand.Seed(time.Now().UnixNano())
}

// RandomInt 生成随机整数，在min和max之间
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString 生成随机字符串
func RandomString(n int) string {
	// 声明一个字符串构造器
	var sb strings.Builder
	k := len(alphabet)
	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

// RandomOwner 生成一个用户名字
func RandomOwner() string {
	return RandomString(6)
}

// RandomMoney 生成一个随机整数 money
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// RandomCurrency 随机币种
func RandomCurrency() string {
	currencies := []string{
		"RMB", "USD", "EUD",
	}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}
