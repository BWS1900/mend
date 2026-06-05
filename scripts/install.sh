#!/usr/bin/env bash
set -euo pipefail

# mend installer
# Usage: curl -sSfL https://raw.githubusercontent.com/BWS1900/mend/main/scripts/install.sh | sh
#
# Environment variables:
#   MEND_VERSION   - version to install (default: latest)
#   MEND_INSTALL   - install location (default: /usr/local/bin or ~/.local/bin)
#   MEND_GO_INSTALL - if set, use `go install` instead of downloading a binary

require() {
  command -v "$1" >/dev/null 2>&1 || { echo "mend-install: $1 required" >&2; exit 1; }
}

if [[ -n "${MEND_GO_INSTALL:-}" ]]; then
  require go
  echo "Installing mend with go install..."
  go install "github.com/BWS1900/mend@${MEND_VERSION:-latest}"
  exit 0
fi

require curl
require tar

repo="BWS1900/mend"
version="${MEND_VERSION:-latest}"

case "$(uname -s)" in
  Linux)  os="linux" ;;
  Darwin) os="darwin" ;;
  *) echo "mend-install: unsupported OS: $(uname -s)" >&2; exit 1 ;;
esac

case "$(uname -m)" in
  x86_64|amd64) arch="amd64" ;;
  arm64|aarch64) arch="arm64" ;;
  *) echo "mend-install: unsupported arch: $(uname -m)" >&2; exit 1 ;;
esac

if [[ "$version" == "latest" ]]; then
  url="https://github.com/${repo}/releases/latest/download/mend_${version}_${os}_${arch}.tar.gz"
else
  base="https://github.com/${repo}/releases/download/${version}"
  url="${base}/mend_${version}_${os}_${arch}.tar.gz"
fi

if [[ -z "${MEND_INSTALL:-}" ]]; then
  if [[ -w "/usr/local/bin" ]]; then
    MEND_INSTALL="/usr/local/bin"
  else
    MEND_INSTALL="${HOME}/.local/bin"
    mkdir -p "$MEND_INSTALL"
  fi
fi

tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT

echo "Downloading mend $version ($os/$arch)..."
curl -sSfL "$url" -o "$tmp/mend.tgz"
tar -xzf "$tmp/mend.tgz" -C "$tmp"
install -m 0755 "$tmp/mend" "$MEND_INSTALL/mend"

echo "Installed: $MEND_INSTALL/mend"
"$MEND_INSTALL/mend" --version
