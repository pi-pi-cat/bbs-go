#!/bin/sh
set -eu

cd "$(dirname "$0")"

./load-images.sh

mkdir -p docker-data/mysql docker-data/data docker-data/logs docker-data/uploads

docker compose -f docker-compose.yaml up -d

echo "bbs-go is starting. Open http://127.0.0.1:${BBSGO_WEB_PORT:-3000}"
