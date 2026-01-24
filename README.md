# 名字生成器 (Name Generator)

这是一个用于生成虚构人物名字的工具。它首先是一个可以被其他项目使用的库，同时还提供了 `namegen` 命令行工具用于生成名字。

## 功能特性

- 支持多种文化背景的名字生成
- 提供命令行工具和HTTP API
- 包含30多种语言和文化的名字数据
- 复杂邮箱前缀和密码生成功能

## 快速开始

### 构建

使用以下脚本构建二进制文件：

```bash
# 构建所有组件（推荐）
make build

# 或者分别构建
make build-cli    # 构建命令行工具 -> build/namegen
make build-api    # 构建API服务 -> build/namegen-api
```

### 使用命令行工具

```bash
# 生成一个英文名字
make run-cli

# 生成5个中文女性名字
make run-cli ARGS="-o chinese -g female -n 5"

# 或者直接运行
./build/namegen
./build/namegen -o chinese -g female -n 5

# 查看所有支持的名字类型
./build/namegen -l
```

### 启动API服务

```bash
# 使用Gin框架启动（推荐，有详细日志）
make run-api-bg

# 前台启动API服务
make run-api

# 或者直接运行
./build/namegen-api
```

**特性:**
- ✅ 详细的请求日志（时间、方法、路径、状态码、延迟、客户端IP）
- ✅ 更好的错误处理和恢复机制
- ✅ 更丰富的中间件支持
- ✅ 更高的性能和稳定性

## API文档

### 生成名字

**请求:**
```
GET /api/v1/names
```

**参数:**

| 参数 | 类型 | 默认值 | 描述 |
|------|------|--------|------|
| origin | string | english | 名字的起源/国家，如: english, chinese, russian等 |
| gender | string | both | 性别，可选值: male, female, both |
| count | int | 1 | 返回的名字数量，最大100 |
| mode | string | full | 名字生成模式: full(完整名字), firstname(仅名), lastname(仅姓) |
| normalize | boolean | true | 是否将特殊字符标准化为基本拉丁字母 |

**示例请求:**
```
GET /api/v1/names?origin=chinese&gender=female&count=2
```

**响应:**
```json
[
  {
    "name": "Hua Chen",
    "first_name": "Hua",
    "last_name": "Chen",
    "gender": "female",
    "origin": "chinese"
  },
  {
    "name": "Mei Li",
    "first_name": "Mei",
    "last_name": "Li",
    "gender": "female",
    "origin": "chinese"
  }
]
```

### 生成个人资料

**请求:**
```
GET /api/v1/profile
```

**参数:**

| 参数 | 类型 | 默认值 | 描述 |
|------|------|--------|------|
| origin | string | chinese | 名字的起源/国家，如: english, chinese, russian等 |
| gender | string | both | 性别，可选值: male, female, both |
| count | int | 1 | 返回的资料数量，最大100 |
| normalize | boolean | true | 是否将特殊字符标准化为基本拉丁字母 |

**示例请求:**
```
GET /api/v1/profile?origin=french&gender=both&count=1
```

**响应:**
```json
{
  "first_name": "Marie",
  "last_name": "Dubois",
  "gender": "both",
  "origin": "french"
}
```

**批量请求示例:**
```
GET /api/v1/profile?origin=chinese&count=3
```

**批量响应:**
```json
[
  {
    "first_name": "小红",
    "last_name": "李",
    "gender": "both",
    "origin": "chinese"
  },
  {
    "first_name": "伟明",
    "last_name": "王",
    "gender": "both",
    "origin": "chinese"
  },
  {
    "first_name": "美丽",
    "last_name": "张",
    "gender": "both",
    "origin": "chinese"
  }
]
```

### 生成完整用户档案

**请求:**
```
GET /api/v1/full-profile
```

**参数:**

| 参数 | 类型 | 默认值 | 描述 |
|------|------|--------|------|
| origin | string | chinese | 名字的起源/国家 |
| gender | string | both | 性别，可选值: male, female, both |
| domain | string | outlook.com | 邮箱域名 |

**示例请求:**
```
GET /api/v1/full-profile?origin=french&domain=gmail.com
```

**响应格式:**
```
email@domain.com----password----lastname----firstname----country----birthdate
```

**示例响应:**
```
marie.dubois1985@gmail.com----P@ssw0rd!----Dubois----Marie----法国----1985-03-15
```

### API使用示例

```bash
# 生成单个个人资料（适合你的Python代码）
curl "http://localhost:8080/api/v1/profile?origin=french&count=1"

# 生成多个中文个人资料
curl "http://localhost:8080/api/v1/profile?origin=chinese&count=3"

# 生成英文男性资料
curl "http://localhost:8080/api/v1/profile?origin=english&gender=male"

# 生成完整用户档案（返回纯文本格式）
curl "http://localhost:8080/api/v1/full-profile?origin=french&domain=gmail.com"

# 简单接口示例
curl "http://localhost:8080/api/v1/generate-name?origin=french"
curl "http://localhost:8080/api/v1/generate-email-prefix"
curl "http://localhost:8080/api/v1/generate-profile?origin=chinese&domain=outlook.com"
```

## 简单API接口（统一使用 /api/v1 前缀，完全匹配你的Python代码）

### 生成名字

**请求:**
```
GET /generate-name
```

**参数:**
- `gender`: 性别（male/female/both，默认: both）
- `origin`: 名字来源（默认: chinese）

**响应:**
```json
{
  "firstname": "小明",
  "lastname": "王"
}
```

**示例:**
```bash
curl "http://localhost:8080/api/v1/generate-name?gender=both&origin=chinese"
```

