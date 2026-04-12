$BASE_URL = "http://localhost:7000"
$API_V1 = "$BASE_URL/api/v1"
$TEST_DIR = "d:\lwdd\code\毕设\webGL-720yun\test"
$DATA_DIR = "d:\lwdd\code\毕设\webGL-720yun\参考\都江堰"
$PANORAMA_FILE = "都江堰景区大门_3576761_sphere.jpg"

Write-Host "=== Upload and Slice Module Test ===" -ForegroundColor Green

Write-Host "`n1. Login..." -ForegroundColor Cyan
$loginResp = curl.exe -s "$API_V1/auth/login" -X POST -H "Content-Type: application/json" -d "@$TEST_DIR\login.json"
$TOKEN = ($loginResp | Select-String -Pattern '"access_token":"([^"]+)"' -AllMatches).Matches.Groups[1].Value
if (-not $TOKEN) {
    Write-Host "Login failed" -ForegroundColor Red
    exit 1
}
Write-Host "Login successful" -ForegroundColor Green

Write-Host "`n2. Get Spaces..." -ForegroundColor Cyan
$spacesResp = curl.exe -s "$API_V1/resource/spaces?page=1" -H "Authorization: Bearer $TOKEN"
$spaceMatches = $spacesResp | Select-String -Pattern '"id":(\d+)' -AllMatches
if ($spaceMatches.Matches.Count -gt 0) {
    $SPACE_ID = $spaceMatches.Matches[0].Groups[1].Value
} else {
    Write-Host "No spaces found, creating one..." -ForegroundColor Yellow
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
    } | ConvertTo-Json -Depth 10
    $spaceBody | Out-File -FilePath "$TEST_DIR\space_req.json" -Encoding utf8
    $spaceResp = curl.exe -s "$API_V1/resource/spaces" -X POST -H "Content-Type: application/json; charset=utf-8" -H "Authorization: Bearer $TOKEN" -d "@$TEST_DIR\space_req.json"
    $spaceMatch = $spaceResp | Select-String -Pattern '"id":(\d+)' -AllMatches
    if ($spaceMatch.Matches.Count -gt 0) {
        $SPACE_ID = $spaceMatch.Matches[0].Groups[1].Value
    }
}
Write-Host "Using Space ID: $SPACE_ID" -ForegroundColor Green

Write-Host "`n3. Calculate file MD5..." -ForegroundColor Cyan
$testFile = Join-Path $DATA_DIR $PANORAMA_FILE
$fileInfo = Get-Item $testFile
$FILE_SIZE = $fileInfo.Length
$md5 = [System.Security.Cryptography.MD5]::Create()
$stream = [System.IO.File]::OpenRead($testFile)
$hash = $md5.ComputeHash($stream)
$stream.Close()
$FILE_MD5 = [System.BitConverter]::ToString($hash).Replace("-", "").ToLower()
$sizeMB = [math]::Round($FILE_SIZE / 1MB, 2)
Write-Host "File: $PANORAMA_FILE, Size: $sizeMB MB, MD5: $FILE_MD5" -ForegroundColor Green

Write-Host "`n4. Init Upload..." -ForegroundColor Cyan
$initBody = @{
    filename = $PANORAMA_FILE
    file_size = $FILE_SIZE
    file_hash = $FILE_MD5
} | ConvertTo-Json -Depth 10
$initBody | Out-File -FilePath "$TEST_DIR\init_req.json" -Encoding utf8
$initResp = curl.exe -s "$API_V1/upload/init" -X POST -H "Content-Type: application/json; charset=utf-8" -H "Authorization: Bearer $TOKEN" -d "@$TEST_DIR\init_req.json"
Write-Host "Init response: $initResp" -ForegroundColor Gray

$instantMatch = $initResp | Select-String -Pattern '"instant":(true|false)' -AllMatches
if ($instantMatch.Matches.Count -gt 0) {
    $isInstant = $instantMatch.Matches[0].Groups[1].Value
} else {
    $isInstant = "false"
}

