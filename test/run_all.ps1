Import-Module "$PSScriptRoot\common.ps1" -Force

Log-Info "#######################################"
Log-Info "#                                     #"
Log-Info "#      WebGL-720yun API 全模块测试      #"
Log-Info "#                                     #"
Log-Info "#######################################"

$scriptDir = $PSScriptRoot
& "$scriptDir\auth_test.ps1"
& "$scriptDir\user_test.ps1"
& "$scriptDir\admin_test.ps1"
& "$scriptDir\upload_test.ps1"

Log-Info "#######################################"
Log-Info "#            全部测试完成             #"
Log-Info "#######################################"
