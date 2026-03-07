package api

import (
	crand "crypto/rand" // 设置别名：用于生成安全密码
	"fmt"
	"math/big"
	mrand "math/rand" // 设置别名：用于普通的随机选择
	"regexp"
	"strconv"
	"strings"
	"time"
)

// 国家映射
var COUNTRY_MAPPING = map[string]string{
	"chinese":  "中国", "french":   "法国", "german":   "德国",
	"spanish":  "西班牙", "italian":  "意大利", "english":  "英国",
	"japanese": "日本", "korean":   "韩国", "russian":  "俄罗斯",
	"arabic":   "阿拉伯", "dutch":    "荷兰", "fantasy":  "幻想",
}

// 随性词库矩阵
var particles = map[string][]string{
	"chinese":  {"da", "xiao", "lao", "ah", "wo", "super", "pro"},
	"english":  {"the", "iam", "real", "thisis", "just", "mr", "top"},
	"french":   {"le", "la", "petit", "grand", "vrai", "ici"},
	"japanese": {"san", "chan", "kun", "neo", "mega"},
	"russian":  {"pro", "best", "real", "top", "cyber"},
}
var universalVibes = []string{"cool", "chill", "vibes", "hyper", "cyber", "neon", "pure", "zero", "prime", "pixel"}
var lifeFields = []string{"lab", "studio", "dev", "design", "art", "code", "life", "hub", "zone", "flow"}
var connectors = []string{".", "_", "-", ""}

func cleanName(name string) string {
	reg := regexp.MustCompile(`[\s\-'\.]`)
	return reg.ReplaceAllString(name, "")
}

func cleanEmailPrefix(s string) string {
	s = strings.ToLower(s)
	reg := regexp.MustCompile(`[^a-z0-9._-]`)
	s = reg.ReplaceAllString(s, "")
	multiSymbolReg := regexp.MustCompile(`[._-]{2,}`)
	s = multiSymbolReg.ReplaceAllString(s, ".")
	s = strings.Trim(s, "._-")
	if len(s) < 3 {
		s += "user"
	}
	return s
}

// generateEmailPrefix 生成复杂的邮箱前缀
func generateEmailPrefix(firstName, lastName, origin string) (string, string) {
	// 1. 先进行拉丁化与基础清理
	formalFn := strings.ToLower(cleanName(NormalizeToBasicLatin(firstName)))
	ln := strings.ToLower(cleanName(NormalizeToBasicLatin(lastName)))

	// 【重头戏】2. 尝试获取昵称
	// 设定 70% 的概率使用昵称（显得更随性），30% 的概率保留全拼
	fn := formalFn
	if mrand.Float32() < 0.90 {
		fn = GetRandomNickname(formalFn)
	}

	// 3. 安全托底
	if len(fn) == 0 { fn = "u" }
	if len(ln) == 0 { ln = "user" }

	// 4. 生成生日信息
	currentYear := time.Now().Year()
	birthYear := mrand.Intn(currentYear-22-(currentYear-35)+1) + (currentYear - 35)
	birthMonth := mrand.Intn(12) + 1
	birthDay := mrand.Intn(28) + 1
	birthDate := fmt.Sprintf("%d-%02d-%02d", birthYear, birthMonth, birthDay)

	birthYY := strconv.Itoa(birthYear)[2:]
	currentYY := strconv.Itoa(currentYear)[2:]
	
	// 注意：因为 fn 可能是缩写（比如 William 变成了 Bill），这里的首字母也会很真实地变成 b
	initial := fmt.Sprintf("%c%c", fn[0], ln[0])

	// 5. 获取随性元素
	pList := particles[origin]
	if len(pList) == 0 {
		pList = []string{"me", "my", "im"}
	}
	particle := pList[mrand.Intn(len(pList))]
	vibe := universalVibes[mrand.Intn(len(universalVibes))]
	field := lifeFields[mrand.Intn(len(lifeFields))]
	conn := connectors[mrand.Intn(len(connectors))]

	// 6. 合并所有模式
	classicPatterns := []string{
		fmt.Sprintf("%s.%s%s", fn, ln, birthYY),
		fmt.Sprintf("%s%s%s", fn, ln, birthYY),
		fmt.Sprintf("%s.%s%s", ln, fn, birthYY),
		fmt.Sprintf("%s%s%s", ln, fn, birthYY),
		fmt.Sprintf("%s.%s%s", fn, ln, currentYY),
		fmt.Sprintf("%s%s%s", fn, ln, currentYY),
		fmt.Sprintf("%s_%s%s", fn, ln, birthYY),
		fmt.Sprintf("%s_%s%s", ln, fn, birthYY),
		fmt.Sprintf("%s_%s_%s", initial, fn, ln),
		fmt.Sprintf("%s_%s_%s", initial, ln, fn),
		fmt.Sprintf("%s-%s%s", fn, ln, birthYY),
		fmt.Sprintf("%s-%s%s", ln, fn, birthYY),
		fmt.Sprintf("%s%s%s", fn, birthYY, ln),
		fmt.Sprintf("%s%s%s", ln, birthYY, fn),
	}

	vibePatterns := []string{
		fmt.Sprintf("%s%s%s", particle, conn, fn),      
		fmt.Sprintf("%s%s%s", fn, conn, vibe),          
		fmt.Sprintf("%s%s%s", ln, conn, field),         
		fmt.Sprintf("%s%s%s%s", fn, conn, ln, birthYY), 
		fmt.Sprintf("%s%s%d", vibe, conn, mrand.Intn(999)), 
	}

	allPatterns := append(classicPatterns, vibePatterns...)
	rawPrefix := allPatterns[mrand.Intn(len(allPatterns))]

	// 7. 返回前执行最终清洗
	return cleanEmailPrefix(rawPrefix), birthDate
}

// generatePassword 生成安全的随机密码
func generatePassword(length int) string {
	if length <= 0 {
		length = 12
	}
	characters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var password strings.Builder
	for i := 0; i < length; i++ {
		idx, _ := crand.Int(crand.Reader, big.NewInt(int64(len(characters))))
		password.WriteByte(characters[idx.Int64()])
	}
	return password.String()
}