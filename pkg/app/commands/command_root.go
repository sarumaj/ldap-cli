package commands

import (
	"crypto/tls"
	"fmt"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	supererrors "github.com/sarumaj/go-super/errors"
	apputil "github.com/sarumaj/ldap-cli/pkg/app/util"
	auth "github.com/sarumaj/ldap-cli/pkg/lib/auth"
	logrus "github.com/sirupsen/logrus"
	cobra "github.com/spf13/cobra"
)

var rootFlags struct {
	Address        string
	AuthType       string
	BindParameters auth.BindParameters
	Debug          bool
	DialOptions    auth.DialOptions
	DisableTLS     bool
}
var address string
var authType string
var bindParameters = &auth.BindParameters{}
var debug bool
var dialOptions = &auth.DialOptions{}
var disableTLS bool

var rootCmd = func() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:              "ldap-cli",
		Short:            "ldap-cli is cross-platform compatible client application based on the lightweight directory access control (LDAP)",
		PersistentPreRun: rootPersistentPreRun,
		Run:              rootRun,
		Example:          `ldap-cli --user "DOMAIN\\user" --password "password" --url "ldaps://example.com:636" <command>`,
		Version:          internalVersion,
	}

	flags := rootCmd.PersistentFlags()
	flags.BoolVar(&debug, "debug", false, "Set log level to debug")

	// dial options
	flags.UintVar(&dialOptions.MaxRetries, "max-retries", 3, "Specify number of retries")
	flags.Int64Var(&dialOptions.SizeLimit, "size-limit", -1, "Specify query size limit (-1: unlimited)")
	flags.DurationVar(&dialOptions.TimeLimit, "timeout", 10*time.Minute, "Specify query timeout")
	flags.BoolVar(&disableTLS, "disable-tls", false, "Disable TLS (not recommended)")

	// bind parameters
	flags.StringVar(&address, "url", auth.URL{Scheme: auth.LDAP, Host: "localhost", Port: auth.LDAP_RW}.String(), "Provide address of the directory server")
	flags.StringVar(&authType, "auth-type", auth.UNAUTHENTICATED.String(), fmt.Sprintf("Set authentication schema (supported: [%v])", strings.Join(auth.ListSupportedAuthTypes(true), ", ")))
	flags.StringVar(&bindParameters.Domain, "domain", "", fmt.Sprintf("Set domain (required for %s authentication schema)", auth.NTLM))
	flags.StringVar(&bindParameters.Password, "password", "", fmt.Sprintf("Set password (will be ignored if authentication schema is set to %s)", auth.UNAUTHENTICATED))
	flags.StringVar(&bindParameters.User, "user", "", fmt.Sprintf("Set username (will be ignored if authentication schema is set to %s)", auth.UNAUTHENTICATED))

	rootCmd.AddCommand(getCmd, versionCmd)

	return rootCmd
}()

func rootPersistentPreRun(cmd *cobra.Command, _ []string) {
	if debug {
		apputil.Logger.SetLevel(logrus.DebugLevel)
	}

	if _ = dialOptions.SetURL(address); dialOptions.URL.Scheme == auth.LDAPS {
		_ = dialOptions.SetTLSConfig(&tls.Config{InsecureSkipVerify: disableTLS})
	}

	switch _ = bindParameters.SetType(auth.TypeFromString(authType)); {

	case len(bindParameters.User)*len(bindParameters.Password) != 0 && bindParameters.AuthType == auth.UNAUTHENTICATED:
		_ = bindParameters.SetType(auth.SIMPLE)

	case len(bindParameters.User)*len(bindParameters.Password)*len(bindParameters.Domain) != 0 && bindParameters.AuthType == auth.UNAUTHENTICATED:
		_ = bindParameters.SetType(auth.NTLM)

	}
}

func rootRun(cmd *cobra.Command, args []string) {
	var x string
	supererrors.Except(survey.AskOne(&survey.Select{
		Message: "Select command from below:",
		Options: []string{getCmd.Name()},
		Default: getCmd.Name(),
	}, &x))

	var child *cobra.Command
	switch x {

	case getCmd.Name():
		child = getCmd

	default:
		return
	}

	child.PersistentPreRun(child, nil)
	child.Run(child, nil)
}

// Execute executes the root command.
func Execute(version, buildDate string) {
	internalVersion, internalBuildDate = version, buildDate

	apputil.Logger.Debugf("Version: %s, build date: %s, executable path: %s", internalVersion, internalBuildDate, apputil.GetExecutablePath())

	if err := rootCmd.Execute(); err != nil {
		apputil.Logger.Debugf("Execution failed: %v", err)
	}
}
