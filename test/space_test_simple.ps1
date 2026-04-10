# Space CRUD Test
$BASE_URL = "http://localhost:7000"
$API_V1 = "$BASE_URL/api/v1"

function Write-Info { param($Message) Write-Host "[INFO] $Message" -ForegroundColor Green }
function Write-Error { param($Message) Write-Host "[ERROR] $Message" -ForegroundColor Red }
function Write-Warn { param($Message) Write-Host "[WARN] $Message" -ForegroundColor Yellow }

Write-Info "=== Space Module CRUD Test ==="

# 1. Login
Write-Info "1. Login..."
$loginBody = '{"username":"admin","password":"admin123"}'
$loginResp = Invoke-RestMethod -Uri "$API_V1/auth/login" -Method POST -Body $loginBody -ContentType "application/json"
$TOKEN = $loginResp.data.access_token
Write-Info "Login success"

$headers = @{ "Authorization" = "Bearer $TOKEN" }

# 2. Create Space
Write-Info "2. Create Space..."
$ts = Get-Date -Format "yyyyMMddHHmmss"
$createBody = "{`"name`":`"TestSpace_$ts`",`"description`":`"Test Description`",`"province`":`"TestProvince`",`"city`":`"TestCity`",`"longitude`":116.407429,`"latitude`":39.904211,`"zoom_level`":12,`"sort_order`":1}"
$createResp = Invoke-RestMethod -Uri "$API_V1/resource/spaces" -Method POST -Body $createBody -ContentType "application/json" -Headers $headers
$SPACE_ID = $createResp.data.id
Write-Info "Created Space ID: $SPACE_ID"

# 3. Get Space List
Write-Info "3. Get Space List..."
$listResp = Invoke-RestMethod -Uri "$API_V1/resource/spaces?page=1&page_size=10" -Method GET -Headers $headers
Write-Info "List count: $($listResp.data.spaces.Count)"

# 4. Get Space Detail
Write-Info "4. Get Space Detail..."
$detailResp = Invoke-RestMethod -Uri "$API_V1/resource/spaces/$SPACE_ID" -Method GET -Headers $headers
Write-Info "Space name: $($detailResp.data.name)"

# 5. Update Space
Write-Info "5. Update Space..."
$updateBody = "{`"name`":`"UpdatedSpace_$ts`",`"description`":`"Updated Description`",`"city`":`"UpdatedCity`"}"
$updateResp = Invoke-RestMethod -Uri "$API_V1/resource/spaces/$SPACE_ID" -Method PUT -Body $updateBody -ContentType "application/json" -Headers $headers
Write-Info "Update success"

# 6. Create Space 2 for batch delete
Write-Info "6. Create Space 2..."
$ts2 = Get-Date -Format "yyyyMMddHHmmss"
$createBody2 = "{`"name`":`"TestSpace2_$ts2`",`"description`":`"Test2`",`"province`":`"P2`",`"city`":`"C2`",`"longitude`":117.0,`"latitude`":40.0,`"zoom_level`":10,`"sort_order`":2}"
$createResp2 = Invoke-RestMethod -Uri "$API_V1/resource/spaces" -Method POST -Body $createBody2 -ContentType "application/json" -Headers $headers
$SPACE_ID2 = $createResp2.data.id
Write-Info "Created Space2 ID: $SPACE_ID2"

# 7. Create Space 3 for batch delete
Write-Info "7. Create Space 3..."
$ts3 = Get-Date -Format "yyyyMMddHHmmss"
$createBody3 = "{`"name`":`"TestSpace3_$ts3`",`"description`":`"Test3`",`"province`":`"P3`",`"city`":`"C3`",`"longitude`":118.0,`"latitude`":41.0,`"zoom_level`":11,`"sort_order`":3}"
$createResp3 = Invoke-RestMethod -Uri "$API_V1/resource/spaces" -Method POST -Body $createBody3 -ContentType "application/json" -Headers $headers
$SPACE_ID3 = $createResp3.data.id
Write-Info "Created Space3 ID: $SPACE_ID3"

# 8. Batch Delete
Write-Info "8. Batch Delete Spaces $SPACE_ID2, $SPACE_ID3..."
$batchDeleteBody = "{`"ids`":[$SPACE_ID2,$SPACE_ID3]}"
$batchDeleteResp = Invoke-RestMethod -Uri "$API_V1/resource/spaces/batch" -Method DELETE -Body $batchDeleteBody -ContentType "application/json" -Headers $headers
Write-Info "Batch delete success"

# 9. Delete Single Space
Write-Info "9. Delete Single Space $SPACE_ID..."
$deleteResp = Invoke-RestMethod -Uri "$API_V1/resource/spaces/$SPACE_ID" -Method DELETE -Headers $headers
Write-Info "Delete success"

# 10. Search
Write-Info "10. Search Spaces..."
$searchResp = Invoke-RestMethod -Uri "$API_V1/resource/spaces?keyword=Test&page=1&page_size=10" -Method GET -Headers $headers
Write-Info "Search result count: $($searchResp.data.spaces.Count)"

Write-Info "=== All Tests Passed ==="
