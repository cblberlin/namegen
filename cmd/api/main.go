package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/ironarachne/namegen/api"
)

func main() {
	// 解析命令行参数
	port := flag.String("port", "8080", "API服务的端口号")
	randomSeed := flag.String("s", "none", "可选的随机数生成器种子（字母数字）")
	flag.Parse()

	// 初始化随机数生成器
	if *randomSeed == "none" {
		rand.Seed(time.Now().UnixNano())
	} else {
		// namegen.RandomSeedFromString(*randomSeed)
		// TODO: 如果需要固定种子可以取消注释上面的代码
		rand.Seed(time.Now().UnixNano())
	}

	// 启动服务器
	fmt.Printf("启动名字生成API服务，监听端口: %s\n", *port)
	fmt.Printf("API文档:\n")
	
	fmt.Printf("  获取名字: http://localhost:%s/api/v1/names?origin=english&gender=male&count=5&mode=full\n", *port)
	fmt.Printf("  查看可用名字起源: http://localhost:%s/api/v1/origins\n", *port)
	fmt.Printf("\n可用参数:\n")
	fmt.Printf("  origin: 名字的起源/国家，如: english, chinese, russian等\n")
	fmt.Printf("  gender: 性别，可选值: male, female, both(默认)\n")
	fmt.Printf("  count: 返回的名字数量，默认为1，最大100\n")
	fmt.Printf("  mode: 名字生成模式，可选值: full(完整名字), firstname(仅名), lastname(仅姓)\n")
	fmt.Printf("  normalize: 是否将特殊字符标准化为基本拉丁字母，可选值: false, true(默认)\n")
	
	// 启动API服务
	log.Fatal(api.StartServer(*port))
} 