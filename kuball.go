package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	// Define a keyword flag
	keyword := flag.String("k", "", "Keyword to filter contexts")
	flag.Parse()

	// Check if there are command-line arguments
	if flag.NArg() < 1 {
		fmt.Println("Usage: kuball [-k keyword] <kubectl-args>")
		os.Exit(1)
	}

	// Get the kubectl command arguments to run in each context
	kubectlArgs := flag.Args()

	// Get the list of context names
	contexts, err := getContexts()
	if err != nil {
		fmt.Printf("Error getting contexts: %v\n", err)
		os.Exit(1)
	}

	// Iterate through the context names and run kubectl commands
	for _, context := range contexts {
		// If a keyword is specified and the context does not contain the keyword, skip it
		if *keyword != "" && !strings.Contains(context, *keyword) {
			continue
		}

		err := useContext(context)
		if err != nil {
			fmt.Printf("Error switching to context '%s': %v\n", context, err)
			continue
		} else {
			fmt.Printf("Context: %s\n", context)
		}

		err = runKubectlCommand(kubectlArgs)
		if err != nil {
			fmt.Printf("Error running 'kubectl' command: %v\n", err)
		}
	}
}

// getContexts retrieves the list of available contexts
func getContexts() ([]string, error) {
	contextsOutput, err := runCommand("kubectl", "config", "get-contexts", "-o", "name")
	if err != nil {
		return nil, err
	}
	contexts := strings.Split(strings.TrimSpace(contextsOutput), "\n")
	return contexts, nil
}

// useContext sets the active context
func useContext(context string) error {
	_, err := runCommand("kubectl", "config", "use-context", context)
	return err
}

// runKubectlCommand runs a kubectl command with specified arguments
func runKubectlCommand(args []string) error {
	cmd := exec.Command("kubectl", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// runCommand executes a shell command and returns the combined output
func runCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}
