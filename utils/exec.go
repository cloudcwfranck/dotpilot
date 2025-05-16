package utils

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// ExecuteCommand executes a command and returns its output
func ExecuteCommand(command string, args ...string) (string, error) {
	Logger.Debug().Msgf("Executing command: %s %s", command, strings.Join(args, " "))
	
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()
	
	return string(output), err
}

// PromptYesNo asks the user for a yes/no answer
func PromptYesNo(question string) bool {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%s [y/n]: ", question)
		response, err := reader.ReadString('\n')
		if err != nil {
			Logger.Error().Err(err).Msg("Error reading input")
			return false
		}

		response = strings.ToLower(strings.TrimSpace(response))
		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}

		fmt.Println("Please answer with 'y' or 'n'")
	}
}
