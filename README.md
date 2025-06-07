# Symlinker

Symlinker is a simple command-line tool for managing symbolic links based on a specified configuration file. It allows you to easily create symlinks for configuration files used in various applications such as shells, version control, and editors.

## Features

- Create symbolic links as specified in a configuration file.
- Supports environment variable expansion in paths.
- Safety features with a dry-run option to preview changes before executing them.

## Setup

### Prerequisites

- Go 1.24.3 or higher
- Ensure that your environment variables such as `$HOME`, `$USER`, and other necessary variables are properly configured.

### Installation

1. Clone this repository:
   ```bash
   git clone <repository-url>
   cd symlinker
   ```

2. Use the following make commands to manage the module:
   ```bash
   make install   # Download dependencies
   make update    # Update dependencies
   make tidy      # Clean up the go.mod file
   ```

## Configuration

The symlinks are configured in the `symlinker.conf` file. The format for each line is:
```
<symlink_path> <actual_path>
```
You can reference environment variables within the path.

### Example Configuration

```plaintext
# ZSH config
$HOME/.zshrc $TOOLS_DIR/.zshrc
$HOME/.gitconfig $TOOLS_DIR/git/.gitconfig
```

## Usage

Run the symlinker command from the terminal:

```bash
symlinker [flags] [config-file]
```

### Flags

- `--dry-run`: Show what would be done without making changes.
- `--help`: Show help message.

### Examples

- Setup symlinks with a default configuration:
  ```bash
  symlinker
  ```

- Setup symlinks with a custom configuration file:
  ```bash
  symlinker symlinker.conf
  ```

- Preview changes before applying:
  ```bash
  symlinker --dry-run
  ```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
