// file: entity_workflow.go
package workflows

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

// EntityWorkflow represents a long-running workflow that maintains state.
func EntityWorkflow(ctx workflow.Context) error {
	logger := workflow.GetLogger(ctx)
	state := "initial"
	shutdownRequested := false

	// Signal channels.
	updateStateCh := workflow.GetSignalChannel(ctx, "updateState")
	shutdownCh := workflow.GetSignalChannel(ctx, "shutdown")

	// For sending a response back to a caller, we assume the caller
	// passed its workflow ID (and optionally runID) along with the update signal.
	var callerWorkflowID string

	// Run indefinitely until a shutdown signal is received.
	for !shutdownRequested {
		selector := workflow.NewSelector(ctx)

		// Handle updateState signal.
		selector.AddReceive(updateStateCh, func(c workflow.ReceiveChannel, more bool) {
			var signalData struct {
				State    string
				CallerID string
				RunID    string
			}
			c.Receive(ctx, &signalData)
			state = signalData.State
			logger.Info("EntityWorkflow: State updated", "state", state)

			// Optionally, if the update signal payload includes the caller's workflow ID,
			// capture it here. For a simple example, we assume it's already set.
			// In a real scenario, you might have a composite struct for the signal data.

			// Send back a response signal to the caller.
			// Ensure you have a valid callerWorkflowID (and optionally runID).

			callerWorkflowID = signalData.CallerID

			if callerWorkflowID != "" {
				responseValue := "State updated to " + state
				err := workflow.SignalExternalWorkflow(
					ctx,
					signalData.CallerID, // Caller workflow ID
					signalData.RunID,    // Caller runID if known, or leave empty
					"responseSignal",
					responseValue,
				).Get(ctx, nil)
				if err != nil {
					logger.Error("EntityWorkflow: Failed to signal caller", "Error", err)
				} else {
					logger.Info("EntityWorkflow: Response signal sent", "response", responseValue)
				}
			}
		})

		// Handle shutdown signal.
		selector.AddReceive(shutdownCh, func(c workflow.ReceiveChannel, more bool) {
			var dummy string
			c.Receive(ctx, &dummy)
			logger.Info("EntityWorkflow: Shutdown signal received")
			shutdownRequested = true
		})

		// Timer for periodic logging.
		selector.AddFuture(workflow.NewTimer(ctx, 5*time.Second), func(f workflow.Future) {
			logger.Info("EntityWorkflow: Current state", "state", state)
		})

		selector.Select(ctx)
	}

	logger.Info("EntityWorkflow: Workflow shutting down")
	return nil
}
