Import-Module "$PSScriptRoot\common.ps1" -Force

Log-Info "=== 开始普通用户模块测试 ==="

$SUFFIX = [DateTimeOffset]::Now.ToUnixTimeSeconds()
$TEST_USER = "student_$SUFFIX"
$TEST_PASS = "password123"
$STUDENT_ID = "S$SUFFIX"

Log-Info "测试学生注册 (Username: $TEST_USER, StudentID: $STUDENT_ID)..."
Check-Status -Method "POST" -Url "$script:API_V1/user/register" -ExpectedStatus 200 -Data "{`"username`":`"$TEST_USER`",`"password`":`"$TEST_PASS`",`"email`":`"$TEST_USER@example.com`",`"role_id`":3,`"student_id`":`"$STUDENT_ID`"}"

$TEACHER_USER = "teacher_$SUFFIX"
Log-Info "测试教师注册 (Username: $TEACHER_USER)..."
Check-Status -Method "POST" -Url "$script:API_V1/user/register" -ExpectedStatus 200 -Data "{`"username`":`"$TEACHER_USER`",`"password`":`"$TEST_PASS`",`"email`":`"$TEACHER_USER@example.com`",`"role_id`":2}"

$TOKEN = Login -Username $TEST_USER -Password $TEST_PASS

if ($TOKEN) {
    $AUTH_HEADER = "Authorization: Bearer $TOKEN"

    Log-Info "测试获取个人资料..."
    Check-Status -Method "GET" -Url "$script:API_V1/user/profile" -ExpectedStatus 200 -Data "" -Headers @($AUTH_HEADER)

    Log-Info "测试更新个人资料..."
    Check-Status -Method "PUT" -Url "$script:API_V1/user/profile" -ExpectedStatus 200 -Data "{`"nickname`":`"学生_$SUFFIX`"}" -Headers @($AUTH_HEADER)

    Log-Info "测试修改密码..."
    Check-Status -Method "POST" -Url "$script:API_V1/user/change-password" -ExpectedStatus 200 -Data "{`"old_password`":`"$TEST_PASS`",`"new_password`":`"newpass_$SUFFIX`"}" -Headers @($AUTH_HEADER)

    $NEW_TOKEN = Login -Username $TEST_USER -Password "newpass_$SUFFIX"
    if ($NEW_TOKEN) {
        Log-Info "新密码登录成功"
        $AUTH_HEADER = "Authorization: Bearer $NEW_TOKEN"
    }
    else {
        Log-Error "新密码登录失败"
    }

    Log-Info "测试退出登录..."
    Check-Status -Method "POST" -Url "$script:API_V1/user/logout" -ExpectedStatus 200 -Data "" -Headers @($AUTH_HEADER)
}
else {
    Log-Error "测试用户登录失败，跳过后续测试"
}

Log-Info "=== 普通用户模块测试完成 ==="
