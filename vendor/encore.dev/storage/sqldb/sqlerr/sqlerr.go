// Package sqlerr provides a set of common error codes for SQL datbaases.
//
// Not all error codes are supported by all databases, and multiple
// underlying database errors may map to the same error code if the
// meaning overlaps.
//
// In future releases additional error codes may be added to this
// package, and the mapping of underlying database errors to these
// error codes may change.
package sqlerr

// Code describes a specific type of database error.
// The value Other is reported when an error does not map to any of the defined codes.
//
// Not all error codes are supported by all databases, and multiple
// underlying database errors may map to the same error code if the
// meaning overlaps.
type Code string

const (
	// Other is reported when an error does not map to any of the defined codes.
	// Consult the underlying database error for more information.
	Other Code = "other"

	// NotNullViolation is reported when a not null constraint would be violated.
	NotNullViolation Code = "not_null_violation"

	// ForeignKeyViolation is reported when a foreign key constraint would be violated.
	ForeignKeyViolation Code = "foreign_key_violation"

	// UniqueViolation is reported when a unique constraint would be violated.
	UniqueViolation Code = "unique_violation"

	// CheckViolation is reported when a check constraint would be violated.
	CheckViolation Code = "check_violation"

	// ExcludeViolation is reported when an exclusion constraint would be violated.
	ExcludeViolation Code = "exclude_violation"

	// TransactionFailed is reported when running a command in a failed transaction,
	// due to some previous command failure.
	TransactionFailed Code = "transaction_failed"

	// DeadlockDetected is reported when a deadlock is detected.
	// Deadlock detection is done on a best-effort basis and not all deadlocks
	// can be detected.
	DeadlockDetected Code = "deadlock_detected"

	// TooManyConnections is reported when the database rejects a connection request
	// due to reaching the maximum number of connections.
	// This is different from blocking waiting on a connection pool.
	TooManyConnections Code = "too_many_connections"
)

// Severity defines the severity of a database error.
type Severity string

const (
	SeverityError   Severity = "ERROR"
	SeverityFatal   Severity = "FATAL"
	SeverityPanic   Severity = "PANIC"
	SeverityWarning Severity = "WARNING"
	SeverityNotice  Severity = "NOTICE"
	SeverityDebug   Severity = "DEBUG"
	SeverityInfo    Severity = "INFO"
	SeverityLog     Severity = "LOG"
)
