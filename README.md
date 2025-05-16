# DotPilot - Dotfiles Management Tool

DotPilot is a cross-platform CLI tool for managing and syncing dotfiles across multiple machines with environment-specific overrides. It uses a Git-backed system to track changes to dotfiles, supports scoped environments, and includes machine-specific configurations.

## Features

- Track dotfiles in a Git repository
- Support for scoped environments (common, environment-specific, machine-specific)
- Automatic Git integration (commit, push, pull)
- Safe sync with backups and diff prompts
- Pre/post hooks for automation
- Package manager integration based on OS detection

## Installation

### Prerequisites

- Go 1.20 or later
- Git
- For SOPS encryption: GPG and Mozilla SOPS

### From Source

1. Clone this repository
   ```bash
   git clone https://github.com/yourusername/dotpilot.git
   cd dotpilot
   ```

2. Build the binary with Go:
   ```bash
   go build -o dotpilot
   ```

3. Move the binary to a location in your PATH:
   ```bash
   # Linux/macOS
   sudo mv dotpilot /usr/local/bin/
   
   # Windows (run in PowerShell as Administrator)
   # Move to a directory in your PATH
   ```

4. Verify installation:
   ```bash
   dotpilot --version
   ```

### Development Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/dotpilot.git
   cd dotpilot
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Test the progress indicators:
   ```bash
   # Run the simple progress indicator test
   go run test_progress.go
   
   # Run the comprehensive demo (great for recording a demo GIF)
   go run demo.go
   ```

4. Run the test suite:
   ```bash
   # Run all tests
   go test ./...
   
   # Run just the progress indicator tests
   go test ./utils -run TestProgress
   ```

5. Create a screencast or GIF of the progress indicators:
   ```bash
   # On macOS with Homebrew
   brew install asciinema
   asciinema rec -t "DotPilot Progress Indicators" demo.cast
   asciinema play demo.cast
   
   # Convert to GIF (requires gifski)
   brew install gifski
   asciicast2gif demo.cast demo.gif
   ```

## Usage

### Initialize DotPilot

To initialize DotPilot with a remote repository:

```bash
# Initialize with a remote repository
dotpilot init --remote https://github.com/username/dotfiles.git --env dev

# Force reinitialization
dotpilot init --remote https://github.com/username/dotfiles.git --force

# Skip package installation and hooks
dotpilot init --remote https://github.com/username/dotfiles.git --skip-packages --skip-hooks
```

### Track Files

To track files or directories in DotPilot:

```bash
# Track a file
dotpilot track ~/.zshrc

# Track a file in a specific environment
dotpilot track ~/.vimrc --env dev

# Track a directory
dotpilot track ~/.config/nvim

# Track a file in the machine-specific environment
dotpilot track ~/.bashrc --env machine
```

### Sync Dotfiles

To sync dotfiles between machines:

```bash
# Sync with default options (pull, apply, push)
dotpilot sync

# Sync without pushing changes
dotpilot sync --no-push

# Dry run to see what would be done
dotpilot sync --dry-run

# Skip backups and diff prompts
dotpilot sync --no-backup --no-diff-prompt

# Sync with advanced conflict resolution
dotpilot sync --resolve-conflicts --strategy=interactive
```

### Bootstrap a Machine

To apply dotfiles and run setup scripts on a new machine:

```bash
# Apply all dotfiles and run setup scripts
dotpilot bootstrap

# Skip common dotfiles
dotpilot bootstrap --skip-common

# Skip environment-specific dotfiles
dotpilot bootstrap --skip-env

# Skip machine-specific dotfiles
dotpilot bootstrap --skip-machine

# Skip setup scripts
dotpilot bootstrap --skip-setup-scripts

# Force overwrite existing files
dotpilot bootstrap --force
```

## Resolve Conflicts

To detect and resolve conflicts between local files and tracked dotfiles:

```bash
# Resolve conflicts interactively (default)
dotpilot resolve

# Keep local versions for all conflicts
dotpilot resolve --strategy=keep-local

# Keep remote versions for all conflicts
dotpilot resolve --strategy=keep-remote

# Attempt three-way merge using a merge tool
dotpilot resolve --strategy=merge

# Keep both versions (backing up the local one)
dotpilot resolve --strategy=backup-both
```

### Check Status

To check the status of your dotfiles:

```bash
dotpilot status
```

## Repository Structure

DotPilot organizes your dotfiles in the following structure:

