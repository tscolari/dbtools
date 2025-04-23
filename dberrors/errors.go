package dberrors

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ToStatusErr converts an error returned from database operations
// into a status.Error.
func ToStatusErr(err error, msgs ...string) error {
	if err == nil {
		return nil
	}

	msgs = append(msgs, "%v")
	message := strings.Join(msgs, ":")

	switch actualErr := err.(type) {
	case *pgconn.PgError:
		return fromPgConn(actualErr, message)

	default:
		code := codes.Internal

		if errors.Is(err, sql.ErrNoRows) {
			code = codes.NotFound
		}

		return status.Errorf(code, message, err.Error())
	}
}

func fromPgConn(err *pgconn.PgError, message string) error {

	switch errType := string(err.Code[0] + err.Code[1]); errType {
	// 08: Connection Exception
	case "08":
		return status.Errorf(codes.Unavailable, message+": %v", err.Message)

	// 22: Data Exception
	case "22":
		return status.Errorf(codes.InvalidArgument, message+": %v", err.Message)

	// 23: Integrity Constraint Violation
	case "23":
		return status.Errorf(codes.FailedPrecondition, message+": %v", err.Message)

	// 28: Invalid Authorization Specification
	case "28":
		return status.Errorf(codes.PermissionDenied, message+": %v", err.Message)

	// 42: Syntax Error or Access Rule Violation
	// XX: Internal errors
	case "42", "XX":
		return status.Errorf(codes.Internal, message+": %v", err.Message)

	}

	return status.Errorf(codes.Unavailable, message+": %v", err.Message)
}
