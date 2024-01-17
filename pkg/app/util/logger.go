package util

import (
	"fmt"
	"os"
	"strings"

	color "github.com/fatih/color"
	supererrors "github.com/sarumaj/go-super/errors"
	auth "github.com/sarumaj/ldap-cli/pkg/lib/auth"
	client "github.com/sarumaj/ldap-cli/pkg/lib/client"
	logrus "github.com/sirupsen/logrus"
	tracerr "github.com/ztrue/tracerr"
)

// App logger (default format JSON)
var Logger = func() *logrus.Logger {
	l := logrus.New()
	l.SetLevel(logrus.WarnLevel)
	l.SetOutput(Stdout())
	l.SetFormatter(&logrus.JSONFormatter{
		DisableHTMLEscape: true,
		PrettyPrint:       true,
	})

	supererrors.RegisterCallback(func(err error) {
		l.SetOutput(Stderr())

		if l.Level >= logrus.DebugLevel {
			err = tracerr.Wrap(err)

			var frames []string
			for _, frame := range err.(tracerr.Error).StackTrace() {
				switch ctx := frame.String(); {

				case
					strings.Contains(ctx, "supererrors.Except"),
					strings.Contains(ctx, "runtime.main()"),
					strings.Contains(ctx, "runtime.goexit()"):

					continue

				default:
					frames = append(frames, frame.String())

				}
			}

			l.WithField("stack", frames).Fatalln(err)
		}

		l.SetFormatter(&logrus.TextFormatter{
			ForceColors:            IsColorEnabled(),
			DisableTimestamp:       true,
			DisableLevelTruncation: true,
		})
		l.Fatalln(err)
	})

	return l
}()

// Log fields in JSON format mode (default)
type Fields = logrus.Fields

// GetFieldsForBind produces log fields describing a bind request
func GetFieldsForBind(bindParameters *auth.BindParameters, dialOptions *auth.DialOptions) logrus.Fields {
	fields := make(logrus.Fields)
	if bindParameters == nil {
		fields["bindParameters"] = nil

	} else {
		fields["bindParameters.AuthType"] = bindParameters.AuthType.String()
		fields["bindParameters.Domain"] = bindParameters.Domain
		fields["bindParameters.User"] = bindParameters.User
		fields["bindParameters.PasswordProvided"] = len(bindParameters.Password) > 0

	}

	if dialOptions == nil {
		fields["dialOptions"] = nil

	} else {
		fields["dialOptions.MaxRetries"] = dialOptions.MaxRetries
		fields["dialOptions.SizeLimit"] = dialOptions.SizeLimit
		fields["dialOptions.TLSEnabled"] = dialOptions.TLSConfig != nil && !dialOptions.TLSConfig.InsecureSkipVerify
		fields["dialOptions.TimeLimit"] = dialOptions.TimeLimit
		fields["dialOptions.URL"] = dialOptions.URL.String()

	}

	return fields
}

// GetFieldsForSearch produces log fields describing a search request
func GetFieldsForSearch(searchArguments *client.SearchArguments) logrus.Fields {
	fields := make(logrus.Fields)
	if searchArguments == nil {
		fields["searchArguments"] = nil

	} else {
		fields["searchArguments.Attributes"] = searchArguments.Attributes.ToAttributeList()
		fields["searchArguments.Filter"] = searchArguments.Filter.String()
		fields["searchArguments.Path"] = searchArguments.Path

	}

	return fields
}

// PrintColors prints a string with colors and exits with provided code
func PrintlnAndExit(code int, format string, a ...any) {
	_ = supererrors.ExceptFn(supererrors.W(fmt.Fprintln(Stderr(), PrintColors(color.RedString, format, a...))))
	os.Exit(code)
}
