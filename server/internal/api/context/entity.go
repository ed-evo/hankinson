package api_context

import (
	"context"
	"net/http"
	"strconv"

	"github.com/ed-evo/hankinson/server/internal/orm"
	"github.com/ed-evo/hankinson/server/pkg/api/api_errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type EntityContextHelper[T any] struct {
	ParamName  string
	ContextKey any
	Fetcher    orm.EntityFetcher[T]
}

func (h *EntityContextHelper[T]) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			paramVal := chi.URLParam(r, h.ParamName)
			id, err := strconv.Atoi(paramVal)
			if err != nil {
				render.Render(w, r, api_errors.ErrInvalidRequest(err))
				return
			}

			entity, err := h.Fetcher.Fetch(r.Context(), id)
			if err != nil {
				render.Render(w, r, api_errors.ErrNotFound)
				return
			}

			// Store entity in context
			ctx := context.WithValue(r.Context(), h.ContextKey, entity)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (h *EntityContextHelper[T]) JsonHandler(
	logic func(r *http.Request, entity *T) (any, error),
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		val := r.Context().Value(h.ContextKey)
		entity, ok := val.(*T)
		if !ok {
			render.Render(w, r, api_errors.ErrNotFound)
			return
		}

		response, err := logic(r, entity)
		if err != nil {
			render.Render(w, r, api_errors.ErrInvalidRequest(err))
			return
		}

		render.JSON(w, r, response)
	}
}

func (h *EntityContextHelper[T]) Process(
	logic func(r *http.Request, entity *T) error,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		val := r.Context().Value(h.ContextKey)
		entity, ok := val.(*T)
		if !ok {
			render.Render(w, r, api_errors.ErrNotFound)
			return
		}

		err := logic(r, entity)
		if err != nil {
			render.Render(w, r, api_errors.ErrInvalidRequest(err))
			return
		}
	}
}
