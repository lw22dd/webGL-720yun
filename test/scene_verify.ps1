# ============================================================
# Scene Data Verification Script
# Check if scene exists in database and MinIO
# ============================================================

param(
    [string]$SceneCode = "doujiangyanhoumen-yejing-616591bb"
)

$ErrorActionPreference = "Continue"

function Write-Info { param($Message) Write-Host "[INFO] $Message" -ForegroundColor Green }
function Write-Err { param($Message) Write-Host "[ERROR] $Message" -ForegroundColor Red }
function Write-Warn { param($Message) Write-Host "[WARN] $Message" -ForegroundColor Yellow }
function Write-Header { param($Message) Write-Host "`n========== $Message ==========" -ForegroundColor Cyan }

Write-Header "Scene Data Verification - $SceneCode"

# 1. Check if backend is running
Write-Header "1. Backend Service Check"
try {
    $health = Invoke-RestMethod -Uri "http://localhost:7000/health" -Method GET -TimeoutSec 5
    Write-Info "Backend is running: $($health.status)"
} catch {
    Write-Err "Backend is not running. Please start the backend service first."
    exit 1
}

# 2. Check scene list API
Write-Header "2. Check Scene List API"
Write-Info "Fetching scene list from backend..."

# Try to get scene list (may need auth, but let's try public endpoint first)
$sceneListUrl = "http://localhost:7000/api/v1/resource/scenes?page=1&page_size=100"
try {
    # First try without auth
    $scenesResp = Invoke-RestMethod -Uri $sceneListUrl -Method GET -TimeoutSec 10
    Write-Info "Found $($scenesResp.data.total) scenes in database"
    
    # Check if our scene exists
    $targetScene = $scenesResp.data.scenes | Where-Object { $_.scene_code -eq $SceneCode }
    if ($targetScene) {
        Write-Info "Scene '$SceneCode' found in database:"
        Write-Info "   ID: $($targetScene.id)"
        Write-Info "   Title: $($targetScene.title)"
        Write-Info "   Status: $($targetScene.status)"
        Write-Info "   Slice Status: $($targetScene.slice_status)"
        Write-Info "   Space ID: $($targetScene.space_id)"
    } else {
        Write-Err "Scene '$SceneCode' NOT found in database!"
        Write-Warn "Available scenes:"
        $scenesResp.data.scenes | Select-Object -First 10 | ForEach-Object {
            Write-Host "   - $($_.scene_code) ($($_.title))" -ForegroundColor Gray
        }
    }
} catch {
    $statusCode = $_.Exception.Response.StatusCode.value__
    if ($statusCode -eq 401) {
        Write-Warn "Scene list API requires authentication"
        Write-Info "Please check database manually or provide auth token"
    } else {
        Write-Err "Failed to fetch scene list: $_"
    }
}

# 3. Check MinIO connection
Write-Header "3. MinIO Connection Check"
Write-Info "MinIO is typically at localhost:9000"
Write-Info "Please verify MinIO is running and check the following:"
Write-Info ""
Write-Info "1. Open MinIO Console: http://localhost:9000"
Write-Info "2. Check bucket: 'panorama' (or your configured bucket name)"
Write-Info "3. Navigate to path: spaces/{spaceName}/tiles/$SceneCode/"
Write-Info ""
Write-Info "Expected tile structure:"
Write-Info "   spaces/{spaceName}/tiles/$SceneCode/cubemap/"
Write-Info "   ├── px/"
Write-Info "   │   ├── level_0/"
Write-Info "   │   │   ├── tile_0_0.jpg"
Write-Info "   │   │   ├── tile_0_1.jpg"
Write-Info "   │   │   ├── tile_1_0.jpg"
Write-Info "   │   │   └── tile_1_1.jpg"
Write-Info "   │   └── level_1/"
Write-Info "   │       └── ... (16 tiles)"
Write-Info "   ├── nx/"
Write-Info "   ├── py/"
Write-Info "   ├── ny/"
Write-Info "   ├── pz/"
Write-Info "   └── nz/"

# 4. Check backend logs
Write-Header "4. Backend Log Analysis"
Write-Info "Check backend logs for slicing task status:"
Write-Info ""
Write-Info "Look for these log patterns:"
Write-Info "   - 'Slicing task started for scene: $SceneCode'"
Write-Info "   - 'Slicing completed for scene: $SceneCode'"
Write-Info "   - 'Slicing failed for scene: $SceneCode'"
Write-Info "   - 'Error processing panorama'"
Write-Info ""
Write-Info "Common issues:"
Write-Info "   1. Slice status is 'pending' or 'slicing' - task not completed"
Write-Info "   2. Slice status is 'failed' - check error logs"
Write-Info "   3. No slicing task created - scene created without panorama file"

# 5. Provide fix suggestions
Write-Header "5. Fix Suggestions"

Write-Info "If scene does not exist in database:"
Write-Info "   1. Create the scene with correct scene_code"
Write-Info "   2. Upload panorama file during creation"
Write-Info "   3. Wait for slicing task to complete"

Write-Info ""
Write-Info "If scene exists but slice_status is not 'completed':"
Write-Info "   1. Check slicing worker logs"
Write-Info "   2. Re-trigger slicing task if needed"
Write-Info "   3. Verify source panorama file exists in MinIO"

Write-Info ""
Write-Info "If scene exists and slice_status is 'completed' but tiles missing:"
Write-Info "   1. Check MinIO bucket permissions"
Write-Info "   2. Verify space name/slug is correct"
Write-Info "   3. Check if tiles were uploaded to wrong path"

Write-Header "Verification Complete"
Write-Info "Next steps:"
Write-Info "   1. Check database for scene record"
Write-Info "   2. Check MinIO for tile files"
Write-Info "   3. Check backend logs for slicing errors"
