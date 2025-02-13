package api

import (
	"context"
)

type Edge interface {
	Common
	Device
	Validate
	DataSync
	CarfileOperation
	WaitQuiet(ctx context.Context) error //perm:read
	// ExternalServiceAddress check service address with different scheduler server
	// if behind nat, service address maybe different
	ExternalServiceAddress(ctx context.Context, schedulerURL string) (string, error) //perm:write
	// UserNATTravel build connection for user
	UserNATTravel(ctx context.Context, userServiceAddress string) error //perm:write
}
