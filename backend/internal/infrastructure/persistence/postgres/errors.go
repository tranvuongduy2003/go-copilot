package postgres

import "errors"

var (
	ErrPoolNotInitialized    = errors.New("database pool not initialized")
	ErrConnectionCancelled   = errors.New("connection cancelled")
	ErrConnectionFailed      = errors.New("failed to connect to database")
	ErrParseConfig           = errors.New("failed to parse connection string")
	ErrBeginTransaction      = errors.New("failed to begin transaction")
	ErrCommitTransaction     = errors.New("failed to commit transaction")
	ErrRollbackTransaction   = errors.New("failed to rollback transaction")
)

type DBError struct {
	Op  string
	Err error
}

func (e *DBError) Error() string {
	if e.Err != nil {
		return e.Op + ": " + e.Err.Error()
	}
	return e.Op
}

func (e *DBError) Unwrap() error {
	return e.Err
}

func newDBError(op string, err error) *DBError {
	return &DBError{Op: op, Err: err}
}
