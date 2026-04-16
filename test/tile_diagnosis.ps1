# ============================================================
# Tile API Diagnosis Script
# Diagnose tile loading 404 issues
# ============================================================

param(
    [string]$SceneCode = "doujiangyanhoumen-yejing-616591bb",
    [string]$BaseUrl = "http://localhost:7000"
)

$API_V1 = "$BaseUrl/api/v1"

function Write-Info { param($Message) Write-Host "[INFO] $Message" -ForegroundColor Green }
function Write-Err { param($Message) Write-Host "[ERROR] $Message" -ForegroundColor Red }
function Write-Warn { param($Message) Write-Host "[WARN] $Message" -ForegroundColor Yellow }
function Write-Header { param($Message) Write-Host "`n========== $Message ==========" -ForegroundColor Cyan }
function Write-Debug { param($Message) Write-Host "[DEBUG] $Message" -ForegroundColor Gray }

Write-Header "Tile API Diagnosis - Scene: $SceneCode"

# 1. Backend Health Check
Write-Header "1. Backend Health Check"
try {
    $health = Invoke-RestMethod -Uri "$BaseUrl/health" -Method GET -TimeoutSec 5
    Write-Info "Backend service is running"
    Write-Info "   Status: $($health.status)"
    Write-Info "   Time: $($health.time)"
} catch {
    Write-Err "Backend service unavailable: $_"
    Write-Err "   Please ensure backend is running at $BaseUrl"
    exit 1
}

# 2. Test HEAD request (frontend preload method)
Write-Header "2. Test HEAD Request (Frontend Preload)"
$headTestUrl = "$API_V1/res/tiles/$SceneCode/px/1/0/0"
Write-Info "Test URL: $headTestUrl"
try {
    $headResp = Invoke-WebRequest -UseBasicParsing -Uri $headTestUrl -Method HEAD -TimeoutSec 10
    Write-Info "HEAD request successful"
    Write-Info "   Status: $($headResp.StatusCode)"
    Write-Info "   Content-Type: $($headResp.Headers['Content-Type'])"
    Write-Info "   Content-Length: $($headResp.Headers['Content-Length'])"
} catch {
    $statusCode = $_.Exception.Response.StatusCode.value__
    Write-Err "HEAD request failed"
    Write-Err "   Status: $statusCode"
    Write-Err "   Error: $_"
    
    Write-Warn "   Trying GET request for detailed error..."
    try {
        $getResp = Invoke-WebRequest -UseBasicParsing -Uri $headTestUrl -Method GET -TimeoutSec 10
        Write-Info "   GET Response: $($getResp.Content)"
    } catch {
        $getError = $_.ErrorDetails.Message
        if ($getError) {
            Write-Err "   GET Error Details: $getError"
        }
    }
}

# 3. Test all Level 0 tiles
Write-Header "3. Test Level 0 Tiles (2x2=4 tiles per face)"
$faces = @("px", "nx", "py", "ny", "pz", "nz")
$level0Stats = @{
    Total = 0
    Success = 0
    Failed = 0
    Results = @()
}

foreach ($face in $faces) {
    for ($y = 0; $y -lt 2; $y++) {
        for ($x = 0; $x -lt 2; $x++) {
            $level0Stats.Total++
            $url = "$API_V1/res/tiles/$SceneCode/$face/0/$y/$x"
            try {
                $resp = Invoke-WebRequest -UseBasicParsing -Uri $url -Method HEAD -TimeoutSec 5
                $level0Stats.Success++
                $level0Stats.Results += "OK $face/0/$y/$x"
            } catch {
                $level0Stats.Failed++
                $statusCode = $_.Exception.Response.StatusCode.value__
                $level0Stats.Results += "FAIL $face/0/$y/$x ($statusCode)"
            }
        }
    }
}

Write-Info "Level 0 Statistics:"
Write-Info "   Total: $($level0Stats.Total)"
Write-Info "   Success: $($level0Stats.Success)"
Write-Info "   Failed: $($level0Stats.Failed)"
Write-Info "   Success Rate: $([math]::Round($level0Stats.Success / $level0Stats.Total * 100, 1))%"

if ($level0Stats.Failed -gt 0) {
    Write-Warn "Failed tiles:"
    $level0Stats.Results | Where-Object { $_ -like "*FAIL*" } | ForEach-Object { Write-Warn "   $_" }
}

# 4. Test Level 1 tiles (the level with errors)
Write-Header "4. Test Level 1 Tiles (4x4=16 tiles per face)"
$level1Stats = @{
    Total = 0
    Success = 0
    Failed = 0
    SampleErrors = @()
}

