// file: caller_workflow.go
package workflows

import (
	"go.temporal.io/sdk/workflow"
)

// CallerWorkflow sends an updateState signal to the TargetWorkflow and then waits for a response.
func CallerWorkflow(ctx workflow.Context) error {
	logger := workflow.GetLogger(ctx)

	// Define the target workflow's IDs.
	targetWorkflowID := "entity_workflow_id"
	targetRunID := "" // If known, include the RunID; otherwise, it can be empty.

	// Prepare the payload for the updateState signal.
	signalData := "UpdatedState"
	// Send the updateState signal to the target workflow.
	err := workflow.SignalExternalWorkflow(
		ctx,
		targetWorkflowID,
		targetRunID,
		"updateState",
		signalData,
	).Get(ctx, nil)
	if err != nil {
		logger.Error("CallerWorkflow: Failed to signal target", "Error", err)
		return err
	}
	logger.Info("CallerWorkflow: updateState signal sent")

	// Now, wait for a response signal.
	responseCh := workflow.GetSignalChannel(ctx, "responseSignal")
	var response string
	responseCh.Receive(ctx, &response)
	logger.Info("CallerWorkflow: Received response", "response", response)

	// Continue with additional processing if needed.
	return nil
}
