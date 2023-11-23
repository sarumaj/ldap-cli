package util

import (
	"errors"
	"fmt"
	"strings"

	survey "github.com/AlecAivazis/survey/v2"
	core "github.com/AlecAivazis/survey/v2/core"
	terminal "github.com/AlecAivazis/survey/v2/terminal"
	cobra "github.com/spf13/cobra"
)

func AskBool(cmd *cobra.Command, flagName string, args *[]string) error {
	f := cmd.Flag(flagName)
	if f == nil {
		return fmt.Errorf("flag %q not defined", flagName)
	}

	var discard string
	err := survey.AskOne(&survey.Select{
		Message: f.Usage + ":",
		Options: []string{"true", "false", "skip"},
		Default: "skip",
	}, &discard, survey.WithValidator(func(answer any) error {
		switch answerOption := answer.(core.OptionAnswer); answerOption.Value {

		case "true", "false":
			*args = append(*args, "--"+flagName+"="+answerOption.Value)
			return nil

		default:
			return nil

		}
	}))

	if errors.Is(err, terminal.InterruptErr) {
		PrintlnAndExit("Aborted")
	}

	return err
}

func AskCommand(cmd *cobra.Command, def *cobra.Command) (*cobra.Command, error) {
	var options []string
	for _, child := range cmd.Commands() {
		options = append(options, child.Name())
	}

	var x string
	err := survey.AskOne(&survey.Select{
		Message: "Select command from below:",
		Options: options,
		Default: def.Name(),
		Description: func(value string, index int) string {
			for _, child := range cmd.Commands() {
				if child.Name() == value {
					return child.Short
				}
			}

			return ""
		},
	}, &x)

	if errors.Is(err, terminal.InterruptErr) {
		PrintlnAndExit("Aborted")
	}

	if err != nil {
		return nil, err
	}

	for _, c := range cmd.Commands() {
		if c.Name() == x {
			return c, nil
		}
	}

	return def, nil
}

func AskString(cmd *cobra.Command, flagName string, args *[]string) error {
	f := cmd.Flag(flagName)
	if f == nil {
		return fmt.Errorf("flag %q not defined", flagName)
	}

	var discard string
	err := survey.AskOne(&survey.Input{
		Message: f.Usage + ":",
		Default: "empty",
	}, &discard)

	if errors.Is(err, terminal.InterruptErr) {
		PrintlnAndExit("Aborted")
	}

	if err != nil {
		return err
	}

	if discard == "empty" || discard == "" {
		return nil
	}

	*args = append(*args, "--"+flagName, discard)
	return nil
}

func AskStrings(cmd *cobra.Command, flagName string, options, def []string, args *[]string) error {
	f := cmd.Flag(flagName)
	if f == nil {
		return fmt.Errorf("flag %q not defined", flagName)
	}

	if len(def) == 1 {
		var discard string
		if err := survey.AskOne(&survey.Select{
			Message: f.Usage + ":",
			Options: options,
			Default: def[0],
		}, &discard); err != nil {
			return err
		}

		*args = append(*args, "--"+flagName, discard)
		return nil
	}

	var discard []string
	err := survey.AskOne(&survey.MultiSelect{
		Message: f.Usage + ":",
		Options: options,
		Default: def,
	}, &discard)

	if errors.Is(err, terminal.InterruptErr) {
		PrintlnAndExit("Aborted")
	}

	if err != nil {
		return err
	}

	*args = append(*args, "--"+flagName, strings.Join(discard, ","))
	return nil
}
