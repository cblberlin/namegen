package main


import (
	"flag"
	"log"
	"math/rand"
	"time"

	"github.com/ironarachne/namegen/router"
)

func main() {
	port := flag.String("port", "8080", "API服务的端口号")
	randomSeed := flag.String("s", "none", "可选的随机数生成器种子（字母数字）")
	flag.Parse()

	if *randomSeed == "none" {
		rand.Seed(time.Now().UnixNano())
	} else {
		rand.Seed(time.Now().UnixNano())
	}

	// 从 router 启动服务
	log.Fatal(router.StartServer(*port))
}