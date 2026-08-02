package api_context

import (
	"net/http"
	"strconv"

	"github.com/ed-evo/hankinson/server/internal/orm"
	"github.com/ed-evo/hankinson/server/pkg/api/api_errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type EntityContextHelper[T any, Q orm.Querier[T]] struct {
	ParamName  string
	ContextKey any
	Fetcher    orm.EntityFetcher[T, Q]
}

func (h *EntityContextHelper[T, Q]) fetch(r *http.Request, mods ...orm.FetcherMidifier[Q]) (*T, error) {
	paramVal := chi.URLParam(r, h.ParamName)
	id, err := strconv.Atoi(paramVal)
	if err != nil {
		return nil, err
	}
	return h.Fetcher.Fetch(r.Context(), id, mods...)
}

func (h *EntityContextHelper[T, Q]) JsonHandler(
	logic func(r *http.Request, entity *T) (any, error),
	mods ...orm.FetcherMidifier[Q],
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		entity, err := h.fetch(r, mods...)

		if err != nil {
			render.Render(w, r, api_errors.ErrInvalidRequest(err))
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

func (h *EntityContextHelper[T, Q]) Process(
	logic func(r *http.Request, entity *T) error,
	mods ...orm.FetcherMidifier[Q],
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		entity, err := h.fetch(r, mods...)

		if err != nil {
			render.Render(w, r, api_errors.ErrInvalidRequest(err))
			return
		}

		if err := logic(r, entity); err != nil {
			render.Render(w, r, api_errors.ErrInvalidRequest(err))
			return
		}
	}
}
