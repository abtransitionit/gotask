// File gotask/phase/show.go
package workflow

import (
	"fmt"

	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/phase"
)

// Name: TestVarInCtx
//
// Description: It is intended to receive a var in the context and display it
//
// Parameters:
//
//   - ctx: The context for the phase.
//   - logger: The logger to use for printing messages.
//   - workflow: The workflow to be executed.
//   - tiers: The tiers of phases to be executed.
//   - toolConfig: The configuration for the tool.
//
// Returns:
//
//   - A string containing the workflow configuration and execution plan.
//   - An error if the phase fails to execute.
//
// Notes:
//
//   - It is intended to receive a var in the context and display it
//   - It is a proof of concept that phase can pass any var to library task
func TestVarInCtx(workflow *phase.Workflow, logger logx.Logger) (string, error) {

	// check parameters
	if workflow == nil {
		return "", fmt.Errorf("workflow cannot be nil")
	}

	if logger == nil {
		return "", fmt.Errorf("logger cannot be nil")
	}

	return "", nil
}
