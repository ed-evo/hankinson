package orm

import (
	"context"

	"github.com/ed-evo/hankinson/server/ent"
	"github.com/ed-evo/hankinson/server/ent/argomento"
)

type Querier[T any] interface {
	Only(ctx context.Context) (*T, error)
}

type FetcherMidifier[Q any] func(q Q) Q

type EntityFetcher[T any, Q Querier[T]] interface {
	Fetch(ctx context.Context, id int, mods ...FetcherMidifier[Q]) (*T, error)
}

func applyMods[Q any](q Q, mods ...FetcherMidifier[Q]) Q {
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

func (f CapitoloFetcher) Fetch(ctx context.Context, id int, mods ...FetcherMidifier[*ent.CapitoloQuery]) (*ent.Capitolo, error) {
	return f.DB.Capitolo.Get(ctx, id)
}

type ArgomentoFetcher struct {
	DB *ent.Client
}

func (f ArgomentoFetcher) Fetch(ctx context.Context, id int, mods ...FetcherMidifier[*ent.ArgomentoQuery]) (*ent.Argomento, error) {
	q := applyMods(f.DB.Debug().Argomento.Query().Where(argomento.ID(id)), mods...)

	return q.Only(ctx)
}

type DomandaFetcher struct {
	DB *ent.Client
}

func (f DomandaFetcher) Fetch(ctx context.Context, id int, mods ...FetcherMidifier[*ent.DomandaQuery]) (*ent.Domanda, error) {
	return f.DB.Domanda.Get(ctx, id)
}

type QuesitoFetcher struct {
	DB *ent.Client
}

func (f QuesitoFetcher) Fetch(ctx context.Context, id int, mods ...FetcherMidifier[*ent.QuesitoEsameQuery]) (*ent.QuesitoEsame, error) {
	return f.DB.QuesitoEsame.Get(ctx, id)
}

type EsameFetcher struct {
	DB *ent.Client
}

func (f EsameFetcher) Fetch(ctx context.Context, id int, mods ...FetcherMidifier[*ent.EsameQuery]) (*ent.Esame, error) {
	return f.DB.Esame.Get(ctx, id)
}
