package framework

import (
	"go.temporal.io/sdk/workflow"
)

func ZebraflowExample(ctx workflow.Context, request lease.Lease) error {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 5,
	}

	ctx = workflow.WithActivityOptions(ctx, options)

	var result string

	// activity: lease requested and leasable;
	err := workflow.ExecuteActivity(ctx, ProcessLease, request).Get(ctx, &result)

	return result, err
}
