@echo off

REM 启动后端服务 cd backend\cmd ; go run main.go
start "后端服务" cmd /k "cd backend\cmd && go run main.go"

REM 等待后端服务初始化
timeout /t 3 /nobreak >nul

REM 启动前端服务 cd frontend ; npm run dev
start "前端服务" cmd /k "cd frontend && npm run dev"

echo 前后端服务已启动，请查看对应的命令窗口
pause
