#!/bin/sh
set -e

# Official installer script for Whiskey
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/sxijyoti/whiskey/main/scripts/install.sh | sh
#   curl -fsSL https://raw.githubusercontent.com/sxijyoti/whiskey/main/scripts/install.sh | sh -s -- v0.1.0

REPO="sxijyoti/whiskey"
INSTALL_DIR="${BIN_DIR:-$HOME/.local/bin}"

log_info() {
    printf '\033[0;34m[info]\033[0m %s\n' "$1"
}

log_success() {
    printf '\033[0;32m✔\033[0m %s\n' "$1"
}

log_error() {
    printf '\033[0;31m[error]\033[0m %s\n' "$1" >&2
}

# 1. OS Detection
OS_RAW="$(uname -s)"
case "$OS_RAW" in
    Linux|linux)
        OS="linux"
        OS_CAP="Linux"
        ;;
    Darwin|darwin)
        OS="darwin"
        OS_CAP="Darwin"
        ;;
    *)
        log_error "Unsupported operating system: $OS_RAW"
        exit 1
        ;;
esac

# 2. Architecture Detection
ARCH_RAW="$(uname -m)"
case "$ARCH_RAW" in
    x86_64|amd64)
        ARCH="amd64"
        ARCH_ALT="x86_64"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ARCH_ALT="arm64"
        ;;
    *)
        log_error "Unsupported architecture: $ARCH_RAW"
        exit 1
        ;;
esac

# Helper: Fetch Release Metadata JSON.
# Exits with an error message on curl failure rather than continuing silently.
fetch_release_json() {
    _tag="$1"
    if [ -n "$_tag" ]; then
        _api_url="https://api.github.com/repos/${REPO}/releases/tags/${_tag}"
    else
        _api_url="https://api.github.com/repos/${REPO}/releases/latest"
    fi
    _body=$(curl -fsSL "$_api_url" 2>&1) || {
        log_error "Failed to fetch release metadata from ${_api_url}"
        log_error "curl output: ${_body}"
        exit 1
    }
    printf '%s' "$_body"
}

# Helper: Parse Release Tag from API JSON.
# Uses jq when available; falls back to grep+sed with BRE (POSIX-compatible).
parse_tag_name() {
    _json="$1"
    if command -v jq >/dev/null 2>&1; then
        printf '%s' "$_json" | jq -r '.tag_name // empty' 2>/dev/null || true
    else
        # BRE sed is POSIX; -E (ERE) is not — use BRE syntax here.
        printf '%s' "$_json" \
            | grep '"tag_name":' \
            | head -n 1 \
            | sed 's/.*"tag_name": *"\([^"]*\)".*/\1/'
    fi
}

# Helper: Parse Asset Download URL from API JSON.
# Matches on OS and architecture; returns the first .tar.gz download URL.
parse_asset_url() {
    _json="$1"
    _os="$2"
    _arch="$3"
    _arch_alt="$4"

    if command -v jq >/dev/null 2>&1; then
        printf '%s' "$_json" | jq -r \
            --arg os "$_os" \
            --arg arch "$_arch" \
            --arg arch_alt "$_arch_alt" \
            '.assets[]?
            | select(
                (.name | ascii_downcase | contains($os)) and
                ((.name | ascii_downcase | contains($arch)) or (.name | ascii_downcase | contains($arch_alt))) and
                (.name | endswith(".tar.gz"))
            )
            | .browser_download_url' 2>/dev/null | head -n 1 || true
    else
        # POSIX fallback: extract all browser_download_url values, then filter
        # by OS and architecture using a case statement.
        # Avoid the `while ... done | head` subshell pattern (breaks return/set -e).
        _urls=$(printf '%s' "$_json" \
            | grep '"browser_download_url":' \
            | sed 's/.*"browser_download_url": *"\([^"]*\)".*/\1/')

        _result=""
        while IFS= read -r _url; do
            # Lowercase comparison via case — no tr or subshell needed.
            case "$_url" in
                *"${_os}"*".tar.gz"|*"${OS_CAP}"*".tar.gz")
                    case "$_url" in
                        *"${_arch}"*|*"${_arch_alt}"*)
                            _result="$_url"
                            break
                            ;;
                    esac
                    ;;
            esac
        done <<EOF
$_urls
EOF
        printf '%s' "$_result"
    fi
}

# 3. Determine Version
VERSION_INPUT="$1"
RELEASE_JSON=""

if [ -n "$VERSION_INPUT" ]; then
    TAG_VERSION="$VERSION_INPUT"
    case "$TAG_VERSION" in
        v*) RAW_VERSION="${TAG_VERSION#v}" ;;
        *)  RAW_VERSION="$TAG_VERSION"
            TAG_VERSION="v$RAW_VERSION"
            ;;
    esac
    RELEASE_JSON=$(fetch_release_json "$TAG_VERSION")
else
    RELEASE_JSON=$(fetch_release_json "")
    TAG_VERSION=$(parse_tag_name "$RELEASE_JSON")

    if [ -z "$TAG_VERSION" ]; then
        log_error "Unable to detect latest release version. Please specify a version explicitly, e.g.: sh install.sh v0.1.0"
        exit 1
    fi
    RAW_VERSION="${TAG_VERSION#v}"
