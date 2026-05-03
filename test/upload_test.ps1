Import-Module "$PSScriptRoot\common.ps1" -Force

$BASE_URL = "http://localhost:7000"
$API_V1 = "$BASE_URL/api/v1"

Write-Host "=== Upload and Slice Module Test ===" -ForegroundColor Cyan

Write-Host "1. Login..."
$LOGIN_RESP = Send-Request -Method "POST" -Url "$API_V1/auth/login" -Data "{`"username`":`"admin`",`"password`":`"admin123`"}"
if ($LOGIN_RESP -match '"access_token`":`"([^`"]+)"') {
    $TOKEN = $matches[1]
    Write-Host "Token: $($TOKEN.Substring(0, [Math]::Min(50, $TOKEN.Length)))..."
}
else {
    Write-Host "Login failed!" -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "2. Get Spaces..."
$SPACES_RESP = Send-Request -Method "GET" -Url "$API_V1/resource/spaces?page=1" -Headers @("Authorization: Bearer $TOKEN")
if ($SPACES_RESP -match '"id`":(\d+)') {
    $SPACE_ID = $matches[1]
    Write-Host "Space ID: $SPACE_ID"
}
else {
    Write-Host "Failed to get spaces" -ForegroundColor Red
}

Write-Host ""
Write-Host "3. Test Init Upload..."
$INIT_RESP = Send-Request -Method "POST" -Url "$API_V1/upload/init" -Data "{`"filename`":`"test.jpg`",`"file_size`":1048576,`"file_hash`":`"abc123def45678901234567890123456`"}" -Headers @("Authorization: Bearer $TOKEN")
Write-Host "Init: $INIT_RESP"

Write-Host ""
Write-Host "Test completed!" -ForegroundColor Green
