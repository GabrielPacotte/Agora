package posthttp

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/GabrielPacotte/Agora/internal/domain"
	httpcommon "github.com/GabrielPacotte/Agora/internal/http/common"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := chi.URLParam(r, "postId")
	postID, err := uuid.Parse(idStr)
	if err != nil {
		httpcommon.WriteError(w, http.StatusBadRequest, httpcommon.ErrCodeInvalidInput, "invalid id")
		return
	}

	requesterID, ok := httpcommon.UserIDFromContext(ctx)
	if !ok {
		httpcommon.WriteError(w, http.StatusUnauthorized, httpcommon.ErrCodeAuthRequired, "unauthorized")
		return
	}

	var req UpdatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpcommon.WriteError(w, http.StatusBadRequest, httpcommon.ErrCodeInvalidInput, "invalid JSON body")
		return
	}

	tags, err := mapStringsToTags(req.Tags)
	if err != nil {
		httpcommon.WriteError(w, http.StatusBadRequest, httpcommon.ErrCodeInvalidInput, "invalid tags")
		return
	}

	stances, err := mapUpdateStances(req.Stances)
	if err != nil {
		httpcommon.WriteError(w, http.StatusBadRequest, httpcommon.ErrCodeInvalidInput, "invalid stances")
		return
	}

	updated, err := h.postService.Update(ctx, requesterID, postID, req.Title, req.Content, tags, stances)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrForbidden):
			httpcommon.WriteError(w, http.StatusForbidden, httpcommon.ErrCodeForbidden, "forbidden")
		case errors.Is(err, domain.ErrNotFound):
			httpcommon.WriteError(w, http.StatusNotFound, httpcommon.ErrCodeNotFound, "post not found")
		case errors.Is(err, domain.ErrEmptyPostTitle),
			errors.Is(err, domain.ErrPostTitleTooLong),
			errors.Is(err, domain.ErrPostContentTooLong),
			errors.Is(err, domain.ErrMinimalStanceAmount),
			errors.Is(err, domain.ErrEmptyStanceLabel),
			errors.Is(err, domain.ErrInvalidColor):
			httpcommon.WriteError(w, http.StatusBadRequest, httpcommon.ErrCodeInvalidInput, err.Error())
		default:
			httpcommon.WriteError(w, http.StatusInternalServerError, httpcommon.ErrCodeInternal, "internal server error")
		}
		return
	}

	resp := mapPostToResponse(updated)
	httpcommon.WriteJSON(w, http.StatusOK, resp)
}

func mapStringsToTags(in []string) ([]domain.Tag, error) {
	if len(in) == 0 {
		return nil, nil
	}
	out := make([]domain.Tag, 0, len(in))
	for _, s := range in {
		t, err := domain.NewTag(s)
		if err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, nil
}

func mapUpdateStances(in []UpdateStanceRequest) ([]domain.Stance, error) {
	if len(in) == 0 {
		return nil, nil
	}
	out := make([]domain.Stance, 0, len(in))
	for _, s := range in {
		id, err := uuid.Parse(s.ID)
		if err != nil {
			return nil, err
		}
		c, err := domain.NewColor(s.Color)
		if err != nil {
			return nil, err
		}
		st, err := domain.NewStance(id, s.Label, c, s.Description)
		if err != nil {
			return nil, err
		}
		out = append(out, st)
	}
	return out, nil
}
