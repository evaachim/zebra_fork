package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

const TaskQueue = "TASK_QUEUE"

func printResults(result string, workflowID, runID string) {
	fmt.Printf("\nWorkflowID: %s RunID: %s\n", workflowID, runID)
	fmt.Printf("\n%s\n\n", result)
}

func ZebraflowExecutable(ctx workflow.Context, email string) (string, error) {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 500,
	}

	ctx = workflow.WithActivityOptions(ctx, options)

	var result string

	// if findUser(store, email) != nil {
	// do login stuff.

	// activity: login;
	err := workflow.ExecuteActivity(ctx, ProcessLogin, email).Get(ctx, &result)
	return result, err

	//}

	// return result, nil
}

func ProcessLogin(ctx context.Context, email string) (string, error) {
	loginAdapter()

	return "Login Processed", nil
}

func firstClient(wg *sync.WaitGroup) {
	defer wg.Done()
	// Create the client object just once per process
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer c.Close()

	// This worker hosts both Workflow and Activity functions
	w := worker.New(c, TaskQueue, worker.Options{})
	//workflow.

	// activity.
	w.RegisterWorkflow(ZebraflowExecutable)

	// Start listening to the Task Queue
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	}

	wg.Done()
}

func secondClient(wg *sync.WaitGroup) {
	defer wg.Done()

	time.Sleep(time.Second * 2)

	// Create the client object just once per process
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer c.Close()

	options := client.StartWorkflowOptions{
		ID:        "workflow",
		TaskQueue: TaskQueue,
	}

	email := "someone@somewhere.something"

	// The  Workflow
	we, err := c.ExecuteWorkflow(context.Background(), options, ZebraflowExecutable, email)
	if err != nil {
		log.Fatalln("unable to complete Workflow", err)
	}

	// Get the results
	var greeting string
	err = we.Get(context.Background(), &greeting)
	if err != nil {
		log.Fatalln("unable to get Workflow result", err)
	}

	printResults(greeting, we.GetID(), we.GetRunID())

}
