#!/bin/sh
set -eu

# Script to update the Homebrew Formula by substituting checksums and version into the template.
# Usage: ./update_formula.sh <tag> <checksums_file> <template_file> <output_file>

if [ "$#" -ne 4 ]; then
    echo "Usage: $0 <tag> <checksums_file> <template_file> <output_file>" >&2
    exit 1
fi

TAG=$1
CHECKSUMS_FILE=$2
TEMPLATE_FILE=$3
OUTPUT_FILE=$4

# Strip 'v' prefix for the version number used inside Homebrew
VERSION=${TAG#v}

if [ ! -f "$CHECKSUMS_FILE" ]; then
    echo "Error: Checksums file not found at $CHECKSUMS_FILE" >&2
    exit 1
fi

if [ ! -f "$TEMPLATE_FILE" ]; then
    echo "Error: Template file not found at $TEMPLATE_FILE" >&2
    exit 1
fi

# Function to get checksum from file
get_checksum() {
    local filename=$1
    local checksum
    checksum=$(grep "$filename" "$CHECKSUMS_FILE" | cut -d' ' -f1)
    if [ -z "$checksum" ]; then
        echo "Error: Could not find checksum for $filename in $CHECKSUMS_FILE" >&2
        exit 1
    fi
    echo "$checksum"
}

# Retrieve checksums
DARWIN_AMD64_SHA256=$(get_checksum "whiskey_${VERSION}_darwin_amd64.tar.gz")
DARWIN_ARM64_SHA256=$(get_checksum "whiskey_${VERSION}_darwin_arm64.tar.gz")
LINUX_AMD64_SHA256=$(get_checksum "whiskey_${VERSION}_linux_amd64.tar.gz")
LINUX_ARM64_SHA256=$(get_checksum "whiskey_${VERSION}_linux_arm64.tar.gz")

# Generate the output file from template using sed
sed \
    -e "s/__VERSION__/${VERSION}/g" \
    -e "s/__TAG__/${TAG}/g" \
    -e "s/__DARWIN_AMD64_SHA256__/${DARWIN_AMD64_SHA256}/g" \
    -e "s/__DARWIN_ARM64_SHA256__/${DARWIN_ARM64_SHA256}/g" \
    -e "s/__LINUX_AMD64_SHA256__/${LINUX_AMD64_SHA256}/g" \
    -e "s/__LINUX_ARM64_SHA256__/${LINUX_ARM64_SHA256}/g" \
    "$TEMPLATE_FILE" > "$OUTPUT_FILE"

echo "Successfully generated Homebrew Formula at $OUTPUT_FILE for version $VERSION ($TAG)"
