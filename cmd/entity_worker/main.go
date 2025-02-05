// file: main.go
package main

import (
	"fmt"

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

	err = w.Run(worker.InterruptCh())
	if err != nil {
		fmt.Println("Unable to start worker", err)
	}
}
