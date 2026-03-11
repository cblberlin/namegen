package v2

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
	"regexp"

	"github.com/ironarachne/namegen"
	v1service "github.com/ironarachne/namegen/service/v1" // 引用 V1
)

var (
	pinyinWords []string
	commonWords []string
	loadOnce    sync.Once
)

const (
	EmailMinLen = 8
	EmailMaxLen = 18
)

// 初始化资源
func initResources() {
	loadOnce.Do(func() {
		pinyinWords = loadFile("email/onewordpinyin.txt")
		commonWords = loadFile("email/words_3_8.txt")
	})
}

func loadFile(path string) []string {
	var results []string
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("⚠️ 无法加载: %s\n", path)
		return []string{"user"}
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		w := strings.ToLower(strings.TrimSpace(scanner.Text()))
		if w != "" {
			results = append(results, w)
		}
	}
	return results
}

// -------------------------
// 工具函数 
// -------------------------

func randSep() string {
	seps := []string{".", "_", "-"}
	return seps[rand.Intn(len(seps))]
}

func randDigits(n int) string {
	var sb strings.Builder
	for i := 0; i < n; i++ {
		fmt.Fprintf(&sb, "%d", rand.Intn(10))
	}
	return sb.String()
}

func randShortID(n int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func generatePassword(minLen, maxLen int) string {
	n := rand.Intn(maxLen-minLen+1) + minLen
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func safeSlice(s string, length int) string {
	// 转换为 rune 切片处理，防止截断多字节字符时发生 panic
	runes := []rune(s)
	if len(runes) > length {
		return string(runes[:length])
	}
	return s
}

func fitLength(parts []string, sep string) string {
	s := strings.Join(parts, sep)
	if len(s) < EmailMinLen {
		s += randDigits(EmailMinLen - len(s))
	} else if len(s) > EmailMaxLen {
		s = safeSlice(s, EmailMaxLen)
	}
	return s
}

// -------------------------
// 四大邮箱策略 
// -------------------------

func buildPinyinWord() string {
	initResources()
	py := safeSlice(pinyinWords[rand.Intn(len(pinyinWords))], 7)
	wd := safeSlice(commonWords[rand.Intn(len(commonWords))], 5)
	
	if rand.Float32() < 0.5 {
		return fitLength([]string{py, wd}, randSep())
	}
	return fitLength([]string{wd, py}, randSep())
}

func buildWordDigits() string {
	initResources()
	wd := safeSlice(commonWords[rand.Intn(len(commonWords))], 8)
	return fitLength([]string{wd, randDigits(rand.Intn(3) + 1)}, randSep())
}


// ✨ 修改点：加入 Normalize 和特殊字符清洗
func buildNameDigits(firstName string) string {
	// 1. 先进行拉丁化 (例如把 "你好" 转成 "Ni Hao", "ü" 转成 "u")
	normName := v1service.NormalizeToBasicLatin(firstName)
	
	// 2. 清除拉丁化后可能出现的空格、连字符等（保证邮箱前缀合法）
	cleanName := strings.ReplaceAll(normName, " ", "")
	cleanName = strings.ReplaceAll(cleanName, "-", "")
	cleanName = strings.ReplaceAll(cleanName, "'", "")
	
	// 3. 转小写并安全截取
	name := safeSlice(strings.ToLower(cleanName), 8)
	
	return fitLength([]string{name, randDigits(rand.Intn(3) + 1)}, randSep())
}

func strictCleanPrefix(s string) string {
	// 只允许小写字母、数字和三个合法分隔符
	reg := regexp.MustCompile(`[^a-z0-9._-]`)
	s = reg.ReplaceAllString(strings.ToLower(s), "")
	
	// 防止出现连续的分隔符，比如 a..b 变成 a.b
	multiSymbolReg := regexp.MustCompile(`[._-]{2,}`)
	s = multiSymbolReg.ReplaceAllString(s, ".")
	
	// 掐头去尾的分隔符
	return strings.Trim(s, "._-")
}

// -------------------------
// 主导出函数
// -------------------------

func GenerateV2Profile(origin, gender, domain string) (string, string, string, string, string, string) {
	initResources()

	// 1. 获取基础名字 (保留原语言用于输出 Profile)
	generator := namegen.NameGeneratorFromType(origin, gender)
	firstName, _ := generator.FirstName(gender)
	lastName, _ := generator.LastName()

	// 2. 选择邮箱生成策略
	var prefix string
	strategy := rand.Float32()
	
	if strategy < 0.30 { 
		prefix = buildPinyinWord()
	} else if strategy < 0.60 { 
		prefix = buildWordDigits()
	} else if strategy < 0.85 { 
		prefix = buildNameDigits(firstName) // 这里传入原始 firstName，内部会 normalize
	} else { 
		prefix = fitLength([]string{randShortID(rand.Intn(5) + 6)}, "")
	}
	prefix = strictCleanPrefix(prefix)

	// ✨ 新增：清洗后可能导致前缀变短（甚至为空），这里做个安全兜底
	if len(prefix) < EmailMinLen {
		prefix += randDigits(EmailMinLen - len(prefix))
	} else if len(prefix) > EmailMaxLen {
		prefix = safeSlice(prefix, EmailMaxLen)
	}

	email := fmt.Sprintf("%s@%s", prefix, domain)
	
	// 3. 密码生成 (10-14位，大小写+数字)
	password := generatePassword(10, 14)

	// 4. 生日生成 (18-45岁)
	year := time.Now().Year() - (rand.Intn(28) + 18)
	birthday := fmt.Sprintf("%d-%02d-%02d", year, rand.Intn(12)+1, rand.Intn(28)+1)

	// 5. 获取国家名称
	country := v1service.GetCountryName(origin)

	// 注意：返回的 firstName 和 lastName 依然是原生字符，方便注册填表
	return email, password, firstName, lastName, country, birthday
}