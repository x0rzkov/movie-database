package user

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/dqkcode/movie-database/internal/app/types"
)

type (
	service interface {
		Register(ctx context.Context, req RegisterRequest) (string, error)
		Update(ctx context.Context, req UpdateInfoRequest) error
		FindUserById(ctx context.Context, id string) (*User, error)
		DeleteUser(ctx context.Context, id string) error
		// ChangePassword(ctx context.Context, req ChangePasswordRequest) error
	}
	Handler struct {
		srv service
	}
)

func NewHandler(srv service) *Handler {
	return &Handler{
		srv,
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, err := h.srv.Register(r.Context(), req)

	if err == ErrDBQuery {
		types.ResponseJson(w, "", types.Normal().Internal)
		return
	} else if err == ErrUserAlreadyExist {
		types.ResponseJson(w, "", types.User().DuplicateEmail)
		return
	} else if err == ErrCreateUserFailed {
		types.ResponseJson(w, "", types.User().CreateFailed)
		return
	}
	data := map[string]string{
		"id": id,
	}
	types.ResponseJson(w, data, types.User().Created)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	var req UpdateInfoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.srv.Update(r.Context(), req)
	if err != nil {
		types.ResponseJson(w, "", types.User().UpdateFailed)
		return
	}

	data := map[string]string{
		"id": r.Context().Value("user").(*types.UserInfo).ID,
	}

	types.ResponseJson(w, data, types.Normal().Success)
}

func (h *Handler) FindUserById(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		types.ResponseJson(w, "", types.Normal().BadRequest)
		return
	}
	user, err := h.srv.FindUserById(r.Context(), id)
	if err != nil {
		types.ResponseJson(w, "", types.User().UserNotFound)
		return
	}
	data := map[string]interface{}{
		"data": user,
	}
	types.ResponseJson(w, data, types.Normal().Success)
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		types.ResponseJson(w, "", types.Normal().BadRequest)
		return
	}

	err := h.srv.DeleteUser(r.Context(), id)
	if err != nil {
		types.ResponseJson(w, "", types.User().DeleteFailed)
		return
	}

	types.ResponseJson(w, "", types.Normal().Success)
}
