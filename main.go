package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	git_transport "github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/urfave/cli/v2"
)

func main() {
	// Test
	app := &cli.App{
		Name:                 "permalinker",
		Usage:                "Generates a permalink URL with the given file and (optional) line number.",
		Description:          "The command locates the git repository which contains the given file. The generated URL combines the following data:\n- The repository HTTPS URL (based on remote \"origin\")\n- The repository HEAD commit hash\n- The file path relative to the repository\n- The (optional) file line number.",
		Version:              "0.1.0",
		HideHelpCommand:      true,
		ArgsUsage:            "<file> [line-number]",
		Action:               permalinker,
		EnableBashCompletion: true,
		// Flags: []cli.Flag{
		// 	&cli.StringFlag{
		// 		Name:    "commit",
		// 		Aliases: []string{"c"},
		// 		Usage:   "Override commit hash used in generated URL",
		// 	},
		// },
	}

	if err := app.Run(os.Args); err != nil {
		log.SetFlags(0)
		log.SetOutput(os.Stderr)
		log.Fatal(err)
	}
}

func permalinker(cCtx *cli.Context) error {
	fileName := cCtx.Args().Get(0)
	if fileName == "" {
		return cli.ShowAppHelp(cCtx)
	}
	lineNumber := cCtx.Args().Get(1)
	if lineNumber != "" {
		_, err := strconv.ParseUint(lineNumber, 10, 64)
		if err != nil {
			return err
		}
	}

	opts := &git.PlainOpenOptions{
		DetectDotGit: true,
	}
	repo, err := git.PlainOpenWithOptions(fileName, opts)
	if err != nil {
		return err
	}

	repoURL, err := getRepoURL(repo)
	if err != nil {
		return err
	}

	repoHeadCommitHash, err := getRepoHeadCommitHash(repo)
	if err != nil {
		return err
	}

	fileNameRelativePath, err := getPathRelativeToGitRepo(fileName)
	if err != nil {
		return err
	}

	fmt.Printf("%s/blob/%s/%s\n", repoURL, repoHeadCommitHash, convertPathForURL(fileNameRelativePath, lineNumber))
	return nil
}

func getRepoURL(repo *git.Repository) (string, error) {
	repoConfig, err := repo.Config()
	if err != nil {
		return "", err
	}
	remoteConfig, ok := repoConfig.Remotes["origin"]
	if !ok {
		return "", errors.New("unable to find remote \"origin\"")
	}
	repoConfigURL := remoteConfig.URLs[0]
	repoURL, err := git_transport.NewEndpoint(repoConfigURL)
	if err != nil {
		return "", err
	}
	repoURL.Protocol = "https"
	repoURL.User = ""
	repoURL.Port = 0
	repoURL.Path = strings.TrimSuffix(repoURL.Path, git.GitDirName)
	return repoURL.String(), nil
}

func getRepoHeadCommitHash(repo *git.Repository) (plumbing.Hash, error) {
	pRef, err := repo.Head()
	if err != nil {
		return plumbing.ZeroHash, err
	}
	return pRef.Hash(), nil
}

func getPathRelativeToGitRepo(path string) (string, error) {
	fileNameAbsPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	repoAbsPath, err := getGitDirPath(fileNameAbsPath)
	if err != nil {
		return "", err
	}
	return filepath.Rel(repoAbsPath, fileNameAbsPath)
}

func getGitDirPath(absPath string) (string, error) {
	originalPath := absPath
	stat, err := os.Stat(absPath)
	if err != nil {
		return "", err
	}
	if !stat.IsDir() {
		absPath = filepath.Dir(absPath)
	}
	for {
		_, err := os.Stat(filepath.Join(absPath, git.GitDirName))
		if err == nil {
			// no error; stop
			break
		}
		if !os.IsNotExist(err) {
			// unknown error; stop
			return "", err
		}
		if dir := filepath.Dir(absPath); dir != absPath {
			absPath = dir
			continue
		}
		return "", errors.New("unable to find git repo in parent directories of " + originalPath)
	}
	return absPath, nil
}

func convertPathForURL(path string, lineNumber string) string {
	path = filepath.ToSlash(path)
	if lineNumber != "" {
		// Ensure UI is able to link to the line number by showing the plain
		// version of the file.
		if strings.HasSuffix(path, ".md") {
			path = path + "?plain=1"
		}
		path = path + "#L" + lineNumber
	}
	return path
}
