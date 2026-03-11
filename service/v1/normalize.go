package v1

import (
	anyascii "github.com/anyascii/go"
)

// NormalizeToBasicLatin 使用 AnyASCII 将任意文字转换为基本拉丁字母(ASCII)
// 它支持 100 多种语言，包括中文、俄文、阿拉伯文以及各种带重音的西欧字符
func NormalizeToBasicLatin(input string) string {
	if input == "" {
		return ""
	}
	
	// anyascii.Transliterate 会自动处理 NFD 规范化、重音去除以及非拉丁字符的转写
	// 例如："你好" -> "Ni Hao", "Æneid" -> "Aeneid", "ü" -> "u"
	return anyascii.Transliterate(input)
}