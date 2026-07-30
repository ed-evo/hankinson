package orm

import (
	"context"

	"github.com/ed-evo/hankinson/server/ent"
)

type EntityFetcher[T any] interface {
	Fetch(ctx context.Context, id int) (*T, error)
}

type CapitoloFetcher struct {
	DB *ent.Client
}

func (f CapitoloFetcher) Fetch(ctx context.Context, id int) (*ent.Capitolo, error) {
	return f.DB.Capitolo.Get(ctx, id)
}

type ArgomentoFetcher struct {
	DB *ent.Client
}

func (f ArgomentoFetcher) Fetch(ctx context.Context, id int) (*ent.Argomento, error) {
	return f.DB.Argomento.Get(ctx, id)
}

type DomandaFetcher struct {
	DB *ent.Client
}

func (f DomandaFetcher) Fetch(ctx context.Context, id int) (*ent.Domanda, error) {
	return f.DB.Domanda.Get(ctx, id)
}

type QuesitoFetcher struct {
	DB *ent.Client
}

func (f QuesitoFetcher) Fetch(ctx context.Context, id int) (*ent.QuesitoEsame, error) {
	return f.DB.QuesitoEsame.Get(ctx, id)
}
