package router

import (
	"fmt"

	"github.com/gin-gonic/gin"

	// 引入你刚刚补全的 v1 handler
	v1handler "github.com/ironarachne/namegen/handler/v1"
	// 这里预留 v2 的引入，等你写好 v2 handler 后取消注释
	v2handler "github.com/ironarachne/namegen/handler/v2"
)

// StartServer 启动 Gin 服务器
func StartServer(port string) error {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// 自定义日志格式
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("[%s] %s %s %d %s %s\n",
			param.TimeStamp.Format("2006/01/02 15:04:05"),
			param.Method, param.Path, param.StatusCode,
			param.Latency, param.ClientIP,
		)
	}))
	r.Use(gin.Recovery())

	// API 根组
	api := r.Group("/api")

	// === V1 路由组 (你的原始逻辑) ===
	v1 := api.Group("/v1")
	{
		v1.GET("/names", v1handler.GenerateName)
		v1.GET("/origins", v1handler.ListOrigins)
		v1.GET("/profile", v1handler.GenerateProfile)
		v1.GET("/full-profile", v1handler.GenerateFullProfile)
		
		// 对应你原本的 simple 系列接口
		v1.GET("/generate-email-prefix", v1handler.GenerateEmailPrefixSimple)
		v1.GET("/generate-profile", v1handler.GenerateProfileSimple)
	}

	// === V2 路由组 (预留给新逻辑) ===
	v2 := api.Group("/v2")
	{
		// 现在 V2 有了真正的处理函数
		v2.GET("/generate-profile", v2handler.GenerateFullProfileV2)
	}

	fmt.Printf("🚀 API服务器启动成功，监听端口: %s\n", port)
	return r.Run(":" + port)
}