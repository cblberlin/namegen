#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import requests
import sys

# 设置API地址
API_URL = "http://51.210.103.207:8080/api/v1/names"

# 要测试的origin列表
ORIGINS = [
    "english", "italian", "spanish", "french", "polish", 
    "dutch", "german", "irish", "russian", "chinese"
]

def test_origin(origin):
    url = f"{API_URL}/api/v1/names"
    headers = {
        "Accept": "application/json"
    }
    params = {
        "origin": origin,
        "gender": "male",
        "count": 1
    }
    
    try:
        response = requests.get(url, headers=headers, params=params, timeout=10)
        if response.status_code == 200:
            data = response.json()
            print(f"Origin: {origin}")
            print(f"返回的名字: {data['name']}")
            print(f"返回的origin: {data['origin']}")
            print(f"请求的origin与返回的是否一致: {'是' if data['origin'] == origin else '否'}")
            print("-" * 50)
            return True
        else:
            print(f"Origin: {origin} - 请求失败，状态码: {response.status_code}")
            print(f"错误消息: {response.text}")
            print("-" * 50)
            return False
    except Exception as e:
        print(f"Origin: {origin} - 请求异常: {str(e)}")
        print("-" * 50)
        return False

def test_all_origins():
    print("开始测试所有origin...")
    print("=" * 50)
    
    success_count = 0
    for origin in ORIGINS:
        if test_origin(origin):
            success_count += 1
    
    print("=" * 50)
    print(f"测试完成: {success_count}/{len(ORIGINS)} 个origin测试成功")

def main():
    if len(sys.argv) > 1:
        origin = sys.argv[1]
        test_origin(origin)
    else:
        test_all_origins()

if __name__ == "__main__":
    main() 