foreach ($face in $faces) {
    for ($y = 0; $y -lt 4; $y++) {
        for ($x = 0; $x -lt 4; $x++) {
            $level1Stats.Total++
            $url = "$API_V1/res/tiles/$SceneCode/$face/1/$y/$x"
            try {
                $resp = Invoke-WebRequest -UseBasicParsing -Uri $url -Method HEAD -TimeoutSec 5
                $level1Stats.Success++
            } catch {
                $level1Stats.Failed++
                if ($level1Stats.SampleErrors.Count -lt 5) {
                    $statusCode = $_.Exception.Response.StatusCode.value__
                    $level1Stats.SampleErrors += "$face/1/$y/$x ($statusCode)"
                }
            }
        }
    }
}

Write-Info "Level 1 Statistics:"
Write-Info "   Total: $($level1Stats.Total)"
Write-Info "   Success: $($level1Stats.Success)"
Write-Info "   Failed: $($level1Stats.Failed)"
Write-Info "   Success Rate: $([math]::Round($level1Stats.Success / $level1Stats.Total * 100, 1))%"

if ($level1Stats.Failed -gt 0) {
    Write-Warn "Failed samples (first 5):"
    $level1Stats.SampleErrors | ForEach-Object { Write-Warn "   $_" }
}

# 5. Test preview image
Write-Header "5. Test Preview Image"
$previewUrl = "$API_V1/res/previews/$SceneCode"
Write-Info "URL: $previewUrl"
try {
    $previewResp = Invoke-WebRequest -UseBasicParsing -Uri $previewUrl -Method HEAD -TimeoutSec 5
    Write-Info "Preview image exists"
    Write-Info "   Content-Length: $($previewResp.Headers['Content-Length']) bytes"
} catch {
    $statusCode = $_.Exception.Response.StatusCode.value__
    Write-Warn "Preview image not found (Status: $statusCode)"
}

# 6. Test complete cubemap faces
Write-Header "6. Test Complete Cubemap Faces"
$cubemapStats = @{
    Success = 0
    Failed = 0
    Results = @()
}

foreach ($face in $faces) {
    $url = "$API_V1/res/cubemap/$SceneCode/$face"
    try {
        $resp = Invoke-WebRequest -UseBasicParsing -Uri $url -Method HEAD -TimeoutSec 5
        $cubemapStats.Success++
        $cubemapStats.Results += "OK $face"
    } catch {
        $cubemapStats.Failed++
        $statusCode = $_.Exception.Response.StatusCode.value__
        $cubemapStats.Results += "FAIL $face ($statusCode)"
    }
}

Write-Info "Cubemap Face Statistics:"
Write-Info "   Success: $($cubemapStats.Success) / 6"
Write-Info "   Results: $($cubemapStats.Results -join ' | ')"

# 7. Diagnosis Suggestions
Write-Header "7. Diagnosis Suggestions"

if ($level0Stats.Failed -eq $level0Stats.Total -and $level1Stats.Failed -eq $level1Stats.Total) {
    Write-Err "All tiles are inaccessible. Possible causes:"
    Write-Err "   1. Scene '$SceneCode' does not exist or has been deleted"
    Write-Err "   2. Scene slicing task is incomplete or failed"
    Write-Err "   3. No tile data in MinIO storage"
    Write-Err ""
    Write-Warn "Recommended actions:"
    Write-Warn "   1. Check database: SELECT * FROM res_scenes WHERE scene_code = '$SceneCode'"
    Write-Warn "   2. Check slice_status field (should be 'completed')"
    Write-Warn "   3. Check MinIO for tile files"
} elseif ($level0Stats.Success -gt 0 -and $level1Stats.Failed -gt 0) {
    Write-Warn "Level 0 partially successful, but Level 1 has many failures. Possible causes:"
    Write-Warn "   1. Slicing task incomplete (only Level 0 done)"
    Write-Warn "   2. Errors during slicing, some levels not generated"
    Write-Warn ""
    Write-Warn "Recommended actions:"
    Write-Warn "   1. Check slicing task status and logs"
    Write-Warn "   2. Re-trigger slicing task"
} elseif ($level0Stats.Success -eq $level0Stats.Total -and $level1Stats.Success -eq $level1Stats.Total) {
    Write-Info "All tiles are accessible"
    Write-Info "   If frontend still shows 404, check:"
    Write-Info "   1. Frontend request URL is correct"
    Write-Info "   2. CORS issues"
    Write-Info "   3. Browser cache issues"
}

# 8. MinIO Path Reference
Write-Header "8. MinIO Storage Path Reference"
Write-Info "Tile path format:"
Write-Info "   spaces/{spaceName}/tiles/{sceneCode}/cubemap/{face}/level_{level}/tile_{y}_{x}.jpg"
Write-Info ""
Write-Info "Example paths:"
Write-Info "   spaces/dujiangyan/tiles/$SceneCode/cubemap/px/level_0/tile_0_0.jpg"
Write-Info "   spaces/dujiangyan/tiles/$SceneCode/cubemap/px/level_1/tile_1_0.jpg"
Write-Info ""
Write-Info "Please check if these files exist in MinIO"

Write-Header "Diagnosis Complete"
