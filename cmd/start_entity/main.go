// file: main.go
package main

import (
	"context"
	"fmt"
	"time"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"github.com/lukegriffith/entity-workflow-example/workflows" // Update with your actual module path.
)

func main() {
	// Create a Temporal client.
	c, err := client.NewClient(client.Options{})
	if err != nil {
		fmt.Println("Unable to create Temporal client", err)
		return
	}
	defer c.Close()

	// Create and start a worker for the task queue.
	w := worker.New(c, "SimpleEntityTaskQueue", worker.Options{})
	w.RegisterWorkflow(workflows.EntityWorkflow)
	w.RegisterWorkflow(workflows.CallerWorkflow)

	err = w.Start()
	if err != nil {
		fmt.Println("Unable to start worker", err)
		return
	}
	defer w.Stop()

	// Start the workflow.
	workflowOptions := client.StartWorkflowOptions{
		ID:        "entity_workflow_id",
		TaskQueue: "SimpleEntityTaskQueue",
		// Optionally, set timeouts like WorkflowRunTimeout.
	}
	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, workflows.EntityWorkflow)
	if err != nil {
		fmt.Println("Unable to execute workflow", err)
		return
	}

	fmt.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())

	// For demonstration, wait 15 seconds before sending a signal.
	time.Sleep(15 * time.Second)
	// Optionally, send a shutdown signal to stop the workflow gracefully.
	// Uncomment the following block to send the shutdown signal.
	/*
		time.Sleep(10 * time.Second) // Wait a bit before shutting down.
		err = c.SignalWorkflow(context.Background(), we.GetID(), we.GetRunID(), "shutdown", "shutdown")
		if err != nil {
			fmt.Println("Error sending shutdown signal", err)
		} else {
			fmt.Println("Sent shutdown signal")
		}
	*/

	// Optionally, wait for the workflow to complete if a shutdown signal is sent.
	// Note: If the workflow is not shut down, it will run indefinitely.
	err = we.Get(context.Background(), nil)
	if err != nil {
		fmt.Println("Workflow failed", err)
	} else {
		fmt.Println("Workflow completed.")
	}
}
