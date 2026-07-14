#!/bin/bash
# ==============================================================
# Call Recording Upload Script
# ==============================================================
# Triggered by Asterisk MixMonitor after call ends.
# Uploads recording to MinIO and publishes event to RabbitMQ.
#
# Usage in ARI (Go backend sets this via MixMonitor):
#   MixMonitor(filename.wav,b,/usr/local/bin/upload_recording.sh ^{UNIQUEID} ^{CALLERID(num)} ^{EXTEN})
# ==============================================================

set -euo pipefail

RECORDING_FILE="$1"
UNIQUE_ID="${2:-unknown}"
CALLER="${3:-unknown}"
CALLEE="${4:-unknown}"

# MinIO configuration
MINIO_ALIAS="crm"
MINIO_ENDPOINT="${MINIO_ENDPOINT:-http://minio:9000}"
MINIO_USER="${MINIO_ROOT_USER:-crm_storage_admin}"
MINIO_PASS="${MINIO_ROOT_PASSWORD:-change-me}"
BUCKET="${MINIO_BUCKET_RECORDINGS:-call-recordings}"

# RabbitMQ configuration
RABBITMQ_URL="${RABBITMQ_URL:-http://rabbitmq:15672}"
RABBITMQ_USER="${RABBITMQ_USER:-crm_broker}"
RABBITMQ_PASS="${RABBITMQ_PASSWORD:-change-me}"

# Date-based path for organization
DATE_PATH=$(date +%Y/%m/%d)
S3_KEY="${DATE_PATH}/${UNIQUE_ID}.wav"

echo "[UPLOAD] Processing recording: ${RECORDING_FILE}"
echo "[UPLOAD] Call: ${CALLER} -> ${CALLEE} (ID: ${UNIQUE_ID})"

# ─── Step 1: Configure MinIO client ──────────────────────────
mc alias set ${MINIO_ALIAS} ${MINIO_ENDPOINT} ${MINIO_USER} ${MINIO_PASS} --api S3v4 2>/dev/null

# Ensure bucket exists
mc mb --ignore-existing ${MINIO_ALIAS}/${BUCKET} 2>/dev/null

# ─── Step 2: Upload recording to MinIO ───────────────────────
if [ -f "${RECORDING_FILE}" ]; then
    mc cp "${RECORDING_FILE}" "${MINIO_ALIAS}/${BUCKET}/${S3_KEY}"
    echo "[UPLOAD] ✅ Uploaded to s3://${BUCKET}/${S3_KEY}"
else
    echo "[UPLOAD] ❌ Recording file not found: ${RECORDING_FILE}"
    exit 1
fi

# ─── Step 3: Publish event to RabbitMQ ────────────────────────
PAYLOAD=$(cat <<EOF
{
    "event": "recording.uploaded",
    "unique_id": "${UNIQUE_ID}",
    "caller": "${CALLER}",
    "callee": "${CALLEE}",
    "s3_bucket": "${BUCKET}",
    "s3_key": "${S3_KEY}",
    "timestamp": "$(date -u +%Y-%m-%dT%H:%M:%SZ)"
}
EOF
)

# Publish via RabbitMQ HTTP API (Management Plugin)
curl -s -u "${RABBITMQ_USER}:${RABBITMQ_PASS}" \
    -H "Content-Type: application/json" \
    -X POST "${RABBITMQ_URL}/api/exchanges/%2f/recording.events/publish" \
    -d "{
        \"properties\": {\"content_type\": \"application/json\"},
        \"routing_key\": \"recording.uploaded\",
        \"payload\": $(echo ${PAYLOAD} | jq -Rs .),
        \"payload_encoding\": \"string\"
    }" && echo "[UPLOAD] ✅ Event published to RabbitMQ" \
       || echo "[UPLOAD] ⚠️  Failed to publish event (non-critical)"

# ─── Step 4: Clean up local file ─────────────────────────────
rm -f "${RECORDING_FILE}"
echo "[UPLOAD] 🗑️  Local file cleaned up"
echo "[UPLOAD] ✅ Done"
