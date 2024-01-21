package commands

import (
	"crypto/tls"
	"errors"
	"fmt"
	"strings"
	"time"

	survey "github.com/AlecAivazis/survey/v2"
	terminal "github.com/AlecAivazis/survey/v2/terminal"
	supererrors "github.com/sarumaj/go-super/errors"
	apputil "github.com/sarumaj/ldap-cli/pkg/app/util"
	auth "github.com/sarumaj/ldap-cli/pkg/lib/auth"
	libutil "github.com/sarumaj/ldap-cli/pkg/lib/util"
	logrus "github.com/sirupsen/logrus"
	cobra "github.com/spf13/cobra"
)

// rootFlags holds the command line flags for the root command.
var rootFlags struct {
	address        string `flag:"url"`
	authType       string `flag:"auth-type"`
	bindParameters auth.BindParameters
	debug          int `flag:"Debug"`
	dialOptions    auth.DialOptions
	disableTLS     bool `flag:"disable-tls"`
}

// rootCmd represents the base command when called without any subcommands.
var rootCmd = func() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:              "ldap-cli",
		Short:            "ldap-cli is cross-platform compatible client application based on the lightweight directory access control (LDAP)",
		PersistentPreRun: rootPersistentPreRun,
		Run:              rootRun,
		Example:          `ldap-cli --user "DOMAIN\\user" --password "password" --url "ldaps://example.com:636" <command>`,
		Version:          versionFlags.internalVersion,
	}

	flags := rootCmd.PersistentFlags()
	flags.CountVarP(&rootFlags.debug, "debug", "v", "Set log level to debug (-v for verbose, -vv for trace)")

	// dial options
	flags.UintVar(&rootFlags.dialOptions.MaxRetries, "max-retries", 3, "Specify number of retries")
	flags.Int64Var(&rootFlags.dialOptions.SizeLimit, "size-limit", 2000, "Specify query size limit (-1: unlimited)")
	flags.DurationVar(&rootFlags.dialOptions.TimeLimit, "timeout", 10*time.Minute, "Specify query timeout")
	flags.BoolVar(&rootFlags.disableTLS, "disable-tls", false, "Disable TLS (not recommended)")

	// bind parameters
	flags.StringVar(&rootFlags.address, "url", auth.URL{Scheme: auth.LDAP, Host: "localhost", Port: auth.LDAP_RW}.String(), "Provide address of the directory server")
	flags.StringVar(&rootFlags.authType, "auth-type", auth.UNAUTHENTICATED.String(), fmt.Sprintf("Set authentication schema (supported: [%v])", strings.Join(auth.ListSupportedAuthTypes(true), ", ")))
	flags.StringVar(&rootFlags.bindParameters.Domain, "domain", "", fmt.Sprintf("Set domain (required for %s authentication schema)", auth.NTLM))
	flags.StringVar(&rootFlags.bindParameters.Password, "password", "", fmt.Sprintf("Set password (will be ignored if authentication schema is set to %s)", auth.UNAUTHENTICATED))
	flags.StringVar(&rootFlags.bindParameters.User, "user", "", fmt.Sprintf("Set username (will be ignored if authentication schema is set to %s)", auth.UNAUTHENTICATED))

	rootCmd.AddCommand(editCmd, getCmd, versionCmd)

	return rootCmd
}()

