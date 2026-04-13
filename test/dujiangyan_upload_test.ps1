$BASE_URL = "http://localhost:7000"
$API_V1 = "$BASE_URL/api/v1"
$TEST_DATA_DIR = "d:\lwdd\code\毕设\webGL-720yun\参考\都江堰"

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

Write-Info "=== Dujiangyan Panorama Upload Test ==="
Write-Info "Test data directory: $TEST_DATA_DIR"

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

Write-Info "3. Create Dujiangyan Space..."
$SPACE_NAME = "Dujiangyan"
try {
    $spacesResp = Invoke-RestMethod -Uri "$API_V1/resource/spaces?page=1&page_size=100" -Method GET -Headers $HEADERS
    $existingSpace = $spacesResp.data.spaces | Where-Object { $_.name -eq $SPACE_NAME } | Select-Object -First 1
    
    if ($existingSpace) {
        $SPACE_ID = $existingSpace.id
        Write-Info "Space already exists, ID: $SPACE_ID"
    } else {
        $spaceBody = @{
            name = $SPACE_NAME
            slug = "dujiangyan-test"
            description = "Dujiangyan Test Space"
            province = "Sichuan"
            city = "Chengdu"
            longitude = 103.611379
            latitude = 31.001752
            zoom_level = 15
            sort_order = 1
            status = 1
        } | ConvertTo-Json
        
        $spaceResp = Invoke-RestMethod -Uri "$API_V1/resource/spaces" -Method POST -Body $spaceBody -ContentType "application/json" -Headers $HEADERS
        $SPACE_ID = $spaceResp.data.id
        Write-Info "Created space, ID: $SPACE_ID"
    }
} catch {
    Write-Err "Failed to create space: $_"
    exit 1
}

$panoramaFiles = Get-ChildItem -Path $TEST_DATA_DIR -Filter "*_sphere.jpg" | Select-Object -First 3
Write-Info "Found $($panoramaFiles.Count) panorama files to test"

$uploadedScenes = @()

foreach ($file in $panoramaFiles) {
    Write-Info ""
    Write-Info "=== Processing: $($file.Name) ==="
    
    $FILE_SIZE = $file.Length
    $FILE_MD5 = Get-FileMD5 -FilePath $file.FullName
    $SCENE_CODE = $file.BaseName -replace ".*_(\d+)_sphere", '$1'
    $SCENE_TITLE = $file.BaseName -replace "_\d+_sphere", ''
    
    Write-Info "File: $($file.Name)"
    Write-Info "Size: $([math]::Round($FILE_SIZE / 1MB, 2)) MB"
    Write-Info "MD5: $FILE_MD5"
    Write-Info "Scene Code: $SCENE_CODE"
    
    Write-Info "Init upload..."
    try {
        $initBody = @{
            space_id = [int]$SPACE_ID
            filename = $file.Name
            file_size = $FILE_SIZE
            file_hash = $FILE_MD5
        } | ConvertTo-Json
        
        $initResp = Invoke-RestMethod -Uri "$API_V1/upload/init" -Method POST -Body $initBody -ContentType "application/json" -Headers $HEADERS
        
        if ($initResp.data.instant -eq $true) {
            Write-Info "Instant upload (file already exists)"
            $FILE_ID = $initResp.data.file_id
        } else {
            $UPLOAD_ID = $initResp.data.upload_id
            $TOTAL_CHUNKS = $initResp.data.total_chunks
            $CHUNK_SIZE = $initResp.data.chunk_size
            Write-Info "Upload initialized: ID=$UPLOAD_ID, Chunks=$TOTAL_CHUNKS"
            
            Write-Info "Uploading chunks..."
            $fileBytes = [System.IO.File]::ReadAllBytes($file.FullName)
            
            for ($i = 0; $i -lt $TOTAL_CHUNKS; $i++) {
                $startIdx = $i * $CHUNK_SIZE
                $endIdx = [Math]::Min(($i + 1) * $CHUNK_SIZE, $fileBytes.Length)
                $chunkBytes = $fileBytes[$startIdx..($endIdx - 1)]
                
                $chunkFile = "$env:TEMP\chunk_$i.tmp"
                [System.IO.File]::WriteAllBytes($chunkFile, $chunkBytes)
                
                $chunkMD5 = [System.BitConverter]::ToString([System.Security.Cryptography.MD5]::Create().ComputeHash($chunkBytes)).Replace("-", "").ToLower()
                
                $progress = [math]::Round(($i + 1) / $TOTAL_CHUNKS * 100, 1)
                Write-Host "  Chunk $($i + 1)/$TOTAL_CHUNKS ($progress %)"
                
                $curlCmd = "curl.exe -s -X POST `"$API_V1/upload/chunk`" -H `"Authorization: Bearer $TOKEN`" -F `"upload_id=$UPLOAD_ID`" -F `"chunk_index=$i`" -F `"chunk_hash=$chunkMD5`" -F `"chunk_data=@$chunkFile`""
                $chunkRespJson = Invoke-Expression $curlCmd
                
                Remove-Item $chunkFile -Force -ErrorAction SilentlyContinue
            }
            
            Write-Info "Completing upload..."
            $completeBody = @{
                upload_id = $UPLOAD_ID
                file_hash = $FILE_MD5
            } | ConvertTo-Json
            
            $completeResp = Invoke-RestMethod -Uri "$API_V1/upload/complete" -Method POST -Body $completeBody -ContentType "application/json" -Headers $HEADERS
            $FILE_ID = $completeResp.data.file_id
            Write-Info "Upload completed, File ID: $FILE_ID"
        }
    } catch {
        Write-Err "Upload failed: $_"
        continue
    }
    
    Write-Info "Creating scene..."
    try {
        $sceneBody = @{
            space_id = [int]$SPACE_ID
            title = $SCENE_TITLE
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
            Write-Info "Scene created: ID=$SCENE_ID, TaskID=$TASK_ID"
            $uploadedScenes += @{
                File = $file.Name
                SceneID = $SCENE_ID
                TaskID = $TASK_ID
            }
        } else {
            Write-Err "Create scene failed: $($sceneResp.message)"
        }
    } catch {
        Write-Err "Create scene failed: $_"
    }
}

Write-Info ""
Write-Info "=== Waiting for slice tasks to complete ==="
Start-Sleep -Seconds 5

Write-Info ""
Write-Info "=== Checking slice status ==="
foreach ($scene in $uploadedScenes) {
    try {
        $sceneDetail = Invoke-RestMethod -Uri "$API_V1/resource/scenes/$($scene.SceneID)" -Method GET -Headers $HEADERS
        Write-Info "$($scene.File): slice_status=$($sceneDetail.data.slice_status)"
    } catch {
        Write-Err "Failed to get scene $($scene.SceneID): $_"
    }
}

Write-Info ""
Write-Info "=== Test Summary ==="
Write-Info "Space ID: $SPACE_ID"
Write-Info "Uploaded scenes: $($uploadedScenes.Count)"
foreach ($scene in $uploadedScenes) {
    Write-Info "  - $($scene.File) -> Scene ID: $($scene.SceneID)"
}

Write-Info ""
Write-Info "=== Checking MinIO Storage ==="
