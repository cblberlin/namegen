package api

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// 国家映射
var COUNTRY_MAPPING = map[string]string{
	"chinese":  "中国",
	"french":   "法国",
	"german":   "德国",
	"spanish":  "西班牙",
	"italian":  "意大利",
	"english":  "英国",
	"japanese": "日本",
	"korean":   "韩国",
	"russian":  "俄罗斯",
	"arabic":   "阿拉伯",
}

// cleanName 清理外国人名中的特殊字符和空格
func cleanName(name string) string {
	// 移除空格、连字符、撇号、点号等特殊字符
	reg := regexp.MustCompile(`[\s\-'\.]`)
	cleaned := reg.ReplaceAllString(name, "")
	return cleaned
}

// generateEmailPrefix 生成复杂的邮箱前缀（复刻Python版本的逻辑）
func generateEmailPrefix(firstName, lastName string) (string, string) {
	// 清理外国人名中的特殊字符和空格
	firstName = strings.ToLower(cleanName(firstName))
	lastName = strings.ToLower(cleanName(lastName))

	// 生成模拟生日（22-35岁之间的随机生日）
	currentYear := time.Now().Year()
	birthYear := rand.Intn(currentYear-22-(currentYear-35)+1) + (currentYear - 35) // 22-35岁
	birthMonth := rand.Intn(12) + 1
	birthDay := rand.Intn(28) + 1 // 使用28避免月份天数问题

	// 保存生日信息用于返回
	birthDate := fmt.Sprintf("%d-%02d-%02d", birthYear, birthMonth, birthDay)

	// 生日相关的日期格式
	birthYY := strconv.Itoa(birthYear)[2:]

	// 当前年份（用于注册年份）
	currentYY := strconv.Itoa(currentYear)[2:]

	// 生成首字母
	initial := fmt.Sprintf("%s%s", string(firstName[0]), string(lastName[0]))

	// 邮箱前缀模式（复刻Python版本的patterns）
	patterns := []string{
		// 带生日年份的模式 (YY)
		fmt.Sprintf("%s.%s%s", firstName, lastName, birthYY),
		fmt.Sprintf("%s%s%s", firstName, lastName, birthYY),
		fmt.Sprintf("%s.%s%s", lastName, firstName, birthYY),
		fmt.Sprintf("%s%s%s", lastName, firstName, birthYY),

		// 带注册年份的模式 (YY)
		fmt.Sprintf("%s.%s%s", firstName, lastName, currentYY),
		fmt.Sprintf("%s%s%s", firstName, lastName, currentYY),
		fmt.Sprintf("%s.%s%s", lastName, firstName, currentYY),
		fmt.Sprintf("%s%s%s", lastName, firstName, currentYY),

		// 使用下划线的模式
		fmt.Sprintf("%s_%s%s", firstName, lastName, birthYY),
		fmt.Sprintf("%s_%s%s", lastName, firstName, birthYY),
		fmt.Sprintf("%s_%s%s", firstName, lastName, currentYY),
		fmt.Sprintf("%s_%s%s", lastName, firstName, currentYY),
		fmt.Sprintf("%s_%s_%s", initial, firstName, lastName),
		fmt.Sprintf("%s_%s_%s", initial, lastName, firstName),

		// 使用连字符的模式
		fmt.Sprintf("%s-%s%s", firstName, lastName, birthYY),
		fmt.Sprintf("%s-%s%s", lastName, firstName, birthYY),
		fmt.Sprintf("%s-%s%s", firstName, lastName, currentYY),
		fmt.Sprintf("%s-%s%s", lastName, firstName, currentYY),

		// 年份在中间的模式
		fmt.Sprintf("%s%s%s", firstName, birthYY, lastName),
		fmt.Sprintf("%s%s%s", lastName, birthYY, firstName),
		fmt.Sprintf("%s%s%s", firstName, currentYY, lastName),
		fmt.Sprintf("%s%s%s", lastName, currentYY, firstName),
	}

	// 随机选择一个模式
	pattern := patterns[rand.Intn(len(patterns))]

	return pattern, birthDate
}

// generatePassword 生成密码（复刻Python版本）
func generatePassword(length int) string {
	if length <= 0 {
		length = 12
	}

	characters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var password strings.Builder
	for i := 0; i < length; i++ {
		password.WriteByte(characters[rand.Intn(len(characters))])
	}
	return password.String()
}