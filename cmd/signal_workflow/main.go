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

	ctx := context.Background()

	// Start the caller workflow.

	for i := 0; i < 20; i++ {
		callerWorkflowOptions := client.StartWorkflowOptions{
			ID:        fmt.Sprintf("caller-workflow-id-%d", i),
			TaskQueue: "SimpleEntityTaskQueue",
		}
		we, err := c.ExecuteWorkflow(
			ctx,
			callerWorkflowOptions,
			workflows.CallerWorkflow,
			"entity_workflow_id",
			fmt.Sprintf("%d", i),
		)
		if err != nil {
			fmt.Println("Failed to start CallerWorkflow", err)
			return
		}

		fmt.Println("CallerWorkflow started. WorkflowID:", we.GetID(), "RunID:", we.GetRunID())

		// Optionally wait for completion.
		//err = we.Get(ctx, nil)
		//if err != nil {
		//	fmt.Println("CallerWorkflow completed with error", err)
		//} else {
		//	fmt.Println("CallerWorkflow completed successfully")
		//}
	}

	fmt.Scanln()
	for char := 'a'; char < 'a'+20; char++ {
		callerWorkflowOptions := client.StartWorkflowOptions{
			ID:        fmt.Sprintf("caller-workflow-id-%c", char),
			TaskQueue: "SimpleEntityTaskQueue",
		}

		we, err := c.ExecuteWorkflow(
			ctx,
			callerWorkflowOptions,
			workflows.CallerWorkflow,
			"entity_workflow_id",
			// Convert the rune (char) to a string for the workflow argument
			fmt.Sprintf("%c", char),
		)
		if err != nil {
			fmt.Println("Failed to start CallerWorkflow:", err)
			return
		}

		fmt.Println("CallerWorkflow started. WorkflowID:", we.GetID(), "RunID:", we.GetRunID())

		// Optionally wait for completion.
		// err = we.Get(ctx, nil)
		// if err != nil {
		//     fmt.Println("CallerWorkflow completed with error:", err)
		// } else {
		//     fmt.Println("CallerWorkflow completed successfully")
		// }
	}

}
