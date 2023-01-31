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

func ZebraflowLogin(ctx workflow.Context, email string) (string, error) {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 5,
	}

	print("\nI had been here")

	ctx = workflow.WithActivityOptions(ctx, options)

	// do login stuff.
	print("\nI was here")

	var result string

	// activity: login;
	err := workflow.ExecuteActivity(ctx, ProcessLogin, email).Get(ctx, &result)
	print("\nI am here")

	return result, err
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
	w.RegisterWorkflow(ZebraflowLogin)
	w.RegisterActivity(ProcessLogin)

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
	var result string
	err = we.Get(context.Background(), &result)

	if err != nil {
		log.Fatalln("unable to get Workflow result", err)
	}

	printResults(result, we.GetID(), we.GetRunID())

}
