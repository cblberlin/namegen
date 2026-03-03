package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ironarachne/namegen"
)

// NameResponse 表示API返回的名字结构
type NameResponse struct {
	Name      string `json:"name"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Gender    string `json:"gender"`
	Origin    string `json:"origin"`
}

// 错误响应结构
type ErrorResponse struct {
	Error string `json:"error"`
}

// ProfileResponse 表示API返回的profile结构
type ProfileResponse struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Gender    string `json:"gender"`
	Origin    string `json:"origin"`
}

// FullProfileResponse 表示完整的用户档案响应
type FullProfileResponse struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	LastName   string `json:"lastname"`
	FirstName  string `json:"firstname"`
	Country    string `json:"country"`
	BirthDate  string `json:"birth_date"`
	ProfileStr string `json:"profile_str"` // 原始格式字符串
}



// StartServer 使用Gin框架启动API服务器
func StartServer(port string) error {
	gin.SetMode(gin.ReleaseMode) // 生产模式

	r := gin.New()

	// 添加中间件
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("[%s] %s %s %d %s %s\n",
			param.TimeStamp.Format("2006/01/02 15:04:05"),
			param.Method,
			param.Path,
			param.StatusCode,
			param.Latency,
			param.ClientIP,
		)
	}))
	r.Use(gin.Recovery())

	// API路由组 - 统一使用 /api/v1 前缀
	v1 := r.Group("/api/v1")
	{
		v1.GET("/names", generateNameHandlerGin)
		v1.GET("/origins", listOriginsHandlerGin)
		v1.GET("/profile", generateProfileHandlerGin)
		v1.GET("/full-profile", generateFullProfileHandlerGin)

		// 简单接口也统一使用 /api/v1 前缀
		v1.GET("/generate-name", generateNameSimpleHandlerGin)
		v1.GET("/generate-email-prefix", generateEmailPrefixSimpleHandlerGin)
		v1.GET("/generate-profile", generateProfileSimpleHandlerGin)
	}

	fmt.Printf("🚀 API服务器启动成功，监听端口: %s\n", port)
	fmt.Printf("📖 API文档 (所有接口统一使用 /api/v1 前缀):\n")
	fmt.Printf("  获取名字: http://localhost:%s/api/v1/names?origin=english&gender=male&count=5&mode=full\n", port)
	fmt.Printf("  生成个人资料: http://localhost:%s/api/v1/profile?origin=chinese&count=1\n", port)
	fmt.Printf("  生成完整档案: http://localhost:%s/api/v1/full-profile?origin=french&domain=outlook.com\n", port)
	fmt.Printf("  查看可用名字起源: http://localhost:%s/api/v1/origins\n", port)
	fmt.Printf("  生成名字: http://localhost:%s/api/v1/generate-name?gender=both&origin=chinese\n", port)
	fmt.Printf("  生成邮箱前缀: http://localhost:%s/api/v1/generate-email-prefix\n", port)
	fmt.Printf("  生成完整档案: http://localhost:%s/api/v1/generate-profile?gender=both&origin=chinese&domain=outlook.com\n", port)

	return r.Run(":" + port)
}

// Gin版本的名字生成处理函数
func generateNameHandlerGin(c *gin.Context) {
	origin := c.DefaultQuery("origin", "english")
	gender := c.DefaultQuery("gender", "both")
	mode := c.DefaultQuery("mode", "full")

	countStr := c.DefaultQuery("count", "1")
	count, err := strconv.Atoi(countStr)
	if err != nil || count < 1 {
		count = 1
	}
	if count > 100 {
		count = 100
	}

	normalizeStr := c.DefaultQuery("normalize", "true")
	normalize := normalizeStr != "false" && normalizeStr != "0" && normalizeStr != "no"

	// 生成名字
	generator := namegen.NameGeneratorFromType(origin, gender)
	var responses []NameResponse

	for i := 0; i < count; i++ {
		response := NameResponse{
			Gender: gender,
			Origin: origin,
		}

		switch mode {
		case "full":
			firstName, err1 := generator.FirstName(gender)
			lastName, err2 := generator.LastName()
			if err1 == nil && err2 == nil {
				if normalize {
					firstName = NormalizeToBasicLatin(firstName)
					lastName = NormalizeToBasicLatin(lastName)
				}
				response.FirstName = firstName
				response.LastName = lastName
				response.Name = firstName + " " + lastName
			} else {
				c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "名字生成失败"})
				return
			}
		case "firstname":
			response.FirstName, err = generator.FirstName(gender)
			if err == nil && normalize {
				response.FirstName = NormalizeToBasicLatin(response.FirstName)
			}
			response.Name = response.FirstName
		case "lastname":
			response.LastName, err = generator.LastName()
			if err == nil && normalize {
				response.LastName = NormalizeToBasicLatin(response.LastName)
			}
			response.Name = response.LastName
		default:
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "不支持的生成模式: " + mode})
			return
		}

		if err != nil {
			if err == namegen.ErrorEmptyItems {
				c.JSON(http.StatusBadRequest, ErrorResponse{Error: "不支持的名字起源: " + origin})
			} else {
				c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "名字生成失败: " + err.Error()})
			}
			return
		}

		responses = append(responses, response)
	}

	if count == 1 {
		c.JSON(http.StatusOK, responses[0])
	} else {
		c.JSON(http.StatusOK, responses)
	}
}

// Gin版本的起源列表处理函数
func listOriginsHandlerGin(c *gin.Context) {
	origins := namegen.GetSupportedOrigins()
	c.JSON(http.StatusOK, gin.H{
		"origins": origins,
	})
}

// Gin版本的profile生成处理函数
func generateProfileHandlerGin(c *gin.Context) {
	origin := c.DefaultQuery("origin", "chinese")
	gender := c.DefaultQuery("gender", "both")

	countStr := c.DefaultQuery("count", "1")
	count, err := strconv.Atoi(countStr)
	if err != nil || count < 1 {
		count = 1
	}
	if count > 100 {
		count = 100
	}

	normalizeStr := c.DefaultQuery("normalize", "true")
	normalize := normalizeStr != "false" && normalizeStr != "0" && normalizeStr != "no"

	generator := namegen.NameGeneratorFromType(origin, gender)
	var responses []ProfileResponse

	for i := 0; i < count; i++ {
		response := ProfileResponse{
			Gender: gender,
			Origin: origin,
		}

		firstName, err1 := generator.FirstName(gender)
		lastName, err2 := generator.LastName()
		if err1 == nil && err2 == nil {
			if normalize {
				firstName = NormalizeToBasicLatin(firstName)
				lastName = NormalizeToBasicLatin(lastName)
			}
			response.FirstName = firstName
			response.LastName = lastName
		} else {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "无法生成完整名字"})
			return
		}

		responses = append(responses, response)
	}

	if count == 1 {
		c.JSON(http.StatusOK, responses[0])
	} else {
		c.JSON(http.StatusOK, responses)
	}
}


// Gin版本的完整档案生成处理函数
func generateFullProfileHandlerGin(c *gin.Context) {
	origin := c.DefaultQuery("origin", "chinese")
	gender := c.DefaultQuery("gender", "both")
	domain := c.DefaultQuery("domain", "outlook.com")

	// 生成姓名
	generator := namegen.NameGeneratorFromType(origin, gender)
	firstName, err1 := generator.FirstName(gender)
	lastName, err2 := generator.LastName()

	if err1 != nil || err2 != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "无法生成名字"})
		return
	}

	// 生成邮箱前缀和生日
	emailPrefix, birthDate := generateEmailPrefix(firstName, lastName)

	// 生成密码
	password := generatePassword(12)

	// 获取国家名称
	country := COUNTRY_MAPPING[origin]
	if country == "" {
		country = origin
	}

	// 构建邮箱
	email := fmt.Sprintf("%s@%s", emailPrefix, domain)

	// 按照指定格式构建profile字符串
	profileStr := fmt.Sprintf("%s----%s----%s----%s----%s----%s",
		email, password, lastName, firstName, country, birthDate)

	// 返回纯文本格式（按照用户要求）
	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.String(http.StatusOK, profileStr)
}

// 简单名字生成处理函数（对应Python的/generate-name）
func generateNameSimpleHandlerGin(c *gin.Context) {
	gender := c.DefaultQuery("gender", "both")
	origin := c.DefaultQuery("origin", "chinese")

	// 生成姓名
	generator := namegen.NameGeneratorFromType(origin, gender)
	firstName, err1 := generator.FirstName(gender)
	lastName, err2 := generator.LastName()

	if err1 != nil || err2 != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "无法生成名字"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"firstname": firstName,
		"lastname":  lastName,
	})
}

// 简单邮箱前缀生成处理函数（对应Python的/generate-email-prefix）
func generateEmailPrefixSimpleHandlerGin(c *gin.Context) {
	firstname := c.Query("firstname")
	lastname := c.Query("lastname")

	// 如果没有提供名字，则生成默认的
	if firstname == "" || lastname == "" {
		generator := namegen.NameGeneratorFromType("chinese", "both")
		var err error
		firstname, err = generator.FirstName("both")
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "无法生成名字"})
			return
		}
		lastname, err = generator.LastName()
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "无法生成名字"})
			return
		}
	}

	// 生成邮箱前缀和生日
	emailPrefix, birthDate := generateEmailPrefix(firstname, lastname)

	c.JSON(http.StatusOK, gin.H{
		"email_prefix": emailPrefix,
		"birth_date":   birthDate,
	})
}

// 简单档案生成处理函数（对应Python的/generate-profile）
func generateProfileSimpleHandlerGin(c *gin.Context) {
	gender := c.DefaultQuery("gender", "both")
	origin := c.DefaultQuery("origin", "chinese")
	domain := c.DefaultQuery("domain", "outlook.com")
	
	// 1. 处理 count 参数
	countStr := c.DefaultQuery("count", "1")
	count, err := strconv.Atoi(countStr)
	if err != nil || count < 1 {
		count = 1
	}
	if count > 100 { // 设置上限，防止恶意请求
		count = 100
	}

	generator := namegen.NameGeneratorFromType(origin, gender)
	var results []string

	// 2. 循环生成
	for i := 0; i < count; i++ {
		firstName, err1 := generator.FirstName(gender)
		lastName, err2 := generator.LastName()

		if err1 != nil || err2 != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "无法生成名字"})
			return
		}

		// 生成邮箱前缀和生日
		emailPrefix, birthDate := generateEmailPrefix(firstName, lastName)

		// 生成密码
		password := generatePassword(12)

		// 获取国家名称
		country := COUNTRY_MAPPING[origin]
		if country == "" {
			country = origin
		}

		// 按照指定格式构建单条 profile 字符串
		profileStr := fmt.Sprintf("%s@%s----%s----%s----%s----%s----%s",
			emailPrefix, domain, password, lastName, firstName, country, birthDate)
		
		results = append(results, profileStr)
	}

	// 3. 将结果用换行符拼接并返回纯文本
	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.String(http.StatusOK, strings.Join(results, "\n"))
}