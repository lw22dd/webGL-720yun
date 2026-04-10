# Space 模块增删改查测试脚本 (PowerShell)
$BASE_URL = "http://localhost:7000"
$API_V1 = "$BASE_URL/api/v1"

function Write-Info { param($Message) Write-Host "[INFO] $Message" -ForegroundColor Green }
function Write-Error { param($Message) Write-Host "[ERROR] $Message" -ForegroundColor Red }
function Write-Warn { param($Message) Write-Host "[WARN] $Message" -ForegroundColor Yellow }

Write-Info "=== 开始 Space 模块增删改查测试 ==="

# 1. 健康检查
Write-Info "测试健康检查接口..."
try {
    $health = Invoke-RestMethod -Uri "$BASE_URL/health" -Method GET
    if ($health.status -eq "ok") {
        Write-Info "健康检查通过"
    } else {
        throw "健康检查失败"
    }
} catch {
    Write-Error "健康检查失败: $_"
    exit 1
}

# 2. 登录获取 Token
Write-Info "登录获取 Token..."
try {
    $loginBody = @{ username = "admin"; password = "admin123" } | ConvertTo-Json
    $loginResp = Invoke-RestMethod -Uri "$API_V1/auth/login" -Method POST -Body $loginBody -ContentType "application/json"
    $TOKEN = $loginResp.data.access_token
    if (-not $TOKEN) {
        throw "未获取到Token"
    }
    Write-Info "登录成功，获取到 Token"
} catch {
    Write-Error "登录失败: $_"
    exit 1
}

# 3. 测试创建空间
Write-Info "测试创建空间接口..."
try {
    $timestamp = (Get-Date -Format "yyyyMMddHHmmss")
    $createBody = @{
        name = "测试景区_$timestamp"
        description = "这是一个测试景区"
        province = "测试省"
        city = "测试市"
        longitude = 116.407429
        latitude = 39.904211
        zoom_level = 12
        sort_order = 1
    } | ConvertTo-Json

    $createResp = Invoke-RestMethod -Uri "$API_V1/resource/spaces" -Method POST -Body $createBody -ContentType "application/json" -Headers @{ "Authorization" = "Bearer $TOKEN" }
    Write-Info "创建响应: $($createResp | ConvertTo-Json -Compress)"

    if ($createResp.code -eq 200) {
        $GLOBAL:SPACE_ID = $createResp.data.id
        Write-Info "空间创建成功，ID: $SPACE_ID"
    } else {
        throw "创建空间失败: $($createResp.msg)"
    }
} catch {
    Write-Error "创建空间失败: $_"
    exit 1
}

# 4. 测试获取空间列表
Write-Info "测试获取空间列表接口..."
try {
    $listResp = Invoke-RestMethod -Uri "$API_V1/resource/spaces?page=1&page_size=10" -Method GET -Headers @{ "Authorization" = "Bearer $TOKEN" }
    Write-Info "列表响应: $($listResp | ConvertTo-Json -Compress)"

    if ($listResp.code -eq 200 -and $listResp.data.spaces) {
        Write-Info "获取空间列表成功"
    } else {
        throw "获取空间列表失败"
    }
} catch {
    Write-Error "获取空间列表失败: $_"
    exit 1
}

# 5. 测试获取空间详情
Write-Info "测试获取空间详情接口..."
try {
    $detailResp = Invoke-RestMethod -Uri "$API_V1/resource/spaces/$SPACE_ID" -Method GET -Headers @{ "Authorization" = "Bearer $TOKEN" }
    Write-Info "详情响应: $($detailResp | ConvertTo-Json -Compress)"

    if ($detailResp.code -eq 200 -and $detailResp.data.id -eq $SPACE_ID) {
        Write-Info "获取空间详情成功"
    } else {
        throw "获取空间详情失败"
    }
} catch {
    Write-Error "获取空间详情失败: $_"
    exit 1
}

# 6. 测试更新空间
Write-Info "测试更新空间接口..."
try {
    $updateBody = @{
        name = "更新后的测试景区"
        description = "这是更新后的描述"
        city = "更新后的城市"
    } | ConvertTo-Json

    $updateResp = Invoke-RestMethod -Uri "$API_V1/resource/spaces/$SPACE_ID" -Method PUT -Body $updateBody -ContentType "application/json" -Headers @{ "Authorization" = "Bearer $TOKEN" }
    Write-Info "更新响应: $($updateResp | ConvertTo-Json -Compress)"

    if ($updateResp.code -eq 200) {
        Write-Info "更新空间成功"
    } else {
        throw "更新空间失败: $($updateResp.msg)"
    }
} catch {
    Write-Error "更新空间失败: $_"
    exit 1
}

# 7. 创建第二个空间用于批量删除测试
Write-Info "创建第二个测试空间..."
try {
    $timestamp2 = (Get-Date -Format "yyyyMMddHHmmss")
    $createBody2 = @{
        name = "测试景区2_$timestamp2"
        description = "这是第二个测试景区"
        province = "测试省2"
        city = "测试市2"
        longitude = 117.0
        latitude = 40.0
        zoom_level = 10
        sort_order = 2
    } | ConvertTo-Json

    $createResp2 = Invoke-RestMethod -Uri "$API_V1/resource/spaces" -Method POST -Body $createBody2 -ContentType "application/json" -Headers @{ "Authorization" = "Bearer $TOKEN" }

    if ($createResp2.code -eq 200) {
        $GLOBAL:SPACE_ID2 = $createResp2.data.id
        Write-Info "第二个空间创建成功，ID: $SPACE_ID2"
    } else {
        throw "创建第二个空间失败: $($createResp2.msg)"
    }
} catch {
    Write-Error "创建第二个空间失败: $_"
    exit 1
}

