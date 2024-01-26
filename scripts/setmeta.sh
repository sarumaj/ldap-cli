#!/bin/bash
set -e

version="$(git describe --tags "$(git rev-list --tags --max-count=1)")"
buildDate="$(date -u "+%Y-%m-%d %H:%M:%S UTC")"

echo "got version: ${version}, buildDate: ${buildDate}"

sed -E "s/^(var Version = \")[^\"]*(\".*)\$/\1${version}\2/" < cmd/ldap-cli/main.go | \
sed -E "s/^(var BuildDate = \")[^\"]*(\".*)\$/\1${buildDate}\2/" > cmd/ldap-cli/main.go.new

mv cmd/ldap-cli/main.go.new cmd/ldap-cli/main.go
echo "updated cmd/ldap-cli/main.go"
