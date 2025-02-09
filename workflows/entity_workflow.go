// file: entity_workflow.go
package workflows

import (
	"fmt"
	"go.temporal.io/sdk/workflow"
	"os"
	"time"
)

type signalData struct {
	State    string
	CallerID string
	RunID    string
}

// EntityWorkflow represents a long-running workflow that maintains state.
func EntityWorkflow(ctx workflow.Context) error {

	logger := workflow.GetLogger(ctx)
	file, err := os.OpenFile("/tmp/workflowlog.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.Error("EntityWorkflow: unable to lock log file")
	}
	defer file.Close()
	stateQueue := []signalData{}
	shutdownRequested := false

	// Signal channels.
	updateStateCh := workflow.GetSignalChannel(ctx, "updateState")
	shutdownCh := workflow.GetSignalChannel(ctx, "shutdown")

	// Run indefinitely until a shutdown signal is received.
	for !shutdownRequested {
		selector := workflow.NewSelector(ctx)
		// Handle updateState signal.
		selector.AddReceive(updateStateCh, func(c workflow.ReceiveChannel, more bool) {
			var sigData signalData
			c.Receive(ctx, &sigData)
			stateQueue = append(stateQueue, sigData)
			logger.Info("EntityWorkflow: Appending to state queue", "state", sigData)
		})

		// Handle shutdown signal.
		selector.AddReceive(shutdownCh, func(c workflow.ReceiveChannel, more bool) {
			var dummy string
			c.Receive(ctx, &dummy)
			logger.Info("EntityWorkflow: Shutdown signal received")
			shutdownRequested = true
		})

		// Timer for periodic logging.
		selector.AddFuture(workflow.NewTimer(ctx, 10*time.Second), func(f workflow.Future) {
			for i, v := range stateQueue {
				logger.Info("EntityWorkflow: Processing stateQueue", "queue pos", i, "state", v)
				responseValue := "State updated to " + v.State
				sigData := v
				err := workflow.SignalExternalWorkflow(
					ctx,
					sigData.CallerID, // Caller workflow ID
					sigData.RunID,    // Caller runID if known, or leave empty
					"responseSignal",
					responseValue,
				).Get(ctx, nil)
				if err != nil {
					logger.Error("EntityWorkflow: Failed to signal caller", "Error", err)
				} else {
					logger.Info("EntityWorkflow: Response signal sent", "response", responseValue)
				}
				if _, err := file.WriteString(fmt.Sprintf("%s\n", v)); err != nil {
					logger.Error(fmt.Sprintf("%v", err))
				}
				workflow.Sleep(ctx, 3*time.Second)
			}
			stateQueue = []signalData{}
		})
		selector.Select(ctx)
	}
	logger.Info("EntityWorkflow: Workflow shutting down")
	return nil
}
