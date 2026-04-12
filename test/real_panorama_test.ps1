param(
    [string]$PanoramaFile = "dujiangyan-gate.jpg"
)

$BASE_URL = "http://localhost:7000"
$API_V1 = "$BASE_URL/api/v1"
$DATA_DIR = "d:\lwdd\code\毕设\webGL-720yun\参考\都江堰"
$TEST_IMAGE = Join-Path $DATA_DIR $PanoramaFile

function LogInfo { 
    param($msg) 
    $prefix = [char]73, [char]78, [char]70, [char]79, [char]32
    Write-Host (-join $prefix + $msg) -ForegroundColor Green 
}
function LogErr { 
    param($msg) 
    $prefix = [char]69, [char]82, [char]82, [char]79, [char]82, [char]32
    Write-Host (-join $prefix + $msg) -ForegroundColor Red 
}

function Get-FileMD5 {
    param([string]$FilePath)
    $md5 = [System.Security.Cryptography.MD5]::Create()
    $stream = [System.IO.File]::OpenRead($FilePath)
    $hash = $md5.ComputeHash($stream)
    $stream.Close()
    return [System.BitConverter]::ToString($hash).Replace("-", "").ToLower()
}

LogInfo "Real Panorama Upload and Slice Test"
LogInfo "Using panorama $PanoramaFile"

if (-not (Test-Path $TEST_IMAGE)) {
    LogErr "Panorama file not found $TEST_IMAGE"
    exit 1
}

LogInfo "1. Health Check..."
try {
    $health = Invoke-RestMethod -Uri "$BASE_URL/health" -Method GET
    LogInfo "Health check passed"
} catch {
    LogErr "Health check failed"
    exit 1
}

LogInfo "2. Login..."
try {
    $loginBody = @{ username = "admin"; password = "admin123" } | ConvertTo-Json -Depth 10
    $loginResp = Invoke-RestMethod -Uri "$API_V1/auth/login" -Method POST -Body $loginBody -ContentType "application/json; charset=utf-8"
    $TOKEN = $loginResp.data.access_token
    $HEADERS = @{ "Authorization" = "Bearer $TOKEN" }
    LogInfo "Login successful"
} catch {
    $errDetail = $_.Exception.Message
    LogErr "Login failed: $errDetail"
    exit 1
}

LogInfo "3. Get/Create Space..."
try {
    $uri = "$API_V1/resource/spaces?page=1"
    $spacesResp = Invoke-RestMethod -Uri $uri -Method GET -Headers $HEADERS
    if ($spacesResp.data.spaces.Count -eq 0) {
        LogInfo "No space found, creating one..."
        $spaceBody = @{
            name = "dujiangyan"
            slug = "dujiangyan-test"
            description = "Dujiangyan panorama tour"
            province = "Sichuan"
            city = "Chengdu"
            longitude = 103.6146
            latitude = 31.0047
            zoom_level = 15
            sort_order = 1
            status = 1
        } | ConvertTo-Json
        $spaceResp = Invoke-RestMethod -Uri "$API_V1/resource/spaces" -Method POST -Body $spaceBody -ContentType "application/json" -Headers $HEADERS
        $SPACE_ID = $spaceResp.data.id
    } else {
        $SPACE_ID = $spacesResp.data.spaces[0].id
    }
    LogInfo "Using Space ID $SPACE_ID"
} catch {
    $errDetail = $_.Exception.Message
    LogErr "Failed to get/create space: $errDetail"
    exit 1
}

$fileInfo = Get-Item $TEST_IMAGE
$FILE_SIZE = $fileInfo.Length
$FILE_MD5 = Get-FileMD5 -FilePath $TEST_IMAGE
$sizeMB = [math]::Round($FILE_SIZE / 1MB, 2)
LogInfo "File size $sizeMB MB"
LogInfo "MD5 $FILE_MD5"

LogInfo "4. Init Upload..."
try {
    $initBody = @{
        filename = $PanoramaFile
        file_size = $FILE_SIZE
        file_hash = $FILE_MD5
    } | ConvertTo-Json
    
    $initResp = Invoke-RestMethod -Uri "$API_V1/upload/init" -Method POST -Body $initBody -ContentType "application/json" -Headers $HEADERS
    
    if ($initResp.data.instant -eq $true) {
        LogInfo "Instant upload - file already exists"
        $FILE_ID = $initResp.data.file_id
    } else {
        $UPLOAD_ID = $initResp.data.upload_id
        $TOTAL_CHUNKS = $initResp.data.total_chunks
        $CHUNK_SIZE = $initResp.data.chunk_size
        LogInfo "Upload initialized Chunks=$TOTAL_CHUNKS"
    }
} catch {
    LogErr "Init upload failed"
    exit 1
}

