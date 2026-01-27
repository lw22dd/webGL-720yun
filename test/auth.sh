#!/bin/bash
source "$(dirname "$0")/common.sh"

log_info "=== 开始认证模块测试 ==="

# 1. 健康检查
log_info "测试健康检查接口..."
check_status "GET" "${BASE_URL}/health" 200 ""

# 2. 默认管理员登录
log_info "测试默认管理员登录..."
TOKEN=$(login "admin" "admin123")
if [ ! -z "$TOKEN" ]; then
    log_info "登录成功，获取到 Token"
else
    log_error "登录失败"
    exit 1
fi

# 3. 错误密码登录
log_info "测试错误密码登录..."
check_status "POST" "${API_V1}/auth/login" 401 "{\"username\":\"admin\",\"password\":\"wrongpass\"}"

# 4. 刷新令牌 (需要先登录获取 Refresh Token)
log_info "测试刷新令牌..."
RESP=$(send_request "POST" "${API_V1}/auth/login" "{\"username\":\"admin\",\"password\":\"admin123\"}" "Content-Type: application/json")
REFRESH_TOKEN=$(echo "$RESP" | grep -o '"refresh_token":"[^"]*' | cut -d'"' -f4)

if [ ! -z "$REFRESH_TOKEN" ]; then
    check_status "POST" "${API_V1}/auth/refresh" 200 "{\"refresh_token\":\"$REFRESH_TOKEN\"}"
else
    log_warn "未获取到 Refresh Token，跳过刷新测试"
fi

log_info "=== 认证模块测试完成 ==="
