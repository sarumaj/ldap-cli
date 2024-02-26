/*
ldap-cli is cross-platform compatible client application based on the lightweight directory access control (LDAP).
The CLI application is written in Go and uses the Cobra framework.
It is mainly intended for administrators and developers who want to query and edit directory objects.
It is highly adopted to work with Microsoft Active Directory, although openldap implementations are being used for testing.

Usage:

	ldap-cli [flags]
	ldap-cli [command]

Examples:
ldap-cli --user "DOMAIN\\user" --password "password" --url "ldaps://example.com:636" <command>

Available Commands:

	completion  Generate the autocompletion script for the specified shell
	edit        Edit a directory object
	get         Get a directory object
	help        Help about any command
	version     Display version information

Flags:

	    --auth-type string   Set authentication schema (supported: ["MD5", "NTLM", "SIMPLE", "UNAUTHENTICATED"]) (default "UNAUTHENTICATED")
	-v, --debug              Set log level to debug (-v for verbose, -vv for trace)
	    --disable-tls        Disable TLS (not recommended)
	    --domain string      Set domain (required for NTLM authentication schema)
	-h, --help               help for ldap-cli
	    --max-retries uint   Specify number of retries (default 3)
	    --password string    Set password (will be ignored if authentication schema is set to UNAUTHENTICATED)
	    --size-limit int     Specify query size limit (-1: unlimited) (default 2000)
	    --timeout duration   Specify query timeout (default 10m0s)
	    --url string         Provide address of the directory server (default "ldap://localhost:389")
	    --user string        Set username (will be ignored if authentication schema is set to UNAUTHENTICATED)

Use "ldap-cli [command] --help" for more information about a command.
*/
package main

import (
	commands "github.com/sarumaj/ldap-cli/v2/pkg/app/commands"
	apputil "github.com/sarumaj/ldap-cli/v2/pkg/app/util"
)

// Version holds the application version.
// It gets filled automatically at build time.
var Version = "v2.6.7"

// BuildDate holds the date and time at which the application was build.
// It gets filled automatically at build time.
var BuildDate = "2024-02-26 19:57:30 UTC"

func main() {
	apputil.Logger.Debugf("version: %q, build date: %q", Version, BuildDate)
	commands.Execute(Version, BuildDate)
}
