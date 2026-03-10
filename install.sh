#!/bin/sh
# install.sh — One-line installer for PortView
# Usage: curl -sSL https://raw.githubusercontent.com/Vinay-Madarkhandi/portview/main/install.sh | bash
#
# Installs the latest release binary to /usr/local/bin (or ~/.local/bin as fallback).

set -e

REPO="Vinay-Madarkhandi/portview"
BINARY="portview"
RELEASES_URL="https://github.com/${REPO}/releases/latest/download"

# --- Helpers ---

info()  { printf "\033[1;34m→\033[0m %s\n" "$1"; }
ok()    { printf "\033[1;32m✓\033[0m %s\n" "$1"; }
err()   { printf "\033[1;31m✗ ERROR:\033[0m %s\n" "$1" >&2; exit 1; }

need_cmd() {
    if ! command -v "$1" > /dev/null 2>&1; then
        err "Required command '$1' not found. Please install it and try again."
    fi
}

# --- Detect OS ---

detect_os() {
    case "$(uname -s)" in
        Linux*)  echo "linux"  ;;
        Darwin*) echo "darwin" ;;
        *)       err "Unsupported operating system: $(uname -s). PortView supports Linux and macOS." ;;
    esac
}

# --- Detect Architecture ---

detect_arch() {
    case "$(uname -m)" in
        x86_64|amd64)   echo "amd64" ;;
        aarch64|arm64)  echo "arm64" ;;
        *)              err "Unsupported architecture: $(uname -m). PortView supports amd64 and arm64." ;;
    esac
}

# --- Determine install directory ---

detect_install_dir() {
    if [ -w /usr/local/bin ]; then
        echo "/usr/local/bin"
    elif command -v sudo > /dev/null 2>&1; then
        echo "/usr/local/bin"
    else
        mkdir -p "${HOME}/.local/bin"
        echo "${HOME}/.local/bin"
    fi
}

# --- Main ---

main() {
    echo ""
    info "PortView Installer"
    echo ""

    need_cmd curl
    need_cmd uname

    OS=$(detect_os)
    ARCH=$(detect_arch)
    INSTALL_DIR=$(detect_install_dir)

    ASSET="${BINARY}-${OS}-${ARCH}"
    DOWNLOAD_URL="${RELEASES_URL}/${ASSET}"

    info "Detected: ${OS}/${ARCH}"
    info "Downloading ${ASSET}..."

    TMPDIR=$(mktemp -d)
    trap 'rm -rf "${TMPDIR}"' EXIT

    HTTP_CODE=$(curl -sL -o "${TMPDIR}/${BINARY}" -w "%{http_code}" "${DOWNLOAD_URL}")

    if [ "${HTTP_CODE}" != "200" ]; then
        err "Download failed (HTTP ${HTTP_CODE}). Check that a release exists for ${OS}/${ARCH} at:\n  ${DOWNLOAD_URL}"
    fi

    chmod +x "${TMPDIR}/${BINARY}"

    info "Installing to ${INSTALL_DIR}/${BINARY}..."

    if [ "${INSTALL_DIR}" = "/usr/local/bin" ] && [ ! -w /usr/local/bin ]; then
        sudo mv "${TMPDIR}/${BINARY}" "${INSTALL_DIR}/${BINARY}"
    else
        mv "${TMPDIR}/${BINARY}" "${INSTALL_DIR}/${BINARY}"
    fi

    ok "PortView installed successfully!"
    echo ""
    info "Run it with:"
    echo "    ${BINARY}"
    echo "    sudo ${BINARY}   # recommended, for full process info"
    echo ""

    # Warn if install dir is not in PATH
    case ":${PATH}:" in
        *":${INSTALL_DIR}:"*) ;;
        *)
            echo "  ⚠  ${INSTALL_DIR} is not in your PATH."
            echo "  Add it with:"
            echo "      export PATH=\"${INSTALL_DIR}:\${PATH}\""
            echo ""
            ;;
    esac
}

main

