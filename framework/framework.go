package framework

import (
	"go.temporal.io/sdk/workflow"
)

func ZebraflowExample(ctx workflow.Context, state CartState) error {
	logger := workflow.GetLogger(ctx)

	err := workflow.SetQueryHandler(ctx, "getState", func(input []byte) (LeaseState, error) {
		return state, nil
	})
	if err != nil {
		logger.Info("SetQueryHandler failed.", "Error", err)
		return err
	}

	// channel := workflow.GetSignalChannel(ctx, "stateMessages")
	selector := workflow.NewSelector(ctx)

	selector.AddReceive(channel, func(c workflow.ReceiveChannel, _ bool) {
		var signal interface{}
		c.Receive(ctx, &signal)
		state.Items = append(state.Items, CartItem{ProductId: 0, Quantity: 1})
	})

	for {
		selector.Select(ctx)
	}

	return nil
}
