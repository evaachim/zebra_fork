package framework

// Note that the workflow in this file assumes that the Ability Zebra PR was merged.
// Functions in that PR are used here.

import (
	"context"
	"fmt"
)

var ErrLeaseNotProcessed = errors.New(`lease request not processed`)

func ProcessLease(ctx context.Context, request Lease.Lease) (string, error) {
	if request.CanLease() && request.IsFree() {
		// do lease stuff.

		return "Lease Processed", nil
	}

	return "Lease Not Processed", ErrLeaseNotProcessed
}
