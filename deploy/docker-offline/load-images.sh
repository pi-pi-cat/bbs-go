#!/bin/sh
set -eu

cd "$(dirname "$0")"

if [ ! -d images ]; then
  echo "images directory not found" >&2
  exit 1
fi

for image_tar in images/*.tar; do
  if [ ! -f "$image_tar" ]; then
    echo "no image tar files found in images/" >&2
    exit 1
  fi

  echo "Loading $image_tar"
  docker load -i "$image_tar"
done

echo "Docker images loaded."
