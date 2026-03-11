package v1

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	// 引用根目录下的 namegen.go 和 random.go (因为它们在根目录)
	"github.com/ironarachne/namegen"

	// 引用 models 文件夹里的结构体
	"github.com/ironarachne/namegen/models"

	// 引用 service/v1 文件夹里的业务逻辑
	v1service "github.com/ironarachne/namegen/service/v1"
)

// GenerateFullProfile 生成完整的个人信息字符串 (email----password----...)
func GenerateFullProfile(c *gin.Context) {
	origin := c.DefaultQuery("origin", "chinese")
	gender := c.DefaultQuery("gender", "both")
	domain := c.DefaultQuery("domain", "outlook.com")

	generator := namegen.NameGeneratorFromType(origin, gender)
	firstName, err1 := generator.FirstName(gender)
	lastName, err2 := generator.LastName()

	if err1 != nil || err2 != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "无法生成名字"})
		return
	}

	emailPrefix, birthDate := v1service.GenerateEmailPrefix(firstName, lastName, origin)
	password := v1service.GeneratePassword(12)

	country := v1service.COUNTRY_MAPPING[origin]
	if country == "" {
		country = origin
	}

	email := fmt.Sprintf("%s@%s", emailPrefix, domain)
	profileStr := fmt.Sprintf("%s----%s----%s----%s----%s----%s", email, password, lastName, firstName, country, birthDate)

	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.String(http.StatusOK, profileStr)
}

// GenerateName 处理 /names 接口，支持批量生成和拉丁化
func GenerateName(c *gin.Context) {
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

	generator := namegen.NameGeneratorFromType(origin, gender)
	var responses []models.NameResponse

	for i := 0; i < count; i++ {
		response := models.NameResponse{Gender: gender, Origin: origin}
		switch mode {
		case "full":
			fn, err1 := generator.FirstName(gender)
			ln, err2 := generator.LastName()
			if err1 == nil && err2 == nil {
				if normalize {
					fn = v1service.NormalizeToBasicLatin(fn)
					ln = v1service.NormalizeToBasicLatin(ln)
				}
				response.FirstName = fn
				response.LastName = ln
				response.Name = fn + " " + ln
			} else {
				c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "名字生成失败"})
				return
			}
		case "firstname":
			fn, err := generator.FirstName(gender)
			if err == nil && normalize {
				fn = v1service.NormalizeToBasicLatin(fn)
			}
			response.FirstName = fn
			response.Name = fn
		case "lastname":
			ln, err := generator.LastName()
			if err == nil && normalize {
				ln = v1service.NormalizeToBasicLatin(ln)
			}
			response.LastName = ln
			response.Name = ln
		default:
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "不支持的生成模式: " + mode})
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

// ListOrigins 返回支持的国家/种族列表
func ListOrigins(c *gin.Context) {
	origins := namegen.GetSupportedOrigins()
	c.JSON(http.StatusOK, gin.H{"origins": origins})
}

// GenerateProfile 生成基础的 Profile 信息 (First + Last Name)
func GenerateProfile(c *gin.Context) {
	origin := c.DefaultQuery("origin", "chinese")
	gender := c.DefaultQuery("gender", "both")
	countStr := c.DefaultQuery("count", "1")
	count, err := strconv.Atoi(countStr)
	if err != nil || count < 1 {
		count = 1
	}

	normalizeStr := c.DefaultQuery("normalize", "true")
	normalize := normalizeStr != "false" && normalizeStr != "0" && normalizeStr != "no"

	generator := namegen.NameGeneratorFromType(origin, gender)
	var responses []models.ProfileResponse

	for i := 0; i < count; i++ {
		fn, err1 := generator.FirstName(gender)
		ln, err2 := generator.LastName()
		if err1 == nil && err2 == nil {
			if normalize {
				fn = v1service.NormalizeToBasicLatin(fn)
				ln = v1service.NormalizeToBasicLatin(ln)
			}
			responses = append(responses, models.ProfileResponse{
				FirstName: fn,
				LastName:  ln,
				Gender:    gender,
				Origin:    origin,
			})
		}
	}

	if count == 1 {
		c.JSON(http.StatusOK, responses[0])
	} else {
		c.JSON(http.StatusOK, responses)
	}
}

// GenerateEmailPrefixSimple 只生成邮箱前缀和日期
func GenerateEmailPrefixSimple(c *gin.Context) {
	firstname := c.Query("firstname")
	lastname := c.Query("lastname")
	origin := c.DefaultQuery("origin", "chinese")

	if firstname == "" || lastname == "" {
		generator := namegen.NameGeneratorFromType(origin, "both")
		firstname, _ = generator.FirstName("both")
		lastname, _ = generator.LastName()
	}

	emailPrefix, birthDate := v1service.GenerateEmailPrefix(firstname, lastname, origin)
	c.JSON(http.StatusOK, gin.H{"email_prefix": emailPrefix, "birth_date": birthDate})
}

// GenerateProfileSimple 批量生成纯文本格式的 Profile (常用于导出)
func GenerateProfileSimple(c *gin.Context) {
	gender := c.DefaultQuery("gender", "both")
	origin := c.DefaultQuery("origin", "chinese")
	domain := c.DefaultQuery("domain", "outlook.com")

	countStr := c.DefaultQuery("count", "1")
	count, err := strconv.Atoi(countStr)
	if err != nil || count < 1 {
		count = 1
	}

	generator := namegen.NameGeneratorFromType(origin, gender)
	var results []string

	for i := 0; i < count; i++ {
		fn, err1 := generator.FirstName(gender)
		ln, err2 := generator.LastName()

		if err1 == nil && err2 == nil {
			emailPrefix, birthDate := v1service.GenerateEmailPrefix(fn, ln, origin)
			password := v1service.GeneratePassword(12)
			country := v1service.COUNTRY_MAPPING[origin]
			if country == "" {
				country = origin
			}

			line := fmt.Sprintf("%s@%s----%s----%s----%s----%s----%s", emailPrefix, domain, password, ln, fn, country, birthDate)
			results = append(results, line)
		}
	}

	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.String(http.StatusOK, strings.Join(results, "\n"))
}