### 生成邮箱前缀

**请求:**
```
GET /api/v1/generate-email-prefix
```

**参数:**
- `firstname`: 可选，指定名字（如果不提供会自动生成）
- `lastname`: 可选，指定姓氏（如果不提供会自动生成）

**响应:**
```json
{
  "email_prefix": "marie.dubois85",
  "birth_date": "1999-03-15"
}
```

**支持的邮箱前缀模式:**
- `firstname.lastname` + 年份
- `lastname.firstname` + 年份
- `firstname_lastname` + 年份
- `firstname-lastname` + 年份
- 年份在姓名中间的组合
- 多种年份格式（2位/4位）

**示例:**
```bash
curl "http://localhost:8080/api/v1/generate-email-prefix"
```

### 生成完整档案

**请求:**
```
GET /generate-profile
```

**参数:**
- `gender`: 性别（male/female/both，默认: both）
- `origin`: 名字来源（默认: chinese）
- `domain`: 邮箱域名（默认: outlook.com）

**响应格式:**
```
email@domain----password----lastname----firstname----country----birthdate
```

**示例:**
```bash
curl "http://localhost:8080/api/v1/generate-profile?gender=both&origin=french&domain=gmail.com"
```

**响应示例:**
```
marie_dubois85@gmail.com----Abc123XyZ----Dubois----Marie----法国----1999-03-15
```

### Python客户端示例

```python
import asyncio
import httpx

BASE_URL = "http://localhost:8080"

async def get_names(count=1, gender="both", origin="chinese"):
    """调用generate-name API endpoint"""
    url = BASE_URL + "/api/v1/generate-name"

    params = {
        "gender": gender,
        "origin": origin
    }

    async with httpx.AsyncClient() as client:
        response = await client.get(url, params=params)
        data = response.json()
        return data["firstname"], data["lastname"]

async def generate_profile(gender="both", origin="chinese", domain="outlook.com"):
    """调用generate-profile API endpoint"""
    url = BASE_URL + "/api/v1/generate-profile"

    params = {
        "gender": gender,
        "origin": origin,
        "domain": domain
    }

    async with httpx.AsyncClient() as client:
        response = await client.get(url, params=params)
        return response.text

async def test_get_names(gender="both", origin="chinese"):
    firstname, lastname = await get_names(gender=gender, origin=origin)
    print(firstname, lastname)

if __name__ == "__main__":
    async def main():
        await test_get_names(origin="french")
        profile = await generate_profile(origin="chinese", domain="gmail.com")
        print("Profile:", profile)
    asyncio.run(main())
```

### 获取可用名字起源

**请求:**
```
GET /api/v1/origins
```

**响应:**
```json
{
  "origins": [
    "anglosaxon", "dutch", "dwarf", "elf", "english",
    "estonian", "fantasy", "finnish", "french", "german",
    "greek", "hindu", "indonesian", "irish", "italian",
    "japanese", "korean", "mayan", "mongolian", "nepalese",
    "norwegian", "portuguese", "russian", "spanish", "swedish",
    "thai", "ukrainian", "somalia", "arabic", "hawaiian",
    "turkish", "serbian", "nigerian", "polish", "chinese",
    "european"
  ]
}
```

## 开发指南

### 添加新的名字生成器

添加新的名字生成器需要4个步骤：

**注意:** 关于名字格式的优秀指南请参考: [https://www.fbiic.gov/public/2008/nov/Naming_practice_guide_UK_2006.pdf](https://www.fbiic.gov/public/2008/nov/Naming_practice_guide_UK_2006.pdf)

#### 第一步: 添加名字列表全局变量

创建新的.go文件用于存放名字列表。在文件中创建三个全局变量 - 男性名、女性名和姓氏。命名要符合规范（参考现有文件）。

#### 第二步: 将新生成器添加到映射

在namegen.go的NameGeneratorFromType函数中，将新的名字生成器添加到映射中。它应该引用第一步创建的全局变量。

#### 第三步: 添加名字类型到CLI程序

在cmd/namegen/main.go中，将名字类型添加到`nameLists`数组中。

#### 第四步: 构建和测试

运行`./build.sh`并测试新的名字生成器。

## 支持的名字起源

目前支持以下文化背景的名字生成：

**欧洲名字:**
- anglosaxon - 盎格鲁撒克逊名字
- dutch - 荷兰名字
- english - 英文名字
- estonian - 爱沙尼亚名字
- finnish - 芬兰名字
- french - 法国名字
- german - 德国名字
- greek - 希腊名字
- icelandic - 冰岛名字
- irish - 爱尔兰名字
- italian - 意大利名字
- norwegian - 挪威名字
- polish - 波兰名字
- portuguese - 葡萄牙名字
- russian - 俄罗斯名字
- spanish - 西班牙名字
- swedish - 瑞典名字

**亚洲名字:**
- arabic - 阿拉伯名字
- chinese - 中文名字
- hindu - 印度名字
- indonesian - 印尼名字
- japanese - 日本名字
- korean - 韩国名字
- mongolian - 蒙古名字
- nepalese - 尼泊尔名字
- thai - 泰国名字
- turkish - 土耳其名字

**其他文化:**
- hawaiian - 夏威夷名字
- maori - 毛利名字
- mayan - 马雅名字
- nigerian - 尼日利亚名字
- serbian - 塞尔维亚名字
- somalia - 索马里名字
- ukrainian - 乌克兰名字

**虚构世界名字:**
- dwarf - 矮人名字
- elf - 精灵名字
- fantasy - 奇幻世界名字
- european - 欧洲中性名字
