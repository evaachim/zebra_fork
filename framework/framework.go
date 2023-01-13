package framework

import (
	"context"
	"time"

	"go.temporal.io/sdk/workflow"
	"go.uber.org/zap"
)

const ApplicationName = "Zebra Resource Organization Tool"

func zebraWorkflow(ctx workflow.Context, name string) error {
	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute,
		StartToCloseTimeout:    time.Minute,
		HeartbeatTimeout:       time.Second * 20,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	logger := workflow.GetLogger(ctx)
	logger.Info("zebra workflow started")
	var helloworldResult string
	err := workflow.ExecuteActivity(ctx, helloWorldActivity, name).Get(ctx, &helloworldResult)
	if err != nil {
		logger.Error("Activity failed.", zap.Error(err))
		return err
	}

	return nil
}

func helloWorldActivity(ctx context.Context, name string) (string, error) {
	return "Hello Zebra !", nil
}
