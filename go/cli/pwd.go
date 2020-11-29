package cli

import (
	"almost-scrum/core"
	"errors"
	"fmt"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
)

func validate(input string) error {
	if 0 < len(input) && len(input) < 3 {
		return errors.New("Password must have more than 2 characters")
	}
	return nil
}

func processPwd(args []string) {
	if len(args) == 0 {
		color.Red("User is required")
		usage()
		return
	}

	user := args[0]
	prompt := promptui.Prompt{
		Label:    "Password",
		Validate: validate,
		Mask:     '*',
	}

	password, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	core.SetPassword(user, password)

}
