package service

import (
	"dock/internal/dockerClient"
	"dock/internal/util"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/charmbracelet/huh"
)

// Validation functions
func validateName(input string) error {
	if len(input) == 0 {
		return fmt.Errorf("Name cannot be empty")
	}

	pattern := `^[a-zA-Z0-9][a-zA-Z0-9_.-]*$`
	re := regexp.MustCompile(pattern)

	if !re.MatchString(input) {
		return fmt.Errorf("Invalid name. Only letters, numbers, underscores, dots, and hyphens allowed. No spaces.")
	}
	return nil
}

func validateUsername(input string) error {
	if len(input) == 0 {
		return fmt.Errorf("Username cannot be empty")
	}

	pattern := `^[a-z0-9][a-z0-9_]{3,29}$`
	re := regexp.MustCompile(pattern)

	if !re.MatchString(input) {
		return fmt.Errorf("Invalid Docker username. Must be 4-30 chars, lowercase, and contain only letters, digits, and underscores.")
	}
	return nil
}

func validatePassword(input string) error {
	if len(input) == 0 {
		return fmt.Errorf("Password cannot be empty")
	}

	pattern := `^[^\s:/@?#\[\]&=%]+$`
	re := regexp.MustCompile(pattern)

	if !re.MatchString(input) {
		return fmt.Errorf("Invalid PostgreSQL password. Avoid spaces and special URI characters (:/@?#[]&=).")
	}
	return nil
}

func validateDBName(input string) error {
	if len(input) == 0 {
		return fmt.Errorf("Database name cannot be empty")
	}

	pattern := `^[a-zA-Z_][a-zA-Z0-9_]{0,62}$`
	re := regexp.MustCompile(pattern)

	if !re.MatchString(input) {
		return fmt.Errorf("Invalid PostgreSQL database name. Must start with a letter or underscore, contain only letters, numbers, and underscores, and be 1-63 characters long.")
	}
	return nil
}

func validatePort(input string) error {
	if len(input) == 0 {
		return fmt.Errorf("Port cannot be empty")
	}

	pattern := `^([1-9][0-9]{0,3}|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5])$`
	re := regexp.MustCompile(pattern)

	if !re.MatchString(input) {
		return fmt.Errorf("Invalid PostgreSQL port. Must be a number between 1 and 65535.")
	}

	num, err := strconv.Atoi(input)
	if err != nil || num < 1 || num > 65535 {
		return fmt.Errorf("Invalid PostgreSQL port. Must be between 1 and 65535.")
	}

	return nil
}

// Main function to run the PostgreSQL creation prompts
func NewPostgresCmd() {
	fmt.Println("🚀 Let's Create a new PostgreSQL DB inside Docker")
	actionType := ""
	options := []huh.Option[string]{
		{Value: "compose", Key: "Copy a Compose file to create a PostgreSQL DB"},
		{Value: "raw", Key: "Create a new PostgreSQL DB with raw Docker commands"},
	}

	huh.NewSelect[string]().Title("What do you want to do?").Options(options...).Value(&actionType).Run()

	if actionType == "compose" {
		NewPostgresComposeCmd()
		return
	}

	// Run prompts
	var (
		name         string
		username     string
		password     string
		confirmPass  string
		port         string
		dbName       string
		showPassword bool
	)

	var form = huh.NewForm(
		huh.NewGroup(huh.NewInput().Title("Container Name").Prompt("? ").Validate(validateName).Suggestions([]string{"postgres", "pg"}).Value(&name)).Title("PostgreSQL Configuration"),
		huh.NewGroup(
			huh.NewInput().Title("Username").Prompt("? ").Validate(validateUsername).Value(&username).Suggestions([]string{"postgres", "pg", "user", "pguser"}),
			huh.NewInput().EchoMode(huh.EchoModePassword).Prompt("? ").Title("Password").Validate(validatePassword).Value(&password),
			huh.NewInput().EchoMode(huh.EchoModePassword).Prompt("? ").Title("Confirm Password").Validate(func(s string) error {
				if password != s {
					return fmt.Errorf("Passwords do not match")
				}
				return nil
			}).Value(&confirmPass)).Title("Credentials"),
		huh.NewGroup(
			huh.NewInput().Title("Database Name").Prompt("? ").Validate(validateDBName).Suggestions([]string{"postgres", "pg", "database", "pgdatabase", "db"}).Value(&dbName),
			huh.NewInput().Suggestions([]string{"5432"}).Prompt("? ").Title("Port").Validate(validatePort).Value(&port)).Title("Database"),
		huh.NewGroup(huh.NewConfirm().Title("You want to see the password?").Value(&showPassword)).Title("Customize"),
	)

	err := form.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if name == "" {
		fmt.Println("❗ Container name cannot be empty")
		os.Exit(1)
	}

	fmt.Println("")
	// Display the final values
	fmt.Println("\n✅ PostgreSQL Configuration:")
	fmt.Printf("📦 Container Name: %s\n", name)
	fmt.Printf("👤 Username: %s\n", username)
	if showPassword {
		fmt.Printf("🔑 Password: %s\n", password)
	}
	fmt.Printf("🗄️  Database: %s\n", dbName)
	fmt.Printf("🔌 Port: %s\n", port)

	cs, err := dockerClient.CreatePostgresContainer(dockerClient.Options{
		Name:     name,
		User:     username,
		Password: password,
		DB:       dbName,
		Port:     port,
		Image:    "postgres:latest",
	})

	if err != nil {
		fmt.Printf("Error While Creating Postgres Container: %s\n", err)
		return
	}

	if showPassword {
		fmt.Printf("Connection String: %s\n", cs)
	} else {
		fmt.Printf("postgres://%s:%s@%s:%s/%s", username, "*********", "127.0.0.1", port, dbName)
	}
}

func NewPostgresComposeCmd() {
	templateFilePath := "./templates/postgres.yml"

	// Check if the file exists
	if _, err := os.Stat(templateFilePath); os.IsNotExist(err) {
		fmt.Printf("Template file not found: %s\n", templateFilePath)
		os.Exit(1)
	}

	// Read the file
	file, err := os.ReadFile(templateFilePath)
	if err != nil {
		fmt.Printf("Error reading template file: %s\n", err)
		os.Exit(1)
	}

	// Parse the YAML

	util.CopyToClipboard(string(file))

	fmt.Println("📋 Template file copied to clipboard")
	fmt.Println("💡 You can paste this file into your project's docker-compose.yml file")
	fmt.Println(" — And then run `docker-compose up -d` to start the PostgreSQL container")
	fmt.Println("🚀 Happy coding!")
}
