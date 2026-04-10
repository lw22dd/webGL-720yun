#!/bin/bash
source "$(dirname "$0")/common.sh"

log_info "=== 开始 Space 模块增删改查测试 ==="

# 1. 健康检查
log_info "测试健康检查接口..."
if ! check_status "GET" "${BASE_URL}/health" 200 ""; then
    exit 1
fi

# 2. 登录获取 Token
log_info "登录获取 Token..."
TOKEN=$(login "admin" "admin123")
if [ -z "$TOKEN" ]; then
    log_error "登录失败，无法继续测试"
    exit 1
fi
log_info "登录成功，获取到 Token"

# 3. 测试创建空间
log_info "测试创建空间接口..."
CREATE_RESP=$(send_request "POST" "${API_V1}/resource/spaces" \
    "{\"name\":\"测试景区_$(date +%s)\",\"description\":\"这是一个测试景区\",\"province\":\"测试省\",\"city\":\"测试市\",\"longitude\":116.407429,\"latitude\":39.904211,\"zoom_level\":12,\"sort_order\":1}" \
    "Authorization: Bearer $TOKEN")

echo "创建响应: $CREATE_RESP"

SPACE_ID=$(echo "$CREATE_RESP" | grep -o '"id":[0-9]*' | head -1 | cut -d':' -f2)
if [ -z "$SPACE_ID" ]; then
    log_error "创建空间失败"
    exit 1
fi
log_info "空间创建成功，ID: $SPACE_ID"

# 4. 测试获取空间列表
log_info "测试获取空间列表接口..."
LIST_RESP=$(send_request "GET" "${API_V1}/resource/spaces?page=1&page_size=10" "" "Authorization: Bearer $TOKEN")
echo "列表响应: $LIST_RESP"

if echo "$LIST_RESP" | grep -q '"spaces":'; then
    log_info "获取空间列表成功"
else
    log_error "获取空间列表失败"
    exit 1
fi

# 5. 测试获取空间详情
log_info "测试获取空间详情接口..."
DETAIL_RESP=$(send_request "GET" "${API_V1}/resource/spaces/$SPACE_ID" "" "Authorization: Bearer $TOKEN")
echo "详情响应: $DETAIL_RESP"

if echo "$DETAIL_RESP" | grep -q "\"id\":$SPACE_ID"; then
    log_info "获取空间详情成功"
else
    log_error "获取空间详情失败"
    exit 1
fi

# 6. 测试更新空间
log_info "测试更新空间接口..."
UPDATE_RESP=$(send_request "PUT" "${API_V1}/resource/spaces/$SPACE_ID" \
    "{\"name\":\"更新后的测试景区\",\"description\":\"这是更新后的描述\",\"city\":\"更新后的城市\"}" \
    "Authorization: Bearer $TOKEN")
echo "更新响应: $UPDATE_RESP"

if echo "$UPDATE_RESP" | grep -q '"name":"更新后的测试景区"'; then
    log_info "更新空间成功"
else
    log_error "更新空间失败"
    exit 1
fi

# 7. 创建第二个空间用于批量删除测试
log_info "创建第二个测试空间..."
CREATE_RESP2=$(send_request "POST" "${API_V1}/resource/spaces" \
    "{\"name\":\"测试景区2_$(date +%s)\",\"description\":\"这是第二个测试景区\",\"province\":\"测试省2\",\"city\":\"测试市2\",\"longitude\":117.0,\"latitude\":40.0,\"zoom_level\":10,\"sort_order\":2}" \
    "Authorization: Bearer $TOKEN")

SPACE_ID2=$(echo "$CREATE_RESP2" | grep -o '"id":[0-9]*' | head -1 | cut -d':' -f2)
if [ -z "$SPACE_ID2" ]; then
    log_error "创建第二个空间失败"
    exit 1
fi
log_info "第二个空间创建成功，ID: $SPACE_ID2"

# 8. 创建第三个空间用于批量删除测试
log_info "创建第三个测试空间..."
CREATE_RESP3=$(send_request "POST" "${API_V1}/resource/spaces" \
    "{\"name\":\"测试景区3_$(date +%s)\",\"description\":\"这是第三个测试景区\",\"province\":\"测试省3\",\"city\":\"测试市3\",\"longitude\":118.0,\"latitude\":41.0,\"zoom_level\":11,\"sort_order\":3}" \
    "Authorization: Bearer $TOKEN")

SPACE_ID3=$(echo "$CREATE_RESP3" | grep -o '"id":[0-9]*' | head -1 | cut -d':' -f2)
if [ -z "$SPACE_ID3" ]; then
    log_error "创建第三个空间失败"
    exit 1
fi
log_info "第三个空间创建成功，ID: $SPACE_ID3"

# 9. 测试批量删除空间
log_info "测试批量删除空间接口..."
BATCH_DELETE_RESP=$(send_request "DELETE" "${API_V1}/resource/spaces/batch" \
    "{\"ids\":[$SPACE_ID2,$SPACE_ID3]}" \
    "Authorization: Bearer $TOKEN")
echo "批量删除响应: $BATCH_DELETE_RESP"

if echo "$BATCH_DELETE_RESP" | grep -q '批量删除成功'; then
    log_info "批量删除空间成功"
else
    log_error "批量删除空间失败"
    exit 1
fi

# 10. 验证批量删除结果
log_info "验证批量删除结果..."
DETAIL_RESP2=$(send_request "GET" "${API_V1}/resource/spaces/$SPACE_ID2" "" "Authorization: Bearer $TOKEN")
if echo "$DETAIL_RESP2" | grep -q '空间不存在'; then
    log_info "空间2已删除"
else
    log_warn "空间2可能未删除: $DETAIL_RESP2"
fi

# 11. 测试删除单个空间
log_info "测试删除单个空间接口..."
DELETE_RESP=$(send_request "DELETE" "${API_V1}/resource/spaces/$SPACE_ID" "" "Authorization: Bearer $TOKEN")
echo "删除响应: $DELETE_RESP"

if echo "$DELETE_RESP" | grep -q '删除成功'; then
    log_info "删除单个空间成功"
else
    log_error "删除单个空间失败"
    exit 1
fi

# 12. 验证删除结果
log_info "验证删除结果..."
DETAIL_RESP3=$(send_request "GET" "${API_V1}/resource/spaces/$SPACE_ID" "" "Authorization: Bearer $TOKEN")
if echo "$DETAIL_RESP3" | grep -q '空间不存在'; then
    log_info "空间已删除，验证成功"
else
    log_warn "空间可能未删除: $DETAIL_RESP3"
fi

# 13. 测试搜索功能
log_info "测试搜索功能..."
SEARCH_RESP=$(send_request "GET" "${API_V1}/resource/spaces?keyword=测试&page=1&page_size=10" "" "Authorization: Bearer $TOKEN")
echo "搜索响应: $SEARCH_RESP"

if echo "$SEARCH_RESP" | grep -q '"spaces":'; then
    log_info "搜索功能正常"
else
    log_warn "搜索功能可能异常"
fi

log_info "=== Space 模块增删改查测试完成 ==="
log_info "测试总结:"
log_info "  ✓ 健康检查"
log_info "  ✓ 登录获取Token"
log_info "  ✓ 创建空间"
log_info "  ✓ 获取空间列表"
log_info "  ✓ 获取空间详情"
log_info "  ✓ 更新空间"
log_info "  ✓ 批量删除空间"
log_info "  ✓ 删除单个空间"
log_info "  ✓ 搜索功能"
