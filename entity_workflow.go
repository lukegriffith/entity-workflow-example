// file: entity_workflow.go
package workflows

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

// EntityWorkflow is a long-running workflow that maintains internal state.
// It listens for signals to update its state or shut itself down.
func EntityWorkflow(ctx workflow.Context) error {
	logger := workflow.GetLogger(ctx)

	// Initialize internal state.
	state := "initial"
	shutdownRequested := false

	// Obtain signal channels.
	updateStateCh := workflow.GetSignalChannel(ctx, "updateState")
	shutdownCh := workflow.GetSignalChannel(ctx, "shutdown")

	// Run indefinitely until a shutdown signal is received.
	for !shutdownRequested {
		// Create a selector to wait concurrently on timer and signals.
		selector := workflow.NewSelector(ctx)

		// Signal handler for updating state.
		selector.AddReceive(updateStateCh, func(c workflow.ReceiveChannel, more bool) {
			var newState string
			c.Receive(ctx, &newState)
			state = newState
			logger.Info("State updated", "state", state)
		})

		// Signal handler for shutdown.
		selector.AddReceive(shutdownCh, func(c workflow.ReceiveChannel, more bool) {
			// Consume the shutdown signal (data can be ignored).
			var dummy string
			c.Receive(ctx, &dummy)
			logger.Info("Shutdown signal received")
			shutdownRequested = true
		})

		// Timer to simulate work and periodically log the state.
		selector.AddFuture(workflow.NewTimer(ctx, 5*time.Second), func(f workflow.Future) {
			logger.Info("Current state", "state", state)
		})

		// Block until one of the events occurs.
		selector.Select(ctx)
	}

	logger.Info("Workflow shutting down.")
	return nil
}
