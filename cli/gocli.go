// File gotask/cli/gocli.go
package cli

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/abtransitionit/gocore/logx"
)

// BuildProject builds a Go project. It conforms to the PhaseFunc signature.
// The `cmd` arguments are expected to be:
// 1. projectPath (string)
// 2. outputDir (string)
func BuildProject(ctx context.Context, l logx.Logger, targets []Target, cmd ...string) (string, error) {
	// The logger is now passed in explicitly.
	logger := l

	if len(cmd) < 2 {
		return "", fmt.Errorf("BuildProject requires at least two arguments: projectPath and outputDir")
	}

	projectPath := cmd[0]
	outputDir := cmd[1]

	// 1. Get the name of the project from the path.
	projectName := filepath.Base(projectPath)
	outputFile := filepath.Join(outputDir, projectName)

	logger.Infof("Building project '%s' from path: %s", projectName, projectPath)
	logger.Infof("Saving artifact to: %s", outputFile)

	// 2. Execute the `go build` command.
	buildArgs := []string{"build", "-o", outputFile, projectPath}
	logger.Infof("Executing command: go %s", buildArgs)

	buildCmd := exec.CommandContext(ctx, "go", buildArgs...)
	buildCmd.Stdout = logger.Info
	buildCmd.Stderr = logger.Error

	if err := buildCmd.Run(); err != nil {
		return "", fmt.Errorf("build failed: %w", err)
	}

	logger.Infof("Successfully built project to: %s", outputFile)

	return outputFile, nil
}
