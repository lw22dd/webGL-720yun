$BASE_URL = "http://localhost:7000"
$API_V1 = "$BASE_URL/api/v1"
$TEST_DATA_DIR = "$PSScriptRoot\test_data"
$TEST_IMAGE = "$TEST_DATA_DIR\test_panorama_2k.jpg"

function Write-Info { param($Message) Write-Host "[INFO] $Message" -ForegroundColor Green }
function Write-Err { param($Message) Write-Host "[ERROR] $Message" -ForegroundColor Red }
function Write-Warn { param($Message) Write-Host "[WARN] $Message" -ForegroundColor Yellow }

function Get-FileMD5 {
    param([string]$FilePath)
    $md5 = [System.Security.Cryptography.MD5]::Create()
    $stream = [System.IO.File]::OpenRead($FilePath)
    $hash = $md5.ComputeHash($stream)
    $stream.Close()
    return [System.BitConverter]::ToString($hash).Replace("-", "").ToLower()
}

Write-Info "=== Upload & Slice Module Test ==="

if (-not (Test-Path $TEST_IMAGE)) {
    Write-Err "Test image not found: $TEST_IMAGE"
    Write-Info "Run create_test_panorama.ps1 first to generate test images"
    exit 1
}

Write-Info "1. Health Check..."
try {
    $health = Invoke-RestMethod -Uri "$BASE_URL/health" -Method GET
    Write-Info "Health check passed: $($health.status)"
} catch {
    Write-Err "Health check failed: $_"
    Write-Err "Please ensure the backend server is running on $BASE_URL"
    exit 1
}

Write-Info "2. Login..."
try {
    $loginBody = @{ username = "admin"; password = "admin123" } | ConvertTo-Json
    $loginResp = Invoke-RestMethod -Uri "$API_V1/auth/login" -Method POST -Body $loginBody -ContentType "application/json"
    $TOKEN = $loginResp.data.access_token
    $HEADERS = @{ "Authorization" = "Bearer $TOKEN" }
    Write-Info "Login successful"
} catch {
    Write-Err "Login failed: $_"
    exit 1
}

Write-Info "3. Create Test Space..."
try {
    $spaceName = "Network Test Space $(Get-Random)"
    $spaceSlug = "net-test-$(Get-Random)"
    $spaceBody = @{
        name = $spaceName
        slug = $spaceSlug
        description = "Temporary space for network test"
        province = "Beijing"
        city = "Beijing"
        longitude = 116.407429
        latitude = 39.904211
        zoom_level = 12
        sort_order = 1
        status = 1
    } | ConvertTo-Json
    $spaceResp = Invoke-RestMethod -Uri "$API_V1/resource/spaces" -Method POST -Body $spaceBody -ContentType "application/json" -Headers $HEADERS
    $SPACE_ID = $spaceResp.data.id
    $SPACE_SLUG = $spaceResp.data.slug
    Write-Info "Created Space ID: $SPACE_ID, Slug: $SPACE_SLUG"
} catch {
    Write-Err "Failed to create space: $_"
    exit 1
}

$fileInfo = Get-Item $TEST_IMAGE
$FILE_SIZE = $fileInfo.Length
$FILE_MD5 = Get-FileMD5 -FilePath $TEST_IMAGE
Write-Info "Test file: $TEST_IMAGE"
Write-Info "File size: $FILE_SIZE bytes"
Write-Info "File MD5: $FILE_MD5"

Write-Info "4. Test Init Upload..."
try {
    $initBody = @{
        space_id = [int]$SPACE_ID
        filename = "test_panorama.jpg"
        file_size = $FILE_SIZE
        file_hash = $FILE_MD5
    } | ConvertTo-Json
    
    $initResp = Invoke-RestMethod -Uri "$API_V1/upload/init" -Method POST -Body $initBody -ContentType "application/json" -Headers $HEADERS
    Write-Info "Init response: $($initResp | ConvertTo-Json -Compress)"
    
    if ($initResp.data.instant -eq $true) {
        Write-Info "Instant upload (file already exists)"
        $FILE_ID = $initResp.data.file_id
    } else {
        $UPLOAD_ID = $initResp.data.upload_id
        $TOTAL_CHUNKS = $initResp.data.total_chunks
        $CHUNK_SIZE = $initResp.data.chunk_size
        Write-Info "Upload initialized: ID=$UPLOAD_ID, Chunks=$TOTAL_CHUNKS, ChunkSize=$CHUNK_SIZE"
    }
} catch {
    Write-Err "Init upload failed: $_"
    exit 1
}