fi

log_info "Detected platform: ${OS}/${ARCH}"
log_info "Target version: ${TAG_VERSION}"

# 4. Existing Installation Detection
EXISTING_BIN=""
if [ -x "${INSTALL_DIR}/whiskey" ]; then
    EXISTING_BIN="${INSTALL_DIR}/whiskey"
elif command -v whiskey >/dev/null 2>&1; then
    EXISTING_BIN="$(command -v whiskey)"
fi

if [ -n "$EXISTING_BIN" ]; then
    INSTALLED_VER_RAW=$("$EXISTING_BIN" version 2>/dev/null || "$EXISTING_BIN" --version 2>/dev/null || true)
    INSTALLED_VER=$(printf '%s' "$INSTALLED_VER_RAW" | sed 's/[Ww]hiskey //g' | tr -d ' \n\r')

    if [ "$EXISTING_BIN" = "${INSTALL_DIR}/whiskey" ] && {
        [ "$INSTALLED_VER" = "$TAG_VERSION" ] || \
        [ "$INSTALLED_VER" = "$RAW_VERSION" ] || \
        [ "v$INSTALLED_VER" = "$TAG_VERSION" ]
    }; then
        log_info "Whiskey ${TAG_VERSION} is already installed at ${EXISTING_BIN}."
        exit 0
    fi

    if [ -n "$INSTALLED_VER" ]; then
        log_info "Existing installation detected (${INSTALLED_VER} at ${EXISTING_BIN})."
        log_info "Updating Whiskey from ${INSTALLED_VER} to ${TAG_VERSION}..."
    else
        log_info "Existing installation detected at ${EXISTING_BIN}."
        log_info "Updating Whiskey to ${TAG_VERSION}..."
    fi
else
    log_info "Installing Whiskey ${TAG_VERSION}..."
fi

# 5. Download Release Asset
TMP_DIR=$(mktemp -d 2>/dev/null || mktemp -d -t 'whiskey-install')
cleanup() {
    rm -rf "$TMP_DIR"
}
trap cleanup EXIT INT TERM

DOWNLOAD_URL=$(parse_asset_url "$RELEASE_JSON" "$OS" "$ARCH" "$ARCH_ALT")

if [ -z "$DOWNLOAD_URL" ]; then
    log_error "No matching release asset found for ${TAG_VERSION} (${OS}/${ARCH})."
    log_error "Please check available assets at https://github.com/${REPO}/releases/tag/${TAG_VERSION}"
    exit 1
fi

ARCHIVE_PATH="${TMP_DIR}/whiskey.tar.gz"
log_info "Downloading from: ${DOWNLOAD_URL}"
if ! curl -fsSL -o "$ARCHIVE_PATH" "$DOWNLOAD_URL"; then
    log_error "Failed to download release archive from ${DOWNLOAD_URL}."
    exit 1
fi

# 6. Extract & Install
log_info "Extracting release archive..."
tar -xzf "$ARCHIVE_PATH" -C "$TMP_DIR"

if [ ! -f "${TMP_DIR}/whiskey" ]; then
    log_error "Extracted archive did not contain 'whiskey' binary."
    exit 1
fi

mkdir -p "$INSTALL_DIR"
DEST_PATH="${INSTALL_DIR}/whiskey"
mv "${TMP_DIR}/whiskey" "$DEST_PATH"
chmod +x "$DEST_PATH"

# 7. Verify & Print Summary
if "$DEST_PATH" version >/dev/null 2>&1 || \
   "$DEST_PATH" --version >/dev/null 2>&1 || \
   "$DEST_PATH" --help >/dev/null 2>&1; then
    log_success "Whiskey ${TAG_VERSION} installed successfully."
else
    log_error "Installation completed, but the Whiskey binary failed to execute."
    exit 1
fi

printf '\n'
printf 'Try:\n'
printf '\n'
printf '    whiskey --help\n'
printf '\n'
printf 'or\n'
printf '\n'
printf '    whiskey build .\n'
printf '\n'
printf 'Documentation:\n'
printf 'https://github.com/sxijyoti/whiskey\n'

# 8. PATH Warning
case ":$PATH:" in
    *":${INSTALL_DIR}:"*)
        ;;
    *)
        printf '\n'
        log_info "NOTE: ${INSTALL_DIR} is currently not in your \$PATH."
        printf 'To run '\''whiskey'\'' directly from anywhere, add it to your PATH by running:\n'
        printf '\n'
        printf '    export PATH="%s:$PATH"\n' "$INSTALL_DIR"
        printf '\n'
        printf 'To make this permanent, append the line above to your shell configuration file\n'
        printf '(e.g., ~/.bashrc, ~/.zshrc, or ~/.profile), or in CI/CD:\n'
        printf '\n'
        printf '    echo "%s" >> $GITHUB_PATH\n' "$INSTALL_DIR"
        printf '\n'
        ;;
esac
