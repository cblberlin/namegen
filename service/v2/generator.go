package v2

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/ironarachne/namegen"
	v1service "github.com/ironarachne/namegen/service/v1" // 引用 V1 的国家映射
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
// 工具函数 (完美复刻 Python 工具库)
// -------------------------

func randSep() string {
	seps := []string{"", "", "", ".", "_", "-"}
	return seps[rand.Intn(len(seps))]
}

func randDigits(n int) string {
	var sb strings.Builder
	for i := 0; i < n; i++ {
		fmt.Fprintf(&sb, "%d", rand.Intn(10))
	}
	return sb.String()
}

// 对应 Python 的 rand_short_id (仅小写+数字)
func randShortID(n int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// 对应 Python 的 generate_password (大小写+数字)
func generatePassword(minLen, maxLen int) string {
	n := rand.Intn(maxLen-minLen+1) + minLen
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// 安全截取字符串，模拟 Python 的切片 s[:max_len]
func safeSlice(s string, length int) string {
	if len(s) > length {
		return s[:length]
	}
	return s
}

func fitLength(parts []string, sep string) string {
	s := strings.Join(parts, sep)
	// 过短补充数字
	if len(s) < EmailMinLen {
		s += randDigits(EmailMinLen - len(s))
	} else if len(s) > EmailMaxLen {
		// 过长直接截断
		s = s[:EmailMaxLen]
	}
	return s
}

// -------------------------
// 四大邮箱策略 (修复了截取长度和分隔符)
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
	// Python 中这里使用了 rand_sep()
	return fitLength([]string{wd, randDigits(rand.Intn(3) + 1)}, randSep())
}

func buildNameDigits(firstName string) string {
	name := safeSlice(strings.ToLower(firstName), 8)
	// Python 中这里也使用了 rand_sep()
	return fitLength([]string{name, randDigits(rand.Intn(3) + 1)}, randSep())
}

// -------------------------
// 主导出函数
// -------------------------

func GenerateV2Profile(origin, gender, domain string) (string, string, string, string, string, string) {
	initResources()

	// 1. 获取基础名字
	generator := namegen.NameGeneratorFromType(origin, gender)
	firstName, _ := generator.FirstName(gender)
	lastName, _ := generator.LastName()

	// 2. 选择邮箱生成策略
	var prefix string
	strategy := rand.Float32()
	
	if strategy < 0.30 { // 30% 权重
		prefix = buildPinyinWord()
	} else if strategy < 0.60 { // 30% 权重
		prefix = buildWordDigits()
	} else if strategy < 0.85 { // 25% 权重
		prefix = buildNameDigits(firstName)
	} else { // 15% 权重
		prefix = fitLength([]string{randShortID(rand.Intn(5) + 6)}, "")
	}

	email := fmt.Sprintf("%s@%s", prefix, domain)
	
	// 3. 密码生成 (10-14位，大小写+数字)
	password := generatePassword(10, 14)

	// 4. 生日生成 (18-45岁)
	year := time.Now().Year() - (rand.Intn(28) + 18)
	birthday := fmt.Sprintf("%d-%02d-%02d", year, rand.Intn(12)+1, rand.Intn(28)+1)

	// 5. 获取国家名称
	country := v1service.GetCountryName(origin)

	return email, password, firstName, lastName, country, birthday
}