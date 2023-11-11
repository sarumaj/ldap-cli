[![test_and_report](https://github.com/sarumaj/ldap-cli/actions/workflows/test_and_report.yml/badge.svg)](https://github.com/sarumaj/ldap-cli/actions/workflows/test_and_report.yml)
[![build_and_release](https://github.com/sarumaj/ldap-cli/actions/workflows/build_and_release.yml/badge.svg)](https://github.com/sarumaj/ldap-cli/actions/workflows/build_and_release.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/sarumaj/ldap-cli)](https://goreportcard.com/report/github.com/sarumaj/ldap-cli)
[![Maintainability](https://img.shields.io/codeclimate/maintainability-percentage/sarumaj/ldap-cli.svg)](https://codeclimate.com/github/sarumaj/ldap-cli/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/147f265284b27931c2d2/test_coverage)](https://codeclimate.com/github/sarumaj/ldap-cli/test_coverage)

---

# ldap-cli

**ldap-cli** is a cross-platform compatible LDAP-based command-line interface allowing ETL operations on Directory objects using LDAP Data Interchange Format (LDIF).

Developed as Computer Science Project for IU (www.iu-fernstudium.de).

## Installation

Download executable binary from the [release section](https://github.com/sarumaj/ldap-cli/releases/latest), e.g.:

```console
$ latest=$(curl -fsSI https://github.com/sarumaj/ldap-cli/releases/latest | grep -i location | sed 's/\r$//' | awk -F/ '{print $NF}') && \
  curl -fsSL "https://github.com/sarumaj/ldap-cli/releases/download/${latest}/ldap-cli_${latest}_linux-amd64" -o /usr/local/bin/ldap-cli
$ chmod +x /usr/local/bin/ldap-cli
```

Alternatively, build from source code (requires go 1.21.x runtime):

```console
$ git clone https://github.com/sarumaj/ldap-cli
$ cd ldap-cli
$ go build \
    -trimpath \
    -ldflags="-s -w -X 'main.Version=$(git describe --tags --abbrev=0)' -X 'main.BuildDate=$(date -u "+%Y-%m-%d %H:%M:%S UTC")' -extldflags=-static" \
    -tags="osusergo netgo static_build" \
    -o /usr/local/bin/ldap-cli \
    "cmd/ldap-cli/main.go"
$ chmod +x /usr/local/bin/ldap-cli
```

## Features

- [x] LDAP/LDAPS authentication
  - [x] SIMPLE BIND
  - [x] NTLM (**not tested**)
  - [x] UNAUTHENTICATED
  - [x] MD5 (**not tested**)
- [x] Search directory objects
  - [x] Track progress of search operations
  - [x] Search users with options (user-id, enabled, expired, memberOf)
  - [x] Search groups with options (group-id)
  - [x] Search by providing custom LDAP filter
    - [x] Parse and validate filter syntax
    - [x] Register lexical aliases
- [x] Edit directory objects
  - [x] Edit group members
    - [x] Support arbitrary membership attribute
    - [x] LDIF edit mode
  - [x] Edit user's password
    - [x] Support arbitrary password attribute
    - [ ] Option to pass old password (**not available in LDIF mode**)
    - [x] LDIF edit mode
  - [x] Edit custom objects
    - [x] LDIF edit mode
- [x] Interactive mode operandi
  - [x] Ask for inputs
  - [x] Utilize ANSI code sequences if available
- [x] Support multiple output format (CSV, LDIF, YAML)

## Usage

```console
$ ldap-cli --help

>> ldap-cli is cross-platform compatible client application based on the lightweight directory access control (LDAP)
>>
>> Usage:
>>   ldap-cli [flags]
>>   ldap-cli [command]
>>
>> Examples:
>> ldap-cli --user "DOMAIN\\user" --password "password" --url "ldaps://example.com:636" <command>
>>
>> Available Commands:
>>   completion  Generate the autocompletion script for the specified shell
>>   edit        Edit a directory object
>>   get         Get a directory object
>>   help        Help about any command
>>   version     Display version information
>>
>> Flags:
>>       --auth-type string   Set authentication schema (supported: ["MD5", "NTLM", "SIMPLE", "UNAUTHENTICATED"]) (default "UNAUTHENTICATED")
>>       --debug              Set log level to debug
>>       --disable-tls        Disable TLS (not recommended)
>>       --domain string      Set domain (required for NTLM authentication schema)
>>   -h, --help               help for ldap-cli
>>       --max-retries uint   Specify number of retries (default 3)
>>       --password string    Set password (will be ignored if authentication schema is set to UNAUTHENTICATED)
>>       --size-limit int     Specify query size limit (-1: unlimited) (default 2000)
>>       --timeout duration   Specify query timeout (default 10m0s)
>>       --url string         Provide address of the directory server (default "ldap://localhost:389")
>>       --user string        Set username (will be ignored if authentication schema is set to UNAUTHENTICATED)
>>
>> Use "ldap-cli [command] --help" for more information about a command.
```
