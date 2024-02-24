package main

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"testing"

	godog "github.com/cucumber/godog"
	supererrors "github.com/sarumaj/go-super/errors"
	commands "github.com/sarumaj/ldap-cli/v2/pkg/app/commands"
	apputil "github.com/sarumaj/ldap-cli/v2/pkg/app/util"
	libutil "github.com/sarumaj/ldap-cli/v2/pkg/lib/util"
)

type contextEntity[T any] struct{}

func (c contextEntity[T]) Add(ctx context.Context, values any) context.Context {
	args, _ := ctx.Value(c).([]string)

	switch v := values.(type) {

	case string:
		args = append(args, v)

	case []string:
		args = append(args, v...)

	case *godog.Table:
		for _, row := range v.Rows {
			if len(row.Cells) < 2 {
				continue
			}

			key, values := row.Cells[0].Value, row.Cells[1:]
			for _, value := range values {
				args = append(args, key, value.Value)
			}
		}

	}

	return context.WithValue(ctx, c, args)
}

func (c contextEntity[T]) Get(ctx context.Context) T {
	v, _ := ctx.Value(c).(T)
	return v
}

func (c contextEntity[T]) Set(ctx context.Context, value T) context.Context {
	return context.WithValue(ctx, c, value)
}

func InitializeScenario(scx *godog.ScenarioContext) {
	args := contextEntity[[]string]{}
	code := contextEntity[int]{}
	output := contextEntity[string]{}
	stdOut, stdErr := os.Stdout, os.Stderr

	scx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		var cancelFunc context.CancelCauseFunc
		ctx, cancelFunc = context.WithCancelCause(ctx)

		libutil.Exit = func(exitCode int) {
			code.Set(ctx, exitCode)
		}

		supererrors.RegisterCallback(func(err error) {
			if err != nil && !libutil.ErrorIs(err, io.EOF) {
				cancelFunc(err)
			}
		})

		apputil.Logger.SetOutput(io.Discard)
		supererrors.Except(os.Setenv("NO_COLOR", "true"))

		return args.Add(ctx, []string{}), nil
	})

	scx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		supererrors.RestoreCallback()
		libutil.Exit, os.Stdout, os.Stderr = os.Exit, stdOut, stdErr
		apputil.Logger.SetOutput(os.Stdout)
		supererrors.Except(os.Unsetenv("NO_COLOR"))

		return ctx, err
	})

	scx.Given(
		`^I perform bind with following parameters:$`,
		func(ctx context.Context, params *godog.Table) (context.Context, error) {
			if params == nil {
				return ctx, nil
			}

			return args.Add(ctx, params), nil
		},
	)

	scx.Given(
		`^I plan to execute "([^"]*)" command with following parameters:$`,
		func(ctx context.Context, command string, data *godog.Table) (context.Context, error) {
			ctx = args.Add(ctx, command)
			if data != nil {
				ctx = args.Add(ctx, data)
			}

			return ctx, nil
		},
	)

	scx.When(
		`^I execute the application$`,
		func(ctx context.Context) (context.Context, error) {
			reader, writer := supererrors.ExceptFn2(supererrors.W2(os.Pipe()))
			os.Stdout, os.Stderr = writer, writer

			commands.Execute(Version, BuildDate, args.Get(ctx)...)

			supererrors.Except(writer.Close())
			buffer := bytes.NewBuffer(nil)
			_ = supererrors.ExceptFn(supererrors.W(io.Copy(buffer, reader)))

			ctx = output.Set(ctx, buffer.String())
			return ctx, nil
		},
	)

	scx.Then(
		`^I expect the output to be:$`,
		func(ctx context.Context, expectedOutput *godog.DocString) (context.Context, error) {
			got := bufio.NewScanner(strings.NewReader(output.Get(ctx)))
			got.Split(bufio.ScanLines)
			expected := bufio.NewScanner(strings.NewReader(expectedOutput.Content))
			expected.Split(bufio.ScanLines)

			var errs []error
			for got.Scan() && expected.Scan() {
				if !regexp.MustCompile(expected.Text()).MatchString(got.Text()) {
					errs = append(errs, fmt.Errorf("expected: %q, got: %q", expected.Text(), got.Text()))
				}
			}

			return ctx, errors.Join(errs...)
		},
	)

	scx.Then(
		`^I expect the exit code to be (\d+)$`,
		func(ctx context.Context, expectedExitCode int) (context.Context, error) {
			got := code.Get(ctx)
			if got != expectedExitCode {
				return ctx, fmt.Errorf("expected: %d, got: %d", expectedExitCode, got)
			}

			return ctx, nil
		},
	)
}

func TestFeatures(t *testing.T) {
	libutil.SkipOAT(t)

	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:        "pretty",
			Paths:         []string{"../../oat/features"},
			Randomize:     -1,
			StopOnFailure: true,
			TestingT:      t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
