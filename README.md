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

1. Clone this repository
2. Build the binary with Go:
   ```
   go build -o dotpilot
   ```
3. Move the binary to a location in your PATH:
   ```
   mv dotpilot /usr/local/bin/
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

DotPilot can securely store sensitive configuration files using encryption:

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

## Advanced Features

### Conflict Resolution

DotPilot provides advanced conflict resolution strategies for handling file conflicts:

- `interactive`: Prompt for each conflict with options
- `keep-local`: Always keep local versions
- `keep-remote`: Always use tracked versions
- `merge`: Attempt to merge with a diff tool
- `backup-both`: Keep both versions

This helps to safely handle conflicting changes that might occur when syncing across multiple machines.

## License

MIT

