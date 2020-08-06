package repos

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/keegancsmith/sqlf"

	"github.com/sourcegraph/sourcegraph/internal/db/basestore"
	"github.com/sourcegraph/sourcegraph/internal/workerutil"
)

func newSyncWorker(ctx context.Context, handle *basestore.TransactableHandle, numWorkers int) *workerutil.Worker {
	store := workerutil.NewStore(handle, workerutil.StoreOptions{
		TableName:         "external_service_sync_jobs",
		ViewName:          "external_service_sync_jobs_with_next_sync_at",
		Scan:              scanSyncJob,
		OrderByExpression: sqlf.Sprintf("next_sync_at"),
		ColumnExpressions: syncJobColumns,
		StalledMaxAge:     10 * time.Second,
		// Zero for now as we expect errors to be transient
		MaxNumResets: 0,
	})

	return workerutil.NewWorker(ctx, store, workerutil.WorkerOptions{
		Name:        "sync_worker",
		Handler:     &syncHandler{},
		NumHandlers: numWorkers,
		Interval:    1 * time.Minute,
	})
}

var syncJobColumns = []*sqlf.Query{
	sqlf.Sprintf("id"),
	sqlf.Sprintf("id"),
	sqlf.Sprintf("state"),
	sqlf.Sprintf("failure_message"),
	sqlf.Sprintf("started_at"),
	sqlf.Sprintf("finished_at"),
	sqlf.Sprintf("process_after"),
	sqlf.Sprintf("num_resets"),
	sqlf.Sprintf("external_service_id"),
	sqlf.Sprintf("next_sync_at"),
}

func scanSyncJob(rows *sql.Rows, err error) (workerutil.Record, bool, error) {
	return nil, false, errors.New("TODO")
}

type syncJob struct {
	ID                int
	ExternalServiceID int
}

func (s *syncJob) RecordID() int {
	return s.ID
}

type syncHandler struct{}

func (h *syncHandler) Handle(ctx context.Context, tx workerutil.Store, record workerutil.Record) error {
	//myStore := h.myStore.With(tx) // combine store handles
	//myRecord := record.(MyType)   // convert type of record
	// do processing ...
	return errors.New("TODO")
}
