#!/usr/bin/env bash
set -e

REPO="tofunmiadewuyi/summon"
PROGRAM="summon"

OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

if [ "$OS" != "darwin" ]; then
  echo "$PROGRAM is macOS only"
  exit 1
fi

case "$ARCH" in
  x86_64) ARCH="amd64" ;;
  arm64|aarch64) ARCH="arm64" ;;
  *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

VERSION=$(curl -s https://api.github.com/repos/$REPO/releases/latest | grep tag_name | cut -d '"' -f4)

FILENAME="${PROGRAM}_${VERSION}_${OS}_${ARCH}.zip"
URL="https://github.com/$REPO/releases/download/$VERSION/$FILENAME"

echo "Installing ${PROGRAM} $VERSION for $OS/$ARCH..."

TMP_DIR=$(mktemp -d)
curl -fsSL "$URL" -o "$TMP_DIR/$FILENAME"
unzip -q "$TMP_DIR/$FILENAME" -d "$TMP_DIR"
chmod +x "$TMP_DIR/${PROGRAM}_${VERSION}_${OS}_${ARCH}"
mv "$TMP_DIR/${PROGRAM}_${VERSION}_${OS}_${ARCH}" "$TMP_DIR/$PROGRAM"

INSTALL_DIR="/usr/local/bin"

if [ ! -w "$INSTALL_DIR" ]; then
  INSTALL_DIR="$HOME/.local/bin"
  mkdir -p "$INSTALL_DIR"
  echo "Installing to $INSTALL_DIR (no sudo access)"

  SHELL_RC=""
  if [ -f "$HOME/.zshrc" ]; then
    SHELL_RC="$HOME/.zshrc"
  elif [ -f "$HOME/.bashrc" ]; then
    SHELL_RC="$HOME/.bashrc"
  elif [ -f "$HOME/.profile" ]; then
    SHELL_RC="$HOME/.profile"
  fi

  if [ -n "$SHELL_RC" ]; then
    if ! grep -q 'HOME/.local/bin' "$SHELL_RC" 2>/dev/null; then
      echo '' >> "$SHELL_RC"
      echo 'export PATH="$HOME/.local/bin:$PATH"' >> "$SHELL_RC"
      echo "Added ~/.local/bin to PATH in $SHELL_RC"
    fi
  fi
fi

mv "$TMP_DIR/$PROGRAM" "$INSTALL_DIR/$PROGRAM"

echo "Installed to $INSTALL_DIR/$PROGRAM"

if [ "$INSTALL_DIR" = "$HOME/.local/bin" ]; then
  echo "Reload your shell or run: export PATH=\"\$HOME/.local/bin:\$PATH\""
fi

echo "Run: $PROGRAM"
