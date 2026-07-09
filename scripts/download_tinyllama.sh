#!/bin/bash
# scripts/download_tinyllama.sh
# Downloads the TinyLlama model and tokenizer files from the Hugging Face Hub

set -e

MODEL_REPO="TinyLlama/TinyLlama-1.1B-Chat-v1.0"
OUT_DIR="models/tinyllama"

mkdir -p "$OUT_DIR"

echo "Downloading tokenizer.json..."
curl -L -o "$OUT_DIR/tokenizer.json" "https://huggingface.co/$MODEL_REPO/resolve/main/tokenizer.json"

echo "Downloading model.safetensors..."
curl -L -o "$OUT_DIR/model.safetensors" "https://huggingface.co/$MODEL_REPO/resolve/main/model.safetensors"

echo "Download complete! Files saved to $OUT_DIR"
