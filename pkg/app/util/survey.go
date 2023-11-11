package util

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	survey "github.com/AlecAivazis/survey/v2"
	core "github.com/AlecAivazis/survey/v2/core"
	terminal "github.com/AlecAivazis/survey/v2/terminal"
	ldif "github.com/go-ldap/ldif"
	cobra "github.com/spf13/cobra"
)

var splitByNewLineRegex = regexp.MustCompile("\r?\n")

// Ask to provide a boolean value for given flag, while offering the option to deny answer (skip)
func AskBool(cmd *cobra.Command, flagName string, args *[]string, opts ...survey.AskOpt) (bool, error) {
	f := cmd.Flag(flagName)
	if f == nil {
		return false, fmt.Errorf("flag %q not defined", flagName)
	}

	description := f.Usage
	if f.DefValue != "" {
		description += " (" + f.DefValue + ")"
	}
	description += ":"

	var discard string
	var set bool
	err := survey.AskOne(&survey.Select{
		Message: description,
		Options: []string{"true", "false", "skip"},
		Default: "skip",
	}, &discard, append(opts, survey.WithValidator(func(answer any) error {
		switch answerOption := answer.(core.OptionAnswer); answerOption.Value {

		case "true", "false":
			*args, set = append(*args, "--"+flagName+"="+answerOption.Value), true

		}

		return nil
	}))...)

	if errors.Is(err, terminal.InterruptErr) {
		PrintlnAndExit("Aborted")
	}

	if err != nil {
		return false, err
	}

	return set, nil
}

// Ask to select a sub-command while providing a default sub-command
func AskCommand(cmd *cobra.Command, def *cobra.Command, opts ...survey.AskOpt) (*cobra.Command, error) {
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
	}, &x, opts...)

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

// Ask to modify an object in LDAP Interchange Data Format
func AskLDAPDataInterchangeFormat(requests *ldif.LDIF, editor string) (bool, error) {
	before, err := ldif.Marshal(requests)
	if err != nil {
		return false, err
	}

	var after string
	err = survey.AskOne(&survey.Editor{
		Message:       "Modify object in LDAP Interchange Format (LDIF)",
		FileName:      "*.ldif",
		Default:       before,
		AppendDefault: true,
		Editor:        editor,
	}, &after)

	if errors.Is(err, terminal.InterruptErr) {
		PrintlnAndExit("Aborted")
	}

	if err != nil {
		return false, err
	}

	err = ldif.Unmarshal(strings.NewReader(after), requests)
	if err != nil {
		return false, err
	}

	return true, nil
}

// Ask to provide answers for given flag in multiline text mode (each option is separated by a new line)
func AskMultiline(cmd *cobra.Command, flagName string, args *[]string, opts ...survey.AskOpt) (bool, error) {
	f := cmd.Flag(flagName)
	if f == nil {
		return false, fmt.Errorf("flag %q not defined", flagName)
	}

	defValue := "empty"
	if f.DefValue != "" {
		defValue = f.DefValue
	}

	var discard string
	err := survey.AskOne(&survey.Multiline{
		Message: f.Usage + ":",
		Default: defValue,
	}, &discard, opts...)

	if errors.Is(err, terminal.InterruptErr) {
		PrintlnAndExit("Aborted")
	}

	if err != nil {
		return false, err
	}

	switch discard {

	case "", "empty", f.DefValue:
		return false, nil

	}

	for _, entry := range splitByNewLineRegex.Split(discard, -1) {
		if entry == "" {
			continue
		}

		*args = append(*args, "--"+flagName, entry)
	}

	return true, nil
}

// Ask a string value for given flag, while offering a default value.
// If default value is not provided, the default value of the flag will be used.
// Password mode (sensitive input mode) is supported
func AskString(cmd *cobra.Command, flagName string, args *[]string, password bool, def string, opts ...survey.AskOpt) (bool, error) {
	f := cmd.Flag(flagName)
	if f == nil {
		return false, fmt.Errorf("flag %q not defined", flagName)
	}

	var discard string
	var prompt survey.Prompt
	if password {
		prompt = &survey.Password{Message: f.Usage + ":"}

	} else {
		defValue := "empty"
		if def != "" {
			defValue = def
		} else if f.DefValue != "" {
			defValue = f.DefValue
		}

		prompt = &survey.Input{
			Message: f.Usage + ":",
			Default: defValue,
		}
	}

	err := survey.AskOne(prompt, &discard, opts...)

	if errors.Is(err, terminal.InterruptErr) {
		PrintlnAndExit("Aborted")
	}

	if err != nil {
		return false, err
	}

	switch discard {

	case "", "empty", f.DefValue:
		return false, nil

	}

	*args = append(*args, "--"+flagName, discard)
	return true, nil
}

// Ask to select one or many option from many fro given flag, default options may be provided.
// If just one default option is being provided, a selection of only one option is foreseen,
// otherwise multiple selection is possible.
func AskStrings(cmd *cobra.Command, flagName string, options, def []string, args *[]string, opts ...survey.AskOpt) (bool, error) {
	f := cmd.Flag(flagName)
	if f == nil {
		return false, fmt.Errorf("flag %q not defined", flagName)
	}

	if len(def) == 1 {
		var discard string
		if err := survey.AskOne(&survey.Select{
			Message: f.Usage + ":",
			Options: options,
			Default: def[0],
		}, &discard, opts...); err != nil {

			return false, err
		}

		*args = append(*args, "--"+flagName, discard)
		return true, nil
	}

	var discard []string
	err := survey.AskOne(&survey.MultiSelect{
		Message: f.Usage + ":",
		Options: options,
		Default: def,
	}, &discard, opts...)

	if errors.Is(err, terminal.InterruptErr) {
		PrintlnAndExit("Aborted")
	}

	if err != nil {
		return false, err
	}

	for _, entry := range discard {
		*args = append(*args, "--"+flagName, entry)
	}

	return true, nil
}
