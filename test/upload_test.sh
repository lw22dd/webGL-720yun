#!/bin/bash
<<<<<<< HEAD
source "$(dirname "$0")/common.sh"

log_info "=== 开始分片上传模块测试 ==="

# 1. 健康检查
log_info "测试健康检查接口..."
check_status "GET" "${BASE_URL}/health" 200 ""

# 2. 登录获取 Token
log_info "登录获取 Token..."
TOKEN=$(login "admin" "admin123")
if [ -z "$TOKEN" ]; then
    log_error "登录失败，无法继续测试"
    exit 1
fi
log_info "登录成功，获取到 Token"

# 3. 创建测试空间
log_info "创建测试空间..."
SPACE_RESP=$(send_request "POST" "${API_V1}/resource/spaces" \
    "{\"name\":\"分片上传测试空间\",\"slug\":\"chunk-upload-test\",\"description\":\"用于测试分片上传功能\",\"province\":\"测试省\",\"city\":\"测试市\",\"longitude\":116.407429,\"latitude\":39.904211,\"zoom_level\":12,\"sort_order\":1,\"status\":1}" \
    "Authorization: Bearer $TOKEN")

SPACE_ID=$(echo "$SPACE_RESP" | grep -o '"id":[0-9]*' | head -1 | cut -d':' -f2)
if [ -z "$SPACE_ID" ]; then
    log_error "创建空间失败"
    exit 1
fi
log_info "空间创建成功，ID: $SPACE_ID"

# 4. 测试初始化上传接口
log_info "测试初始化上传接口..."
INIT_RESP=$(send_request "POST" "${API_V1}/upload/init" \
    "{\"file_name\":\"test_panorama.jpg\",\"file_size\":10485760,\"file_md5\":\"abc123def45678901234567890123456\",\"space_id\":$SPACE_ID,\"scene_code\":\"test_scene_001\",\"title\":\"测试场景\"}" \
    "Authorization: Bearer $TOKEN")

echo "初始化响应: $INIT_RESP"

UPLOAD_ID=$(echo "$INIT_RESP" | grep -o '"upload_id":"[^"]*' | cut -d'"' -f4)
if [ -z "$UPLOAD_ID" ]; then
    log_error "初始化上传失败"
    # 清理空间
    send_request "DELETE" "${API_V1}/resource/spaces/$SPACE_ID" "" "Authorization: Bearer $TOKEN"
    exit 1
fi
log_info "初始化上传成功，Upload ID: $UPLOAD_ID"

# 检查是否秒传
SKIP_UPLOAD=$(echo "$INIT_RESP" | grep -o '"skip_upload":[^,}]*' | cut -d':' -f2)
if [ "$SKIP_UPLOAD" = "true" ]; then
    log_info "文件已存在，秒传成功"
    # 清理并退出
    send_request "DELETE" "${API_V1}/resource/spaces/$SPACE_ID" "" "Authorization: Bearer $TOKEN"
    log_info "=== 分片上传测试完成 ==="
    exit 0
fi

TOTAL_CHUNKS=$(echo "$INIT_RESP" | grep -o '"total_chunks":[0-9]*' | cut -d':' -f2)
log_info "总分片数: $TOTAL_CHUNKS"

# 5. 测试分片上传接口（模拟上传2个分片）
log_info "测试分片上传接口..."

# 创建测试分片文件（1MB 的测试数据）
TEST_CHUNK_FILE="/tmp/test_chunk_$(date +%s).bin"
dd if=/dev/zero of="$TEST_CHUNK_FILE" bs=1M count=1 2>/dev/null

# 上传第一个分片
log_info "上传分片 0..."
CHUNK_RESP=$(curl -s -X POST "${API_V1}/upload/chunk" \
    -H "Authorization: Bearer $TOKEN" \
    -F "upload_id=$UPLOAD_ID" \
    -F "chunk_index=0" \
    -F "chunk_data=@$TEST_CHUNK_FILE")

echo "分片上传响应: $CHUNK_RESP"

