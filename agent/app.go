package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: agenttest <repository-path>")
		os.Exit(2)
	}

	repoPath := os.Args[1]
	stagingPath := "/tmp/staging"

	target, err := fetchTargetCommit(repoPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "fetch target commit:", err)
		os.Exit(1)
	}

	fmt.Println("target commit:", target)

	if err := reconcileConfigs(repoPath, target, stagingPath); err != nil {
		fmt.Fprintln(os.Stderr, "stage target configuration:", err)
		os.Exit(1)
	}

	fmt.Println("staged configuration at:", stagingPath)

	if err := validateStaged(stagingPath); err != nil {
		fmt.Fprintln(os.Stderr, "validate staged configuration:", err)
		os.Exit(1)
	}

	fmt.Println("staged configuration is valid")
}

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

func managedPathsChanged(repoPath string, target string, applied string) (diff bool, err error) {
	cmd := exec.Command("git", "diff", "--quiet", applied, target, "--", "compose.yaml", "nginx/conf.d/")
	cmd.Dir = repoPath
	err = cmd.Run()
	if err == nil {
		return false, nil
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) && exitErr.ExitCode() == 1 {
		return true, nil
	}
	return false, err
}

func reconcileConfigs(repoPath, target, stagingPath string) error {
	if err := os.RemoveAll(stagingPath); err != nil {
		return fmt.Errorf("remove staging directory: %w", err)
	}

	if err := os.MkdirAll(stagingPath, 0o755); err != nil {
		return fmt.Errorf("create staging directory: %w", err)
	}

	archiveFile, err := os.CreateTemp("", "seal-configs-*.tar")
	if err != nil {
		return fmt.Errorf("create temporary archive: %w", err)
	}

	archivePath := archiveFile.Name()

	if err := archiveFile.Close(); err != nil {
		os.Remove(archivePath)
		return fmt.Errorf("close temporary archive: %w", err)
	}

	defer os.Remove(archivePath)

	archiveCmd := exec.Command(
		"git",
		"archive",
		"--format=tar",
		"--output="+archivePath,
		target,
		"compose.yaml",
		"nginx/conf.d/",
	)
	archiveCmd.Dir = repoPath

	if output, err := archiveCmd.CombinedOutput(); err != nil {
		return fmt.Errorf(
			"archive target configuration: %w: %s",
			err,
			strings.TrimSpace(string(output)),
		)
	}

	tarCmd := exec.Command(
		"tar",
		"-xf",
		archivePath,
		"-C",
		stagingPath,
	)

	if output, err := tarCmd.CombinedOutput(); err != nil {
		return fmt.Errorf(
			"extract target configuration: %w: %s",
			err,
			strings.TrimSpace(string(output)),
		)
	}

	return nil
}

// we validate compos config and also start a temporary nginx container to validate
func validateStaged(stagingPath string) error {
	composeCmd := exec.Command(
		"docker",
		"compose",
		"-f",
		stagingPath+"/compose.yaml",
		"config",
		"--quiet",
	)
	output, err := composeCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf(
			"validate staged compose: %w: %s",
			err,
			strings.TrimSpace(string(output)),
		)
	}

	nginxCmd := exec.Command(
		"docker",
		"run",
		"--rm",
		"--mount",
		"type=bind,src="+filepath.Join(stagingPath, "nginx", "conf.d")+",dst=/etc/nginx/conf.d,readonly",
		"nginx:1.30.3-alpine",
		"nginx",
		"-t",
	)
	output, err = nginxCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf(
			"validate nginx conf: %w: %s",
			err,
			strings.TrimSpace(string(output)),
		)
	}
	return nil
}
