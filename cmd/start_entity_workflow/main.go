// file: main.go
package main

import (
	"context"
	"fmt"

	"go.temporal.io/sdk/client"

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
}
