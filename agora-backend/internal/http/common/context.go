package httpcommon

import (
	"context"

	"github.com/google/uuid"
)

type ctxKey string

const ctxKeyUserID ctxKey = "user_id"

func WithUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, ctxKeyUserID, userID)
}

func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	v := ctx.Value(ctxKeyUserID)
	id, ok := v.(uuid.UUID)
	return id, ok
}
