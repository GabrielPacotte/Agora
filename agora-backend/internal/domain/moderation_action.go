package domain

import (
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

type ModerationDecision string

const (
	ModerationPending ModerationDecision = "pending"
	ReportReviewed    ModerationDecision = "report_reviewed"
	ContentDeleted    ModerationDecision = "content_deleted"
	UserSuspended     ModerationDecision = "user_suspended"
	UserBanished      ModerationDecision = "user_banished"
)

var validModerationDecisions = map[ModerationDecision]struct{}{
	ModerationPending: {},
	ReportReviewed:    {},
	ContentDeleted:    {},
	UserSuspended:     {},
	UserBanished:      {},
}

const maxModerationActionDescriptionLength = 1000

type ModerationAction struct {
	identified
	authorID    uuid.UUID
	decision    ModerationDecision
	reasonIDs   []uuid.UUID
	description string
}

func NewModerationAction(date time.Time, authorID uuid.UUID, decision ModerationDecision, reasonIDs []uuid.UUID, desc string) (*ModerationAction, error) {
	identified, err := newIdentified(date)
	if err != nil {
		return nil, err
	}
	reasonIDs = dedupeReasons(reasonIDs)
	desc = strings.TrimSpace(desc)
	err = validateModerationActionFields(authorID, decision, reasonIDs, desc)
	if err != nil {
		return nil, err
	}
	return &ModerationAction{
		identified:  identified,
		authorID:    authorID,
		reasonIDs:   reasonIDs,
		description: desc,
		decision:    decision,
	}, nil
}

func ExistingModerationAction(id uuid.UUID, createdAt, updatedAt time.Time, authorID uuid.UUID, decision ModerationDecision, reasonIDs []uuid.UUID, desc string) (*ModerationAction, error) {
	reasonIDs = dedupeReasons(reasonIDs)
	desc = strings.TrimSpace(desc)
	err := validateModerationActionFields(authorID, decision, reasonIDs, desc)
	if err != nil {
		return nil, err
	}

	identified, err := existingIdentified(id, createdAt, updatedAt)
	if err != nil {
		return nil, err
	}
	return &ModerationAction{
		identified:  identified,
		authorID:    authorID,
		reasonIDs:   reasonIDs,
		description: desc,
		decision:    decision,
	}, nil
}

func validateModerationActionFields(authorID uuid.UUID, decision ModerationDecision, reasonIDs []uuid.UUID, desc string) error {
	if authorID == uuid.Nil {
		return ErrInvalidID
	}

	if _, ok := validModerationDecisions[decision]; !ok {
		return ErrInvalidModerationActionDecision
	}

	if len(reasonIDs) == 0 {
		return ErrEmptyReasons
	}
	for _, id := range reasonIDs {
		if id == uuid.Nil {
			return ErrInvalidID
		}
	}

	if desc == "" || utf8.RuneCountInString(desc) > maxModerationActionDescriptionLength {
		return ErrModerationActionDescriptionSize
	}
	return nil
}

func dedupeReasons(reasonIDs []uuid.UUID) []uuid.UUID {
	seen := make(map[uuid.UUID]struct{}, len(reasonIDs))
	out := make([]uuid.UUID, 0, len(reasonIDs))

	for _, r := range reasonIDs {
		if _, exists := seen[r]; exists {
			continue
		}
		seen[r] = struct{}{}
		out = append(out, r)
	}

	return out
}

func (a *ModerationAction) AuthorID() uuid.UUID {
	return a.authorID
}

func (a *ModerationAction) Decision() ModerationDecision {
	return a.decision
}

func (a *ModerationAction) ReasonIDs() []uuid.UUID {
	out := make([]uuid.UUID, len(a.reasonIDs))
	copy(out, a.reasonIDs)
	return out
}

func (a *ModerationAction) Description() string {
	return a.description
}
