#!/bin/sh
set -e

REPO="made-purple/clog"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

# Detect OS
OS="$(uname -s)"
case "$OS" in
  Linux)  OS="linux" ;;
  Darwin) OS="darwin" ;;
  *)
    echo "Error: unsupported OS: $OS" >&2
    exit 1
    ;;
esac

# Detect architecture
ARCH="$(uname -m)"
case "$ARCH" in
  x86_64|amd64)  ARCH="amd64" ;;
  arm64|aarch64) ARCH="arm64" ;;
  *)
    echo "Error: unsupported architecture: $ARCH" >&2
    exit 1
    ;;
esac

# Check for unsupported combinations
if [ "$OS" = "linux" ] && [ "$ARCH" = "arm64" ]; then
  echo "Error: linux/arm64 is not currently supported" >&2
  exit 1
fi

# Get latest version
echo "Fetching latest release..."
VERSION="$(curl -sSf "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"v?([^"]+)".*/\1/')"
if [ -z "$VERSION" ]; then
  echo "Error: could not determine latest version" >&2
  exit 1
fi
echo "Latest version: v${VERSION}"

# Download
FILENAME="clog_${VERSION}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/${REPO}/releases/download/v${VERSION}/${FILENAME}"

TMPDIR="$(mktemp -d)"
trap 'rm -rf "$TMPDIR"' EXIT

echo "Downloading ${URL}..."
curl -sSfL -o "${TMPDIR}/${FILENAME}" "$URL"

# Extract
tar -xzf "${TMPDIR}/${FILENAME}" -C "$TMPDIR"

# Install
if [ -w "$INSTALL_DIR" ]; then
  mv "${TMPDIR}/clog" "${INSTALL_DIR}/clog"
else
  echo "Installing to ${INSTALL_DIR} (requires sudo)..."
  sudo mv "${TMPDIR}/clog" "${INSTALL_DIR}/clog"
fi

chmod +x "${INSTALL_DIR}/clog"
echo "clog v${VERSION} installed to ${INSTALL_DIR}/clog"

# Install shell completions
CLOG="${INSTALL_DIR}/clog"
SHELL_NAME="$(basename "${SHELL:-}")"

install_completions() {
  case "$SHELL_NAME" in
    bash)
      COMP_DIR="${BASH_COMPLETION_USER_DIR:-${XDG_DATA_HOME:-$HOME/.local/share}/bash-completion/completions}"
      mkdir -p "$COMP_DIR"
      "$CLOG" completion bash > "${COMP_DIR}/clog"
      echo "Bash completions installed to ${COMP_DIR}/clog"
      ;;
    zsh)
      # Use ~/.zsh/completions and ensure it's in fpath via .zshrc
      COMP_DIR="${HOME}/.zsh/completions"
      mkdir -p "$COMP_DIR"
      "$CLOG" completion zsh > "${COMP_DIR}/_clog"
      # Add to fpath if not already present
      if [ -f "${HOME}/.zshrc" ]; then
        if ! grep -q '\.zsh/completions' "${HOME}/.zshrc" 2>/dev/null; then
          printf '\n# clog completions\nfpath=(~/.zsh/completions $fpath)\nautoload -Uz compinit && compinit\n' >> "${HOME}/.zshrc"
          echo "Added completion setup to ~/.zshrc"
        fi
      fi
      echo "Zsh completions installed to ${COMP_DIR}/_clog"
      ;;
    fish)
      COMP_DIR="${XDG_CONFIG_HOME:-$HOME/.config}/fish/completions"
      mkdir -p "$COMP_DIR"
      "$CLOG" completion fish > "${COMP_DIR}/clog.fish"
      echo "Fish completions installed to ${COMP_DIR}/clog.fish"
      ;;
    *)
      echo "Shell completions: could not detect shell (SHELL=$SHELL). Run 'clog completion --help' to install manually."
      return
      ;;
  esac
}

install_completions
