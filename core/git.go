package core

import (
        "fmt"
        "os"
        "path/filepath"
        "time"

        "github.com/dotpilot/utils"
        "github.com/go-git/go-git/v5"
        "github.com/go-git/go-git/v5/plumbing"
        "github.com/go-git/go-git/v5/plumbing/object"
        // "github.com/go-git/go-git/v5/plumbing/transport/http"
)

// RemoteStatus represents the status of the local repository compared to the remote
type RemoteStatus struct {
        Ahead  int
        Behind int
}

// InitializeRepo initializes the dotpilot repository
func InitializeRepo(remoteURL, dotpilotDir, environment string) error {
        // Create directory if it doesn't exist
        if err := os.MkdirAll(dotpilotDir, 0755); err != nil {
                return err
        }

        // Clone repository
        utils.Logger.Debug().Msgf("Cloning repository %s to %s", remoteURL, dotpilotDir)
        _, err := git.PlainClone(dotpilotDir, false, &git.CloneOptions{
                URL:      remoteURL,
                Progress: os.Stdout,
        })

        if err != nil {
                // If the repository doesn't exist, initialize a new one
                if err == git.ErrRepositoryAlreadyExists {
                        utils.Logger.Debug().Msg("Repository already exists, skipping clone")
                } else if err == git.ErrRepositoryNotExists {
                        utils.Logger.Debug().Msg("Remote repository doesn't exist, initializing new one")
                        
                        // Initialize new repo
                        repo, err := git.PlainInit(dotpilotDir, false)
                        if err != nil {
                                return err
                        }

                        // Create default directory structure
                        createDirStructure(dotpilotDir)

                        // Add remote
                        _, err = repo.CreateRemote(&git.RemoteConfig{
                                Name: "origin",
                                URLs: []string{remoteURL},
                        })
                        if err != nil {
                                return err
                        }

                        // Initial commit
                        w, err := repo.Worktree()
                        if err != nil {
                                return err
                        }

                        _, err = w.Add(".")
                        if err != nil {
                                return err
                        }

                        _, err = w.Commit("Initial commit", &git.CommitOptions{
                                Author: &object.Signature{
                                        Name:  "dotpilot",
                                        Email: "dotpilot@local",
                                        When:  time.Now(),
                                },
                        })
                        if err != nil {
                                return err
                        }
                } else {
                        return err
                }
        }

        // Create dotpilotrc file
        return CreateDefaultConfigFile(remoteURL, environment)
}

// createDirStructure creates the default directory structure for dotpilot
func createDirStructure(dotpilotDir string) error {
        // Create common directory
        if err := os.MkdirAll(filepath.Join(dotpilotDir, "common"), 0755); err != nil {
                return err
        }

        // Create envs directory
        if err := os.MkdirAll(filepath.Join(dotpilotDir, "envs", "default"), 0755); err != nil {
                return err
        }

        // Create machine directory with hostname
        hostname, err := os.Hostname()
        if err != nil {
                hostname = "unknown"
        }
        if err := os.MkdirAll(filepath.Join(dotpilotDir, "machine", hostname), 0755); err != nil {
                return err
        }

        // Create README
        readmePath := filepath.Join(dotpilotDir, "README.md")
        readmeContent := `# Dotfiles managed by DotPilot

This repository contains dotfiles managed by the DotPilot tool.

## Structure

- common/ - Files common to all environments
- envs/ - Environment-specific configurations
- machine/ - Machine-specific configurations
`
        if err := os.WriteFile(readmePath, []byte(readmeContent), 0644); err != nil {
                return err
        }

        return nil
}

// CommitChanges commits the changes in the repository with the given message
func CommitChanges(dotpilotDir, message string) error {
        // Open repository
        repo, err := git.PlainOpen(dotpilotDir)
        if err != nil {
                return err
        }

        // Get worktree
        w, err := repo.Worktree()
        if err != nil {
                return err
        }

        // Add all changes
        _, err = w.Add(".")
        if err != nil {
                return err
        }

        // Commit
        _, err = w.Commit(message, &git.CommitOptions{
                Author: &object.Signature{
                        Name:  "dotpilot",
                        Email: "dotpilot@local",
                        When:  time.Now(),
                },
        })
        if err != nil {
                return err
        }

        return nil
}

