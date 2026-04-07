package ui

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"charm.land/huh/v2"
)

// ProfileFormData contains profile values used by interactive profile forms.
type ProfileFormData struct {
	Profile        string
	UserName       string
	UserEmail      string
	UserSigningKey string
}

// SelectProfile prompts the user to select a profile from the provided list.
func SelectProfile(names []string, in io.Reader, out io.Writer) (string, error) {
	if len(names) == 0 {
		return "", fmt.Errorf("There are no available profiles")
	}

	var profile string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select a profile").
				Options(huh.NewOptions(names...)...).
				Value(&profile),
		),
	).WithInput(in).WithOutput(out)

	err := form.Run()
	if err != nil {
		return "", err
	}

	return profile, nil
}

// PromptProfileName prompts the user for a profile name.
func PromptProfileName(in io.Reader, out io.Writer) (string, error) {
	var profile string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter an existing profile name or a new one.").
				Value(&profile).
				Validate(func(input string) error {
					if strings.TrimSpace(input) == "" {
						return fmt.Errorf("profile name cannot be empty")
					}

					return nil
				}),
		),
	).WithInput(in).WithOutput(out)

	err := form.Run()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(profile), nil
}

// PromptProfileFields prompts the user for profile field values.
func PromptProfileFields(initial ProfileFormData, in io.Reader, out io.Writer) (ProfileFormData, error) {
	result := initial

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Value for user.name").Value(&result.UserName),
			huh.NewInput().Title("Value for user.email").Value(&result.UserEmail),
			huh.NewInput().Title("Value for user.signingkey").Value(&result.UserSigningKey),
		),
	).WithInput(in).WithOutput(out)

	err := form.Run()
	if err != nil {
		return ProfileFormData{}, err
	}

	return result, nil
}

// IsAborted reports whether the error represents user-initiated form cancellation.
func IsAborted(err error) bool {
	return errors.Is(err, huh.ErrUserAborted)
}
