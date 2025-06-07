package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Command line flags
var (
	dryRun = flag.Bool("dry-run", false, "Show what would be done without making changes")
	help   = flag.Bool("help", false, "Show help message")
)

// ensureDirExists creates a directory if it doesn't exist
func ensureDirExists(path string, dryRun bool) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if dryRun {
			fmt.Printf("[DRY RUN] Would create directory: %s\n", path)
			return nil
		}
		fmt.Printf("Creating directory: %s\n", path)
		return os.MkdirAll(path, 0755)
	}
	return nil
}

// expandPath expands ALL environment variables in a string
func expandPath(path string) string {
	// Use os.ExpandEnv to expand all environment variables
	// This handles $VAR and ${VAR} syntax automatically
	return os.ExpandEnv(path)
}

// createSymlink creates a symbolic link
func createSymlink(targetPath, symlinkPath string, dryRun bool) error {
	// Check if existing symlink or file exists
	if _, err := os.Lstat(symlinkPath); err == nil {
		if dryRun {
			fmt.Printf("[DRY RUN] Would remove existing: %s\n", symlinkPath)
		} else {
			fmt.Printf("Removing existing: %s\n", symlinkPath)
			if err := os.RemoveAll(symlinkPath); err != nil {
				return fmt.Errorf("error removing existing path: %w", err)
			}
		}
	}

	// Create the symlink
	if dryRun {
		fmt.Printf("[DRY RUN] Would create symlink: %s -> %s\n", symlinkPath, targetPath)
		return nil
	}

	fmt.Printf("Creating symlink: %s -> %s\n", symlinkPath, targetPath)
	return os.Symlink(targetPath, symlinkPath)
}

// setupSymlinks reads a configuration file and creates symlinks
func setupSymlinks(configFilePath string, dryRun bool) error {
	// Check if config file exists
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		return fmt.Errorf("error: Config file not found: %s", configFilePath)
	}

	if dryRun {
		fmt.Printf("[DRY RUN] Would set up symlinks from config: %s\n", configFilePath)
	} else {
		fmt.Printf("Setting up symlinks from config: %s\n", configFilePath)
	}

	// Open the config file
	file, err := os.Open(configFilePath)
	if err != nil {
		return fmt.Errorf("error opening config file: %w", err)
	}
	defer file.Close()

	// Read the file line by line
	scanner := bufio.NewScanner(file)
	lineNumber := 0

	commentRegex := regexp.MustCompile(`^\s*#`)

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()

		// Skip empty lines and comments
		if line == "" || commentRegex.MatchString(line) {
			continue
		}

		// Split line into symlink_path and actual_path
		fields := strings.Fields(line)
		if len(fields) < 2 {
			fmt.Printf("Warning: Invalid line %d in config file: %s\n", lineNumber, line)
			continue
		}

		// Expand environment variables in both paths
		symlinkPath := expandPath(fields[0])
		actualPath := expandPath(fields[1])

		// Skip if either path is empty after expansion
		if symlinkPath == "" || actualPath == "" {
			fmt.Printf("Warning: Invalid paths at line %d in config file: %s\n", lineNumber, line)
			continue
		}

		// Check if expansion actually happened (detect unexpanded variables)
		if strings.Contains(symlinkPath, "$") || strings.Contains(actualPath, "$") {
			fmt.Printf("Warning: Unexpanded environment variables at line %d: %s\n", lineNumber, line)
		}

		// Get the directory of the symlink
		symlinkDir := filepath.Dir(symlinkPath)

		if dryRun {
			// new line for dry run output
			fmt.Printf("\n")
			fmt.Printf("[DRY RUN] Line %d: %s -> %s\n", lineNumber, fields[0], fields[1])
			fmt.Printf("[DRY RUN] Expanded: %s -> %s (dir: %s)\n", symlinkPath, actualPath, symlinkDir)
		}

		// Create symlink directory if it doesn't exist
		if err := ensureDirExists(symlinkDir, dryRun); err != nil {
			return fmt.Errorf("error creating directory %s: %w", symlinkDir, err)
		}

		// Create the symlink
		if err := createSymlink(actualPath, symlinkPath, dryRun); err != nil {
			return fmt.Errorf("error creating symlink at line %d: %w", lineNumber, err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}

	if dryRun {
		fmt.Println("[DRY RUN] Symlink setup complete! (No changes made)")
	} else {
		fmt.Println("Symlink setup complete!")
	}
	return nil
}

func getExecutablePath() (string, error) {
	// Get the executable path
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}

	// Resolve symlinks
	realPath, err := filepath.EvalSymlinks(exe)
	if err != nil {
		return "", err
	}

	return filepath.Dir(realPath), nil
}

func showHelp() {
	fmt.Println("Symlink Manager - Create and manage symlinks from configuration files")
	fmt.Println("\nUsage:")
	fmt.Println("  symlinker [flags] [config-file]")
	fmt.Println("\nFlags:")
	flag.PrintDefaults()
	fmt.Println("\nEnvironment Variable Expansion:")
	fmt.Println("  Supports all environment variables in format $VAR or ${VAR}")
	fmt.Println("  Examples: $HOME, $USER, $DOTFILES_HOME, ${XDG_CONFIG_HOME}")
	fmt.Println("\nRequired Environment Variables:")
	fmt.Println("  Make sure to set the environment variables used in your config file")
	fmt.Println("  Example: export DOTFILES_HOME=\"$HOME/dotfiles\"")
	fmt.Println("\nExamples:")
	fmt.Println("  symlinker                    # Use default config file")
	fmt.Println("  symlinker custom.conf        # Use custom config file")
	fmt.Println("  symlinker --dry-run          # Preview changes without applying")
	fmt.Println("  symlinker --dry-run my.conf  # Preview with custom config")
}

func printEnvironmentInfo(dryRun bool) {
	if dryRun {
		fmt.Println("🔍 DRY RUN MODE - No changes will be made")
		fmt.Println("=" + strings.Repeat("=", 50))
		fmt.Printf("Current environment variables:\n")
		
		// Show commonly used environment variables
		commonVars := []string{"HOME", "USER", "XDG_CONFIG_HOME", "DOTFILES_HOME", "TOOLS_DIR", "NOTES_DIR"}
		for _, varName := range commonVars {
			if value := os.Getenv(varName); value != "" {
				fmt.Printf("  %s: %s\n", varName, value)
			}
		}
		fmt.Println("=" + strings.Repeat("=", 50))
	}
}

func main() {
	// CLI args
	flag.Parse()

	// Show help if requested
	if *help {
		showHelp()
		return
	}

	// Get the executable directory
	execDir, err := getExecutablePath()
	if err != nil {
		fmt.Printf("Error getting executable path: %s\n", err)
		os.Exit(1)
	}

	// Default config file location
	defaultConfigFile := filepath.Join(execDir, "symlinker.conf")

	// Get config file path from remaining arguments
	configFilePath := defaultConfigFile
	if flag.NArg() > 0 {
		configFilePath = flag.Arg(0)
	}

	// Print environment info if dry run
	printEnvironmentInfo(*dryRun)

	// Setup symlinks
	if err := setupSymlinks(configFilePath, *dryRun); err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