// HasUncommittedChanges checks if there are uncommitted changes in the repository
func HasUncommittedChanges(dotpilotDir string) (bool, error) {
        // Open repository
        repo, err := git.PlainOpen(dotpilotDir)
        if err != nil {
                return false, err
        }

        // Get worktree
        w, err := repo.Worktree()
        if err != nil {
                return false, err
        }

        // Get status
        status, err := w.Status()
        if err != nil {
                return false, err
        }

        return !status.IsClean(), nil
}

// PullChanges pulls changes from the remote
func PullChanges(dotpilotDir string) error {
        // Open repository
        repo, err := git.PlainOpen(dotpilotDir)
        if err != nil {
                return err
        }

        // Get worktree
        w, err := repo.Worktree()
        if err != nil {
                return err
        }

        // Pull
        err = w.Pull(&git.PullOptions{
                RemoteName: "origin",
                Progress:   os.Stdout,
        })

        if err != nil && err != git.NoErrAlreadyUpToDate {
                return err
        }

        return nil
}

// PushChanges pushes changes to the remote
func PushChanges(dotpilotDir string) error {
        // Open repository
        repo, err := git.PlainOpen(dotpilotDir)
        if err != nil {
                return err
        }

        // Push
        err = repo.Push(&git.PushOptions{
                RemoteName: "origin",
                Progress:   os.Stdout,
        })

        if err != nil && err != git.NoErrAlreadyUpToDate {
                return err
        }

        return nil
}

// GetGitStatus returns a string representation of the git status
func GetGitStatus(dotpilotDir string) (string, error) {
        // Open repository
        repo, err := git.PlainOpen(dotpilotDir)
        if err != nil {
                return "", err
        }

        // Get worktree
        w, err := repo.Worktree()
        if err != nil {
                return "", err
        }

        // Get status
        status, err := w.Status()
        if err != nil {
                return "", err
        }

        return status.String(), nil
}

// GetRemoteStatus returns the status of the local repository compared to the remote
func GetRemoteStatus(dotpilotDir string) (RemoteStatus, error) {
        result := RemoteStatus{
                Ahead:  0,
                Behind: 0,
        }

        // Open repository
        repo, err := git.PlainOpen(dotpilotDir)
        if err != nil {
                return result, err
        }

        // Get reference to HEAD
        head, err := repo.Head()
        if err != nil {
                return result, err
        }

        // Get remote reference
        remoteRef, err := repo.Reference(plumbing.NewRemoteReferenceName("origin", head.Name().Short()), true)
        if err != nil {
                return result, err
        }

        // Count commits ahead and behind
        revList, err := repo.Log(&git.LogOptions{
                From:  head.Hash(),
                Order: git.LogOrderCommitterTime,
        })
        if err != nil {
                return result, err
        }

        // Count commits ahead
        err = revList.ForEach(func(c *object.Commit) error {
                if c.Hash == remoteRef.Hash() {
                        return fmt.Errorf("stop")
                }
                result.Ahead++
                return nil
        })
        if err != nil && err.Error() != "stop" {
                return result, err
        }

        // Count commits behind
        revList, err = repo.Log(&git.LogOptions{
                From:  remoteRef.Hash(),
                Order: git.LogOrderCommitterTime,
        })
        if err != nil {
                return result, err
        }

        err = revList.ForEach(func(c *object.Commit) error {
                if c.Hash == head.Hash() {
                        return fmt.Errorf("stop")
                }
                result.Behind++
                return nil
        })
        if err != nil && err.Error() != "stop" {
                return result, err
        }

        return result, nil
}

// GetTrackedFiles returns a list of files tracked by dotpilot
func GetTrackedFiles(dotpilotDir string) ([]string, error) {
        var trackedFiles []string

        // Open repository
        repo, err := git.PlainOpen(dotpilotDir)
        if err != nil {
                return nil, err
        }

        // Get HEAD reference
        ref, err := repo.Head()
        if err != nil {
                return nil, err
        }

        // Get commit
        commit, err := repo.CommitObject(ref.Hash())
        if err != nil {
                return nil, err
        }

        // Get tree
        tree, err := commit.Tree()
        if err != nil {
                return nil, err
        }

        // Walk the tree
        err = tree.Files().ForEach(func(f *object.File) error {
                // Skip .git directory and README.md
                if f.Name == "README.md" {
                        return nil
                }

                trackedFiles = append(trackedFiles, f.Name)
                return nil
        })
        if err != nil {
                return nil, err
        }

        return trackedFiles, nil
}
