// file: main.go
package main

import (
	"context"
	"fmt"

	"go.temporal.io/sdk/client"

	workflows "github.com/lukegriffith/entity-workflow-example" // Update with your actual module path.
)

func main() {
	// Create a Temporal client.
	c, err := client.NewClient(client.Options{})
	if err != nil {
		fmt.Println("Unable to create Temporal client", err)
		return
	}
	defer c.Close()

	ctx := context.Background()

	// Start the caller workflow.
	callerWorkflowOptions := client.StartWorkflowOptions{
		ID:        "caller-workflow-id",
		TaskQueue: "SimpleEntityTaskQueue",
	}
	we, err := c.ExecuteWorkflow(
		ctx,
		callerWorkflowOptions,
		workflows.CallerWorkflow,
	)
	if err != nil {
		fmt.Println("Failed to start CallerWorkflow", err)
		return
	}

	fmt.Println("CallerWorkflow started. WorkflowID:", we.GetID(), "RunID:", we.GetRunID())

	// Optionally wait for completion.
	err = we.Get(ctx, nil)
	if err != nil {
		fmt.Println("CallerWorkflow completed with error", err)
	} else {
		fmt.Println("CallerWorkflow completed successfully")
	}
}
