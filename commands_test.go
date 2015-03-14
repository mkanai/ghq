package main

import (
	"flag"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"

	"github.com/codegangsta/cli"
)

import (
	"testing"
	. "github.com/onsi/gomega"
)

func flagSet(name string, flags []cli.Flag) *flag.FlagSet {
	set := flag.NewFlagSet(name, flag.ContinueOnError)

	for _, f := range flags {
		f.Apply(set)
	}
	return set
}

func captureReader(block func()) (*os.File, *os.File, error) {
	rOut, wOut, err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}

	rErr, wErr, err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}

	defer wOut.Close()
	defer wErr.Close()

	var stdout, stderr *os.File
	os.Stdout, stdout = wOut, os.Stdout
	os.Stderr, stderr = wErr, os.Stderr

	defer func() {
		os.Stdout = stdout
		os.Stderr = stderr
	}()

	block()

	wOut.Close()
	wErr.Close()

	return rOut, rErr, nil
}

func capture(block func()) (string, string, error) {
	rOut, rErr, err := captureReader(block)
	if err != nil {
		return "", "", err
	}

	bufOut, err := ioutil.ReadAll(rOut)
	if err != nil {
		return "", "", err
	}

	bufErr, err := ioutil.ReadAll(rErr)
	if err != nil {
		return "", "", err
	}

	return string(bufOut), string(bufErr), nil
}

type _cloneArgs struct {
	remote  *url.URL
	local   string
	shallow bool
    recursive bool
}

type _updateArgs struct {
	local string
}

func withFakeGitBackend(t *testing.T, block func(string, *_cloneArgs, *_updateArgs)) {
	tmpRoot, err := ioutil.TempDir("", "ghq")
	if err != nil {
		t.Fatalf("Could not create tempdir: %s", err)
	}
	defer func() { os.RemoveAll(tmpRoot) }()

	_localRepositoryRoots = []string{tmpRoot}
	defer func() { _localRepositoryRoots = []string{} }()

	var cloneArgs _cloneArgs
	var updateArgs _updateArgs

	var originalGitBackend = GitBackend
	GitBackend = &VCSBackend{
		Clone: func(remote *url.URL, local string, shallow bool, recursive bool) error {
			cloneArgs = _cloneArgs{
				remote:  remote,
				local:   local,
				shallow: shallow,
                recursive: recursive,
			}
			return nil
		},
		Update: func(local string) error {
			updateArgs = _updateArgs{
				local: local,
			}
			return nil
		},
	}
	defer func() { GitBackend = originalGitBackend }()

	block(tmpRoot, &cloneArgs, &updateArgs)
}

func TestCommandGet(t *testing.T) {
	RegisterTestingT(t)

	app := newApp()

	withFakeGitBackend(t, func(tmpRoot string, cloneArgs *_cloneArgs, updateArgs *_updateArgs) {
		localDir := filepath.Join(tmpRoot, "github.com", "motemen", "ghq-test-repo")

		app.Run([]string{"", "get", "motemen/ghq-test-repo"})

		Expect(cloneArgs.remote.String()).To(Equal("https://github.com/motemen/ghq-test-repo"))
		Expect(filepath.ToSlash(cloneArgs.local)).To(Equal(filepath.ToSlash(localDir)))
		Expect(cloneArgs.shallow).To(Equal(false))
        Expect(cloneArgs.recursive).To(Equal(false))
	})

	withFakeGitBackend(t, func(tmpRoot string, cloneArgs *_cloneArgs, updateArgs *_updateArgs) {
		localDir := filepath.Join(tmpRoot, "github.com", "motemen", "ghq-test-repo")

		app.Run([]string{"", "get", "-p", "motemen/ghq-test-repo"})

		Expect(cloneArgs.remote.String()).To(Equal("ssh://git@github.com/motemen/ghq-test-repo"))
		Expect(filepath.ToSlash(cloneArgs.local)).To(Equal(filepath.ToSlash(localDir)))
	})

	withFakeGitBackend(t, func(tmpRoot string, cloneArgs *_cloneArgs, updateArgs *_updateArgs) {
		localDir := filepath.Join(tmpRoot, "github.com", "motemen", "ghq-test-repo")
		// mark as "already cloned", the condition may change later
		os.MkdirAll(filepath.Join(localDir, ".git"), 0755)

		app.Run([]string{"", "get", "-u", "motemen/ghq-test-repo"})

		Expect(updateArgs.local).To(Equal(localDir))
	})

	withFakeGitBackend(t, func(tmpRoot string, cloneArgs *_cloneArgs, updateArgs *_updateArgs) {
		localDir := filepath.Join(tmpRoot, "github.com", "motemen", "ghq-test-repo")
		app.Run([]string{"", "get", "-p", "motemen/ghq-test-repo"})

		Expect(cloneArgs.remote.String()).To(Equal("ssh://git@github.com/motemen/ghq-test-repo"))
		Expect(filepath.ToSlash(cloneArgs.local)).To(Equal(filepath.ToSlash(localDir)))
		Expect(cloneArgs.shallow).To(Equal(false))
        Expect(cloneArgs.recursive).To(Equal(false))
	})

	withFakeGitBackend(t, func(tmpRoot string, cloneArgs *_cloneArgs, updateArgs *_updateArgs) {
		localDir := filepath.Join(tmpRoot, "github.com", "motemen", "ghq-test-repo")

		app.Run([]string{"", "get", "-shallow", "motemen/ghq-test-repo"})

		Expect(cloneArgs.remote.String()).To(Equal("https://github.com/motemen/ghq-test-repo"))
		Expect(filepath.ToSlash(cloneArgs.local)).To(Equal(filepath.ToSlash(localDir)))
		Expect(cloneArgs.shallow).To(Equal(true))
        Expect(cloneArgs.recursive).To(Equal(false))
	})

	withFakeGitBackend(t, func(tmpRoot string, cloneArgs *_cloneArgs, updateArgs *_updateArgs) {
		localDir := filepath.Join(tmpRoot, "github.com", "motemen", "ghq-test-repo")

		app.Run([]string{"", "get", "-recursive", "motemen/ghq-test-repo"})

		Expect(cloneArgs.remote.String()).To(Equal("https://github.com/motemen/ghq-test-repo"))
		Expect(filepath.ToSlash(cloneArgs.local)).To(Equal(filepath.ToSlash(localDir)))
		Expect(cloneArgs.shallow).To(Equal(false))
        Expect(cloneArgs.recursive).To(Equal(true))
	})
}

func TestCommandList(t *testing.T) {
	RegisterTestingT(t)

	_, _, err := capture(func() {
		app := cli.NewApp()
		flagSet := flagSet("list", commandList.Flags)
		c := cli.NewContext(app, flagSet, flagSet)

		doList(c)
	})

	Expect(err).To(BeNil())
}

func TestCommandListUnique(t *testing.T) {
	RegisterTestingT(t)

	_, _, err := capture(func() {
		app := cli.NewApp()
		flagSet := flagSet("list", commandList.Flags)
		flagSet.Parse([]string{"--unique"})
		c := cli.NewContext(app, flagSet, flagSet)

		doList(c)
	})

	Expect(err).To(BeNil())
}

func TestCommandListUnknown(t *testing.T) {
	RegisterTestingT(t)

	_, _, err := capture(func() {
		app := cli.NewApp()
		flagSet := flagSet("list", commandList.Flags)
		flagSet.Parse([]string{"--unknown-flag"})
		c := cli.NewContext(app, flagSet, flagSet)

		doList(c)
	})

	Expect(err).To(BeNil())
}
