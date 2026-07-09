#!/bin/bash
echo "[Client] Initiating generation stream request..."

curl -N -X POST -H "Content-Type: application/json" \
     -d '{"prompt": [1, 512, 1024, 42], "max_tokens": 20, "eos_token_id": 2}' \
     http://localhost:8080/v2/generate

echo -e "\n[Client] Stream complete."
