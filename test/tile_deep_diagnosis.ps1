# ============================================================
# Tile 404 Issue Deep Diagnosis
# Check database, MinIO path, and backend logs
# ============================================================

param(
    [string]$SceneCode = "doujiangyanhoumen-yejing-616591bb",
    [string]$BaseUrl = "http://localhost:7000"
)

$ErrorActionPreference = "Continue"

function Write-Info { param($Message) Write-Host "[INFO] $Message" -ForegroundColor Green }
function Write-Err { param($Message) Write-Host "[ERROR] $Message" -ForegroundColor Red }
function Write-Warn { param($Message) Write-Host "[WARN] $Message" -ForegroundColor Yellow }
function Write-Header { param($Message) Write-Host "`n========== $Message ==========" -ForegroundColor Cyan }
function Write-Debug { param($Message) Write-Host "[DEBUG] $Message" -ForegroundColor Gray }

Write-Header "Tile 404 Deep Diagnosis"

# 1. Check backend health
Write-Header "1. Backend Health"
try {
    $health = Invoke-RestMethod -Uri "$BaseUrl/health" -Method GET -TimeoutSec 5
    Write-Info "Backend is running"
} catch {
    Write-Err "Backend not running: $_"
    exit 1
}

# 2. Test tile request and capture response
Write-Header "2. Test Tile Request"
$testUrl = "$BaseUrl/api/v1/res/tiles/$SceneCode/px/0/0/0"
Write-Info "Requesting: $testUrl"

try {
    $resp = Invoke-WebRequest -UseBasicParsing -Uri $testUrl -Method GET -TimeoutSec 10
    Write-Info "Success! Status: $($resp.StatusCode)"
    Write-Info "Content-Type: $($resp.Headers['Content-Type'])"
    Write-Info "Content-Length: $($resp.Headers['Content-Length'])"
} catch {
    $statusCode = $_.Exception.Response.StatusCode.value__
    Write-Err "Request failed with status: $statusCode"
    
    # Get response body
    try {
        $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
        $reader.BaseStream.Position = 0
        $responseBody = $reader.ReadToEnd()
        Write-Err "Response body: $responseBody"
    } catch {
        Write-Debug "Could not read response body"
    }
}

# 3. Check backend logs for the request
Write-Header "3. Backend Log Analysis"
Write-Info "Looking for tile request logs..."

$logPath = "d:\lwdd\code\毕设\webGL-720yun\backend\logs\app.log"
if (Test-Path $logPath) {
    $recentLogs = Get-Content $logPath -Tail 50
    $tileLogs = $recentLogs | Select-String -Pattern "GetTileStream|scene_code=$SceneCode|objectPath" | Select-Object -Last 10
    
    if ($tileLogs) {
        Write-Info "Found tile-related logs:"
        $tileLogs | ForEach-Object { Write-Debug $_ }
    } else {
        Write-Warn "No tile logs found in recent entries"
    }
    
    # Check for errors
    $errorLogs = $recentLogs | Select-String -Pattern "error|Error|ERROR|failed|Failed" | Select-Object -Last 5
    if ($errorLogs) {
        Write-Warn "Found error logs:"
        $errorLogs | ForEach-Object { Write-Err $_ }
    }
} else {
    Write-Warn "Log file not found at: $logPath"
}

# 4. Database check instructions
Write-Header "4. Database Check Instructions"
Write-Info "Please run these SQL queries to check the scene:"

$sqlQueries = @"
-- Check if scene exists
SELECT 
    id, 
    scene_code, 
    title, 
    slice_status,
    space_id,
    source_url,
    status
FROM res_scene 
WHERE scene_code = '$SceneCode';

-- Check the associated space
SELECT 
    s.id AS scene_id,
    s.scene_code,
    s.slice_status,
    sp.id AS space_id,
    sp.name AS space_name,
    sp.slug AS space_slug
FROM res_scene s
LEFT JOIN res_space sp ON s.space_id = sp.id
WHERE s.scene_code = '$SceneCode';

-- Check all scenes with their space info
SELECT 
    s.scene_code,
    s.slice_status,
    sp.slug AS space_slug
FROM res_scene s
LEFT JOIN res_space sp ON s.space_id = sp.id
ORDER BY s.id DESC
LIMIT 10;
"@

Write-Host $sqlQueries -ForegroundColor Cyan

# 5. MinIO path verification
Write-Header "5. MinIO Path Verification"
Write-Info "Expected MinIO paths based on your structure:"
Write-Info ""
Write-Info "Bucket: resource-scene"
Write-Info "Tile path: spaces/{space_slug}/tiles/$SceneCode/cubemap/{face}/level_{level}/tile_{y}_{x}.jpg"
Write-Info ""
Write-Info "Example paths:"
Write-Info "  spaces/dujiangyan/tiles/$SceneCode/cubemap/px/level_0/tile_0_0.jpg"
Write-Info "  spaces/dujiangyan/tiles/$SceneCode/cubemap/px/level_1/tile_0_0.jpg"
Write-Info ""
Write-Info "Your MinIO structure shows:"
Write-Info "  spaces/dujiangyan/tiles/$SceneCode/cubemap/px.jpg (face image)"
Write-Info "  spaces/dujiangyan/tiles/$SceneCode/cubemap/px/level_0/ (tile directory)"
Write-Info ""
Write-Warn "IMPORTANT: Check if level_0 directory contains tile files!"
Write-Info "Expected files in level_0:"
Write-Info "  tile_0_0.jpg, tile_0_1.jpg, tile_1_0.jpg, tile_1_1.jpg"

# 6. Common issues checklist
Write-Header "6. Common Issues Checklist"

Write-Info "[ ] Scene exists in database with slice_status = 'ready'"
Write-Info "[ ] Space slug matches MinIO folder name (dujiangyan)"
Write-Info "[ ] MinIO bucket name is 'resource-scene' (matches config)"
Write-Info "[ ] Tile files exist in MinIO at correct path"
Write-Info "[ ] Level directories (level_0, level_1, level_2) contain tile files"

# 7. Quick fix suggestions
Write-Header "7. Quick Fix Suggestions"

Write-Warn "If slice_status is not 'ready':"
Write-Info "   UPDATE res_scene SET slice_status = 'ready' WHERE scene_code = '$SceneCode';"

Write-Warn ""
Write-Warn "If space_slug doesn't match:"
Write-Info "   Check res_space.slug for the space associated with this scene"
Write-Info "   UPDATE res_space SET slug = 'dujiangyan' WHERE id = (SELECT space_id FROM res_scene WHERE scene_code = '$SceneCode');"

Write-Warn ""
Write-Warn "If tiles exist but at wrong path:"
Write-Info "   Backend expects: spaces/{slug}/tiles/{sceneCode}/cubemap/{face}/level_{level}/tile_{y}_{x}.jpg"
Write-Info "   Check MinIO console to verify actual path"

Write-Header "Diagnosis Complete"
Write-Info "Next steps:"
Write-Info "1. Run the SQL queries above to check database"
Write-Info "2. Verify MinIO path structure"
Write-Info "3. Check backend logs for detailed error messages"
