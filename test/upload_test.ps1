# Chunk Upload Test Script (PowerShell)
$BASE_URL = "http://localhost:7000"
$API_V1 = "$BASE_URL/api/v1"

function Write-Info { param($Message) Write-Host "[INFO] $Message" -ForegroundColor Green }
function Write-Error { param($Message) Write-Host "[ERROR] $Message" -ForegroundColor Red }
function Write-Warn { param($Message) Write-Host "[WARN] $Message" -ForegroundColor Yellow }

Write-Info "=== Starting Chunk Upload Module Test ==="

# 1. Health Check
Write-Info "Testing health check endpoint..."
try {
    $health = Invoke-RestMethod -Uri "$BASE_URL/health" -Method GET
    Write-Info "Health check passed: $($health.status)"
} catch {
    Write-Error "Health check failed: $_"
    exit 1
}

# 2. Login to get Token
Write-Info "Logging in to get Token..."
try {
    $loginBody = @{ username = "admin"; password = "admin123" } | ConvertTo-Json
    $loginResp = Invoke-RestMethod -Uri "$API_V1/auth/login" -Method POST -Body $loginBody -ContentType "application/json"
    $TOKEN = $loginResp.data.access_token
    Write-Info "Login successful, got Token"
} catch {
    Write-Error "Login failed: $_"
    exit 1
}

# 3. Get existing space
Write-Info "Getting existing space..."
try {
    $spacesResp = Invoke-RestMethod -Uri "$API_V1/resource/spaces?page=1&page_size=1" -Method GET -Headers @{ "Authorization" = "Bearer $TOKEN" }
    if ($spacesResp.data.spaces.Count -eq 0) {
        Write-Error "No spaces found, please create a space first"
        exit 1
    }
    $SPACE_ID = $spacesResp.data.spaces[0].id
    $SPACE_SLUG = $spacesResp.data.spaces[0].slug
    Write-Info "Using existing space, ID: $SPACE_ID, Slug: $SPACE_SLUG"
} catch {
    Write-Error "Failed to get space: $_"
    exit 1
}

# 4. Test init upload endpoint
Write-Info "Testing init upload endpoint..."
$SCENE_CODE = "test_scene_$(Get-Random)"
try {
    $initBody = @{
        file_name = "test_panorama.jpg"
        file_size = 10485760
        file_md5 = "abc123def45678901234567890123456"
        space_id = [int]$SPACE_ID
        scene_code = $SCENE_CODE
        title = "Test Scene"
    } | ConvertTo-Json
    
    $initResp = Invoke-RestMethod -Uri "$API_V1/upload/init" -Method POST -Body $initBody -ContentType "application/json" -Headers @{ "Authorization" = "Bearer $TOKEN" }
    Write-Info "Init response: $($initResp | ConvertTo-Json -Compress)"
    
    if ($initResp.data.skip_upload -eq $true) {
        Write-Info "File already exists, instant upload successful"
        Write-Info "=== Chunk Upload Test Completed ==="
        exit 0
    }
    
    $UPLOAD_ID = $initResp.data.upload_id
    $TOTAL_CHUNKS = $initResp.data.total_chunks
    Write-Info "Init upload successful, Upload ID: $UPLOAD_ID, Total chunks: $TOTAL_CHUNKS"
} catch {
    Write-Error "Init upload failed: $_"
    exit 1
}

# 5. Test chunk upload endpoint using curl.exe
Write-Info "Testing chunk upload endpoint..."

# Create test chunk file (1MB)
$TEST_CHUNK_FILE = "$env:TEMP\test_chunk_$(Get-Random).bin"
$testData = New-Object byte[] (1 * 1024 * 1024)
[System.IO.File]::WriteAllBytes($TEST_CHUNK_FILE, $testData)

# Upload first chunk using curl
Write-Info "Uploading chunk 0..."
try {
    $curlCmd = "curl.exe -s -X POST `"$API_V1/upload/chunk`" -H `"Authorization: Bearer $TOKEN`" -F `"upload_id=$UPLOAD_ID`" -F `"chunk_index=0`" -F `"chunk_data=@$TEST_CHUNK_FILE`""
    $chunkRespJson = Invoke-Expression $curlCmd
    $chunkResp = $chunkRespJson | ConvertFrom-Json
    Write-Info "Chunk upload response: $($chunkResp | ConvertTo-Json -Compress)"
    
    if ($chunkResp.data.chunk_index -eq 0) {
        Write-Info "Chunk 0 uploaded successfully"
    } else {
        throw "Chunk upload returned abnormal response"
    }
} catch {
    Write-Error "Chunk 0 upload failed: $_"
    Remove-Item $TEST_CHUNK_FILE -Force -ErrorAction SilentlyContinue
    # Try to cancel upload
    try { Invoke-RestMethod -Uri "$API_V1/upload/$UPLOAD_ID" -Method DELETE -Headers @{ "Authorization" = "Bearer $TOKEN" } } catch {}
    exit 1
}

# 6. Test get upload status endpoint
Write-Info "Testing get upload status endpoint..."
try {
    $statusResp = Invoke-RestMethod -Uri "$API_V1/upload/status/$UPLOAD_ID" -Method GET -Headers @{ "Authorization" = "Bearer $TOKEN" }
    Write-Info "Status response: $($statusResp | ConvertTo-Json -Compress)"
    Write-Info "Current status: $($statusResp.data.status)"
} catch {
    Write-Error "Get status failed: $_"
}

# 7. Test cancel upload endpoint
Write-Info "Testing cancel upload endpoint..."
try {
    $cancelResp = Invoke-RestMethod -Uri "$API_V1/upload/$UPLOAD_ID" -Method DELETE -Headers @{ "Authorization" = "Bearer $TOKEN" }
    Write-Info "Cancel response: $($cancelResp | ConvertTo-Json -Compress)"
} catch {
    Write-Error "Cancel upload failed: $_"
}

# 8. Test instant upload feature
Write-Info "Testing instant upload feature..."
try {
    $initBody2 = @{
        file_name = "test_panorama.jpg"
        file_size = 10485760
        file_md5 = "abc123def45678901234567890123456"
        space_id = [int]$SPACE_ID
        scene_code = "test_scene_$(Get-Random)"
        title = "Test Scene 2"
    } | ConvertTo-Json
    
    $initResp2 = Invoke-RestMethod -Uri "$API_V1/upload/init" -Method POST -Body $initBody2 -ContentType "application/json" -Headers @{ "Authorization" = "Bearer $TOKEN" }
    Write-Info "Instant upload test response: $($initResp2 | ConvertTo-Json -Compress)"
    
    if ($initResp2.data.skip_upload -eq $true) {
        Write-Info "Instant upload feature working"
    } else {
        Write-Warn "Instant upload may not be working (this is normal without actual file content)"
    }
} catch {
    Write-Error "Instant upload test failed: $_"
}

# Cleanup test file
Remove-Item $TEST_CHUNK_FILE -Force -ErrorAction SilentlyContinue

Write-Info "=== Chunk Upload Module Test Completed ==="
