package commands

import (
	"crypto/tls"
	"fmt"
	"strings"
	"time"

	supererrors "github.com/sarumaj/go-super/errors"
	auth "github.com/sarumaj/ldap-cli/pkg/lib/auth"
	client "github.com/sarumaj/ldap-cli/pkg/lib/client"
	cobra "github.com/spf13/cobra"
)

var searchCmd = func() *cobra.Command {
	bindParameters := &auth.BindParameters{}
	dialOptions := &auth.DialOptions{}

	var address string
	var authType string
	var disableTLS bool

	searchCmd := &cobra.Command{
		Use:     "search",
		Short:   "Search a directory object",
		Example: "ldap-cli search <object>",
		Run: func(*cobra.Command, []string) {
			_ = dialOptions.SetURL(address)
			if dialOptions.URL.Scheme == auth.LDAPS {
				_ = dialOptions.SetTLSConfig(&tls.Config{InsecureSkipVerify: disableTLS})
			}

			_ = bindParameters.SetType(auth.TypeFromString(authType))
			switch {

			case len(bindParameters.User)*len(bindParameters.Password) != 0 && bindParameters.AuthType == auth.UNAUTHENTICATED:
				_ = bindParameters.SetType(auth.SIMPLE)

			case len(bindParameters.User)*len(bindParameters.Password)*len(bindParameters.Domain) != 0 && bindParameters.AuthType == auth.UNAUTHENTICATED:
				_ = bindParameters.SetType(auth.NTLM)

			}

			conn := supererrors.ExceptFn(supererrors.W(auth.Bind(
				bindParameters,
				dialOptions,
			)))

			client.Search(conn, client.SearchArguments{})
		},
	}

	flags := searchCmd.Flags()

	// dial options
	flags.UintVarP(&dialOptions.MaxRetries, "max-retries", "r", 3, "Specify number of retries")
	flags.Int64VarP(&dialOptions.SizeLimit, "size-limit", "s", -1, "Specify query size limit (-1: unlimited)")
	flags.DurationVarP(&dialOptions.TimeLimit, "timeout", "t", 10*time.Second, "Specify query timeout")
	flags.BoolVar(&disableTLS, "disable-tls", false, "Disable TLS (not recommended)")

	// bind parameters
	flags.StringVarP(&authType, "auth-type", "a", auth.UNAUTHENTICATED.String(), fmt.Sprintf("Set authentication schema (supported: [%v])", strings.Join(auth.ListSupportedAuthTypes(true), ", ")))
	flags.StringVarP(&bindParameters.Domain, "domain", "d", "", fmt.Sprintf("Set domain (required for %s authentication schema)", auth.NTLM))
	flags.StringVarP(&bindParameters.Password, "password", "p", "", fmt.Sprintf("Set password (will be ignored if authentication schema is set to %s)", auth.UNAUTHENTICATED))
	flags.StringVar(&address, "url", auth.URL{Scheme: auth.LDAP, Host: "localhost", Port: auth.LDAP_RW}.String(), "Provide address of the directory server")
	flags.StringVarP(&bindParameters.User, "username", "u", "", fmt.Sprintf("Set username (will be ignored if authentication schema is set to %s)", auth.UNAUTHENTICATED))

	return searchCmd
}()
