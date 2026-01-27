#!/bin/bash

# 配置信息
BASE_URL="http://localhost:7000"
API_V1="${BASE_URL}/api/v1"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 打印日志
log_info() { echo -e "${GREEN}[INFO] $1${NC}"; }
log_warn() { echo -e "${YELLOW}[WARN] $1${NC}"; }
log_error() { echo -e "${RED}[ERROR] $1${NC}"; }

# 发送请求并获取响应
# 参数: $1=method, $2=url, $3=data (json string), $4=extra_headers (array)
send_request() {
    local method=$1
    local url=$2
    local data=$3
    shift 3
    local headers=("$@")

    local curl_cmd=("curl" "-s" "-X" "$method")
    curl_cmd+=("-H" "Content-Type: application/json")
    for h in "${headers[@]}"; do
        curl_cmd+=("-H" "$h")
    done
    
    if [ ! -z "$data" ]; then
        curl_cmd+=("-d" "$data")
    fi
    
    curl_cmd+=("$url")
    
    # 执行并返回响应体
    "${curl_cmd[@]}"
}

# 登录并获取令牌
# 参数: $1=username, $2=password
login() {
    local username=$1
    local password=$2
    local resp=$(send_request "POST" "${API_V1}/auth/login" "{\"username\":\"$username\",\"password\":\"$password\"}")
    echo "$resp" | grep -o '"access_token":"[^"]*' | cut -d'"' -f4
}

# 检查 HTTP 状态码
check_status() {
    local method=$1
    local url=$2
    local expected_status=$3
    local data=$4
    shift 4
    local headers=("$@")

    local curl_args=("-s" "-o" "/tmp/resp_body" "-w" "%{http_code}" "-X" "$method")
    curl_args+=("-H" "Content-Type: application/json")
    
    for h in "${headers[@]}"; do
        curl_args+=("-H" "$h")
    done
    
    if [ ! -z "$data" ]; then
        curl_args+=("-d" "$data")
    fi
    
    local status=$(curl "${curl_args[@]}" "$url")

    if [ "$status" -eq "$expected_status" ]; then
        log_info "测试通过: $method $url (预期 $expected_status, 实际 $status)"
        return 0
    else
        log_error "测试失败: $method $url (预期 $expected_status, 实际 $status)"
        log_error "响应内容: $(cat /tmp/resp_body)"
        return 1
    fi
}
