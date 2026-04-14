# ============================================================
# 全景资源读取接口网络测试
# 测试瓦片、预览图、封面的流式读取接口
# ============================================================

$BASE_URL = "http://localhost:7000"
$API_V1 = "$BASE_URL/api/v1"

function Write-Info { param($Message) Write-Host "[INFO] $Message" -ForegroundColor Green }
function Write-Err { param($Message) Write-Host "[ERROR] $Message" -ForegroundColor Red }
function Write-Warn { param($Message) Write-Host "[WARN] $Message" -ForegroundColor Yellow }
function Write-Header { param($Message) Write-Host "`n========== $Message ==========" -ForegroundColor Cyan }

# ------ 已知的测试数据（都江堰景区的场景） ------
# 这些 sceneCode 来自种子数据，如果你的数据不同请修改
$SCENE_CODE = "anlanqiao-18fbf61b"
$SPACE_NAME = "dujiangyan"

Write-Header "全景资源读取接口测试"

# 1. 健康检查
Write-Info "1. Health Check..."
try {
    $health = Invoke-RestMethod -Uri "$BASE_URL/health" -Method GET
    Write-Info "Health check passed: $($health.status)"
} catch {
    Write-Err "Health check failed: $_"
    Write-Err "请确保后端服务已在 $BASE_URL 运行"
    exit 1
}

# ------ 瓦片读取测试 ------
Write-Header "瓦片读取测试 (Tile Stream)"

# 2. 测试 Level 0 瓦片（首屏加载, 2x2=4 tiles per face）
Write-Info "2. Test Level 0 Tile (首屏加载)..."
$tileUrl = "$API_V1/res/tiles/$SCENE_CODE/px/0/0/0"
Write-Info "   URL: $tileUrl"
try {
    $resp = Invoke-WebRequest -UseBasicParsing -Uri $tileUrl -Method GET -UseBasicParsing
    $contentType = $resp.Headers["Content-Type"]
    $contentLength = $resp.Headers["Content-Length"]
    $cacheControl = $resp.Headers["Cache-Control"]
    $etag = $resp.Headers["ETag"]
    
    Write-Info "   Status: $($resp.StatusCode)"
    Write-Info "   Content-Type: $contentType"
    Write-Info "   Content-Length: $contentLength bytes"
    Write-Info "   Cache-Control: $cacheControl"
    Write-Info "   ETag: $etag"
    
    if ($resp.StatusCode -eq 200 -and $contentType -like "*image/jpeg*") {
        Write-Info "   ✅ Level 0 Tile OK"
    } else {
        Write-Err "   ❌ Level 0 Tile FAILED"
    }
} catch {
    Write-Err "   ❌ Level 0 Tile request failed: $_"
}

# 3. 测试 Level 2 瓦片（放大时加载, 8x8=64 tiles per face）
Write-Info "3. Test Level 2 Tile (高精度)..."
$tileUrl2 = "$API_V1/res/tiles/$SCENE_CODE/nx/2/3/5"
Write-Info "   URL: $tileUrl2"
try {
    $resp2 = Invoke-WebRequest -UseBasicParsing -Uri $tileUrl2 -Method GET
    Write-Info "   Status: $($resp2.StatusCode)"
    Write-Info "   Content-Length: $($resp2.Headers["Content-Length"]) bytes"
    
    if ($resp2.StatusCode -eq 200) {
        Write-Info "   ✅ Level 2 Tile OK"
    } else {
        Write-Err "   ❌ Level 2 Tile FAILED"
    }
} catch {
    Write-Err "   ❌ Level 2 Tile request failed: $_"
}

# 4. 测试所有 6 个面
Write-Info "4. Test All 6 Cubemap Faces..."
$faces = @("px", "nx", "py", "ny", "pz", "nz")
$faceResults = @()
foreach ($face in $faces) {
    $faceUrl = "$API_V1/res/tiles/$SCENE_CODE/$face/0/0/0"
    try {
        $faceResp = Invoke-WebRequest -UseBasicParsing -Uri $faceUrl -Method GET
        if ($faceResp.StatusCode -eq 200) {
            $faceResults += "✅ $face"
        } else {
            $faceResults += "❌ $face ($($faceResp.StatusCode))"
        }
    } catch {
        $faceResults += "❌ $face (error)"
    }
}
Write-Info "   Results: $($faceResults -join ' | ')"

# 5. 测试条件请求 (304 Not Modified)
Write-Info "5. Test Conditional Request (If-None-Match → 304)..."
try {
    # 先获取 ETag
    $firstResp = Invoke-WebRequest -UseBasicParsing -Uri "$API_V1/res/tiles/$SCENE_CODE/px/0/0/0" -Method GET
    $etag = $firstResp.Headers["ETag"]
    Write-Info "   Got ETag: $etag"
    
    if ($etag) {
        # 带 If-None-Match 请求
        $headers = @{ "If-None-Match" = $etag }
        try {
            $condResp = Invoke-WebRequest -UseBasicParsing -Uri "$API_V1/res/tiles/$SCENE_CODE/px/0/0/0" -Method GET -Headers $headers
            if ($condResp.StatusCode -eq 304) {
                Write-Info "   ✅ Got 304 Not Modified (conditional request works)"
            } else {
                Write-Warn "   ⚠️ Got $($condResp.StatusCode) instead of 304"
            }
        } catch {
            $statusCode = $_.Exception.Response.StatusCode.value__
            if ($statusCode -eq 304) {
                Write-Info "   ✅ Got 304 Not Modified (conditional request works)"
            } else {
                Write-Warn "   ⚠️ Got status $statusCode instead of 304"
            }
        }
    }
} catch {
    Write-Err "   ❌ Conditional request test failed: $_"
}

