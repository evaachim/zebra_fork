package framework

import (
	"go.temporal.io/sdk/workflow"
)

func ZebraflowExecutable(ctx workflow.Context, request lease.Lease, store zebra.Store, email string) error {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 500,
	}

	ctx = workflow.WithActivityOptions(ctx, options)

	var result string

	// activity: login;
	err := workflow.ExecuteActivity(ctx, ProcessLogin, store, email).Get(ctx, &result)

	// activity: lease requested and leasable;
	err := workflow.ExecuteActivity(ctx, ProcessLease, request).Get(ctx, &result)

	return result, err
}