```
~/.dotpilot/
├── common/                # Files common to all environments
├── envs/                  # Environment-specific configurations
│   ├── default/
│   ├── dev/
│   └── prod/
└── machine/               # Machine-specific configurations
    └── {hostname}/
```

## Configuration Files

DotPilot uses the following special files:

- `preinstall.sh`: Run before package installation
- `postinstall.sh`: Run after package installation
- `postpull.sh`: Run after pulling changes from remote
- `packages.apt`, `packages.brew`, `packages.yay`: Package lists for different package managers

## Secrets Management

DotPilot offers two methods for securely storing sensitive configuration files:

### Basic Encryption

Simple encryption for sensitive files:

```bash
# Add a file as an encrypted secret
dotpilot secrets add ~/.aws/credentials

# Add with a custom name
dotpilot secrets add ~/.ssh/id_rsa --name ssh_key

# List all encrypted secrets
dotpilot secrets list

# Retrieve a decrypted secret
dotpilot secrets get aws_credentials ~/.aws/credentials

# Remove a secret
dotpilot secrets remove aws_credentials
```

DotPilot will use GPG if available on your system, or fall back to AES-256 encryption if GPG is not available.

### Advanced SOPS/GPG Integration

For enhanced security with Mozilla SOPS and GPG:

```bash
# Add a file with SOPS encryption
dotpilot sops add ~/.aws/credentials

# Add with custom name
dotpilot sops add ~/.ssh/id_rsa --name ssh_key

# Add and immediately edit the encrypted file
dotpilot sops add ~/.npmrc --edit

# List SOPS encrypted secrets
dotpilot sops list

# Decrypt and retrieve a secret
dotpilot sops get aws_credentials ~/.aws/credentials

# Edit an encrypted secret directly
dotpilot sops edit aws_credentials

# Remove a SOPS secret
dotpilot sops remove aws_credentials
```

The SOPS integration offers several advantages:
- In-place editing of encrypted files
- Compatible with Mozilla SOPS CLI tool
- Uses your GPG keys for strong encryption
- Human-readable encrypted files (JSON format)
- Supports team-based secret sharing when using multiple GPG keys

Requirements:
- GPG must be installed with a key generated
- SOPS must be installed (https://github.com/mozilla/sops)

## Advanced Features

### Animated Progress Indicators

DotPilot provides animated progress indicators for long-running operations:

```bash
# Normal sync with animated progress
dotpilot sync

# Disable animated progress if needed
dotpilot sync --no-progress

# Progress indicators are also available for SOPS operations
dotpilot sops add ~/.aws/credentials
dotpilot sops get credentials ~/.aws/credentials

# Test the progress indicators
dotpilot test progress
```

The progress indicators provide real-time visual feedback for:
- Git operations (pull, push, commit)
- File sync operations
- Encryption and decryption
- Conflict resolution
- Configuration application

#### Progress Indicator Types

DotPilot implements four styles of animated progress indicators:

1. **Spinner**: A rotating animation that spins continuously
   - Best for: Operations with unknown completion time
   - Used in: Commit operations, hooks execution

2. **Bar**: A horizontal bar that fills from left to right
   - Best for: Operations with measurable progress
   - Used in: Configuration application, file sync operations

3. **Bounce**: A dot that bounces back and forth 
   - Best for: Network operations
   - Used in: Git pull/push operations

4. **Dots**: Text followed by animated dots
   - Best for: Encryption/decryption operations
   - Used in: SOPS file processing

### Conflict Resolution

DotPilot provides advanced conflict resolution strategies for handling file conflicts:

- `interactive`: Prompt for each conflict with options
- `keep-local`: Always keep local versions
- `keep-remote`: Always use tracked versions
- `merge`: Attempt to merge with a diff tool
- `backup-both`: Keep both versions

This helps to safely handle conflicting changes that might occur when syncing across multiple machines.

## Shell Completion

DotPilot provides smart command-line completion for various shells to enhance productivity:

```bash
# Generate bash completion
dotpilot completion bash > ~/.bash_completion.d/dotpilot

# Generate zsh completion
dotpilot completion zsh > "${fpath[1]}/_dotpilot"

# Generate fish completion
dotpilot completion fish > ~/.config/fish/completions/dotpilot.fish

# Generate PowerShell completion
dotpilot completion powershell > dotpilot.ps1
```

This enables context-aware completion for:
- File paths when tracking files
- Available environments
- Conflict resolution strategies
- Secret names when accessing encrypted secrets
- Package management systems

## License

MIT

