Add-Type -AssemblyName System.Drawing

function New-TestPanorama {
    param(
        [string]$OutputPath,
        [int]$Width = 4096,
        [int]$Height = 2048
    )
    
    $bitmap = New-Object System.Drawing.Bitmap($Width, $Height)
    $graphics = [System.Drawing.Graphics]::FromImage($bitmap)
    
    $brush = New-Object System.Drawing.SolidBrush([System.Drawing.Color]::SkyBlue)
    $graphics.FillRectangle($brush, 0, 0, $Width, $Height)
    
    $rand = New-Object Random
    for ($i = 0; $i -lt 100; $i++) {
        $x = $rand.Next($Width)
        $y = $rand.Next($Height)
        $size = $rand.Next(10, 50)
        $color = [System.Drawing.Color]::FromArgb($rand.Next(256), $rand.Next(256), $rand.Next(256))
        $cloudBrush = New-Object System.Drawing.SolidBrush($color)
        $graphics.FillEllipse($cloudBrush, $x, $y, $size, $size / 2)
    }
    
    $font = New-Object System.Drawing.Font("Arial", 48, [System.Drawing.FontStyle]::Bold)
    $textBrush = New-Object System.Drawing.SolidBrush([System.Drawing.Color]::White)
    $format = New-Object System.Drawing.StringFormat
    $format.Alignment = [System.Drawing.StringAlignment]::Center
    $format.LineAlignment = [System.Drawing.StringAlignment]::Center
    $graphics.DrawString("TEST PANORAMA $Width x $Height", $font, $textBrush, $Width / 2, $Height / 2, $format)
    
    $bitmap.Save($OutputPath, [System.Drawing.Imaging.ImageFormat]::Jpeg)
    
    $graphics.Dispose()
    $bitmap.Dispose()
    
    Write-Host "Created test panorama: $OutputPath ($Width x $Height)"
}

$outputDir = "$PSScriptRoot\test_data"
if (-not (Test-Path $outputDir)) {
    New-Item -ItemType Directory -Path $outputDir -Force | Out-Null
}

New-TestPanorama -OutputPath "$outputDir\test_panorama_4k.jpg" -Width 4096 -Height 2048
New-TestPanorama -OutputPath "$outputDir\test_panorama_2k.jpg" -Width 2048 -Height 1024

Write-Host "Test panoramas created successfully!"
