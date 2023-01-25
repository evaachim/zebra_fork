package framework

// Note that the workflow in this file assumes that the Ability Zebra PR was merged.
// Functions in that PR are used here.

import (
	"context"
	"errors"
	"project-safari/zebra"
)

var ErrLeaseNotProcessed = errors.New(`lease request not processed`)
var ErrLoginNotProcessed = errors.New(`login request not processed`)
var ErrQuerryNotProcessed = errors.New(`query request not processed`)

func ProcessLease(ctx context.Context, request Lease.Lease) (string, error) {
	if request.CanLease() && request.IsFree() {
		// do lease stuff.

		return "Lease Processed", nil
	}

	return "Lease Not Processed", ErrLeaseNotProcessed
}

func ProcessLogin(ctx context.Context, store zebra.Store, email string) (string, error) {
	/*if findUser() && request.IsFree() {
		// do login stuff.

		return "Login Processed", nil
	}
	*/

	return "Login Not Processed", ErrLoginNotProcessed
}

func ProcessQuery(ctx context.Context, this string) (string, error) {
	if this != "" {
		// do query stuff.
		return "Querry Processed", nil
	}

	return "Querry Not Processed", ErrQuerryNotProcessed
}

func ProcessPost(ctx context.Context, this string) (string, error) {
	if this != "" {
		// do post stuff.
		// some post stuff already done (such as for the db data).
		return "Post Processed", nil
	}

	return "Post Not Processed", ErrQuerryNotProcessed
}

func ProcessTasks(ctx context.Context, store zebra.Store, email string) (string, error) {
	/*	if findUser() && request.IsFree() {
			// do lease stuff.

			return "Login Processed", nil
		}
	*/

	return "Login Not Processed", ErrLoginNotProcessed
}
