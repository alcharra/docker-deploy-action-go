#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

REPO="alcharra/docker-deploy-action-go"
LATEST_VERSION="v2.1.0"
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH_RAW="$(uname -m)"
ARCH="$ARCH_RAW"
VARIANT=""

case "$ARCH_RAW" in
  x86_64 | amd64) ARCH="amd64" ;;
  i386 | i686) ARCH="386" ;;
  armv5* | armv5l) ARCH="arm"; VARIANT="5" ;;
  armv6* | armv6l) ARCH="arm"; VARIANT="6" ;;
  armv7* | armv7l | armhf) ARCH="arm"; VARIANT="7" ;;
  armv8* | aarch64) ARCH="arm64" ;;
  *) echo "❌ ERROR: Unsupported architecture: $ARCH_RAW"; exit 1 ;;
esac

VERSION_NO_V="${LATEST_VERSION#v}"
FILENAME="docker-deploy-action-go-${VERSION_NO_V}-${OS}-${ARCH}"
[ -n "$VARIANT" ] && FILENAME="${FILENAME}-${VARIANT}"
ARCHIVE="${FILENAME}.tar.gz"

URL="https://github.com/${REPO}/releases/download/${LATEST_VERSION}/${ARCHIVE}"

echo "📦 Downloading from: $URL"
curl -fsSL --retry 5 "$URL" -o "$ARCHIVE"

echo -e "\U0001F9F0 Extracting binary..."
tar -xzf "$ARCHIVE"
chmod +x docker-deploy-action-go*

echo -e "\U0001F680 Running docker-deploy-action-go version $LATEST_VERSION..."
./docker-deploy-action-go* "$@"

rm "$ARCHIVE"
rm docker-deploy-action-go*