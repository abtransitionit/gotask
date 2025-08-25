// File gotask/phase/show.go
package workflow

import (
	"context"
	"fmt"

	corectx "github.com/abtransitionit/gocore/ctx"
	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/phase"
)

// Name: ShowPhase
//
// Description: displays the workflow's configuration and execution plan.
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
//   - It is intended to be used as the first phase in a workflow
func ShowTier(workflow *phase.Workflow, l logx.Logger) (string, error) {

	// check parameters
	if workflow == nil {
		return "", fmt.Errorf("workflow cannot be nil")
	}

	// get phases sorted by tier
	PhaseSortedByTier, err := workflow.TopoSort(context.Background())
	if err != nil {
		return "", err
	}
	// show them
	l.Info("list of sorted phases")
	PhaseSortedByTier.Show(l)

	return "ShowTier phase complete", nil
}

// Name: ShowPhase
//
// Description: displays the workflow's phases
//
// Parameters:
//
//   - logger: The logger to use for printing messages.
//   - workflow: The workflow to be executed.
//
// Returns:
//
//   - An error if any
func ShowPhase(workflow *phase.Workflow, l logx.Logger) error {

	// check parameters
	if workflow == nil {
		return fmt.Errorf("workflow cannot be nil")
	}

	// show all workflow phases
	workflow.Show(l)
	return nil
}

// Name: ShowWorkflow
//
// Description: displays the workflow's phases and execution plan.
//
// Parameters:
//
//   - ctx: The context for the phase.
//   - logger: The logger to use for printing messages.
//   - ...string: The command to be executed.
//
// Returns:
//
//   - An error if any
//
// Notes:
//   - intented for test
//   - the ctx must have received the worflow
func ShowWorkflow(ctx context.Context, logger logx.Logger, targets []phase.Target, cmd ...string) (string, error) {

	logger.Info("From gotask/workflow : Displaying workflow execution plan:")
	// Get the var that was pass to the ctx and convert it
	wrkflw, ok := ctx.Value(corectx.WorkflowKeyId).(*phase.Workflow)
	if !ok || wrkflw == nil {
		return "", fmt.Errorf("from gotask/workflow : failed to get executionID from context")
	}
	wrkflw.Show(logger)
	return "", nil
}
