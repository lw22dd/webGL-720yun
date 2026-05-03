Import-Module "$PSScriptRoot\common.ps1" -Force

Log-Info "=== 开始认证模块测试 ==="

Log-Info "测试健康检查接口..."
Check-Status -Method "GET" -Url "$script:BASE_URL/health" -ExpectedStatus 200 -Data ""

Log-Info "测试默认管理员登录..."
$TOKEN = Login -Username "admin" -Password "admin123"
if ($TOKEN) {
    Log-Info "登录成功，获取到 Token"
}
else {
    Log-Error "登录失败"
    exit 1
}

Log-Info "测试错误密码登录..."
Check-Status -Method "POST" -Url "$script:API_V1/auth/login" -ExpectedStatus 401 -Data "{`"username`":`"admin`",`"password`":`"wrongpass`"}"

Log-Info "测试刷新令牌..."
$RESP = Send-Request -Method "POST" -Url "$script:API_V1/auth/login" -Data "{`"username`":`"admin`",`"password`":`"admin123`"}"
if ($RESP -match '"refresh_token`":`"([^`"]+)"') {
    $REFRESH_TOKEN = $matches[1]
    Check-Status -Method "POST" -Url "$script:API_V1/auth/refresh" -ExpectedStatus 200 -Data "{`"refresh_token`":`"$REFRESH_TOKEN`"}"
}
else {
    Log-Warn "未获取到 Refresh Token，跳过刷新测试"
}

Log-Info "=== 认证模块测试完成 ==="