if (-not $FILE_ID) {
    LogInfo "5. Uploading chunks..."
    
    $fileBytes = [System.IO.File]::ReadAllBytes($TEST_IMAGE)
    
    for ($i = 0; $i -lt $TOTAL_CHUNKS; $i++) {
        $startIdx = $i * $CHUNK_SIZE
        $endIdx = [Math]::Min(($i + 1) * $CHUNK_SIZE, $fileBytes.Length)
        $chunkBytes = $fileBytes[$startIdx..($endIdx - 1)]
        
        $chunkFile = "$env:TEMP\chunk_$i.tmp"
        [System.IO.File]::WriteAllBytes($chunkFile, $chunkBytes)
        
        $chunkMD5 = [System.BitConverter]::ToString([System.Security.Cryptography.MD5]::Create().ComputeHash($chunkBytes)).Replace("-", "").ToLower()
        
        LogInfo "Uploading chunk $i"
        
        $curlCmd = "curl.exe -s -X POST `"$API_V1/upload/chunk`" -H `"Authorization: Bearer $TOKEN`" -F `"upload_id=$UPLOAD_ID`" -F `"chunk_index=$i`" -F `"chunk_hash=$chunkMD5`" -F `"chunk_data=@$chunkFile`""
        $chunkRespJson = Invoke-Expression $curlCmd
        
        Remove-Item $chunkFile -Force -ErrorAction SilentlyContinue
    }

    LogInfo "6. Complete Upload..."
    $completeBody = @{
        upload_id = $UPLOAD_ID
        file_hash = $FILE_MD5
    } | ConvertTo-Json
    
    $completeResp = Invoke-RestMethod -Uri "$API_V1/upload/complete" -Method POST -Body $completeBody -ContentType "application/json" -Headers $HEADERS
    
    if ($completeResp.code -eq 200) {
        $FILE_ID = $completeResp.data.file_id
        LogInfo "Upload completed File ID $FILE_ID"
    } else {
        LogErr "Complete upload failed"
        exit 1
    }
}

if ($FILE_ID) {
    LogInfo "7. Create Scene..."
    $SCENE_CODE = $PanoramaFile -replace ".jpg", ""
    try {
        $sceneBody = @{
            space_id = [int]$SPACE_ID
            title = $SCENE_CODE
            scene_code = $SCENE_CODE
            file_id = $FILE_ID
            panorama_type = "equirectangular"
            initial_fov = 90
            initial_pitch = 0
            initial_yaw = 0
        } | ConvertTo-Json
        
        $sceneResp = Invoke-RestMethod -Uri "$API_V1/resource/scenes" -Method POST -Body $sceneBody -ContentType "application/json" -Headers $HEADERS
        
        if ($sceneResp.code -eq 200) {
            $SCENE_ID = $sceneResp.data.scene_id
            $TASK_ID = $sceneResp.data.task_id
            $SLICE_STATUS = $sceneResp.data.slice_status
            LogInfo "Scene created ID=$SCENE_ID TaskID=$TASK_ID Status=$SLICE_STATUS"
        } else {
            LogErr "Create scene failed"
        }
    } catch {
        LogErr "Create scene failed"
    }

    if ($SCENE_ID) {
        LogInfo "8. Waiting for slice task..."
        for ($i = 0; $i -lt 30; $i++) {
            Start-Sleep -Seconds 1
            $sceneDetailResp = Invoke-RestMethod -Uri "$API_V1/resource/scenes/$SCENE_ID" -Method GET -Headers $HEADERS
            $status = $sceneDetailResp.data.slice_status
            Write-Host "  Status: $status"
            if ($status -eq "ready") {
                break
            } elseif ($status -eq "failed") {
                break
            }
        }

        LogInfo "9. Final Scene Details..."
        $finalResp = Invoke-RestMethod -Uri "$API_V1/resource/scenes/$SCENE_ID" -Method GET -Headers $HEADERS
        LogInfo "Scene ID $($finalResp.data.id)"
        LogInfo "Slice Status $($finalResp.data.slice_status)"
        LogInfo "Tile URL $($finalResp.data.tile_url)"
    }
}

LogInfo "Test Completed"
