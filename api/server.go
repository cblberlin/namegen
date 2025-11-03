package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

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

// StartServer 启动API服务器
func StartServer(port string) error {
	http.HandleFunc("/api/v1/names", generateNameHandler)
	http.HandleFunc("/api/v1/origins", listOriginsHandler)

	return http.ListenAndServe(":"+port, nil)
}

// 名字生成处理函数
func generateNameHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// 获取请求参数
	origin := r.URL.Query().Get("origin")
	if origin == "" {
		origin = "english" // 默认使用英语名
	}

	gender := r.URL.Query().Get("gender")
	if gender == "" {
		gender = "both" // 默认返回任意性别
	}

	countStr := r.URL.Query().Get("count")
	count := 1
	if countStr != "" {
		var err error
		count, err = strconv.Atoi(countStr)
		if err != nil || count < 1 {
			count = 1
		}
		if count > 100 {
			count = 100 // 限制最大请求数量
		}
	}

	// 名字生成模式: full, firstname, lastname
	mode := r.URL.Query().Get("mode")
	if mode == "" {
		mode = "full"
	}
	
	// 是否标准化名字，默认为true
	normalizeParam := r.URL.Query().Get("normalize")
	normalize := normalizeParam != "false" && normalizeParam != "0" && normalizeParam != "no"

	// 生成名字
	generator := namegen.NameGeneratorFromType(origin, gender)
	var responses []NameResponse

	// 遍历生成请求的名字数量
	for i := 0; i < count; i++ {
		response := NameResponse{
			Gender: gender,
			Origin: origin,
		}

		var err error

		// 根据不同模式生成名字
		switch mode {
		case "full":
			firstName, err1 := generator.FirstName(gender)
			lastName, err2 := generator.LastName()
			if err1 == nil && err2 == nil {
				// 标准化处理
				if normalize {
					firstName = NormalizeToBasicLatin(firstName)
					lastName = NormalizeToBasicLatin(lastName)
				}
				response.FirstName = firstName
				response.LastName = lastName
				response.Name = firstName + " " + lastName
			} else {
				err = errors.New("无法生成完整名字")
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
			sendErrorResponse(w, "不支持的生成模式: "+mode, http.StatusBadRequest)
			return
		}

		if err != nil {
			if errors.Is(err, namegen.ErrorEmptyItems) {
				sendErrorResponse(w, "不支持的名字起源: "+origin, http.StatusBadRequest)
			} else {
				sendErrorResponse(w, "名字生成失败: "+err.Error(), http.StatusInternalServerError)
			}
			return
		}

		responses = append(responses, response)
	}

	// 根据count决定返回单个对象还是数组
	if count == 1 {
		json.NewEncoder(w).Encode(responses[0])
	} else {
		json.NewEncoder(w).Encode(responses)
	}
}

// 列出所有可用的名字起源
func listOriginsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	origins := []string{
		"anglosaxon", "dutch", "dwarf", "elf", "english", 
		"estonian", "fantasy", "finnish", "french", "german", 
		"greek", "hindu", "indonesian", "irish", "italian", 
		"japanese", "korean", "mayan", "mongolian", "nepalese", 
		"norwegian", "portuguese", "russian", "spanish", "swedish", 
		"thai", "ukrainian", "somalia", "arabic", "hawaiian", 
		"turkish", "serbian", "nigerian", "polish", "chinese",
		"european",
	}

	json.NewEncoder(w).Encode(map[string][]string{
		"origins": origins,
	})
}

// 辅助函数：发送错误响应
func sendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error: message,
	})
} 