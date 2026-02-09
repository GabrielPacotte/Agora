package posthttp

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/GabrielPacotte/Agora/internal/domain"
	httpcommon "github.com/GabrielPacotte/Agora/internal/http/common"
	"github.com/google/uuid"
)

func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	userID, ok := httpcommon.UserIDFromContext(r.Context())
	if !ok {
		httpcommon.WriteError(w, http.StatusUnauthorized, httpcommon.ErrCodeAuthRequired, "unauthorized")
		return
	}

	var req CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpcommon.WriteError(w, http.StatusBadRequest, httpcommon.ErrCodeInvalidInput, "invalid JSON body")
		return
	}

	tags, err := mapTags(req.Tags)
	if err != nil {
		httpcommon.WriteError(w, http.StatusBadRequest, httpcommon.ErrCodeInvalidInput, "invalid tags: "+err.Error())
		return
	}

	stances, err := mapStances(req.Stances)
	if err != nil {
		httpcommon.WriteError(w, http.StatusBadRequest, httpcommon.ErrCodeInvalidInput, "invalid stances: "+err.Error())
		return
	}

	post, err := h.postService.Create(r.Context(), userID, req.Title, req.Content, tags, stances)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrEmptyPostTitle),
			errors.Is(err, domain.ErrPostTitleTooLong),
			errors.Is(err, domain.ErrPostContentTooLong),
			errors.Is(err, domain.ErrMinimalStanceAmount):
			httpcommon.WriteError(w, http.StatusBadRequest, httpcommon.ErrCodeInvalidInput, err.Error())
			return
		default:
			httpcommon.WriteError(w, http.StatusInternalServerError, httpcommon.ErrCodeInternal, "internal error"+err.Error())
			return
		}
	}

	resp := mapPostToResponse(post)

	httpcommon.WriteJSON(w, http.StatusCreated, resp)
}

func mapTags(input []string) ([]domain.Tag, error) {
	out := make([]domain.Tag, 0, len(input))
	for _, raw := range input {
		tag, err := domain.NewTag(raw)
		if err != nil {
			return nil, err
		}
		out = append(out, tag)
	}
	return out, nil
}

func mapStances(input []CreateStanceRequest) ([]domain.Stance, error) {
	out := make([]domain.Stance, 0, len(input))
	for _, in := range input {
		color, err := domain.NewColor(in.Color)
		if err != nil {
			return nil, err
		}
		stance, err := domain.NewStance(uuid.New(), in.Label, color, in.Description)
		if err != nil {
			return nil, err
		}
		out = append(out, stance)
	}
	return out, nil
}

func mapPostToResponse(p *domain.Post) GetPostResponse {
	tags := p.Tags()
	tagStrings := make([]string, len(tags))
	for i, t := range tags {
		tagStrings[i] = t.String()
	}

	stances := p.Stances()
	stanceDTOs := make([]StanceJSON, len(stances))
	for i, s := range stances {
		stanceDTOs[i] = StanceJSON{
			ID:          s.ID().String(),
			Label:       s.Label(),
			Color:       s.Color().String(),
			Description: s.Description(),
		}
	}

	postJson := PostJSON{
		ID:        p.ID().String(),
		Title:     p.Title(),
		Content:   p.Content(),
		Tags:      tagStrings,
		Stances:   stanceDTOs,
		AuthorID:  p.AuthorID().String(),
		CreatedAt: p.CreatedAt(),
		UpdatedAt: p.UpdatedAt(),
	}

	return GetPostResponse{
		Post: postJson,
	}
}