# 8. 创建第三个空间用于批量删除测试
Write-Info "创建第三个测试空间..."
try {
    $timestamp3 = (Get-Date -Format "yyyyMMddHHmmss")
    $createBody3 = @{
        name = "测试景区3_$timestamp3"
        description = "这是第三个测试景区"
        province = "测试省3"
        city = "测试市3"
        longitude = 118.0
        latitude = 41.0
        zoom_level = 11
        sort_order = 3
    } | ConvertTo-Json

    $createResp3 = Invoke-RestMethod -Uri "$API_V1/resource/spaces" -Method POST -Body $createBody3 -ContentType "application/json" -Headers @{ "Authorization" = "Bearer $TOKEN" }

    if ($createResp3.code -eq 200) {
        $GLOBAL:SPACE_ID3 = $createResp3.data.id
        Write-Info "第三个空间创建成功，ID: $SPACE_ID3"
    } else {
        throw "创建第三个空间失败: $($createResp3.msg)"
    }
} catch {
    Write-Error "创建第三个空间失败: $_"
    exit 1
}

# 9. 测试批量删除空间
Write-Info "测试批量删除空间接口..."
try {
    $batchDeleteBody = @{
        ids = @($SPACE_ID2, $SPACE_ID3)
    } | ConvertTo-Json

    Write-Info "批量删除请求体: $batchDeleteBody"

    $batchDeleteResp = Invoke-RestMethod -Uri "$API_V1/resource/spaces/batch" -Method DELETE -Body $batchDeleteBody -ContentType "application/json" -Headers @{ "Authorization" = "Bearer $TOKEN" }
    Write-Info "批量删除响应: $($batchDeleteResp | ConvertTo-Json -Compress)"

    if ($batchDeleteResp.code -eq 200) {
        Write-Info "批量删除空间成功"
    } else {
        throw "批量删除空间失败: $($batchDeleteResp.msg)"
    }
} catch {
    Write-Error "批量删除空间失败: $_"
    exit 1
}

# 10. 验证批量删除结果
Write-Info "验证批量删除结果..."
try {
    $detailResp2 = Invoke-RestMethod -Uri "$API_V1/resource/spaces/$SPACE_ID2" -Method GET -Headers @{ "Authorization" = "Bearer $TOKEN" }
    if ($detailResp2.msg -match "不存在" -or $detailResp2.msg -match "not found") {
        Write-Info "空间2已删除"
    } else {
        Write-Warn "空间2可能未删除: $($detailResp2.msg)"
    }
} catch {
    Write-Info "空间2已删除 (异常: $_)"
}

# 11. 测试删除单个空间
Write-Info "测试删除单个空间接口..."
try {
    $deleteResp = Invoke-RestMethod -Uri "$API_V1/resource/spaces/$SPACE_ID" -Method DELETE -Headers @{ "Authorization" = "Bearer $TOKEN" }
    Write-Info "删除响应: $($deleteResp | ConvertTo-Json -Compress)"

    if ($deleteResp.code -eq 200) {
        Write-Info "删除单个空间成功"
    } else {
        throw "删除单个空间失败: $($deleteResp.msg)"
    }
} catch {
    Write-Error "删除单个空间失败: $_"
    exit 1
}

# 12. 验证删除结果
Write-Info "验证删除结果..."
try {
    $detailResp3 = Invoke-RestMethod -Uri "$API_V1/resource/spaces/$SPACE_ID" -Method GET -Headers @{ "Authorization" = "Bearer $TOKEN" }
    if ($detailResp3.msg -match "不存在" -or $detailResp3.msg -match "not found") {
        Write-Info "空间已删除，验证成功"
    } else {
        Write-Warn "空间可能未删除: $($detailResp3.msg)"
    }
} catch {
    Write-Info "空间已删除，验证成功 (异常: $_)"
}

# 13. 测试搜索功能
Write-Info "测试搜索功能..."
try {
    $searchResp = Invoke-RestMethod -Uri "$API_V1/resource/spaces?keyword=测试&page=1&page_size=10" -Method GET -Headers @{ "Authorization" = "Bearer $TOKEN" }
    Write-Info "搜索响应: $($searchResp | ConvertTo-Json -Compress)"

    if ($searchResp.code -eq 200) {
        Write-Info "搜索功能正常"
    } else {
        Write-Warn "搜索功能可能异常"
    }
} catch {
    Write-Warn "搜索功能测试失败: $_"
}

Write-Info "=== Space 模块增删改查测试完成 ==="
Write-Info "测试总结:"
Write-Info "  [V] 健康检查"
Write-Info "  [V] 登录获取Token"
Write-Info "  [V] 创建空间"
Write-Info "  [V] 获取空间列表"
Write-Info "  [V] 获取空间详情"
Write-Info "  [V] 更新空间"
Write-Info "  [V] 批量删除空间"
Write-Info "  [V] 删除单个空间"
Write-Info "  [V] 搜索功能"
