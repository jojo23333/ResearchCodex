#!/usr/bin/env bash
# install.sh - Build and install rcodex without requiring root access.
# Re-running this script updates to the latest code automatically.

set -euo pipefail

REPO_URL="${RCODEX_REPO_URL:-https://github.com/jojo/ResearchCodex.git}"
REF="${RCODEX_REF:-main}"
CACHE_DIR="${RCODEX_CACHE_DIR:-$HOME/.cache/rcodex-src}"
INSTALL_DIR="${RCODEX_INSTALL_DIR:-$HOME/.local/bin}"
BINARY_NAME="${RCODEX_BINARY_NAME:-rcodex}"

log() {
  printf '[rcodex-install] %s\n' "$*" >&2
}

ensure_go() {
  if ! command -v go >/dev/null 2>&1; then
    log "Go toolchain not found. Please install Go 1.21+ and re-run."
    exit 1
  fi
}

clone_or_update_repo() {
  if [ -d "$CACHE_DIR/.git" ]; then
    log "Updating existing repository in $CACHE_DIR"
    git -C "$CACHE_DIR" fetch --tags --prune origin
    git -C "$CACHE_DIR" checkout "$REF"
    git -C "$CACHE_DIR" reset --hard "origin/$REF"
  else
    log "Cloning rcodex into $CACHE_DIR"
    mkdir -p "$(dirname "$CACHE_DIR")"
    git clone --branch "$REF" "$REPO_URL" "$CACHE_DIR"
  fi
}

build_and_install() {
  mkdir -p "$INSTALL_DIR"
  log "Building rcodex"
  (cd "$CACHE_DIR" && go build -o "$INSTALL_DIR/$BINARY_NAME" ./cmd/rcodex)
  chmod +x "$INSTALL_DIR/$BINARY_NAME"
  log "rcodex installed to $INSTALL_DIR/$BINARY_NAME"
}

post_install_hint() {
  case ":$PATH:" in
    *":$INSTALL_DIR:"*) ;;
    *)
      log "Add the following to your shell profile to use rcodex:"
      log "  export PATH=\"$INSTALL_DIR:\$PATH\""
      ;;
  esac
}

ensure_go
clone_or_update_repo
build_and_install
post_install_hint

