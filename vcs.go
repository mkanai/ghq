package main

import (
	"net/url"
	"os"
	"path/filepath"

	"github.com/motemen/ghq/utils"
)

// A VCSBackend represents a VCS backend.
type VCSBackend struct {
	// Clones a remote repository to local path.
	Clone func(*url.URL, string, string, bool, bool) error
	// Updates a cloned local repository.
	Update func(string) error
}

var GitBackend = &VCSBackend{
	Clone: func(remote *url.URL, local string, branch string, shallow bool, recursive bool) error {
		dir, _ := filepath.Split(local)
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}

		args := []string{"clone"}
        if branch != "" {
            args = append(args, "--branch", branch)
        }
		if shallow {
			args = append(args, "--depth", "1")
		}
        if recursive {
            args = append(args, "--recursive")
        }
		args = append(args, remote.String(), local)

		return utils.Run("git", args...)
	},
	Update: func(local string) error {
		return utils.RunInDir(local, "git", "pull", "--ff-only")
	},
}

var SubversionBackend = &VCSBackend{
	Clone: func(remote *url.URL, local string, ignoredBranch string, shallow bool, ignoredRecursive bool) error {
		dir, _ := filepath.Split(local)
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}

		args := []string{"checkout"}
		if shallow {
			args = append(args, "--depth", "1")
		}
		args = append(args, remote.String(), local)

		return utils.Run("svn", args...)
	},
	Update: func(local string) error {
		return utils.RunInDir(local, "svn", "update")
	},
}

var GitsvnBackend = &VCSBackend{
	// git-svn seems not supporting shallow clone currently.
	Clone: func(remote *url.URL, local string, ignoredBranch string, ignoredShallow bool, ignoredRecursive bool) error {
		dir, _ := filepath.Split(local)
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}

		return utils.Run("git", "svn", "clone", remote.String(), local)
	},
	Update: func(local string) error {
		return utils.RunInDir(local, "git", "svn", "rebase")
	},
}

var MercurialBackend = &VCSBackend{
	// Mercurial seems not supporting shallow clone currently.
	Clone: func(remote *url.URL, local string, ignoredBranch string, ignoredShallow bool, ignoredRecursive bool) error {
		dir, _ := filepath.Split(local)
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}

		return utils.Run("hg", "clone", remote.String(), local)
	},
	Update: func(local string) error {
		return utils.RunInDir(local, "hg", "pull", "--update")
	},
}
