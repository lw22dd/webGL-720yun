Import-Module "$PSScriptRoot\common.ps1" -Force

Log-Info "=== 开始管理员模块测试 ==="

$ADMIN_TOKEN = Login -Username "admin" -Password "admin123"
if (-not $ADMIN_TOKEN) {
    Log-Error "管理员登录失败"
    exit 1
}
$AUTH_HEADER = "Authorization: Bearer $ADMIN_TOKEN"

Log-Info "测试管理员获取用户列表..."
Check-Status -Method "GET" -Url "$script:API_V1/user/admin/list?page=1&page_size=10" -ExpectedStatus 200 -Data "" -Headers @($AUTH_HEADER)

$SUFFIX = [DateTimeOffset]::Now.ToUnixTimeSeconds()
$NEW_USER = "adm_$SUFFIX"
Log-Info "测试管理员创建教师用户 (Username: $NEW_USER)..."
Check-Status -Method "POST" -Url "$script:API_V1/user/admin/create" -ExpectedStatus 200 -Data "{`"username`":`"$NEW_USER`",`"password`":`"pass123`",`"email`":`"$NEW_USER@example.com`",`"role_id`":2}" -Headers @($AUTH_HEADER)

Log-Info "测试管理员获取 ID=1 的用户信息..."
Check-Status -Method "GET" -Url "$script:API_V1/user/admin/1" -ExpectedStatus 200 -Data "" -Headers @($AUTH_HEADER)

Log-Info "测试管理员更新 ID=1 的昵称..."
Check-Status -Method "PUT" -Url "$script:API_V1/user/admin/1" -ExpectedStatus 200 -Data "{`"nickname`":`"超级管理员旗舰版`"}" -Headers @($AUTH_HEADER)

Log-Info "测试 RBAC 越权拦截 (普通学生访问管理员列表)..."
$S_USER = "rbac_stu_$SUFFIX"
$S_ID = "RS$SUFFIX"
Send-Request -Method "POST" -Url "$script:API_V1/user/register" -Data "{`"username`":`"$S_USER`",`"password`":`"pass123`",`"email`":`"$S_USER@example.com`",`"role_id`":3,`"student_id`":`"$S_ID`"}"
$U_TOKEN = Login -Username $S_USER -Password "pass123"
if ($U_TOKEN) {
    $U_HEADER = "Authorization: Bearer $U_TOKEN"
    Check-Status -Method "GET" -Url "$script:API_V1/user/admin/list" -ExpectedStatus 403 -Data "" -Headers @($U_HEADER)
}
else {
    Log-Error "注册/登录 RBAC 测试用户失败"
}

Log-Info "=== 管理员模块测试完成 ==="
