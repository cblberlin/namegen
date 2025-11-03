package namegen

import (
	"bufio"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	europeanMaleFirstNames   []string
	europeanFemaleFirstNames []string
	europeanLastNames        []string
)

// init 函数在包加载时执行，读取欧洲名字文件
func init() {
	europeanMaleFirstNames, europeanFemaleFirstNames, europeanLastNames = loadEuropeanNames()
}

// loadEuropeanNames 从文本文件中加载欧洲名字
func loadEuropeanNames() (maleFirstNames, femaleFirstNames, lastNames []string) {
	maleFirstNames = make([]string, 0)
	femaleFirstNames = make([]string, 0)
	lastNameSet := make(map[string]bool)
	lastNames = make([]string, 0)

	// 获取当前文件所在目录
	_, filename, _, _ := runtime.Caller(0)
	baseDir := filepath.Dir(filename)

	// 读取女性名字文件
	femaleFile := filepath.Join(baseDir, "欧洲女性名（名+姓）.txt")
	if err := readNameFile(femaleFile, &femaleFirstNames, &lastNameSet); err != nil {
		// 如果文件不存在，返回空切片
		return
	}

	// 读取男性名字文件
	maleFile := filepath.Join(baseDir, "欧洲男性名（名+姓）.txt")
	if err := readNameFile(maleFile, &maleFirstNames, &lastNameSet); err != nil {
		// 如果文件不存在，返回空切片
		return
	}

	// 将去重后的姓转换为切片
	for lastName := range lastNameSet {
		lastNames = append(lastNames, lastName)
	}

	return
}

// readNameFile 读取名字文件并解析名和姓
func readNameFile(filepath string, firstNames *[]string, lastNameSet *map[string]bool) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// 按空格分割，取第一部分作为名，剩余部分作为姓
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			firstName := parts[0]
			lastName := strings.Join(parts[1:], " ") // 处理可能的复合姓

			*firstNames = append(*firstNames, firstName)
			(*lastNameSet)[lastName] = true
		}
	}

	return scanner.Err()
}
