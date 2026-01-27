#!/bin/bash
source "$(dirname "$0")/common.sh"

log_info "=== 开始管理员模块测试 ==="

# 登录管理员
ADMIN_TOKEN=$(login "admin" "admin123")
if [ -z "$ADMIN_TOKEN" ]; then
    log_error "管理员登录失败"
    exit 1
fi
AUTH_HEADER="Authorization: Bearer $ADMIN_TOKEN"

# 1. 获取用户列表
log_info "测试管理员获取用户列表..."
check_status "GET" "${API_V1}/user/admin/list?page=1&page_size=10" 200 "" "$AUTH_HEADER"

# 2. 管理员创建用户 (教师角色的创建)
SUFFIX=$(date +%s)
NEW_USER="adm_${SUFFIX}"
log_info "测试管理员创建教师用户 (Username: $NEW_USER)..."
check_status "POST" "${API_V1}/user/admin/create" 200 "{\"username\":\"$NEW_USER\",\"password\":\"pass123\",\"email\":\"$NEW_USER@example.com\",\"role_id\":2}" "$AUTH_HEADER"

# 3. 获取特定用户信息 (ID=1 是默认管理员)
log_info "测试管理员获取 ID=1 的用户信息..."
check_status "GET" "${API_V1}/user/admin/1" 200 "" "$AUTH_HEADER"

# 4. 更新用户信息 (ID=1)
log_info "测试管理员更新 ID=1 的昵称..."
check_status "PUT" "${API_V1}/user/admin/1" 200 "{\"nickname\":\"超级管理员旗舰版\"}" "$AUTH_HEADER"

# 5. RBAC 权限越权测试
log_info "测试 RBAC 越权拦截 (普通学生访问管理员列表)..."
# 注册并登录一个普通学生
S_USER="rbac_stu_${SUFFIX}"
S_ID="RS${SUFFIX}"
send_request "POST" "${API_V1}/user/register" "{\"username\":\"$S_USER\",\"password\":\"pass123\",\"email\":\"$S_USER@example.com\",\"role_id\":3,\"student_id\":\"$S_ID\"}"
U_TOKEN=$(login "$S_USER" "pass123")
if [ ! -z "$U_TOKEN" ]; then
    U_HEADER="Authorization: Bearer $U_TOKEN"
    check_status "GET" "${API_V1}/user/admin/list" 403 "" "$U_HEADER"
else
    log_error "注册/登录 RBAC 测试用户失败"
fi

log_info "=== 管理员模块测试完成 ==="
