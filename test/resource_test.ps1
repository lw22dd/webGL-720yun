# Resource Module Test Script
$BASE_URL = "http://localhost:7000"
$API_V1 = "$BASE_URL/api/v1"

# Color definitions
$RED = "`e[31m"
$GREEN = "`e[32m"
$YELLOW = "`e[33m"
$NC = "`e[0m"

function Log-Info { param($msg) Write-Host "$GREEN[INFO] $msg$NC" }
function Log-Warn { param($msg) Write-Host "$YELLOW[WARN] $msg$NC" }
function Log-Error { param($msg) Write-Host "$RED[ERROR] $msg$NC" }

function Send-Request {
    param($method, $url, $data, $headers)
    try {
        $params = @{
            Uri = $url
            Method = $method
            ContentType = "application/json"
        }
        if ($headers) {
            $params.Headers = $headers
        }
        if ($data) {
            $params.Body = $data
        }
        $response = Invoke-RestMethod @params
        return $response
    } catch {
        return $_.Exception.Response
    }
}

function Login {
    param($username, $password)
    $body = "{`"username`":`"$username`",`"password`":`"$password`"}"
    $resp = Send-Request -method "POST" -url "$API_V1/auth/login" -data $body
    return $resp.data.access_token
}

function Check-Status {
    param($method, $url, $expected_status, $data, $headers)
    try {
        $params = @{
            Uri = $url
            Method = $method
            ContentType = "application/json"
        }
        if ($headers) {
            $params.Headers = $headers
        }
        if ($data) {
            $params.Body = $data
        }
        $response = Invoke-RestMethod @params
        Log-Info "Test passed: $method $url"
        return $response
    } catch {
        Log-Error "Test failed: $method $url - $($_.Exception.Message)"
        return $null
    }
}

Log-Info "=== Starting Resource Module Test ==="

# Login as admin
$ADMIN_TOKEN = Login -username "admin" -password "admin123"
if (-not $ADMIN_TOKEN) {
    Log-Error "Admin login failed"
    exit 1
}
Log-Info "Admin login successful"
$AUTH_HEADER = @{ "Authorization" = "Bearer $ADMIN_TOKEN" }

$SUFFIX = Get-Date -Format "yyyyMMddHHmmss"

# ========== Space Management Test ==========
Log-Info "=== Space Management Test ==="

# 1. Create Space
$SPACE_NAME = "TestSpace_$SUFFIX"
$SPACE_SLUG = "test_space_$SUFFIX"
Log-Info "Testing create space..."
$spaceBody = "{`"name`":`"$SPACE_NAME`",`"slug`":`"$SPACE_SLUG`",`"description`":`"This is a test space`"}"
$CREATE_SPACE_RESP = Check-Status -method "POST" -url "$API_V1/resource/spaces" -expected_status 200 -data $spaceBody -headers $AUTH_HEADER

if ($CREATE_SPACE_RESP -and $CREATE_SPACE_RESP.data) {
    $SPACE_ID = $CREATE_SPACE_RESP.data.id
    Log-Info "Create space success, ID: $SPACE_ID"
} else {
    Log-Error "Create space failed"
    exit 1
}

# 2. Get Space List
Log-Info "Testing get space list..."
$listUrl = "$API_V1/resource/spaces?page=1&page_size=10"
Check-Status -method "GET" -url $listUrl -expected_status 200 -headers $AUTH_HEADER

# 3. Get Space Detail
Log-Info "Testing get space detail..."
Check-Status -method "GET" -url "$API_V1/resource/spaces/$SPACE_ID" -expected_status 200 -headers $AUTH_HEADER

# 4. Update Space
Log-Info "Testing update space..."
$updateBody = "{`"name`":`"${SPACE_NAME}_Updated`",`"description`":`"Updated description`"}"
Check-Status -method "PUT" -url "$API_V1/resource/spaces/$SPACE_ID" -expected_status 200 -data $updateBody -headers $AUTH_HEADER

# ========== Scene Management Test ==========
Log-Info "=== Scene Management Test ==="

# 5. Create Scene (without file)
$SCENE_TITLE = "TestScene_$SUFFIX"
Log-Info "Testing create scene (without file)..."
$sceneBody = "{`"space_id`":$SPACE_ID,`"title`":`"$SCENE_TITLE`",`"initial_fov`":100,`"initial_pitch`":0,`"initial_yaw`":0,`"longitude`":116.3974,`"latitude`":39.9093,`"sort_order`":0}"
$CREATE_SCENE_RESP = Check-Status -method "POST" -url "$API_V1/resource/scenes" -expected_status 200 -data $sceneBody -headers $AUTH_HEADER

if ($CREATE_SCENE_RESP -and $CREATE_SCENE_RESP.data) {
    $SCENE_ID = $CREATE_SCENE_RESP.data.scene_id
    if (-not $SCENE_ID) {
        $SCENE_ID = $CREATE_SCENE_RESP.data.id
    }
    Log-Info "Create scene success, ID: $SCENE_ID"
} else {
    Log-Error "Create scene failed"
}

# 6. Get Scene List
Log-Info "Testing get scene list..."
$sceneListUrl = "$API_V1/resource/scenes?space_id=$SPACE_ID&page=1&page_size=10"
Check-Status -method "GET" -url $sceneListUrl -expected_status 200 -headers $AUTH_HEADER

# 7. Get Scene Detail, Update, Delete (if scene ID exists)
if ($SCENE_ID) {
    Log-Info "Testing get scene detail..."
    Check-Status -method "GET" -url "$API_V1/resource/scenes/$SCENE_ID" -expected_status 200 -headers $AUTH_HEADER

    Log-Info "Testing update scene..."
    $updateSceneBody = "{`"title`":`"${SCENE_TITLE}_Updated`",`"initial_fov`":120}"
    Check-Status -method "PUT" -url "$API_V1/resource/scenes/$SCENE_ID" -expected_status 200 -data $updateSceneBody -headers $AUTH_HEADER

    Log-Info "Testing delete scene..."
    Check-Status -method "DELETE" -url "$API_V1/resource/scenes/$SCENE_ID" -expected_status 200 -headers $AUTH_HEADER
    Log-Info "Scene deleted successfully"
} else {
    Log-Warn "Skip scene detail, update and delete tests (no valid scene ID)"
}

# ========== Space Graph Data Test ==========
Log-Info "Testing get space graph data..."
Check-Status -method "GET" -url "$API_V1/resource/spaces/$SPACE_ID/graph" -expected_status 200 -headers $AUTH_HEADER

# 8. Delete Space
Log-Info "Testing delete space..."
Check-Status -method "DELETE" -url "$API_V1/resource/spaces/$SPACE_ID" -expected_status 200 -headers $AUTH_HEADER
Log-Info "Space deleted successfully"

Log-Info "=== Resource Module Test Completed ==="
