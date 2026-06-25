#!/bin/sh
set -eu

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
PACKAGE_DIR="$ROOT_DIR/dist/bbs-go-docker-offline"
IMAGES_DIR="$PACKAGE_DIR/images"

BBSGO_IMAGE="${BBSGO_IMAGE:-bbs-go:offline}"
MYSQL_IMAGE="${MYSQL_IMAGE:-mysql:8.4}"
PLATFORM="${DOCKER_DEFAULT_PLATFORM:-linux/amd64}"
RETRY_COUNT="${BBSGO_PACKAGE_RETRY_COUNT:-3}"
SKIP_BUILD="${BBSGO_SKIP_BUILD:-false}"

cd "$ROOT_DIR"

retry() {
  attempt=1
  while :; do
    if "$@"; then
      return 0
    fi

    if [ "$attempt" -ge "$RETRY_COUNT" ]; then
      return 1
    fi

    attempt=$((attempt + 1))
    echo "Command failed; retrying ($attempt/$RETRY_COUNT): $*"
    sleep 5
  done
}

if [ "$SKIP_BUILD" = "true" ]; then
  echo "Skipping build; using existing $BBSGO_IMAGE"
  docker image inspect "$BBSGO_IMAGE" >/dev/null
else
  echo "Building $BBSGO_IMAGE for $PLATFORM"
  retry docker build --pull=false --platform "$PLATFORM" -t "$BBSGO_IMAGE" .
fi

echo "Ensuring $MYSQL_IMAGE exists locally"
if ! docker image inspect "$MYSQL_IMAGE" >/dev/null 2>&1; then
  retry docker pull --platform "$PLATFORM" "$MYSQL_IMAGE"
fi

rm -rf "$PACKAGE_DIR"
mkdir -p "$IMAGES_DIR"

echo "Saving Docker images"
docker save -o "$IMAGES_DIR/bbs-go-offline.tar" "$BBSGO_IMAGE"
docker save -o "$IMAGES_DIR/mysql-8.4.tar" "$MYSQL_IMAGE"
chmod 644 "$IMAGES_DIR/bbs-go-offline.tar" "$IMAGES_DIR/mysql-8.4.tar"

cp deploy/docker-offline/docker-compose.yaml "$PACKAGE_DIR/docker-compose.yaml"
cp deploy/docker-offline/load-images.sh "$PACKAGE_DIR/load-images.sh"
cp deploy/docker-offline/deploy.sh "$PACKAGE_DIR/deploy.sh"
cp deploy/docker-offline/README.md "$PACKAGE_DIR/README.md"

chmod +x "$PACKAGE_DIR/load-images.sh" "$PACKAGE_DIR/deploy.sh"

echo "Offline Docker package created at: $PACKAGE_DIR"
