package api

import (
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// 特殊字符到基本拉丁字母的映射表
var specialCharMap = map[rune]string{
	// 北欧字符
	'æ': "ae", 'Æ': "Ae",
	'œ': "oe", 'Œ': "Oe",
	'ø': "o", 'Ø': "O",
	'å': "a", 'Å': "A",
	
	// 德语字符
	'ä': "a", 'Ä': "A",
	'ö': "o", 'Ö': "O",
	'ü': "u", 'Ü': "U",
	'ß': "ss",
	
	// 冰岛/古英语字符
	'þ': "th", 'Þ': "Th",
	'ð': "d", 'Ð': "D",
	
	// 东欧字符
	'ą': "a", 'Ą': "A",
	'ć': "c", 'Ć': "C",
	'ę': "e", 'Ę': "E",
	'ł': "l", 'Ł': "L",
	'ń': "n", 'Ń': "N",
	'ś': "s", 'Ś': "S",
	'ź': "z", 'Ź': "Z",
	'ż': "z", 'Ż': "Z",
	
	// 土耳其字符
	'ı': "i", 'İ': "I",
	'ğ': "g", 'Ğ': "G",
	'ş': "s", 'Ş': "S",
	
	// 捷克/斯洛伐克字符
	'č': "c", 'Č': "C",
	'ř': "r", 'Ř': "R",
	'š': "s", 'Š': "S",
	'ž': "z", 'Ž': "Z",
	
	// 西班牙/葡萄牙字符
	'ñ': "n", 'Ñ': "N",
	
	// 法语字符已通过变音符号处理
}

// NormalizeToBasicLatin 将特殊字符转换为基本拉丁字母(A-Z, a-z)
func NormalizeToBasicLatin(input string) string {
	if input == "" {
		return ""
	}
	
	// 步骤1: 分解带重音符号的字符（如é->e+´）
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, _ := transform.String(t, input)
	
	// 步骤2: 处理仍存在的特殊字符
	var builder strings.Builder
	
	for _, r := range result {
		if r < 128 {
			// ASCII字符直接保留
			builder.WriteRune(r)
		} else if replacement, exists := specialCharMap[r]; exists {
			// 已知特殊字符替换成映射值
			builder.WriteString(replacement)
		} else if unicode.IsLetter(r) {
			// 其他未知字母字符替换为'x'
			builder.WriteRune('x')
		} else if unicode.IsPunct(r) || unicode.IsSpace(r) || unicode.IsDigit(r) {
			// 保留标点符号、空格和数字
			builder.WriteRune(r)
		}
	}
	
	return builder.String()
} 