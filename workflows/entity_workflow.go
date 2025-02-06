// file: entity_workflow.go
package workflows

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

type SignalData struct {
	CallerID string
	RunID    string
}

type Work struct {
	SignalData SignalData
	WorkData   string
}

// EntityWorkflow represents a long-running workflow that maintains state.
func EntityWorkflow(ctx workflow.Context) error {
	logger := workflow.GetLogger(ctx)
	workQueue := []Work{}
	shutdownRequested := false
	state := "initial"

	// Signal channels.
	signalCh := workflow.GetSignalChannel(ctx, "signal")
	shutdownCh := workflow.GetSignalChannel(ctx, "shutdown")

	// Run indefinitely until a shutdown signal is received.
	for !shutdownRequested {
		selector := workflow.NewSelector(ctx)

		// Handle updateState signal.
		selector.AddReceive(signalCh, func(c workflow.ReceiveChannel, more bool) {
			var signalData struct {
				Data     string
				CallerID string
				RunID    string
			}
			c.Receive(ctx, &signalData)

			workQueue = append(workQueue, Work{
				SignalData: SignalData{
					signalData.CallerID,
					signalData.RunID,
				},
				WorkData: signalData.Data,
			})

			logger.Info("EntityWorkflow: Work Queued")

		})
		// Optionally, if the update signal payload includes the caller's workflow ID,
		// capture it here. For a simple example, we assume it's already set.
		// In a real scenario, you might have a composite struct for the signal data.

		// Send back a response signal to the caller.
		// Ensure you have a valid callerWorkflowID (and optionally runID).

		// Handle shutdown signal.
		selector.AddReceive(shutdownCh, func(c workflow.ReceiveChannel, more bool) {
			var dummy string
			c.Receive(ctx, &dummy)
			logger.Info("EntityWorkflow: Shutdown signal received")
			shutdownRequested = true
		})

		// Timer for periodic logging.
		selector.AddFuture(workflow.NewTimer(ctx, 60*time.Second), func(f workflow.Future) {
			logger.Info("EntityWorkflow: Current state", "state", state)

			for _, w := range workQueue {
				updateData := w.WorkData
				state = updateData
				responseValue := "State updated to " + state

				err := workflow.SignalExternalWorkflow(
					ctx,
					w.SignalData.CallerID, // Caller workflow ID
					w.SignalData.RunID,    // Caller runID if known, or leave empty
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

		selector.Select(ctx)
	}

	logger.Info("EntityWorkflow: Workflow shutting down")
	return nil
}
