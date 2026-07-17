package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main()

// fetches origin/main and returns its current commit SHA.
func fetchTargetCommit(repoPath string) (target string, err error) {
	cmd := exec.Command("git", "fetch", "origin", "main")
	cmd.Dir = repoPath
	err = cmd.Run()
	if err != nil {
		return "", err
	}
	cmd = exec.Command("git", "rev-parse", "--verify", "origin/main")
	cmd.Dir = repoPath
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func handleHead(repoPath string, target string, applied string) error {
	cmd := exec.Command("git", "diff", "--name-only", target, applied, "--", "compose.yaml", "nginx/conf.d/")
	cmd.Dir = repoPath
	err := cmd.Run()
	if err != nil {
		return err
	}
	cmd = exec.Command("rm", "-rf", "/tmp/staging")
	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = exec.Command("mkdir", "-p", "/tmp/staging")
	if err := cmd.Run(); err != nil {
		return err
	}

}