# ------ 错误用例测试 ------
Write-Header "错误用例测试"

# 6. 无效的 face
Write-Info "6. Test Invalid Face..."
try {
    $invalidResp = Invoke-WebRequest -UseBasicParsing -Uri "$API_V1/res/tiles/$SCENE_CODE/xx/0/0/0" -Method GET -ErrorAction Stop
    Write-Err "   ❌ Should have returned 400, got $($invalidResp.StatusCode)"
} catch {
    $statusCode = $_.Exception.Response.StatusCode.value__
    if ($statusCode -eq 400) {
        Write-Info "   ✅ Got 400 Bad Request for invalid face"
    } else {
        Write-Info "   Got status $statusCode"
    }
}

# 7. 无效的 level
Write-Info "7. Test Invalid Level (level=5)..."
try {
    $invalidLevelResp = Invoke-WebRequest -UseBasicParsing -Uri "$API_V1/res/tiles/$SCENE_CODE/px/5/0/0" -Method GET -ErrorAction Stop
    Write-Err "   ❌ Should have returned 400, got $($invalidLevelResp.StatusCode)"
} catch {
    $statusCode = $_.Exception.Response.StatusCode.value__
    if ($statusCode -eq 400) {
        Write-Info "   ✅ Got 400 Bad Request for invalid level"
    } else {
        Write-Info "   Got status $statusCode"
    }
}

# 8. 不存在的 sceneCode
Write-Info "8. Test Non-existent SceneCode..."
try {
    $notFoundResp = Invoke-WebRequest -UseBasicParsing -Uri "$API_V1/res/tiles/nonexistent-scene/px/0/0/0" -Method GET -ErrorAction Stop
    Write-Err "   ❌ Should have returned 404, got $($notFoundResp.StatusCode)"
} catch {
    $statusCode = $_.Exception.Response.StatusCode.value__
    if ($statusCode -eq 404) {
        Write-Info "   ✅ Got 404 Not Found for non-existent scene"
    } else {
        Write-Info "   Got status $statusCode"
    }
}

# ------ 预览图测试 ------
Write-Header "预览图读取测试 (Preview Stream)"

Write-Info "9. Test Preview Image..."
$previewUrl = "$API_V1/res/previews/$SCENE_CODE"
Write-Info "   URL: $previewUrl"
try {
    $previewResp = Invoke-WebRequest -UseBasicParsing -Uri $previewUrl -Method GET
    Write-Info "   Status: $($previewResp.StatusCode)"
    Write-Info "   Content-Type: $($previewResp.Headers["Content-Type"])"
    Write-Info "   Content-Length: $($previewResp.Headers["Content-Length"]) bytes"
    
    if ($previewResp.StatusCode -eq 200) {
        Write-Info "   ✅ Preview OK"
    }
} catch {
    $statusCode = $_.Exception.Response.StatusCode.value__
    Write-Warn "   ⚠️ Preview not found (status $statusCode) — may not have been generated yet"
}

# ------ 封面测试 ------
Write-Header "封面读取测试 (Cover Stream)"

Write-Info "10. Test Cover Image..."
$coverUrl = "$API_V1/res/covers/$SPACE_NAME"
Write-Info "   URL: $coverUrl"
try {
    $coverResp = Invoke-WebRequest -UseBasicParsing -Uri $coverUrl -Method GET
    Write-Info "   Status: $($coverResp.StatusCode)"
    Write-Info "   Content-Type: $($coverResp.Headers["Content-Type"])"
    Write-Info "   Content-Length: $($coverResp.Headers["Content-Length"]) bytes"
    
    if ($coverResp.StatusCode -eq 200) {
        Write-Info "   ✅ Cover OK"
    }
} catch {
    $statusCode = $_.Exception.Response.StatusCode.value__
    Write-Warn "   ⚠️ Cover not found (status $statusCode) — may not have been uploaded yet"
}

# ------ 性能测试：批量瓦片请求 ------
Write-Header "性能测试：批量 Level 0 瓦片请求"

Write-Info "11. Fetching all Level 0 tiles for all 6 faces (6 faces × 4 tiles = 24 requests)..."
$totalTime = [System.Diagnostics.Stopwatch]::StartNew()
$successCount = 0
$failCount = 0
$totalBytes = 0

foreach ($face in $faces) {
    for ($y = 0; $y -lt 2; $y++) {
        for ($x = 0; $x -lt 2; $x++) {
            $url = "$API_V1/res/tiles/$SCENE_CODE/$face/0/$y/$x"
            try {
                $r = Invoke-WebRequest -UseBasicParsing -Uri $url -Method GET
                $successCount++
                $totalBytes += [int]$r.Headers["Content-Length"]
            } catch {
                $failCount++
            }
        }
    }
}
$totalTime.Stop()

Write-Info "   Total Requests: $($successCount + $failCount)"
Write-Info "   Success: $successCount | Failed: $failCount"
Write-Info "   Total Data: $([math]::Round($totalBytes / 1024, 1)) KB"
Write-Info "   Total Time: $($totalTime.ElapsedMilliseconds) ms"
Write-Info "   Avg per Tile: $([math]::Round($totalTime.ElapsedMilliseconds / [math]::Max(1, $successCount + $failCount), 1)) ms"

# ------ 总结 ------
Write-Header "测试总结"
Write-Info "瓦片接口:   /api/v1/res/tiles/:sceneCode/:face/:level/:y/:x"
Write-Info "预览图接口: /api/v1/res/previews/:sceneCode"
Write-Info "封面接口:   /api/v1/res/covers/:spaceName"
Write-Info ""
Write-Info "所有接口无需鉴权，支持强缓存 + ETag 条件请求"
Write-Info "=== 测试完成 ==="
