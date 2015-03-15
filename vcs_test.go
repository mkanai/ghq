package main

import (
	"io/ioutil"
	"net/url"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/motemen/ghq/utils"
	. "github.com/onsi/gomega"
)

func TestGitBackend(t *testing.T) {
	RegisterTestingT(t)

	tempDir, err := ioutil.TempDir("", "ghq-test")
	if err != nil {
		t.Fatal(err)
	}

	localDir := filepath.Join(tempDir, "repo")

	remoteURL, err := url.Parse("https://example.com/git/repo")
	if err != nil {
		t.Fatal(err)
	}

    branch := ""

	commands := []*exec.Cmd{}
	lastCommand := func() *exec.Cmd { return commands[len(commands)-1] }
	utils.CommandRunner = func(cmd *exec.Cmd) error {
		commands = append(commands, cmd)
		return nil
	}

	err = GitBackend.Clone(remoteURL, localDir, branch, false, false)

	Expect(err).NotTo(HaveOccurred())
	Expect(commands).To(HaveLen(1))
	Expect(lastCommand().Args).To(Equal([]string{
		"git", "clone", remoteURL.String(), localDir,
	}))

	err = GitBackend.Clone(remoteURL, localDir, branch, true, false)

	Expect(err).NotTo(HaveOccurred())
	Expect(commands).To(HaveLen(2))
	Expect(lastCommand().Args).To(Equal([]string{
		"git", "clone", "--depth", "1", remoteURL.String(), localDir,
	}))

    err = GitBackend.Clone(remoteURL, localDir, branch, false, true)
    Expect(err).NotTo(HaveOccurred())
    Expect(commands).To(HaveLen(3))
    Expect(lastCommand().Args).To(Equal([]string{
        "git", "clone", "--recursive", remoteURL.String(), localDir,
    }))

    branch = "develop"
    err = GitBackend.Clone(remoteURL, localDir, branch, false, false)

    Expect(err).NotTo(HaveOccurred())
    Expect(commands).To(HaveLen(4))
    Expect(lastCommand().Args).To(Equal([]string{
        "git", "clone", "--branch", branch, remoteURL.String(), localDir,
    }))


	err = GitBackend.Update(localDir)

	Expect(err).NotTo(HaveOccurred())
	Expect(commands).To(HaveLen(5))
	Expect(lastCommand().Args).To(Equal([]string{
		"git", "pull", "--ff-only",
	}))
	Expect(lastCommand().Dir).To(Equal(localDir))
}

func TestSubversionBackend(t *testing.T) {
	RegisterTestingT(t)

	tempDir, err := ioutil.TempDir("", "ghq-test")
	if err != nil {
		t.Fatal(err)
	}

	localDir := filepath.Join(tempDir, "repo")

	remoteURL, err := url.Parse("https://example.com/git/repo")
	if err != nil {
		t.Fatal(err)
	}

    branch := ""

	commands := []*exec.Cmd{}
	lastCommand := func() *exec.Cmd { return commands[len(commands)-1] }
	utils.CommandRunner = func(cmd *exec.Cmd) error {
		commands = append(commands, cmd)
		return nil
	}

	err = SubversionBackend.Clone(remoteURL, localDir, branch, false, false)

	Expect(err).NotTo(HaveOccurred())
	Expect(commands).To(HaveLen(1))
	Expect(lastCommand().Args).To(Equal([]string{
		"svn", "checkout", remoteURL.String(), localDir,
	}))

	err = SubversionBackend.Clone(remoteURL, localDir, branch, true, false)

	Expect(err).NotTo(HaveOccurred())
	Expect(commands).To(HaveLen(2))
	Expect(lastCommand().Args).To(Equal([]string{
		"svn", "checkout", "--depth", "1", remoteURL.String(), localDir,
	}))

	err = SubversionBackend.Clone(remoteURL, localDir, branch, false, true)

	Expect(err).NotTo(HaveOccurred())
	Expect(commands).To(HaveLen(3))
	Expect(lastCommand().Args).To(Equal([]string{
		"svn", "checkout", remoteURL.String(), localDir,
	}))

	err = SubversionBackend.Update(localDir)

	Expect(err).NotTo(HaveOccurred())
	Expect(commands).To(HaveLen(4))
	Expect(lastCommand().Args).To(Equal([]string{
		"svn", "update",
	}))
	Expect(lastCommand().Dir).To(Equal(localDir))
}

