#!/bin/bash
set -e

supported_platforms=(
    darwin-amd64
    darwin-arm64
    freebsd-386
    freebsd-amd64
    freebsd-arm64
    linux-386
    linux-amd64
    linux-arm
    linux-arm64
    windows-386
    windows-amd64
    windows-arm64
)

VERSION="$1"
BUILD_DATE="$(date -u "+%Y-%m-%d %H:%M:%S UTC")"

echo "VERSION=${VERSION} BUILD_DATE=${BUILD_DATE}"

for ((j = 0; j < ${#supported_platforms[@]}; j++)); do
    p="${supported_platforms[$j]}"
    goos="${p%-*}"
    goarch="${p#*-}"

    ext=""
    if [ "$goos" = "windows" ]; then
        ext=".exe"
    fi

    echo "go build ( $((j + 1)) / ${#supported_platforms[@]} ): GOOS=${goos} GOARCH=${goarch} CGO_ENABLED=${CGO_ENABLED:-0} -o dist/ldap-cli_${VERSION}_${p}${ext}"

    GOOS="$goos" GOARCH="$goarch" CGO_ENABLED="${CGO_ENABLED:-0}" go build \
        -trimpath \
        -ldflags="-s -w -X 'main.Version=${VERSION}' -X 'main.BuildDate=${BUILD_DATE}' -extldflags=-static" \
        -tags="osusergo netgo static_build" \
        -o "dist/ldap-cli_${VERSION}_${p}${ext}.uncompressed" \
        "cmd/ldap-cli/main.go"

    (
        upx --best -q -q -v "dist/ldap-cli_${VERSION}_${p}${ext}.uncompressed" -o "dist/ldap-cli_${VERSION}_${p}${ext}" &&
            rm "dist/ldap-cli_${VERSION}_${p}${ext}.uncompressed"
    ) || mv "dist/ldap-cli_${VERSION}_${p}${ext}.uncompressed" "dist/ldap-cli_${VERSION}_${p}${ext}"

    sha256sum "dist/ldap-cli_${VERSION}_${p}${ext}" | awk '{print $1}' >"dist/ldap-cli_${VERSION}_${p}${ext}.sha256"
done
