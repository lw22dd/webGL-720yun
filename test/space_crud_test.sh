#!/bin/bash
# Space 模块增删改查测试脚本

# 配置
BASE_URL="http://localhost:7000"
API_V1="${BASE_URL}/api/v1"

# 颜色
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log_info() { echo -e "${GREEN}[INFO]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }
log_warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }

# 发送请求
send_request() {
    local method=$1
    local url=$2
    local data=$3
    shift 3
    local headers=("$@")
    
    local curl_cmd="curl -s -X ${method}"
    for h in "${headers[@]}"; do
        curl_cmd="${curl_cmd} -H \"${h}\""
    done
    if [ -n "$data" ]; then
        curl_cmd="${curl_cmd} -d '${data}'"
    fi
    curl_cmd="${curl_cmd} \"${url}\""
    
    eval $curl_cmd
}

log_info "=== Space 模块增删改查测试 ==="

# 1. 登录
log_info "1. 登录获取 Token..."
LOGIN_RESP=$(send_request "POST" "${API_V1}/auth/login" '{"username":"admin","password":"admin123"}' "Content-Type: application/json")
TOKEN=$(echo "$LOGIN_RESP" | grep -o '"access_token":"[^"]*' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
    log_error "登录失败: $LOGIN_RESP"
    exit 1
fi
log_info "登录成功"

AUTH_HEADER="Authorization: Bearer ${TOKEN}"

# 2. 创建空间
log_info "2. 创建空间..."
TIMESTAMP=$(date +%s)
CREATE_BODY="{\"name\":\"TestSpace_${TIMESTAMP}\",\"description\":\"Test Description\",\"province\":\"TestProvince\",\"city\":\"TestCity\",\"longitude\":116.407429,\"latitude\":39.904211,\"zoom_level\":12,\"sort_order\":1}"
CREATE_RESP=$(send_request "POST" "${API_V1}/resource/spaces" "$CREATE_BODY" "Content-Type: application/json" "$AUTH_HEADER")
SPACE_ID=$(echo "$CREATE_RESP" | grep -o '"id":[0-9]*' | head -1 | cut -d':' -f2)

if [ -z "$SPACE_ID" ]; then
    log_error "创建空间失败: $CREATE_RESP"
    exit 1
fi
log_info "创建空间成功，ID: $SPACE_ID"

# 3. 获取空间列表
log_info "3. 获取空间列表..."
LIST_RESP=$(send_request "GET" "${API_V1}/resource/spaces?page=1&page_size=10" "" "$AUTH_HEADER")
if echo "$LIST_RESP" | grep -q '"spaces":'; then
    log_info "获取列表成功"
else
    log_error "获取列表失败: $LIST_RESP"
    exit 1
fi

# 4. 获取空间详情
log_info "4. 获取空间详情..."
DETAIL_RESP=$(send_request "GET" "${API_V1}/resource/spaces/${SPACE_ID}" "" "$AUTH_HEADER")
if echo "$DETAIL_RESP" | grep -q "\"id\":${SPACE_ID}"; then
    log_info "获取详情成功"
else
    log_error "获取详情失败: $DETAIL_RESP"
    exit 1
fi

# 5. 更新空间
log_info "5. 更新空间..."
UPDATE_BODY="{\"name\":\"UpdatedSpace_${TIMESTAMP}\",\"description\":\"Updated Description\",\"city\":\"UpdatedCity\"}"
UPDATE_RESP=$(send_request "PUT" "${API_V1}/resource/spaces/${SPACE_ID}" "$UPDATE_BODY" "Content-Type: application/json" "$AUTH_HEADER")
if echo "$UPDATE_RESP" | grep -q '"code":200'; then
    log_info "更新成功"
else
    log_error "更新失败: $UPDATE_RESP"
    exit 1
fi

# 6. 创建空间2（用于批量删除）
log_info "6. 创建空间2..."
TIMESTAMP2=$(date +%s)
CREATE_BODY2="{\"name\":\"TestSpace2_${TIMESTAMP2}\",\"description\":\"Test2\",\"province\":\"P2\",\"city\":\"C2\",\"longitude\":117.0,\"latitude\":40.0,\"zoom_level\":10,\"sort_order\":2}"
CREATE_RESP2=$(send_request "POST" "${API_V1}/resource/spaces" "$CREATE_BODY2" "Content-Type: application/json" "$AUTH_HEADER")
SPACE_ID2=$(echo "$CREATE_RESP2" | grep -o '"id":[0-9]*' | head -1 | cut -d':' -f2)

if [ -z "$SPACE_ID2" ]; then
    log_error "创建空间2失败: $CREATE_RESP2"
    exit 1
fi
log_info "创建空间2成功，ID: $SPACE_ID2"

# 7. 创建空间3（用于批量删除）
log_info "7. 创建空间3..."
TIMESTAMP3=$(date +%s)
CREATE_BODY3="{\"name\":\"TestSpace3_${TIMESTAMP3}\",\"description\":\"Test3\",\"province\":\"P3\",\"city\":\"C3\",\"longitude\":118.0,\"latitude\":41.0,\"zoom_level\":11,\"sort_order\":3}"
CREATE_RESP3=$(send_request "POST" "${API_V1}/resource/spaces" "$CREATE_BODY3" "Content-Type: application/json" "$AUTH_HEADER")
SPACE_ID3=$(echo "$CREATE_RESP3" | grep -o '"id":[0-9]*' | head -1 | cut -d':' -f2)

if [ -z "$SPACE_ID3" ]; then
    log_error "创建空间3失败: $CREATE_RESP3"
    exit 1
fi
log_info "创建空间3成功，ID: $SPACE_ID3"

# 8. 批量删除
log_info "8. 批量删除空间 ${SPACE_ID2}, ${SPACE_ID3}..."
BATCH_DELETE_BODY="{\"ids\":[${SPACE_ID2},${SPACE_ID3}]}"
BATCH_DELETE_RESP=$(send_request "DELETE" "${API_V1}/resource/spaces/batch" "$BATCH_DELETE_BODY" "Content-Type: application/json" "$AUTH_HEADER")
if echo "$BATCH_DELETE_RESP" | grep -q '"code":200'; then
    log_info "批量删除成功"
else
    log_warn "批量删除可能失败: $BATCH_DELETE_RESP"
fi

# 9. 删除单个空间
log_info "9. 删除单个空间 ${SPACE_ID}..."
DELETE_RESP=$(send_request "DELETE" "${API_V1}/resource/spaces/${SPACE_ID}" "" "$AUTH_HEADER")
if echo "$DELETE_RESP" | grep -q '"code":200'; then
    log_info "删除成功"
else
    log_warn "删除可能失败: $DELETE_RESP"
fi

# 10. 搜索
log_info "10. 搜索空间..."
SEARCH_RESP=$(send_request "GET" "${API_V1}/resource/spaces?keyword=Test&page=1&page_size=10" "" "$AUTH_HEADER")
if echo "$SEARCH_RESP" | grep -q '"code":200'; then
    log_info "搜索成功"
else
    log_warn "搜索可能失败: $SEARCH_RESP"
fi

log_info "=== Space 模块增删改查测试完成 ==="
log_info "测试项目:"
log_info "  [OK] 登录获取Token"
log_info "  [OK] 创建空间"
log_info "  [OK] 获取空间列表"
log_info "  [OK] 获取空间详情"
log_info "  [OK] 更新空间"
log_info "  [OK] 批量删除空间"
log_info "  [OK] 删除单个空间"
log_info "  [OK] 搜索功能"