if (-not $FILE_ID) {
    Write-Info "5. Test Chunk Upload (splitting file into $TOTAL_CHUNKS chunks)..."
    
    $fileBytes = [System.IO.File]::ReadAllBytes($TEST_IMAGE)
    $chunkSizeBytes = $CHUNK_SIZE
    
    for ($i = 0; $i -lt $TOTAL_CHUNKS; $i++) {
        $startIdx = $i * $chunkSizeBytes
        $endIdx = [Math]::Min(($i + 1) * $chunkSizeBytes, $fileBytes.Length)
        $chunkBytes = $fileBytes[$startIdx..($endIdx - 1)]
        
        $chunkFile = "$env:TEMP\chunk_$i.tmp"
        [System.IO.File]::WriteAllBytes($chunkFile, $chunkBytes)
        
        $chunkMD5 = [System.BitConverter]::ToString([System.Security.Cryptography.MD5]::Create().ComputeHash($chunkBytes)).Replace("-", "").ToLower()
        
        Write-Info "Uploading chunk $i/$($TOTAL_CHUNKS - 1) ($($chunkBytes.Length) bytes)..."
        
        try {
            $curlCmd = "curl.exe -s -X POST `"$API_V1/upload/chunk`" -H `"Authorization: Bearer $TOKEN`" -F `"upload_id=$UPLOAD_ID`" -F `"chunk_index=$i`" -F `"chunk_hash=$chunkMD5`" -F `"chunk_data=@$chunkFile`""
            $chunkRespJson = Invoke-Expression $curlCmd
            $chunkResp = $chunkRespJson | ConvertFrom-Json
            
            if ($chunkResp.code -eq 200) {
                Write-Info "  Chunk $i uploaded successfully"
            } else {
                Write-Err "  Chunk $i upload failed: $($chunkResp.message)"
            }
        } catch {
            Write-Err "  Chunk $i upload failed: $_"
        }
        
        Remove-Item $chunkFile -Force -ErrorAction SilentlyContinue
    }

    Write-Info "6. Test Upload Status..."
    try {
        $statusResp = Invoke-RestMethod -Uri "$API_V1/upload/status/$UPLOAD_ID" -Method GET -Headers $HEADERS
        Write-Info "Status: $($statusResp.data.status), Progress: $($statusResp.data.percentage)%, Chunks: $($statusResp.data.uploaded_chunks -join ',')"
    } catch {
        Write-Err "Get status failed: $_"
    }

    Write-Info "7. Test Complete Upload..."
    try {
        $completeBody = @{
            upload_id = $UPLOAD_ID
            file_hash = $FILE_MD5
        } | ConvertTo-Json
        
        $completeResp = Invoke-RestMethod -Uri "$API_V1/upload/complete" -Method POST -Body $completeBody -ContentType "application/json" -Headers $HEADERS
        Write-Info "Complete response: $($completeResp | ConvertTo-Json -Compress)"
        
        if ($completeResp.code -eq 200) {
            $FILE_ID = $completeResp.data.file_id
            Write-Info "Upload completed, File ID: $FILE_ID"
        } else {
            Write-Err "Complete upload failed: $($completeResp.message)"
        }
    } catch {
        Write-Err "Complete upload failed: $_"
    }
}

if ($FILE_ID) {
    Write-Info "8. Test Get File Info..."
    try {
        $fileInfoResp = Invoke-RestMethod -Uri "$API_V1/upload/file/$FILE_ID" -Method GET -Headers $HEADERS
        Write-Info "File info: $($fileInfoResp | ConvertTo-Json -Compress)"
    } catch {
        Write-Err "Get file info failed: $_"
    }

    Write-Info "9. Create Scene (triggers slice task)..."
    $SCENE_CODE = "test_scene_$(Get-Random)"
    try {
        $sceneBody = @{
            space_id = [int]$SPACE_ID
            title = "Test Panorama Scene"
            scene_code = $SCENE_CODE
            file_id = $FILE_ID
            panorama_type = "equirectangular"
            initial_fov = 90
            initial_pitch = 0
            initial_yaw = 0
        } | ConvertTo-Json
        
        $sceneResp = Invoke-RestMethod -Uri "$API_V1/resource/scenes" -Method POST -Body $sceneBody -ContentType "application/json" -Headers $HEADERS
        Write-Info "Scene response: $($sceneResp | ConvertTo-Json -Compress)"
        
        if ($sceneResp.code -eq 200) {
            $SCENE_ID = $sceneResp.data.scene_id
            $TASK_ID = $sceneResp.data.task_id
            $SLICE_STATUS = $sceneResp.data.slice_status
            Write-Info "Scene created: ID=$SCENE_ID, TaskID=$TASK_ID, SliceStatus=$SLICE_STATUS"
        } else {
            Write-Err "Create scene failed: $($sceneResp.message)"
        }
    } catch {
        Write-Err "Create scene failed: $_"
    }

    if ($SCENE_ID) {
        Write-Info "10. Check Scene Slice Status..."
        Start-Sleep -Seconds 2
        try {
            $sceneDetailResp = Invoke-RestMethod -Uri "$API_V1/resource/scenes/$SCENE_ID" -Method GET -Headers $HEADERS
            Write-Info "Scene detail: slice_status=$($sceneDetailResp.data.slice_status)"
        } catch {
            Write-Err "Get scene detail failed: $_"
        }

        Write-Info "11. List Scenes..."
        try {
            $listResp = Invoke-RestMethod -Uri "$API_V1/resource/scenes?space_id=$SPACE_ID&page=1&page_size=10" -Method GET -Headers $HEADERS
            Write-Info "Scenes count: $($listResp.data.page_info.total)"
        } catch {
            Write-Err "List scenes failed: $_"
        }
    }
}

Write-Info ""
Write-Info "=== Test Completed ==="
Write-Info "Summary:"
Write-Info "  - Space ID: $SPACE_ID"
Write-Info "  - File ID: $FILE_ID"
Write-Info "  - Scene ID: $SCENE_ID"
Write-Info "  - Task ID: $TASK_ID"
