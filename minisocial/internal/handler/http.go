package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"minisocial/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type HTTPHandler struct {
	svc *service.SocialService
}

func NewHTTPHandler(svc *service.SocialService) *HTTPHandler {
	return &HTTPHandler{svc: svc}
}

func (h *HTTPHandler) RegisterRoutes(r chi.Router) {
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/users", h.RegisterUser)

		// Fвторизация эмулируется через заголовок X-User-ID
		r.With(h.authMiddleware).Route("/", func(r chi.Router) {
			r.Post("/posts", h.CreatePost)
			r.Post("/users/{id}/follow", h.FollowUser)
			r.Post("/posts/{id}/like", h.LikePost)
			r.Get("/feed", h.GetFeed)
		})
	})
}

func (h *HTTPHandler) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userIDStr := r.Header.Get("X-User-ID")
		if _, err := uuid.Parse(userIDStr); err != nil {
			http.Error(w, "Unauthorized: Invalid X-User-ID", http.StatusUnauthorized)
			return
		}
		r = r.WithContext(setUserID(r.Context(), userIDStr))
		next.ServeHTTP(w, r)
	})
}

// RegisterUser godoc
// @Summary Регистрация нового пользователя
// @Description Создает нового пользователя по имени
// @Tags users
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Информация о пользователе"
// @Success 201 {object} domain.User
// @Failure 400 {string} string "Bad request"
// @Router /users [post]
func (h *HTTPHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	user, err := h.svc.Register(r.Context(), req.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	respondJSON(w, http.StatusCreated, user)
}

// CreatePost godoc
// @Summary Создать новый пост
// @Description Добавляет пост в ленту от имени текущего пользователя
// @Tags posts
// @Accept json
// @Produce json
// @Param request body PostRequest true "Текст поста"
// @Param X-User-ID header string true "ID пользователя"
// @Success 201 {object} domain.Post
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Router /posts [post]
func (h *HTTPHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r.Context())
	var req struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	post, err := h.svc.CreatePost(r.Context(), uuid.MustParse(userID), req.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	respondJSON(w, http.StatusCreated, post)
}

// FollowUser godoc
// @Summary Подписаться на пользователя
// @Description Добавляет текущего пользователя в подписчики к указанному ID
// @Tags follows
// @Produce json
// @Param id path string true "ID пользователя, на которого подписываемся"
// @Param X-User-ID header string true "ID пользователя"
// @Success 204 "No Content"
// @Failure 400 {string} string "Bad request or self-follow"
// @Failure 401 {string} string "Unauthorized"
// @Router /users/{id}/follow [post]
func (h *HTTPHandler) FollowUser(w http.ResponseWriter, r *http.Request) {
	followerID := uuid.MustParse(getUserID(r.Context()))
	followingID := uuid.MustParse(chi.URLParam(r, "id"))

	if err := h.svc.Follow(r.Context(), followerID, followingID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// LikePost godoc
// @Summary Поставить лайк посту
// @Description Увеличивает счетчик лайков у указанного поста
// @Tags likes
// @Produce json
// @Param id path string true "ID поста"
// @Param X-User-ID header string true "ID пользователя"
// @Success 204 "No Content"
// @Failure 401 {string} string "Unauthorized"
// @Router /posts/{id}/like [post]
func (h *HTTPHandler) LikePost(w http.ResponseWriter, r *http.Request) {
	userID := uuid.MustParse(getUserID(r.Context()))
	postID := uuid.MustParse(chi.URLParam(r, "id"))

	if err := h.svc.Like(r.Context(), userID, postID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// GetFeed godoc
// @Summary Получить ленту новостей
// @Description Возвращает список постов от пользователей, на которых подписан текущий пользователь
// @Tags feed
// @Produce json
// @Param X-User-ID header string true "ID пользователя"
// @Success 200 {array} domain.Post
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /feed [get]
func (h *HTTPHandler) GetFeed(w http.ResponseWriter, r *http.Request) {
	userID := uuid.MustParse(getUserID(r.Context()))

	posts, err := h.svc.GetFeed(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respondJSON(w, http.StatusOK, posts)
}

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload != nil {
		json.NewEncoder(w).Encode(payload)
	}
}

type contextKey string

const userIDKey contextKey = "userID"

func setUserID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, userIDKey, id)
}

func getUserID(ctx context.Context) string {
	return ctx.Value(userIDKey).(string)
}

type RegisterRequest struct {
	Username string `json:"username" example:"john_doe"`
}

type PostRequest struct {
	Content string `json:"content" example:"Hello world!"`
}
