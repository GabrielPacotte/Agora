package filters

import (
	"time"

	"github.com/GabrielPacotte/Agora/internal/domain"
)

type SearchPostFilter struct {
	Page
	Tags     []domain.Tag
	FromDate *time.Time
	ToDate   *time.Time
}