if ($isInstant -eq "true") {
    Write-Host "Instant upload - file already exists!" -ForegroundColor Yellow
    $fileMatch = $initResp | Select-String -Pattern '"file_id":"([^"]+)"' -AllMatches
    if ($fileMatch.Matches.Count -gt 0) {
        $FILE_ID = $fileMatch.Matches[0].Groups[1].Value
    }
} else {
    $uploadMatch = $initResp | Select-String -Pattern '"upload_id":"([^"]+)"' -AllMatches
    $chunksMatch = $initResp | Select-String -Pattern '"total_chunks":(\d+)' -AllMatches
    $chunkSizeMatch = $initResp | Select-String -Pattern '"chunk_size":(\d+)' -AllMatches
    
    if ($uploadMatch.Matches.Count -gt 0) {
        $UPLOAD_ID = $uploadMatch.Matches[0].Groups[1].Value
    }
    if ($chunksMatch.Matches.Count -gt 0) {
        $TOTAL_CHUNKS = [int]$chunksMatch.Matches[0].Groups[1].Value
    }
    if ($chunkSizeMatch.Matches.Count -gt 0) {
        $CHUNK_SIZE = [int]$chunkSizeMatch.Matches[0].Groups[1].Value
    }
    
    Write-Host "Upload ID: $UPLOAD_ID, Chunks: $TOTAL_CHUNKS" -ForegroundColor Green

    Write-Host "`n5. Uploading chunks..." -ForegroundColor Cyan
    $fileBytes = [System.IO.File]::ReadAllBytes($testFile)
    for ($i = 0; $i -lt $TOTAL_CHUNKS; $i++) {
        $startIdx = $i * $CHUNK_SIZE
        $endIdx = [Math]::Min(($i + 1) * $CHUNK_SIZE, $fileBytes.Length)
        $chunkBytes = $fileBytes[$startIdx..($endIdx - 1)]
        
        $chunkFile = "$env:TEMP\chunk_$i.tmp"
        [System.IO.File]::WriteAllBytes($chunkFile, $chunkBytes)
        
        $chunkMD5 = [System.BitConverter]::ToString([System.Security.Cryptography.MD5]::Create().ComputeHash($chunkBytes)).Replace("-", "").ToLower()
        
        Write-Host "  Uploading chunk $i..." -ForegroundColor Gray
        curl.exe -s "$API_V1/upload/chunk" -X POST -H "Authorization: Bearer $TOKEN" -F "upload_id=$UPLOAD_ID" -F "chunk_index=$i" -F "chunk_hash=$chunkMD5" -F "chunk_data=@$chunkFile" | Out-Null
        
        Remove-Item $chunkFile -Force -ErrorAction SilentlyContinue
    }

    Write-Host "`n6. Complete Upload..." -ForegroundColor Cyan
    $completeBody = @{
        upload_id = $UPLOAD_ID
        file_hash = $FILE_MD5
    } | ConvertTo-Json -Depth 10
    $completeBody | Out-File -FilePath "$TEST_DIR\complete_req.json" -Encoding utf8
    $completeResp = curl.exe -s "$API_V1/upload/complete" -X POST -H "Content-Type: application/json; charset=utf-8" -H "Authorization: Bearer $TOKEN" -d "@$TEST_DIR\complete_req.json"
    Write-Host "Complete response: $completeResp" -ForegroundColor Gray
    
    $fileMatch = $completeResp | Select-String -Pattern '"file_id":"([^"]+)"' -AllMatches
    if ($fileMatch.Matches.Count -gt 0) {
        $FILE_ID = $fileMatch.Matches[0].Groups[1].Value
        Write-Host "Upload completed, File ID: $FILE_ID" -ForegroundColor Green
    } else {
        Write-Host "Complete upload failed" -ForegroundColor Red
    }
}

if ($FILE_ID) {
    Write-Host "`n7. Create Scene..." -ForegroundColor Cyan
    $SCENE_CODE = $PANORAMA_FILE -replace "_sphere.jpg", ""
    $sceneBody = @{
        space_id = [int]$SPACE_ID
        title = $SCENE_CODE
        scene_code = $SCENE_CODE
        file_id = $FILE_ID
        panorama_type = "equirectangular"
        initial_fov = 90
        initial_pitch = 0
        initial_yaw = 0
    } | ConvertTo-Json -Depth 10
    $sceneBody | Out-File -FilePath "$TEST_DIR\scene_req.json" -Encoding utf8
    $sceneResp = curl.exe -s "$API_V1/resource/scenes" -X POST -H "Content-Type: application/json; charset=utf-8" -H "Authorization: Bearer $TOKEN" -d "@$TEST_DIR\scene_req.json"
    Write-Host "Scene response: $sceneResp" -ForegroundColor Gray
    
    $sceneMatch = $sceneResp | Select-String -Pattern '"scene_id":(\d+)' -AllMatches
    $taskMatch = $sceneResp | Select-String -Pattern '"task_id":"([^"]+)"' -AllMatches
    
    if ($sceneMatch.Matches.Count -gt 0) {
        $SCENE_ID = $sceneMatch.Matches[0].Groups[1].Value
    }
    if ($taskMatch.Matches.Count -gt 0) {
        $TASK_ID = $taskMatch.Matches[0].Groups[1].Value
    }
    
    Write-Host "Scene created: ID=$SCENE_ID, TaskID=$TASK_ID" -ForegroundColor Green

    if ($SCENE_ID) {
        Write-Host "`n8. Waiting for slice task..." -ForegroundColor Cyan
        for ($i = 0; $i -lt 30; $i++) {
            Start-Sleep -Seconds 1
            $detailResp = curl.exe -s "$API_V1/resource/scenes/$SCENE_ID" -H "Authorization: Bearer $TOKEN"
            $statusMatch = $detailResp | Select-String -Pattern '"slice_status":"([^"]+)"' -AllMatches
            if ($statusMatch.Matches.Count -gt 0) {
                $status = $statusMatch.Matches[0].Groups[1].Value
                Write-Host "  Status: $status"
                if ($status -eq "ready" -or $status -eq "failed") {
                    break
                }
            }
        }

        Write-Host "`n9. Final Scene Details..." -ForegroundColor Cyan
        $detailResp = curl.exe -s "$API_V1/resource/scenes/$SCENE_ID" -H "Authorization: Bearer $TOKEN"
        $finalStatusMatch = $detailResp | Select-String -Pattern '"slice_status":"([^"]+)"' -AllMatches
        $tileUrlMatch = $detailResp | Select-String -Pattern '"tile_url":"([^"]+)"' -AllMatches
        
        if ($finalStatusMatch.Matches.Count -gt 0) {
            Write-Host "Slice Status: $($finalStatusMatch.Matches[0].Groups[1].Value)" -ForegroundColor Green
        }
        if ($tileUrlMatch.Matches.Count -gt 0) {
            Write-Host "Tile URL: $($tileUrlMatch.Matches[0].Groups[1].Value)" -ForegroundColor Green
        }
    }
}

Write-Host "`n=== Test Completed ===" -ForegroundColor Green
