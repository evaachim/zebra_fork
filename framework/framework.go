package framework

import (
	"project-safari/zebra"
	"project-safari/zebra/model/lease"
	"time"

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

	// activity: queries;
	err := workflow.ExecuteActivity(ctx, ProcessQuery, request).Get(ctx, &result)

	// activity: posts;
	err := workflow.ExecuteActivity(ctx, ProcessPost, request).Get(ctx, &result)

	return result, err
}
