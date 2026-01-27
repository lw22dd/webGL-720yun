#!/bin/bash
source "$(dirname "$0")/common.sh"

log_info "=== 开始普通用户模块测试 ==="

# 准备验证清理（可选，如果数据库允许的话）

# 1. 注册新学生用户 (RoleID=3)
SUFFIX=$(date +%s)
TEST_USER="student_${SUFFIX}"
TEST_PASS="password123"
STUDENT_ID="S${SUFFIX}"

log_info "测试学生注册 (Username: $TEST_USER, StudentID: $STUDENT_ID)..."
check_status "POST" "${API_V1}/user/register" 200 "{\"username\":\"$TEST_USER\",\"password\":\"$TEST_PASS\",\"email\":\"$TEST_USER@example.com\",\"role_id\":3,\"student_id\":\"$STUDENT_ID\"}"

# 2. 注册新教师用户 (RoleID=2)
TEACHER_USER="teacher_${SUFFIX}"
log_info "测试教师注册 (Username: $TEACHER_USER)..."
check_status "POST" "${API_V1}/user/register" 200 "{\"username\":\"$TEACHER_USER\",\"password\":\"$TEST_PASS\",\"email\":\"$TEACHER_USER@example.com\",\"role_id\":2}"

# 登录测试学生
TOKEN=$(login "$TEST_USER" "$TEST_PASS")

if [ ! -z "$TOKEN" ]; then
    AUTH_HEADER="Authorization: Bearer $TOKEN"

    # 3. 获取个人资料
    log_info "测试获取个人资料..."
    check_status "GET" "${API_V1}/user/profile" 200 "" "$AUTH_HEADER"

    # 4. 更新个人资料
    log_info "测试更新个人资料..."
    check_status "PUT" "${API_V1}/user/profile" 200 "{\"nickname\":\"学生_${SUFFIX}\"}" "$AUTH_HEADER"

    # 5. 修改密码
    log_info "测试修改密码..."
    check_status "POST" "${API_V1}/user/change-password" 200 "{\"old_password\":\"$TEST_PASS\",\"new_password\":\"newpass_${SUFFIX}\"}" "$AUTH_HEADER"
    
    # 用新密码重新登录验证
    NEW_TOKEN=$(login "$TEST_USER" "newpass_${SUFFIX}")
    if [ ! -z "$NEW_TOKEN" ]; then
        log_info "新密码登录成功"
        AUTH_HEADER="Authorization: Bearer $NEW_TOKEN"
    else
        log_error "新密码登录失败"
    fi

    # 6. 退出登录
    log_info "测试退出登录..."
    check_status "POST" "${API_V1}/user/logout" 200 "" "$AUTH_HEADER"
else
    log_error "测试用户登录失败，跳过后续测试"
fi

log_info "=== 普通用户模块测试完成 ==="
