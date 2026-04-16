# ============================================================
# MinIO Tile Checker
# Check if tiles exist in MinIO storage
# ============================================================

param(
    [string]$MinioUrl = "http://localhost:9000",
    [string]$Bucket = "panorama",
    [string]$SceneCode = "doujiangyanhoumen-yejing-616591bb",
    [string]$SpaceName = "dujiangyan"
)

function Write-Info { param($Message) Write-Host "[INFO] $Message" -ForegroundColor Green }
function Write-Err { param($Message) Write-Host "[ERROR] $Message" -ForegroundColor Red }
function Write-Warn { param($Message) Write-Host "[WARN] $Message" -ForegroundColor Yellow }
function Write-Header { param($Message) Write-Host "`n========== $Message ==========" -ForegroundColor Cyan }

Write-Header "MinIO Tile Checker"

Write-Info "Configuration:"
Write-Info "   MinIO URL: $MinioUrl"
Write-Info "   Bucket: $Bucket"
Write-Info "   Scene Code: $SceneCode"
Write-Info "   Space Name: $SpaceName"

Write-Header "Manual Verification Steps"

Write-Info "Step 1: Access MinIO Console"
Write-Info "   Open browser: $MinioUrl"
Write-Info "   Default credentials (if not changed):"
Write-Info "     Username: minioadmin"
Write-Info "     Password: minioadmin"

Write-Info ""
Write-Info "Step 2: Navigate to bucket"
Write-Info "   1. Click on '$Bucket' bucket"
Write-Info "   2. Navigate to: spaces/$SpaceName/tiles/$SceneCode/"

Write-Info ""
Write-Info "Step 3: Check expected structure"
Write-Info "   You should see:"
Write-Info "   spaces/$SpaceName/tiles/$SceneCode/"
Write-Info "   └── cubemap/"
Write-Info "       ├── px/"
Write-Info "       │   ├── level_0/ (4 tiles: tile_0_0.jpg, tile_0_1.jpg, tile_1_0.jpg, tile_1_1.jpg)"
Write-Info "       │   ├── level_1/ (16 tiles)"
Write-Info "       │   └── level_2/ (64 tiles)"
Write-Info "       ├── nx/"
Write-Info "       ├── py/"
Write-Info "       ├── ny/"
Write-Info "       ├── pz/"
Write-Info "       └── nz/"

Write-Info ""
Write-Info "Step 4: Check preview image"
Write-Info "   Path: spaces/$SpaceName/previews/$SceneCode/preview.jpg"

Write-Info ""
Write-Info "Step 5: Check source panorama"
Write-Info "   Path: spaces/$SpaceName/sources/$SceneCode/source.jpg"

Write-Header "Using MinIO Client (mc)"

Write-Info "If you have MinIO client installed, run these commands:"
Write-Info ""
Write-Host "mc alias set local $MinioUrl minioadmin minioadmin" -ForegroundColor Cyan
Write-Host "mc ls local/$Bucket/spaces/$SpaceName/tiles/$SceneCode/" -ForegroundColor Cyan
Write-Host "mc ls local/$Bucket/spaces/$SpaceName/tiles/$SceneCode/cubemap/px/level_0/" -ForegroundColor Cyan
Write-Host "mc ls local/$Bucket/spaces/$SpaceName/previews/$SceneCode/" -ForegroundColor Cyan

Write-Header "Common Issues"

Write-Warn "Issue 1: Bucket does not exist"
Write-Info "   Solution: Create bucket '$Bucket' in MinIO console"

Write-Warn "Issue 2: Path does not exist"
Write-Info "   Solution: Scene slicing may not have completed"
Write-Info "   Check: Database res_scenes.slice_status for scene '$SceneCode'"

Write-Warn "Issue 3: Tiles exist but API returns 404"
Write-Info "   Solution: Check if space_name/slug matches"
Write-Info "   The space slug in database should match the folder name in MinIO"

Write-Header "Automated Check (requires mc)"

Write-Info "Attempting to check MinIO with mc command..."
try {
    $result = mc ls "$MinioUrl/$Bucket/spaces/$SpaceName/tiles/$SceneCode/" 2>&1
    if ($LASTEXITCODE -eq 0) {
        Write-Info "Tiles found:"
        Write-Host $result
    } else {
        Write-Warn "mc command not found or path does not exist"
        Write-Info "Please install MinIO client or check manually via console"
    }
} catch {
    Write-Warn "mc command not available"
    Write-Info "Please check MinIO manually using the console at: $MinioUrl"
}

Write-Header "Check Complete"
