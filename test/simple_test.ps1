$BASE_URL = "http://localhost:7000"
$API_V1 = "$BASE_URL/api/v1"

Write-Host "Testing Upload and Slice Module..." -ForegroundColor Green

Write-Host "1. Login..." -ForegroundColor Cyan
$loginResp = curl.exe -s "$API_V1/auth/login" -X POST -H "Content-Type: application/json" -d '{"username":"admin","password":"admin123"}'
Write-Host $loginResp

$TOKEN = ($loginResp | ConvertFrom-Json).data.access_token
Write-Host "Token: $TOKEN" -ForegroundColor Yellow

Write-Host "`n2. Get Spaces..." -ForegroundColor Cyan
$spacesResp = curl.exe -s "$API_V1/resource/spaces?page=1" -H "Authorization: Bearer $TOKEN"
Write-Host $spacesResp

$spacesData = $spacesResp | ConvertFrom-Json
if ($spacesData.data.spaces.Count -gt 0) {
    $SPACE_ID = $spacesData.data.spaces[0].id
    Write-Host "Using Space ID: $SPACE_ID" -ForegroundColor Yellow
} else {
    Write-Host "No spaces found" -ForegroundColor Red
}

Write-Host "`nTest completed!" -ForegroundColor Green
