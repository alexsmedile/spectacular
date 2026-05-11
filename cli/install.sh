#!/usr/bin/env bash
set -euo pipefail

GITHUB_REPO="alexsmedile/spectacular"
INSTALL_DIR="${HOME}/.local/bin"
BINARY_NAME="spectacular"
SCRIPT_URL="https://raw.githubusercontent.com/${GITHUB_REPO}/main/cli/spectacular"

die()  { echo "Error: $*" >&2; exit 1; }
info() { echo "$*"; }

install_binary() {
  mkdir -p "$INSTALL_DIR"

  info "Downloading ${BINARY_NAME}..."

  if command -v curl &>/dev/null; then
    curl -fsSL "$SCRIPT_URL" -o "${INSTALL_DIR}/${BINARY_NAME}" \
      || die "Could not download ${BINARY_NAME} from GitHub. Check your network connection."
  elif command -v wget &>/dev/null; then
    wget -q "$SCRIPT_URL" -O "${INSTALL_DIR}/${BINARY_NAME}" \
      || die "Could not download ${BINARY_NAME} from GitHub. Check your network connection."
  else
    die "Neither curl nor wget found. Install one to proceed."
  fi

  chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
  info "Installed: ${INSTALL_DIR}/${BINARY_NAME}"
}

check_path() {
  if ! echo "$PATH" | tr ':' '\n' | grep -qx "$INSTALL_DIR"; then
    info ""
    info "Note: ${INSTALL_DIR} is not in your PATH."
    info "Add this to your shell profile (~/.zshrc or ~/.bashrc):"
    info ""
    info "  export PATH=\"\$HOME/.local/bin:\$PATH\""
    info ""
  fi
}

main() {
  install_binary
  check_path
  info ""
  info "Done. Run 'spectacular init' inside any project to get started."
}

main "$@"
