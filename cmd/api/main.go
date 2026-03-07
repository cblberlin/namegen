package main

import (
	"flag"
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
		rand.Seed(time.Now().UnixNano())
	}

	// 启动服务器
	log.Fatal(api.StartServer(*port))
} 