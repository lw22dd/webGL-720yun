#!/bin/bash
BASE_URL="http://localhost:7000"
API_V1="$BASE_URL/api/v1"

echo "=== Upload and Slice Module Test ==="

echo "1. Login..."
LOGIN_RESP=$(curl -s "$API_V1/auth/login" -X POST -H "Content-Type: application/json" -d '{"username":"admin","password":"admin123"}')
TOKEN=$(echo $LOGIN_RESP | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)
echo "Token: ${TOKEN:0:50}..."

echo ""
echo "2. Get Spaces..."
SPACES_RESP=$(curl -s "$API_V1/resource/spaces?page=1" -H "Authorization: Bearer $TOKEN")
SPACE_ID=$(echo $SPACES_RESP | grep -o '"id":[0-9]*' | head -1 | cut -d':' -f2)
echo "Space ID: $SPACE_ID"

echo ""
echo "3. Test Init Upload..."
INIT_RESP=$(curl -s "$API_V1/upload/init" -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $TOKEN" -d '{"filename":"test.jpg","file_size":1048576,"file_hash":"abc123def45678901234567890123456"}')
echo "Init: $INIT_RESP"

echo ""
echo "Test completed!"
