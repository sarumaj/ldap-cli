package app

import (
	commands "github.com/sarumaj/ldap-cli/v2/pkg/app/internal/commands"
	apputil "github.com/sarumaj/ldap-cli/v2/pkg/app/internal/util"
)

// Execute is the main entry point for the application.
var Execute = commands.Execute

// Logger is the main logger for the application.
var Logger = apputil.Logger