CHUNK_INDEX=$(echo "$CHUNK_RESP" | grep -o '"chunk_index":[0-9]*' | cut -d':' -f2)
if [ "$CHUNK_INDEX" != "0" ]; then
    log_error "分片 0 上传失败"
    rm -f "$TEST_CHUNK_FILE"
    send_request "DELETE" "${API_V1}/upload/$UPLOAD_ID" "" "Authorization: Bearer $TOKEN"
    send_request "DELETE" "${API_V1}/resource/spaces/$SPACE_ID" "" "Authorization: Bearer $TOKEN"
    exit 1
fi
log_info "分片 0 上传成功"

# 6. 测试查询上传状态接口
log_info "测试查询上传状态接口..."
STATUS_RESP=$(send_request "GET" "${API_V1}/upload/status/$UPLOAD_ID" "" "Authorization: Bearer $TOKEN")
echo "状态响应: $STATUS_RESP"

STATUS=$(echo "$STATUS_RESP" | grep -o '"status":"[^"]*' | cut -d'"' -f4)
if [ -z "$STATUS" ]; then
    log_error "查询状态失败"
else
    log_info "当前状态: $STATUS"
fi

# 7. 测试取消上传接口
log_info "测试取消上传接口..."
CANCEL_RESP=$(send_request "DELETE" "${API_V1}/upload/$UPLOAD_ID" "" "Authorization: Bearer $TOKEN")
echo "取消响应: $CANCEL_RESP"

# 8. 测试秒传功能（再次上传相同MD5的文件）
log_info "测试秒传功能..."
INIT_RESP2=$(send_request "POST" "${API_V1}/upload/init" \
    "{\"file_name\":\"test_panorama.jpg\",\"file_size\":10485760,\"file_md5\":\"abc123def45678901234567890123456\",\"space_id\":$SPACE_ID,\"scene_code\":\"test_scene_002\",\"title\":\"测试场景2\"}" \
    "Authorization: Bearer $TOKEN")

echo "秒传测试响应: $INIT_RESP2"
SKIP_UPLOAD2=$(echo "$INIT_RESP2" | grep -o '"skip_upload":[^,}]*' | cut -d':' -f2)
if [ "$SKIP_UPLOAD2" = "true" ]; then
    log_info "秒传功能正常"
else
    log_warn "秒传功能可能未生效（这是正常的，因为没有实际文件内容）"
fi

# 清理测试文件
rm -f "$TEST_CHUNK_FILE"

# 9. 清理测试空间
log_info "清理测试空间..."
send_request "DELETE" "${API_V1}/resource/spaces/$SPACE_ID" "" "Authorization: Bearer $TOKEN"

log_info "=== 分片上传模块测试完成 ==="
=======
BASE_URL="http://localhost:7000"
API_V1="$BASE_URL/api/v1"

echo "=== Upload and Slice Module Test ==="

echo "1. Login..."
LOGIN_RESP=$(curl -s "$API_V1/auth/login" -X POST -H "Content-Type: application/json" -d '{"username":"admin","password":"admin123"}')
TOKEN=$(echo $LOGIN_RESP | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)
echo "Token: ${TOKEN:0:50}..."

echo ""
echo "2. Get Spaces..."
SPACES_RESP=$(curl -s "$API_V1/resource/spaces?page=1" -H "Authorization: Bearer $TOKEN")
SPACE_ID=$(echo $SPACES_RESP | grep -o '"id":[0-9]*' | head -1 | cut -d':' -f2)
echo "Space ID: $SPACE_ID"

echo ""
echo "3. Test Init Upload..."
INIT_RESP=$(curl -s "$API_V1/upload/init" -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $TOKEN" -d '{"filename":"test.jpg","file_size":1048576,"file_hash":"abc123def45678901234567890123456"}')
echo "Init: $INIT_RESP"

echo ""
echo "Test completed!"
>>>>>>> 8146554307dc850256079e5aa35fb05bb5b6a503
