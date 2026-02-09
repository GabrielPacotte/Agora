package posthttp

import (
	"net/http"
	"strconv"

	"github.com/GabrielPacotte/Agora/internal/domain"
	"github.com/GabrielPacotte/Agora/internal/filters"
	httpcommon "github.com/GabrielPacotte/Agora/internal/http/common"
)

func (h *Handler) GetFeed(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := httpcommon.UserIDFromContext(ctx)
	if !ok {
		httpcommon.WriteError(w, http.StatusUnauthorized, httpcommon.ErrCodeAuthRequired, "unauthorized")
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	page := filters.NewPage(limit, offset)

	posts, err := h.postService.GetFeed(ctx, userID, page)
	if err != nil {
		print(err.Error())
		httpcommon.WriteError(w, http.StatusInternalServerError, httpcommon.ErrCodeInternal, "internal server error")
		return
	}

	out := make([]PostJSON, 0, len(posts))
	for _, p := range posts {
		out = append(out, mapPost(p))
	}

	resp := ListPostsResponse{
		Posts: out,
		Meta:  buildPageMeta(page.Limit(), page.Offset(), len(out)),
	}

	httpcommon.WriteJSON(w, http.StatusOK, resp)
}

func mapPost(p domain.Post) PostJSON {
	tags := make([]string, 0, len(p.Tags()))
	for _, t := range p.Tags() {
		tags = append(tags, string(t))
	}

	stances := make([]StanceJSON, 0, len(p.Stances()))
	for _, s := range p.Stances() {
		stances = append(stances, StanceJSON{
			ID:          s.ID().String(),
			Label:       s.Label(),
			Color:       s.Color().String(),
			Description: s.Description(),
		})
	}

	return PostJSON{
		ID:        p.ID().String(),
		AuthorID:  p.AuthorID().String(),
		Title:     p.Title(),
		Content:   p.Content(),
		Tags:      tags,
		Stances:   stances,
		CreatedAt: p.CreatedAt(),
		UpdatedAt: p.UpdatedAt(),
	}
}

func buildPageMeta(limit, offset, count int) PageMeta {
	meta := PageMeta{
		Limit: limit,
		Count: count,
	}
	if count == limit {
		next := offset + limit
		meta.NextOffset = &next
	}
	return meta
}
