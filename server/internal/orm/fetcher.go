package orm

import (
	"context"

	"github.com/ed-evo/hankinson/server/ent"
	"github.com/ed-evo/hankinson/server/ent/argomento"
	"github.com/ed-evo/hankinson/server/ent/capitolo"
	"github.com/ed-evo/hankinson/server/ent/domanda"
	"github.com/ed-evo/hankinson/server/ent/esame"
	"github.com/ed-evo/hankinson/server/ent/quesitoesame"
)

type Querier[T any] interface {
	Only(ctx context.Context) (*T, error)
	All(ctx context.Context) ([]*T, error)
}

type FetcherModifier[Q any] func(q Q) Q

type EntityFetcher[T any, Q Querier[T]] interface {
	Fetch(ctx context.Context, id int, mods ...FetcherModifier[Q]) (*T, error)
	List(ctx context.Context, mods ...FetcherModifier[Q]) ([]*T, error)
}

func applyMods[Q any](q Q, mods ...FetcherModifier[Q]) Q {
	_q := q
	for _, m := range mods {
		if m == nil {
			continue
		}
		_q = m(_q)
	}
	return _q
}

type CapitoloFetcher struct {
	DB *ent.Client
}

func (f *CapitoloFetcher) Fetch(ctx context.Context, id int, mods ...FetcherModifier[*ent.CapitoloQuery]) (*ent.Capitolo, error) {
	q := applyMods(f.DB.Capitolo.Query().Where(capitolo.ID(id)), mods...)
	return q.Only(ctx)
}

func (f *CapitoloFetcher) List(ctx context.Context, mods ...FetcherModifier[*ent.CapitoloQuery]) ([]*ent.Capitolo, error) {
	q := applyMods(f.DB.Capitolo.Query(), mods...)
	return q.All(ctx)
}

type ArgomentoFetcher struct {
	DB *ent.Client
}

func (f *ArgomentoFetcher) Fetch(ctx context.Context, id int, mods ...FetcherModifier[*ent.ArgomentoQuery]) (*ent.Argomento, error) {
	q := applyMods(f.DB.Debug().Argomento.Query().Where(argomento.ID(id)), mods...)
	return q.Only(ctx)
}

func (f *ArgomentoFetcher) List(ctx context.Context, mods ...FetcherModifier[*ent.ArgomentoQuery]) ([]*ent.Argomento, error) {
	q := applyMods(f.DB.Debug().Argomento.Query(), mods...)
	return q.All(ctx)
}

type DomandaFetcher struct {
	DB *ent.Client
}

func (f *DomandaFetcher) Fetch(ctx context.Context, id int, mods ...FetcherModifier[*ent.DomandaQuery]) (*ent.Domanda, error) {
	q := applyMods(f.DB.Domanda.Query().Where(domanda.ID(id)), mods...)
	return q.Only(ctx)
}
func (f *DomandaFetcher) List(ctx context.Context, mods ...FetcherModifier[*ent.DomandaQuery]) ([]*ent.Domanda, error) {
	q := applyMods(f.DB.Domanda.Query(), mods...)
	return q.All(ctx)
}

type QuesitoFetcher struct {
	DB *ent.Client
}

func (f *QuesitoFetcher) Fetch(ctx context.Context, id int, mods ...FetcherModifier[*ent.QuesitoEsameQuery]) (*ent.QuesitoEsame, error) {
	q := applyMods(f.DB.QuesitoEsame.Query().Where(quesitoesame.ID(id)), mods...)
	return q.Only(ctx)
}
func (f *QuesitoFetcher) List(ctx context.Context, mods ...FetcherModifier[*ent.QuesitoEsameQuery]) ([]*ent.QuesitoEsame, error) {
	q := applyMods(f.DB.QuesitoEsame.Query(), mods...)
	return q.All(ctx)
}

type EsameFetcher struct {
	DB *ent.Client
}

func (f *EsameFetcher) Fetch(ctx context.Context, id int, mods ...FetcherModifier[*ent.EsameQuery]) (*ent.Esame, error) {
	q := applyMods(f.DB.Esame.Query().Where(esame.ID(id)), mods...)
	return q.Only(ctx)
}
func (f *EsameFetcher) List(ctx context.Context, mods ...FetcherModifier[*ent.EsameQuery]) ([]*ent.Esame, error) {
	q := applyMods(f.DB.Esame.Query(), mods...)
	return q.All(ctx)
}

var _ EntityFetcher[ent.Esame, *ent.EsameQuery] = (*EsameFetcher)(nil)