// Runs always before all executions (inherited by child commands provided).
// Intelligently sets bind request options
func rootPersistentPreRun(cmd *cobra.Command, _ []string) {
	if rootFlags.debug > 1 {
		apputil.Logger.SetLevel(logrus.TraceLevel)
	}

	if rootFlags.debug > 0 {
		apputil.Logger.SetLevel(logrus.DebugLevel)
	}

	apputil.Logger.Debugf(
		"Version: %s, build date: %s, executable path: %s, keyring backends: %v",
		versionFlags.internalVersion,
		versionFlags.internalBuildDate,
		libutil.GetExecutablePath(),
		libutil.Config.AllowedBackends,
	)

	if err := rootFlags.bindParameters.FromKeyring(); err != nil {
		apputil.Logger.Debugf("Failed to access keyring: %v", err)
	}

	logger := apputil.Logger.WithFields(apputil.Fields{"command": cmd.CommandPath(), "step": "rootPersistentPreRun"})
	logger.Trace("Executing")

	_ = rootFlags.dialOptions.SetURL(rootFlags.address)
	logger.WithField("dialOptions.URL", rootFlags.dialOptions.URL.String()).Trace("Set")
	if rootFlags.dialOptions.URL.Scheme == auth.LDAPS {
		_ = rootFlags.dialOptions.SetTLSConfig(&tls.Config{InsecureSkipVerify: rootFlags.disableTLS})
		logger.WithField("dialOptions.TLSConfig.InsecureSkipVerify", rootFlags.dialOptions.TLSConfig.InsecureSkipVerify).Trace("Set")
	}

	switch _ = rootFlags.bindParameters.SetType(auth.TypeFromString(rootFlags.authType)); {

	case
		len(rootFlags.bindParameters.User)*len(rootFlags.bindParameters.Password) != 0 &&
			rootFlags.bindParameters.AuthType == auth.UNAUTHENTICATED:

		_ = rootFlags.bindParameters.SetType(auth.SIMPLE)
		logger.WithField("bindParameters.Type", rootFlags.bindParameters.AuthType.String()).Trace("Set")

	case
		len(rootFlags.bindParameters.User)*len(rootFlags.bindParameters.Password)*len(rootFlags.bindParameters.Domain) != 0 &&
			rootFlags.bindParameters.AuthType == auth.UNAUTHENTICATED:

		_ = rootFlags.bindParameters.SetType(auth.NTLM)
		logger.WithField("bindParameters.Type", rootFlags.bindParameters.AuthType.String()).Trace("Set")

	}

	logger.WithFields(apputil.GetFieldsForBind(&rootFlags.bindParameters, &rootFlags.dialOptions)).Trace("Options")
}

// Runs in interactive mode by asking user to provide values for app parameters
func rootRun(cmd *cobra.Command, _ []string) {
	logger := apputil.Logger.WithFields(apputil.Fields{"command": cmd.CommandPath(), "step": "rootRun"})
	logger.Trace("Executing")

	if rootFlags.bindParameters.AuthType == auth.UNAUTHENTICATED {
		var confirm bool
		if err := survey.AskOne(&survey.Confirm{
			Message: "Running in UNAUTHENTICATED mode, proceed?",
			Default: false,
		}, &confirm); errors.Is(err, terminal.InterruptErr) {

			apputil.PrintlnAndExit(1, "Aborted")
		}

		if !confirm {
			var args []string
			_ = supererrors.ExceptFn(supererrors.W(apputil.AskString(cmd, "url", &args, false, "", survey.WithValidator(func(ans interface{}) error {
				_, err := auth.URLFromString(fmt.Sprint(ans))
				return err
			}))))

			_ = supererrors.ExceptFn(supererrors.W(apputil.AskString(cmd, "auth-type", &args, false, "", survey.WithValidator(func(ans interface{}) error {
				if !auth.TypeFromString(fmt.Sprint(ans)).IsValid() {
					return fmt.Errorf("invalid auth type: %q", fmt.Sprint(ans))
				}
				return nil
			}))))

			if len(args) > 0 && args[len(args)-1] == auth.NTLM.String() {
				_ = supererrors.ExceptFn(supererrors.W(apputil.AskString(cmd, "domain", &args, false, "")))
			}

			_ = supererrors.ExceptFn(supererrors.W(apputil.AskString(cmd, "user", &args, false, "")))
			_ = supererrors.ExceptFn(supererrors.W(apputil.AskString(cmd, "password", &args, true, "")))

			supererrors.Except(cmd.ParseFlags(args))
		}
	}

	child := supererrors.ExceptFn(supererrors.W(apputil.AskCommand(cmd, getCmd)))
	if child.PersistentPreRun != nil {
		child.PersistentPreRun(child, nil)
	}
	child.Run(child, nil)
}

// Execute executes the root command.
// For test purpose, command arguments can be provided
func Execute(version, buildDate string, args ...string) {
	versionFlags.internalVersion, versionFlags.internalBuildDate = version, buildDate

	if len(args) > 0 {
		rootCmd.SetArgs(args)
	}

	if err := rootCmd.Execute(); err != nil {
		apputil.Logger.Debugf("Execution failed: %v", err)
	}

	if err := rootFlags.bindParameters.ToKeyring(); err != nil {
		apputil.Logger.Debugf("Failed to access keyring: %v", err)
	}
}
