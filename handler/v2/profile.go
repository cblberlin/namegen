package v2

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	v2service "github.com/ironarachne/namegen/service/v2"
)

func GenerateFullProfileV2(c *gin.Context) {
	// V2 参数：origin 默认 chinese
	origin := c.DefaultQuery("origin", "chinese")
	gender := c.DefaultQuery("gender", "both")
	
	// 处理 domain 逻辑：如果用户没传，我们就还原 Python 脚本里的从 outlook/hotmail 随机选
	domain := c.Query("domain")
	if domain == "" {
		domains := []string{"outlook.com", "hotmail.com"}
		// 设置一下随机种子，确保每次请求能真正随机
		rand.Seed(time.Now().UnixNano()) 
		domain = domains[rand.Intn(len(domains))]
	}
	
	countStr := c.DefaultQuery("count", "1")
	count, _ := strconv.Atoi(countStr)
	if count < 1 { 
		count = 1 
	}

	var results []string
	for i := 0; i < count; i++ {
		// 👉 修改点 1：接收 6 个返回值，增加 country
		email, pwd, fn, ln, country, birth := v2service.GenerateV2Profile(origin, gender, domain)
		
		// 👉 修改点 2：把硬编码的 "China" 替换成动态的 %s 和 country 变量
		line := fmt.Sprintf("%s----%s----%s----%s----%s----%s", email, pwd, ln, fn, country, birth)
		results = append(results, line)
	}

	// 文本形式返回
	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.String(http.StatusOK, strings.Join(results, "\n"))
}