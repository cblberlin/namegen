package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ironarachne/namegen"
)

type NameResponse struct {
	Name      string `json:"name"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Gender    string `json:"gender"`
	Origin    string `json:"origin"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type ProfileResponse struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Gender    string `json:"gender"`
	Origin    string `json:"origin"`
}

type FullProfileResponse struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	LastName   string `json:"lastname"`
	FirstName  string `json:"firstname"`
	Country    string `json:"country"`
	BirthDate  string `json:"birth_date"`
	ProfileStr string `json:"profile_str"`
}

func StartServer(port string) error {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("[%s] %s %s %d %s %s\n",
			param.TimeStamp.Format("2006/01/02 15:04:05"),
			param.Method, param.Path, param.StatusCode,
			param.Latency, param.ClientIP,
		)
	}))
	r.Use(gin.Recovery())

	v1 := r.Group("/api/v1")
	{
		v1.GET("/names", generateNameHandlerGin)
		v1.GET("/origins", listOriginsHandlerGin)
		v1.GET("/profile", generateProfileHandlerGin)
		v1.GET("/full-profile", generateFullProfileHandlerGin)
		v1.GET("/generate-name", generateNameSimpleHandlerGin)
		v1.GET("/generate-email-prefix", generateEmailPrefixSimpleHandlerGin)
		v1.GET("/generate-profile", generateProfileSimpleHandlerGin)
	}

	fmt.Printf("🚀 API服务器启动成功，监听端口: %s\n", port)
	return r.Run(":" + port)
}

func generateNameHandlerGin(c *gin.Context) {
	origin := c.DefaultQuery("origin", "english")
	gender := c.DefaultQuery("gender", "both")
	mode := c.DefaultQuery("mode", "full")

	countStr := c.DefaultQuery("count", "1")
	count, err := strconv.Atoi(countStr)
	if err != nil || count < 1 { count = 1 }
	if count > 100 { count = 100 }

	normalizeStr := c.DefaultQuery("normalize", "true")
	normalize := normalizeStr != "false" && normalizeStr != "0" && normalizeStr != "no"

	generator := namegen.NameGeneratorFromType(origin, gender)
	var responses []NameResponse

	for i := 0; i < count; i++ {
		response := NameResponse{Gender: gender, Origin: origin}
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
			if err == nil && normalize { response.FirstName = NormalizeToBasicLatin(response.FirstName) }
			response.Name = response.FirstName
		case "lastname":
			response.LastName, err = generator.LastName()
			if err == nil && normalize { response.LastName = NormalizeToBasicLatin(response.LastName) }
			response.Name = response.LastName
		default:
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "不支持的生成模式: " + mode})
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "生成失败"})
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

func listOriginsHandlerGin(c *gin.Context) {
	origins := namegen.GetSupportedOrigins()
	c.JSON(http.StatusOK, gin.H{"origins": origins})
}

func generateProfileHandlerGin(c *gin.Context) {
	origin := c.DefaultQuery("origin", "chinese")
	gender := c.DefaultQuery("gender", "both")
	countStr := c.DefaultQuery("count", "1")
	count, err := strconv.Atoi(countStr)
	if err != nil || count < 1 { count = 1 }
	if count > 100 { count = 100 }

	normalizeStr := c.DefaultQuery("normalize", "true")
	normalize := normalizeStr != "false" && normalizeStr != "0" && normalizeStr != "no"

	generator := namegen.NameGeneratorFromType(origin, gender)
	var responses []ProfileResponse

	for i := 0; i < count; i++ {
		response := ProfileResponse{Gender: gender, Origin: origin}
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

	if count == 1 { c.JSON(http.StatusOK, responses[0]) } else { c.JSON(http.StatusOK, responses) }
}

func generateFullProfileHandlerGin(c *gin.Context) {
	origin := c.DefaultQuery("origin", "chinese")
	gender := c.DefaultQuery("gender", "both")
	domain := c.DefaultQuery("domain", "outlook.com")

	generator := namegen.NameGeneratorFromType(origin, gender)
	firstName, err1 := generator.FirstName(gender)
	lastName, err2 := generator.LastName()

	if err1 != nil || err2 != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "无法生成名字"})
		return
	}

	// 传入 origin 参数
	emailPrefix, birthDate := generateEmailPrefix(firstName, lastName, origin)
	password := generatePassword(12)

	country := COUNTRY_MAPPING[origin]
	if country == "" { country = origin }

	email := fmt.Sprintf("%s@%s", emailPrefix, domain)
	profileStr := fmt.Sprintf("%s----%s----%s----%s----%s----%s", email, password, lastName, firstName, country, birthDate)

	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.String(http.StatusOK, profileStr)
}

func generateNameSimpleHandlerGin(c *gin.Context) {
	gender := c.DefaultQuery("gender", "both")
	origin := c.DefaultQuery("origin", "chinese")

	generator := namegen.NameGeneratorFromType(origin, gender)
	firstName, err1 := generator.FirstName(gender)
	lastName, err2 := generator.LastName()

	if err1 != nil || err2 != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "无法生成名字"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"firstname": firstName, "lastname": lastName})
}

func generateEmailPrefixSimpleHandlerGin(c *gin.Context) {
	firstname := c.Query("firstname")
	lastname := c.Query("lastname")
	origin := c.DefaultQuery("origin", "chinese") // 获取 origin

	if firstname == "" || lastname == "" {
		generator := namegen.NameGeneratorFromType(origin, "both")
		var err error
		firstname, err = generator.FirstName("both")
		lastname, err = generator.LastName()
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "无法生成名字"})
			return
		}
	}

	// 传入 origin
	emailPrefix, birthDate := generateEmailPrefix(firstname, lastname, origin)
	c.JSON(http.StatusOK, gin.H{"email_prefix": emailPrefix, "birth_date": birthDate})
}

func generateProfileSimpleHandlerGin(c *gin.Context) {
	gender := c.DefaultQuery("gender", "both")
	origin := c.DefaultQuery("origin", "chinese")
	domain := c.DefaultQuery("domain", "outlook.com")
	
	countStr := c.DefaultQuery("count", "1")
	count, err := strconv.Atoi(countStr)
	if err != nil || count < 1 { count = 1 }
	const maxCount = 100000
	if count > maxCount { count = maxCount }

	generator := namegen.NameGeneratorFromType(origin, gender)
	var results []string

	for i := 0; i < count; i++ {
		firstName, err1 := generator.FirstName(gender)
		lastName, err2 := generator.LastName()

		if err1 != nil || err2 != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "无法生成名字"})
			return
		}

		// 传入 origin
		emailPrefix, birthDate := generateEmailPrefix(firstName, lastName, origin)
		password := generatePassword(12)

		country := COUNTRY_MAPPING[origin]
		if country == "" { country = origin }

		profileStr := fmt.Sprintf("%s@%s----%s----%s----%s----%s----%s", emailPrefix, domain, password, lastName, firstName, country, birthDate)
		results = append(results, profileStr)
	}

	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.String(http.StatusOK, strings.Join(results, "\n"))
}