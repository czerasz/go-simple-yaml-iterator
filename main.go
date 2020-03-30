package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	filePath := flag.String("file", "", "Path to YAML file")
	rootEl := flag.String("root", "", "Root key inside YAML")
	cmdStr := flag.String("command", "", "Command to execute for each item")
	dry := flag.Bool("dry-run", false, "Only pretend to run")
	ver := flag.Bool("version", false, "Show application version")
	flag.Parse()

	if len(os.Args) > 1 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
		flag.PrintDefaults()
	}

	// If --version flag was specified print out version information
	if *ver {
		fmt.Printf("version: %s\n", version)
		fmt.Printf("commit: %s\n", commit)
		fmt.Printf("built at: %s\n", date)
		os.Exit(0)
	}

	if strings.TrimSpace(*filePath) == "" ||
		strings.TrimSpace(*rootEl) == "" ||
		strings.TrimSpace(*cmdStr) == "" {
		flag.PrintDefaults()
		log.Fatalf("Missing flag. Please provide all flags")
	}

	cmdSli := strings.Split(*cmdStr, " ")
	cmdName := cmdSli[0]
	cmdArgs := cmdSli[1:len(cmdSli)]

	dat, err := ioutil.ReadFile(*filePath)
	if err != nil {
		log.Fatalf("Error while reading file: %v", err)
	}

	// Create map for YAML structure
	m := make(map[interface{}]interface{})

	err = yaml.Unmarshal(dat, &m)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// Get current environment variables
	envVars := os.Environ()

	// Iterate through root element items
	for _, item := range m[*rootEl].([]interface{}) {
		newVars := []string{}
		for k, val := range item.(map[string]interface{}) {
			varName := k
			varVal := val.(string)

			// Make parameters environment variable conform
			varName = strings.ReplaceAll(strings.ToUpper(varName), "-", "_")

			// Extend standard environment variables
			additionalEnv := fmt.Sprintf("%s=%s", varName, varVal)
			newVars = append(newVars, additionalEnv)
		}

		if *dry {
			fmt.Printf("%s %s %s\n", strings.Join(newVars, " "), cmdName, strings.Join(cmdArgs, " "))
			continue
		}

		// Prepare command
		cmd := exec.Command(cmdName, cmdArgs...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Env = append(envVars, newVars...)

		// Execute command
		err = cmd.Run()
		if err != nil {
			log.Fatalf("Command execution failed with %s\n", err)
		}
	}
}
