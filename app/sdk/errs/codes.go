package errs

import "encore.dev/beta/errs"

// This set of variables exist so I can move APP layer code from the service
// project over without the need to refactor.
var (
	OK                 = errs.OK
	NoContent          = errs.OK
	Canceled           = errs.Canceled
	Unknown            = errs.Unknown
	InvalidArgument    = errs.InvalidArgument
	DeadlineExceeded   = errs.DeadlineExceeded
	NotFound           = errs.NotFound
	AlreadyExists      = errs.AlreadyExists
	PermissionDenied   = errs.PermissionDenied
	ResourceExhausted  = errs.ResourceExhausted
	FailedPrecondition = errs.FailedPrecondition
	Aborted            = errs.Aborted
	OutOfRange         = errs.OutOfRange
	Unimplemented      = errs.Unimplemented
	Internal           = errs.Internal
	Unavailable        = errs.Unavailable
	DataLoss           = errs.DataLoss
	Unauthenticated    = errs.Unauthenticated
)
