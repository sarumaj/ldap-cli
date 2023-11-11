package commands

import (
	"crypto/tls"
	"fmt"
	"strings"
	"time"

	supererrors "github.com/sarumaj/go-super/errors"
	apputil "github.com/sarumaj/ldap-cli/pkg/app/util"
	auth "github.com/sarumaj/ldap-cli/pkg/lib/auth"
	logrus "github.com/sirupsen/logrus"
	cobra "github.com/spf13/cobra"
)

var rootFlags struct {
	address        string
	authType       string
	bindParameters auth.BindParameters
	debug          bool
	dialOptions    auth.DialOptions
	disableTLS     bool
}

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
	flags.BoolVar(&rootFlags.debug, "debug", false, "Set log level to debug")

	// dial options
	flags.UintVar(&rootFlags.dialOptions.MaxRetries, "max-retries", 3, "Specify number of retries")
	flags.Int64Var(&rootFlags.dialOptions.SizeLimit, "size-limit", -1, "Specify query size limit (-1: unlimited)")
	flags.DurationVar(&rootFlags.dialOptions.TimeLimit, "timeout", 10*time.Minute, "Specify query timeout")
	flags.BoolVar(&rootFlags.disableTLS, "disable-tls", false, "Disable TLS (not recommended)")

	// bind parameters
	flags.StringVar(&rootFlags.address, "url", auth.URL{Scheme: auth.LDAP, Host: "localhost", Port: auth.LDAP_RW}.String(), "Provide address of the directory server")
	flags.StringVar(&rootFlags.authType, "auth-type", auth.UNAUTHENTICATED.String(), fmt.Sprintf("Set authentication schema (supported: [%v])", strings.Join(auth.ListSupportedAuthTypes(true), ", ")))
	flags.StringVar(&rootFlags.bindParameters.Domain, "domain", "", fmt.Sprintf("Set domain (required for %s authentication schema)", auth.NTLM))
	flags.StringVar(&rootFlags.bindParameters.Password, "password", "", fmt.Sprintf("Set password (will be ignored if authentication schema is set to %s)", auth.UNAUTHENTICATED))
	flags.StringVar(&rootFlags.bindParameters.User, "user", "", fmt.Sprintf("Set username (will be ignored if authentication schema is set to %s)", auth.UNAUTHENTICATED))

	rootCmd.AddCommand(getCmd, versionCmd)

	return rootCmd
}()

func rootPersistentPreRun(cmd *cobra.Command, _ []string) {
	if rootFlags.debug {
		apputil.Logger.SetLevel(logrus.DebugLevel)
	}

	if _ = rootFlags.dialOptions.SetURL(rootFlags.address); rootFlags.dialOptions.URL.Scheme == auth.LDAPS {
		_ = rootFlags.dialOptions.SetTLSConfig(&tls.Config{InsecureSkipVerify: rootFlags.disableTLS})
	}

	switch _ = rootFlags.bindParameters.SetType(auth.TypeFromString(rootFlags.authType)); {

	case len(rootFlags.bindParameters.User)*len(rootFlags.bindParameters.Password) != 0 &&
		rootFlags.bindParameters.AuthType == auth.UNAUTHENTICATED:

		_ = rootFlags.bindParameters.SetType(auth.SIMPLE)

	case len(rootFlags.bindParameters.User)*len(rootFlags.bindParameters.Password)*len(rootFlags.bindParameters.Domain) != 0 &&
		rootFlags.bindParameters.AuthType == auth.UNAUTHENTICATED:

		_ = rootFlags.bindParameters.SetType(auth.NTLM)

	}
}

func rootRun(cmd *cobra.Command, args []string) {
	child := supererrors.ExceptFn(supererrors.W(apputil.AskCommand(cmd, getCmd)))
	if child.PersistentPreRun != nil {
		child.PersistentPreRun(child, nil)
	}
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
