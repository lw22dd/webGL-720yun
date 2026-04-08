#!/bin/bash
source "$(dirname "$0")/common.sh"

log_info "#######################################"
log_info "#                                     #"
log_info "#      WebGL-720yun API 全模块测试      #"
log_info "#                                     #"
log_info "#######################################"

# 赋予执行权限
chmod +x "$(dirname "$0")"/*.sh

# 执行各模块测试
"$(dirname "$0")/auth.sh"
"$(dirname "$0")/user.sh"
"$(dirname "$0")/admin.sh"
"$(dirname "$0")/upload_test.sh"

log_info "#######################################"
log_info "#            全部测试完成             #"
log_info "#######################################"
