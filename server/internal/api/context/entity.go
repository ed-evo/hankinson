package api_context

import (
	"net/http"
	"strconv"

	"github.com/ed-evo/hankinson/server/ent"
	"github.com/ed-evo/hankinson/server/internal/orm"
	"github.com/ed-evo/hankinson/server/pkg/api/api_errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type EntityContextHelper[T any, Q orm.Querier[T]] struct {
	ParamName string
	Fetcher   orm.EntityFetcher[T, Q]
}

func (h *EntityContextHelper[T, Q]) fetch(r *http.Request, mods ...orm.FetcherModifier[Q]) (*T, error) {
	paramVal := chi.URLParam(r, h.ParamName)
	id, err := strconv.Atoi(paramVal)
	if err != nil {
		return nil, err
	}
	return h.Fetcher.Fetch(r.Context(), id, mods...)
}

func (h *EntityContextHelper[T, Q]) JsonHandler(
	logic func(r *http.Request, entity *T) (any, error),
	mods ...orm.FetcherModifier[Q],
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		entity, err := h.fetch(r, mods...)

		if err != nil {
			if ent.IsNotFound(err) {
				render.Render(w, r, api_errors.ErrNotFound)
			} else {
				render.Render(w, r, api_errors.ErrInvalidRequest(err))
			}
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

func (h *EntityContextHelper[T, Q]) JsonListHandler(
	logic func(r *http.Request, entities []*T) (any, error),
	mods ...orm.FetcherModifier[Q],
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l, err := h.Fetcher.List(r.Context(), mods...)
		if err != nil {
			render.Render(w, r, api_errors.ErrInternal(err))
			return
		}

		response, err := logic(r, l)
		if err != nil {
			render.Render(w, r, api_errors.ErrInvalidRequest(err))
			return
		}
		render.JSON(w, r, response)
	}
}

func (h *EntityContextHelper[T, Q]) Process(
	logic func(r *http.Request, entity *T) error,
	mods ...orm.FetcherModifier[Q],
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
		w.WriteHeader(http.StatusNoContent)
	}
}
