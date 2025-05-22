#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

REPO="alcharra/docker-deploy-action-go"
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH_RAW="$(uname -m)"
ARCH="$ARCH_RAW"
VARIANT=""

case "$ARCH_RAW" in
  x86_64 | amd64) ARCH="amd64" ;;
  i386 | i686) ARCH="386" ;;
  armv5* | armv5l) ARCH="arm"; VARIANT="v5" ;;
  armv6* | armv6l) ARCH="arm"; VARIANT="v6" ;;
  armv7* | armv7l | armhf) ARCH="arm"; VARIANT="v7" ;;
  armv8* | aarch64) ARCH="arm64" ;;
  *) echo "❌ ERROR: Unsupported architecture: $ARCH_RAW"; exit 1 ;;
esac

FILENAME="deploy-action_${OS}_${ARCH}"
[ -n "$VARIANT" ] && FILENAME="${FILENAME}${VARIANT}"
FILENAME="${FILENAME}_${RELEASE_VERSION}"
[ "$OS" = "windows" ] && FILENAME="${FILENAME}.exe"

TARGET="${GITHUB_ACTION_PATH}/${FILENAME}"

if ! echo "$RELEASE_VERSION" | grep -Eq '^v[0-9]+(\.[0-9]+)*$'; then
  echo "❌ No valid release version provided (got '$RELEASE_VERSION'). Must be a tag like v1.2.3"
  exit 1
fi

echo "📦 Downloading release binary: $FILENAME"
URL="https://github.com/${REPO}/releases/download/${RELEASE_VERSION}/${FILENAME}"
curl -fsSL --retry 5 --keepalive-time 2 "$URL" -o "$TARGET"

chmod +x "$TARGET"
"$TARGET" "$@"
rm -f "$TARGET"