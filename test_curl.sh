#!/bin/sh
curl -N -s -X POST -H "Content-Type: application/json" -d '{
  "prompt": "<|system|>\nYou are a helpful assistant.</s>\n<|user|>\nWhat is the capital of France?</s>\n<|assistant|>\n",
  "max_tokens": 15,
  "eos_token_id": 2
}' http://localhost:8080/v2/generate
