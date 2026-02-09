package filters

type Page struct {
	limit  int
	offset int
}

const (
	defaultLimit = 20
	maxLimit     = 100
	minLimit     = 1
)

func NewPage(limit, offset int) Page {
	return Page{
		limit:  limit,
		offset: offset,
	}
}

func (p Page) Limit() int {
	if p.limit < minLimit {
		return defaultLimit
	}
	if p.limit > maxLimit {
		return maxLimit
	}
	return p.limit
}

func (p Page) Offset() int {
	if p.offset < 0 {
		return 0
	}
	return p.offset
}
