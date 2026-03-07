package api

import (
	"encoding/csv"
	"fmt"
	mrand "math/rand"
	"os"
	"strings"
	"sync"
)

var (
	nicknameDB   map[string][]string
	nicknameOnce sync.Once
)

// loadCSV 初始化读取本地的 CSV 文件
func loadCSV() {
	nicknameDB = make(map[string][]string)

	// 【修改点1】文件名改为你实际的 names.csv
	file, err := os.Open("names.csv")
	if err != nil {
		fmt.Println("⚠️ 未找到 names.csv，将跳过昵称映射")
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1 // 允许不规则列数
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("⚠️ 读取 names.csv 失败:", err)
		return
	}

	count := 0
	for _, row := range records {
		// 你的格式: aaron,has_nickname,erin
		if len(row) >= 3 {
			formal := strings.ToLower(strings.TrimSpace(row[0]))
			formal = strings.TrimPrefix(formal, "\xef\xbb\xbf") // 去除 BOM 头
			nick := strings.ToLower(strings.TrimSpace(row[2]))
			
			if formal != "" && nick != "" {
				nicknameDB[formal] = append(nicknameDB[formal], nick)
				count++
			}
		}
	}
	fmt.Printf("✅ 成功从 names.csv 加载了 %d 条昵称映射规则！\n", count)
}

// GetRandomNickname 获取昵称（支持精准匹配与包含匹配）
func GetRandomNickname(formalName string) string {
	nicknameOnce.Do(loadCSV)

	cleanInput := strings.ToLower(strings.TrimSpace(formalName))
	
	// 【第一道防线】：精准匹配
	// 比如传入 "william"，直接命中
	if nicks, exists := nicknameDB[cleanInput]; exists && len(nicks) > 0 {
		return nicks[mrand.Intn(len(nicks))]
	}

	// 【第二道防线】：包含匹配 (只要名字里包含 CSV 里的某个词)
	// 比如传入 "johnathan", 包含了 "john"
	var possibleNicks []string
	for key, nicks := range nicknameDB {
		// 限制 key 长度 >= 3，防止 "ed" 这种太短的词误伤了 "kennedy"
		if len(key) >= 3 && strings.Contains(cleanInput, key) {
			possibleNicks = append(possibleNicks, nicks...)
		}
	}

	// 如果包含匹配找到了任何候选昵称，把它们聚在一起随机选一个！
	if len(possibleNicks) > 0 {
		return possibleNicks[mrand.Intn(len(possibleNicks))]
	}

	// 如果精准和包含都没命中，直接原样返回清洗后的名字
	return cleanInput
}