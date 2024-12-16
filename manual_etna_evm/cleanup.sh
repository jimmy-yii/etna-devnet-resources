#!/bin/bash

# This script performs cleanup by:
# 1. Removing all docker containers (node0-node4)
# 2. Recursively deleting the ./data directory while preserving any *_key.txt files
# 3. Restoring the preserved *_key.txt files to a fresh ./data directory

set -euo pipefail

docker stop node0 node1 node2 node3 node4 2>/dev/null || true
echo "- Removed all containers"

mkdir -p data_backup
if mv data/*_key.txt data_backup/ 2>/dev/null; then
  echo "- Moved all *_key.txt files to data_backup"
else
  echo "- No *_key.txt files to move"
fi

sudo rm -rf data/*.txt data/*.json data/chains/
echo "- Removed data directory's *.txt and *.json files keeping node keys and data"

mkdir -p data
if mv data_backup/*_key.txt data/ 2>/dev/null; then
  echo "- Restored all *_key.txt files to data"
else
  echo "- No *_key.txt files to restore"
fi
rm -rf data_backup

echo "✅ Cleanup completed"
