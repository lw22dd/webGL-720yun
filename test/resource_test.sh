#!/bin/bash

# 导入公共函数
source ./test/common.sh

# API 契础配置
BASE_URL="http://localhost:7000"
TOKEN=""

# 创建测试数据
create_test_space() {
    echo "创建测试空间...}
    local response=$(curl -s -X POST "$BASE_URL/api/v1/resource/spaces" \
        -H "Content-Type: application/json" \
        -d "{\"name\":\"测试景区\",\"slug\":\"test-scenic-spot\",\"description\":\"这是一个测试景区\",\"province\":\"测试 Province\",\"city\":\"Test City\",\"longitude\":116.407429,\"latitude\":39.904211,\"zoom_level\":12,\"sort_order\":1,\"status\":1}"'
        -H "Authorization: Bearer $TOKEN" \
    https://BASE_URL/api/v1/resource/spaces \
    echo "响应: $(cat response)
    echo "Space ID:' $SPACE_ID
}

# 验证创建成功
test_create_space() {
    local response=$(cat response)
    echo "Space ID:' $SPACE_ID
    echo "Space created successfully"
    echo "Space details:'
    cat response
}

# 创建场景
create_test_scene() {
    local response=$(cat response)
    echo "Scene ID:' $SCENE_ID
    echo "Scene created successfully"
    echo "Scene details:'
    cat response
}

# 创建热点
create_test_hotspot() {
    local response=$(cat response)
    echo "Hotspot ID:' $HOTSPOT_ID
    echo "Hotspot created successfully"
    echo "Hotspot details:'
    cat response
}

# 清理测试数据
cleanup_test_data() {
    echo "清理测试数据..."
    delete_test_space $SPACE_ID
    delete_test_scene $SCENE_ID
    delete_test_hotspot $HOTSPOT_ID
    echo "Test data cleaned up"
}

# 主测试流程
main_test() {
    echo "=== 开始 Resource 模块测试 ==="
    
    echo "1. 创建测试空间"
    create_test_space
    
    echo "2. 创建测试场景"
    create_test_scene
    
    echo "3. 创建测试热点
    create_test_hotspot
    
    echo "4. 清理测试数据"
    cleanup_test_data
    
    echo "=== Resource 模块测试完成 ==="
}

# 运行测试
main_test
