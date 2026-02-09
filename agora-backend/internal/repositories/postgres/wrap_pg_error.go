package postgres

import (
	"context"
	"errors"

	"github.com/GabrielPacotte/Agora/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func wrapPGError(err error) error {
	if err == nil {
		return nil
	}

	// --------------------------------------------------------------------
	// 1. Canceled / timeout
	// --------------------------------------------------------------------
	if errors.Is(err, context.Canceled) {
		return domain.ErrCanceled
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return domain.ErrTimeout
	}

	// --------------------------------------------------------------------
	// 2. No result: Not Found
	// --------------------------------------------------------------------
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ErrNotFound
	}

	// --------------------------------------------------------------------
	// 3. PostgreSQL Errors (pgconn.PgError)
	// --------------------------------------------------------------------
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {

		switch pgErr.Code {

		// ------------------------------------------------------------
		// unique_violation
		// ------------------------------------------------------------
		case "23505":
			if pgErr.ConstraintName == "users_login_key" {
				return domain.ErrAlreadyExists
			}
			return domain.ErrAlreadyExists

		// ------------------------------------------------------------
		// foreign_key_violation
		// ------------------------------------------------------------
		case "23503":
			return domain.ErrForeignKeyViolation

		// ------------------------------------------------------------
		// not_null_violation
		// ------------------------------------------------------------
		case "23502":
			return domain.ErrValidation

		// ------------------------------------------------------------
		// check_violation (CHECK constraint)
		// ------------------------------------------------------------
		case "23514":
			return domain.ErrValidation

		// ------------------------------------------------------------
		// invalid text representation (ex: mauvais UUID)
		// ------------------------------------------------------------
		case "22P02":
			return domain.ErrValidation

		// ------------------------------------------------------------
		// invalid enum value
		// ------------------------------------------------------------
		case "22P03":
			return domain.ErrValidation

		// ------------------------------------------------------------
		// insufficient_privilege (ex. SELECT interdit)
		// ------------------------------------------------------------
		case "42501":
			return domain.ErrForbidden

		case "P0002":
			return domain.ErrNotFound
		}

		return err
	}

	// --------------------------------------------------------------------
	// 4. Go / Networking errors
	// --------------------------------------------------------------------
	return err
}
