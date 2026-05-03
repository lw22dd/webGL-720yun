$script:BASE_URL = "http://localhost:7000"
$script:API_V1 = "$BASE_URL/api/v1"

function Log-Info {
    param([string]$Message)
    Write-Host "[INFO] $Message" -ForegroundColor Green
}

function Log-Warn {
    param([string]$Message)
    Write-Host "[WARN] $Message" -ForegroundColor Yellow
}

function Log-Error {
    param([string]$Message)
    Write-Host "[ERROR] $Message" -ForegroundColor Red
}

function Send-Request {
    param(
        [string]$Method,
        [string]$Url,
        [string]$Data = "",
        [string[]]$Headers = @()
    )

    $headersDict = @{"Content-Type" = "application/json"}
    foreach ($h in $Headers) {
        if ($h -match "(\w+):\s*(.+)") {
            $headersDict[$matches[1]] = $matches[2].Trim()
        }
    }

    $params = @{
        Method = $Method
        Uri = $Url
        Headers = $headersDict
    }

    if ($Data) {
        $params.Body = $Data
    }

    $tempFile = "$env:TEMP\ps_resp_body_$PID.txt"
    $statusCodeFile = "$env:TEMP\ps_status_$PID.txt"

    try {
        $statusCode = [int]([System.Net.HttpWebRequest]::Create($Url)).GetResponse().StatusCode
    } catch {
        $statusCode = [int]$_.Exception.Response.StatusCode
    }

    $webRequest = [System.Net.HttpWebRequest]::Create($Url)
    $webRequest.Method = $Method
    $webRequest.ContentType = "application/json"
    foreach ($key in $headersDict.Keys) {
        if ($key -ne "Content-Type") {
            $webRequest.Headers[$key] = $headersDict[$key]
        }
    }

    if ($Data -and ($Method -eq "POST" -or $Method -eq "PUT")) {
        $bytes = [System.Text.Encoding]::UTF8.GetBytes($Data)
        $webRequest.ContentLength = $bytes.Length
        $reqStream = $webRequest.GetRequestStream()
        $reqStream.Write($bytes, 0, $bytes.Length)
        $reqStream.Close()
    }

    try {
        $response = $webRequest.GetResponse()
        $statusCode = [int]$response.StatusCode
        $stream = $response.GetResponseStream()
        $reader = New-Object System.IO.StreamReader($stream)
        $body = $reader.ReadToEnd()
        $reader.Close()
        $response.Close()

        $body | Out-File -FilePath $tempFile -Encoding UTF8
        $statusCode | Out-File -FilePath $statusCodeFile -Encoding UTF8

        return $body
    }
    catch {
        if ($_.Exception.Response) {
            $statusCode = [int]$_.Exception.Response.StatusCode
            $statusCode | Out-File -FilePath $statusCodeFile -Encoding UTF8

            $stream = $_.Exception.Response.GetResponseStream()
            $reader = New-Object System.IO.StreamReader($stream)
            $body = $reader.ReadToEnd()
            $reader.Close()
            $body | Out-File -FilePath $tempFile -Encoding UTF8

            return $body
        }
        return ""
    }
}

function Login {
    param([string]$Username, [string]$Password)

    $resp = Send-Request -Method "POST" -Url "$script:API_V1/auth/login" -Data "{`"username`":`"$Username`",`"password`":`"$Password`"}"
    if ($resp -match '"access_token`":`"([^`"]+)"') {
        return $matches[1]
    }
    return ""
}

function Check-Status {
    param(
        [string]$Method,
        [string]$Url,
        [int]$ExpectedStatus,
        [string]$Data = "",
        [string[]]$Headers = @()
    )

    $headersDict = @{"Content-Type" = "application/json"}
    foreach ($h in $Headers) {
        if ($h -match "(\w+):\s*(.+)") {
            $headersDict[$matches[1]] = $matches[2].Trim()
        }
    }

    $webRequest = [System.Net.HttpWebRequest]::Create($Url)
    $webRequest.Method = $Method
    $webRequest.ContentType = "application/json"
    foreach ($key in $headersDict.Keys) {
        if ($key -ne "Content-Type") {
            $webRequest.Headers[$key] = $headersDict[$key]
        }
    }

    if ($Data -and ($Method -eq "POST" -or $Method -eq "PUT")) {
        $bytes = [System.Text.Encoding]::UTF8.GetBytes($Data)
        $webRequest.ContentLength = $bytes.Length
        $reqStream = $webRequest.GetRequestStream()
        $reqStream.Write($bytes, 0, $bytes.Length)
        $reqStream.Close()
    }

    $tempFile = "$env:TEMP\ps_resp_body_$PID.txt"
    $statusCode = 0

    try {
        $response = $webRequest.GetResponse()
        $statusCode = [int]$response.StatusCode
        $stream = $response.GetResponseStream()
        $reader = New-Object System.IO.StreamReader($stream)
        $body = $reader.ReadToEnd()
        $reader.Close()
        $response.Close()
        $body | Out-File -FilePath $tempFile -Encoding UTF8
    }
    catch {
        if ($_.Exception.Response) {
            $statusCode = [int]$_.Exception.Response.StatusCode
            try {
                $stream = $_.Exception.Response.GetResponseStream()
                $reader = New-Object System.IO.StreamReader($stream)
                $body = $reader.ReadToEnd()
                $reader.Close()
                $body | Out-File -FilePath $tempFile -Encoding UTF8
            }
            catch {
                "" | Out-File -FilePath $tempFile -Encoding UTF8
            }
        }
    }

    if ($statusCode -eq $ExpectedStatus) {
        Log-Info "测试通过: $Method $Url (预期 $ExpectedStatus, 实际 $statusCode)"
        return $true
    }
    else {
        $respContent = Get-Content $tempFile -Raw -ErrorAction SilentlyContinue
        Log-Error "测试失败: $Method $Url (预期 $ExpectedStatus, 实际 $statusCode)"
        Log-Error "响应内容: $respContent"
        return $false
    }
}

Export-ModuleMember -Function Log-Info, Log-Warn, Log-Error, Send-Request, Login, Check-Status
