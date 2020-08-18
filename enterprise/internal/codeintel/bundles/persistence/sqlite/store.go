package sqlite

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/sourcegraph/sourcegraph/enterprise/internal/codeintel/bundles/persistence"
	"github.com/sourcegraph/sourcegraph/enterprise/internal/codeintel/bundles/persistence/cache"
	"github.com/sourcegraph/sourcegraph/enterprise/internal/codeintel/bundles/persistence/serialization"
	gobserializer "github.com/sourcegraph/sourcegraph/enterprise/internal/codeintel/bundles/persistence/serialization/gob"
	"github.com/sourcegraph/sourcegraph/enterprise/internal/codeintel/bundles/persistence/sqlite/migrate"
	"github.com/sourcegraph/sourcegraph/enterprise/internal/codeintel/bundles/persistence/sqlite/store"
)

type sqliteStore struct {
	filename   string
	cache      cache.DataCache
	store      *store.Store
	closer     func() error
	serializer serialization.Serializer
}

var _ persistence.Store = &sqliteStore{}

func NewStore(ctx context.Context, filename string, cache cache.DataCache) (_ persistence.Store, err error) {
	store, closer, err := store.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			if closeErr := closer(); closeErr != nil {
				err = multierror.Append(err, closeErr)
			}
		}
	}()

	serializer := gobserializer.New()

	if err := migrate.Migrate(ctx, store, serializer); err != nil {
		return nil, err
	}

	return &sqliteStore{
		filename:   filename,
		cache:      cache,
		store:      store,
		closer:     closer,
		serializer: serializer,
	}, nil
}

// func NewWriter(ctx context.Context, filename string) (_ persistence.Writer, err error) {
// 	store, closer, err := store.Open(filename)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer func() {
// 		if err != nil {
// 			if closeErr := closer(); closeErr != nil {
// 				err = multierror.Append(err, closeErr)
// 			}
// 		}
// 	}()

// 	tx, err := store.Transact(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if err := createTables(ctx, tx); err != nil {
// 		return nil, err
// 	}

// 	return &sqliteStore{
// 		store:      tx,
// 		closer:     closer,
// 		serializer: gobserializer.New(),
// 	}, nil
// }
