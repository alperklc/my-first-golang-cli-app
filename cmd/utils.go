package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/manifoldco/promptui"
)

func promptGetInput(pc promptContent) string {
	validate := func(input string) error {
		if len(input) <= 0 {
			return errors.New(pc.errorMsg)
		}
		return nil
	}

	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }} ",
		Valid:   "{{ . | green }} ",
		Invalid: "{{ . | red }} ",
		Success: "{{ . | bold }} ",
	}

	prompt := promptui.Prompt{
		Label:     pc.label,
		Templates: templates,
		Validate:  validate,
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Input: %s\n", result)

	return result
}

func promptGetInputWithDefault(pc promptContent, defaultValue string) string {
	validate := func(input string) error {
		if len(input) <= 0 {
			return errors.New(pc.errorMsg)
		}
		return nil
	}

	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }} ",
		Valid:   "{{ . | green }} ",
		Invalid: "{{ . | red }} ",
		Success: "{{ . | bold }} ",
	}

	prompt := promptui.Prompt{
		Label:     pc.label,
		Templates: templates,
		Validate:  validate,
		Default:   defaultValue,
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Input: %s\n", result)

	return result
}

func promptAddTasks(pc promptContent) string {
	items := []string{}

	var lastEntry string

	for lastEntry == "" || len(items) == 0 {
		prompt := promptui.Prompt{Label: "Type task and hit enter to add. Hit enter after leaving it blank to stop adding."}

		entry, errPrompt := prompt.Run()
		if errPrompt != nil {
			fmt.Printf("Prompt failed %v\n", errPrompt)
			os.Exit(1)
		}

		if entry != "" {
			items = append(items, entry)
		} else {
			log.Infof("Added %d tasks.", len(items))

			return strings.Join(items[:], "|")
		}
	}

	log.Infof("Added %d tasks.", 0)
	return ""
}
