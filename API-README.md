# Name Generator API

This is an HTTP API service for generating names from different countries and cultural backgrounds, built on the ironarachne/namegen library.

## Quick Start

### Build and Run

```bash
# Clone the repository
git clone https://github.com/ironarachne/namegen.git
cd namegen

# Build the API service
go build -o namegen-api ./cmd/api

# Run the API service (default port 8080)
./namegen-api

# Specify custom port
./namegen-api -port 3000

# Enable API key authentication (restrict API access)
./namegen-api -key your_secret_key
```

## API Authentication

If the server has API key authentication enabled, you can provide the API key in one of three ways:

1. **Bearer Token Authentication (Recommended)**:
   ```
   Authorization: Bearer your_secret_key
   ```

2. **Custom Header**:
   ```
   X-API-Key: your_secret_key
   ```

3. **URL Parameter**:
   ```
   ?api_key=your_secret_key
   ```

The Bearer Token authentication method is recommended as it's the standard practice for REST APIs.

## API Documentation

### Generate Names

**Request:**
```
GET /api/v1/names
```

**Headers (if API key is enabled):**
```
Authorization: Bearer your_secret_key
```

**Parameters:**

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| origin | string | english | Name origin/country, e.g., english, chinese, russian |
| gender | string | both | Gender, options: male, female, both |
| count | int | 1 | Number of names to return, maximum 100 |
| mode | string | full | Name generation mode: full(full name), firstname(first name only), lastname(last name only) |
| normalize | boolean | false | Whether to normalize special characters to basic Latin letters |

**Example Request:**
```
GET /api/v1/names?origin=chinese&gender=female&count=2
Authorization: Bearer your_secret_key
```

**Response:**
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

If requesting only one name (default or count=1), a single object is returned instead of an array.

### Get Available Name Origins

**Request:**
```
GET /api/v1/origins
Authorization: Bearer your_secret_key
```

**Response:**
```json
{
  "origins": [
    "anglosaxon", "dutch", "dwarf", "elf", "english", 
    "estonian", "fantasy", "finnish", "french", "german", 
    "greek", "hindu", "indonesian", "irish", "italian", 
    "japanese", "korean", "mayan", "mongolian", "nepalese", 
    "norwegian", "portuguese", "russian", "spanish", "swedish", 
    "thai", "ukrainian", "somalia", "arabic", "hawaiian", 
    "turkish", "serbian", "nigerian", "polish", "chinese"
  ]
}
```

## Usage Examples

### Using curl to Get Names

```bash
# Using Bearer Token (recommended)
curl -H "Authorization: Bearer your_secret_key" "http://localhost:8080/api/v1/names?origin=chinese&gender=male&count=5"

# Using X-API-Key header
curl -H "X-API-Key: your_secret_key" "http://localhost:8080/api/v1/names?origin=russian&mode=firstname"

# Passing API key via URL parameter
curl "http://localhost:8080/api/v1/names?origin=japanese&count=3&api_key=your_secret_key"
```

### Programmatic API Calls

#### Python Example
```python
import requests

# API key
api_key = "your_secret_key"

# Method 1: Using Bearer Token (recommended)
headers = {"Authorization": f"Bearer {api_key}"}
response = requests.get(
    "http://localhost:8080/api/v1/names?origin=japanese&count=5", 
    headers=headers
)
names = response.json()
for name in names:
    print(f"{name['name']} - {name['gender']}")

# Method 2: Using X-API-Key header
headers = {"X-API-Key": api_key}
response = requests.get(
    "http://localhost:8080/api/v1/names?origin=french&gender=female", 
    headers=headers
)

# Method 3: Via URL parameter
params = {
    "origin": "chinese", 
    "count": 5,
    "api_key": api_key
}
response = requests.get("http://localhost:8080/api/v1/names", params=params)
```

#### JavaScript Example
```javascript
// Get French female names
const apiKey = "your_secret_key";

// Method 1: Using Bearer Token (recommended)
fetch("http://localhost:8080/api/v1/names?origin=french&gender=female", {
  headers: {
    "Authorization": `Bearer ${apiKey}`
  }
})
  .then(response => response.json())
  .then(data => {
    console.log(`Name: ${data.name}`);
    console.log(`Gender: ${data.gender}`);
    console.log(`First Name: ${data.first_name}`);
    console.log(`Last Name: ${data.last_name}`);
  });

// Method 2: Using X-API-Key header
fetch("http://localhost:8080/api/v1/names?origin=chinese&gender=male&count=3", {
  headers: {
    "X-API-Key": apiKey
  }
})
  .then(response => response.json())
  .then(data => console.log(data));

// Method 3: Via URL parameter
fetch(`http://localhost:8080/api/v1/names?origin=english&count=2&api_key=${apiKey}`)
  .then(response => response.json())
  .then(data => console.log(data));
```

## Server Deployment

### Docker Deployment (Recommended)

```bash
# Build Docker image
docker build -t namegen-api .

# Run container
docker run -d -p 8080:8080 --name namegen namegen-api -key your_secret_key
```

### Direct Deployment

```bash
# Build API service
go build -o namegen-api ./cmd/api

# Run service in background
nohup ./namegen-api -key your_secret_key > namegen.log 2>&1 &

# Set up as system service (using systemd)
sudo nano /etc/systemd/system/namegen.service
# Add the following content:
# [Unit]
# Description=Name Generator API Service
# After=network.target
#
# [Service]
# User=youruser
# WorkingDirectory=/path/to/namegen
# ExecStart=/path/to/namegen/namegen-api -key your_secret_key
# Restart=on-failure
#
# [Install]
# WantedBy=multi-user.target

# Enable and start service
sudo systemctl enable namegen
sudo systemctl start namegen
```

## Available Name Origins

Currently supports name generation from the following cultural backgrounds:
- anglosaxon - Anglo-Saxon names
- arabic - Arabic names
- chinese - Chinese names
- dutch - Dutch names
- english - English names
- finnish - Finnish names
- french - French names
- german - German names
- greek - Greek names
- hindu - Hindu names
- italian - Italian names
- japanese - Japanese names
- korean - Korean names
- mongolian - Mongolian names
- norwegian - Norwegian names
- polish - Polish names
- portuguese - Portuguese names
- russian - Russian names
- spanish - Spanish names
- swedish - Swedish names
- thai - Thai names
- turkish - Turkish names
- ukrainian - Ukrainian names

And fantasy world names:
- dwarf - Dwarf names
- elf - Elf names
- fantasy - Fantasy world names 