func TestGitsvnBackend(t *testing.T) {
	RegisterTestingT(t)

	tempDir, err := ioutil.TempDir("", "ghq-test")
	if err != nil {
		t.Fatal(err)
	}

	localDir := filepath.Join(tempDir, "repo")

	remoteURL, err := url.Parse("https://example.com/git/repo")
	if err != nil {
		t.Fatal(err)
	}

    branch := ""

	commands := []*exec.Cmd{}
	lastCommand := func() *exec.Cmd { return commands[len(commands)-1] }
	utils.CommandRunner = func(cmd *exec.Cmd) error {
		commands = append(commands, cmd)
		return nil
	}

	err = GitsvnBackend.Clone(remoteURL, localDir, branch, false, false)

	Expect(err).NotTo(HaveOccurred())
	Expect(commands).To(HaveLen(1))
	Expect(lastCommand().Args).To(Equal([]string{
		"git", "svn", "clone", remoteURL.String(), localDir,
	}))

	err = GitsvnBackend.Clone(remoteURL, localDir, branch, true, false)

	Expect(err).NotTo(HaveOccurred())
	Expect(commands).To(HaveLen(2))
	Expect(lastCommand().Args).To(Equal([]string{
		"git", "svn", "clone", remoteURL.String(), localDir,
	}))

	err = GitsvnBackend.Clone(remoteURL, localDir, branch, false, true)

	Expect(err).NotTo(HaveOccurred())
	Expect(commands).To(HaveLen(3))
	Expect(lastCommand().Args).To(Equal([]string{
		"git", "svn", "clone", remoteURL.String(), localDir,
	}))

	err = GitsvnBackend.Update(localDir)

	Expect(err).NotTo(HaveOccurred())
	Expect(commands).To(HaveLen(4))
	Expect(lastCommand().Args).To(Equal([]string{
		"git", "svn", "rebase",
	}))
	Expect(lastCommand().Dir).To(Equal(localDir))
}

func TestMercurialBackend(t *testing.T) {
	RegisterTestingT(t)

	tempDir, err := ioutil.TempDir("", "ghq-test")
	if err != nil {
		t.Fatal(err)
	}

	localDir := filepath.Join(tempDir, "repo")

	remoteURL, err := url.Parse("https://example.com/git/repo")
	if err != nil {
		t.Fatal(err)
	}

    branch := ""

	commands := []*exec.Cmd{}
	lastCommand := func() *exec.Cmd { return commands[len(commands)-1] }
	utils.CommandRunner = func(cmd *exec.Cmd) error {
		commands = append(commands, cmd)
		return nil
	}

	err = MercurialBackend.Clone(remoteURL, localDir, branch, false, false)

	Expect(err).NotTo(HaveOccurred())
	Expect(commands).To(HaveLen(1))
	Expect(lastCommand().Args).To(Equal([]string{
		"hg", "clone", remoteURL.String(), localDir,
	}))

	err = MercurialBackend.Clone(remoteURL, localDir, branch, true, false)

	Expect(err).NotTo(HaveOccurred())
	Expect(commands).To(HaveLen(2))
	Expect(lastCommand().Args).To(Equal([]string{
		"hg", "clone", remoteURL.String(), localDir,
	}))

	err = MercurialBackend.Clone(remoteURL, localDir, branch, false, true)

	Expect(err).NotTo(HaveOccurred())
	Expect(commands).To(HaveLen(3))
	Expect(lastCommand().Args).To(Equal([]string{
		"hg", "clone", remoteURL.String(), localDir,
	}))

	err = MercurialBackend.Update(localDir)

	Expect(err).NotTo(HaveOccurred())
	Expect(commands).To(HaveLen(4))
	Expect(lastCommand().Args).To(Equal([]string{
		"hg", "pull", "--update",
	}))
	Expect(lastCommand().Dir).To(Equal(localDir))
}
