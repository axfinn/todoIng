#!/usr/bin/env bash
set -euo pipefail
API_BASE="http://localhost:5001/api"
EMAIL="admin@example.com"
PASS="admin123"

log(){ echo -e "\n==== $* ====\n"; }

log "Login"
TOKEN=$(curl -s -X POST "$API_BASE/auth/login" \
  -H 'Content-Type: application/json' \
  -d '{"email":"'$EMAIL'","password":"'$PASS'"}' | jq -r '.token')
if [ "$TOKEN" = "null" ] || [ -z "$TOKEN" ]; then echo "Login failed"; exit 1; fi
echo "Token acquired"

AUTH_HDR="Authorization: Bearer $TOKEN"

log "Create task"
CREATE_RES=$(curl -s -X POST "$API_BASE/tasks" -H 'Content-Type: application/json' -H "$AUTH_HDR" \
  -d '{"title":"调试任务1","description":"测试","status":"To Do","priority":"Medium"}')
echo "$CREATE_RES" | jq '.' || true
TASK_ID=$(echo "$CREATE_RES" | jq -r '._id')

log "List tasks"
curl -s -H "$AUTH_HDR" "$API_BASE/tasks" | jq '.[0]' || true

log "Generate report"
START=$(date -u +%Y-%m-%dT00:00:00Z)
END=$(date -u +%Y-%m-%dT23:59:59Z)
REPORT=$(curl -s -X POST "$API_BASE/reports/generate" -H 'Content-Type: application/json' -H "$AUTH_HDR" \
  -d '{"type":"daily","period":"'$(date +%Y-%m-%d)'","startDate":"'$START'","endDate":"'$END'"}')
echo "$REPORT" | jq '{_id,title,statistics}' || true

log "Polish report (placeholder)"
RID=$(echo "$REPORT" | jq -r '._id')
if [ "$RID" != "null" ]; then curl -s -X POST "$API_BASE/reports/$RID/polish" -H "$AUTH_HDR" | jq '{_id,polishedContent: (.polishedContent|length)}'; fi

log "Done"
