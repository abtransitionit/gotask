// File gotask/phase/show.go
package workflow

import (
	"context"
	"fmt"

	gocorectx "github.com/abtransitionit/gocore/ctx"
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
func ShowPhase(workflow *phase.Workflow, l logx.Logger) error {

	// check parameters
	if workflow == nil {
		return fmt.Errorf("workflow cannot be nil")
	}

	// show all workflow phases
	workflow.Show(l)
	return nil
}
func ShowFiltered(workflow *phase.Workflow, l logx.Logger, ctx context.Context, skipPhases []int) error {

	// check parameters
	if workflow == nil {
		return fmt.Errorf("workflow cannot be nil")
	}

	// get phases topoSorted
	PhaseSortedByTier, err := workflow.TopoSort(ctx)
	if err != nil {
		l.ErrorWithStack(err, "failed to sort phases")
		return err
	}
	// filter them
	l.Info("filtered the tiers")
	PhaseFilteredByTier := PhaseSortedByTier.Filter(*workflow, l, skipPhases)

	// show them
	l.Info("list of filtered phases")
	PhaseFilteredByTier.Show(l)
	return nil

}

func ShowWorkflow(ctx context.Context, logger logx.Logger, cmd ...string) (string, error) {
	// Define a custom key type to avoid collisions.
	type contextKey string
	const mxExecutionKey contextKey = "executionID"

	logger.Info("From gotask/workflow : Displaying workflow execution plan:")
	// Get the var that was pass to the ctx and convert it to string
	// execID, ok := ctx.Value(gocorectx.ExecutionIDKey).(string)
	theWrkfl, ok := ctx.Value(gocorectx.WorkflowKey).(*phase.Workflow)
	if !ok || theWrkfl == nil {
		return "", fmt.Errorf("failed to get executionID from context")
	}
	// logger.Infof("From gotask/workflow : wrkfl: %v", theWrkfl)
	theWrkfl.Show(logger)
	// wkf.Show(logger) // Assuming you have a Show() method on your Workflow struct
	// return execID, nil
	return "Workflow plan displayed successfully.", nil
}
