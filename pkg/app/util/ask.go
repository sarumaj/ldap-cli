package util

import (
	"fmt"
	"strings"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/core"
	cobra "github.com/spf13/cobra"
)

func AskBool(cmd *cobra.Command, flagName string, args *[]string) error {
	f := cmd.Flag(flagName)
	if f == nil {
		return fmt.Errorf("flag %q not defined", flagName)
	}

	var discard string
	return survey.AskOne(&survey.Select{
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
}

func AskString(cmd *cobra.Command, flagName string, args *[]string) error {
	f := cmd.Flag(flagName)
	if f == nil {
		return fmt.Errorf("flag %q not defined", flagName)
	}

	var discard string
	if err := survey.AskOne(&survey.Input{
		Message: f.Usage + ":",
		Default: "empty",
	}, &discard); err != nil {
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
	if err := survey.AskOne(&survey.MultiSelect{
		Message: f.Usage + ":",
		Options: options,
		Default: def,
	}, &discard); err != nil {
		return err
	}

	*args = append(*args, "--"+flagName, strings.Join(discard, ","))
	return nil
}
