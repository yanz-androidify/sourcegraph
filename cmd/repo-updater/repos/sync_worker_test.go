package repos_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/sourcegraph/sourcegraph/cmd/repo-updater/repos"
	"github.com/sourcegraph/sourcegraph/internal/extsvc"
	"github.com/sourcegraph/sourcegraph/internal/workerutil"
)

func testSyncWorkerPlumbing(db *sql.DB) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		// Add fake
		repoStore := repos.NewDBStore(db, sql.TxOptions{
			Isolation: sql.LevelSerializable,
		})
		testSvc := &repos.ExternalService{
			Kind:        extsvc.KindGitHub,
			DisplayName: "TestService",
			Config:      "{}",
		}

		// Create external service
		err := repoStore.UpsertExternalServices(ctx, testSvc)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("Test service created, ID: %d", testSvc.ID)

		// Add item to queue
		result, err := db.ExecContext(ctx, `insert into external_service_sync_jobs (external_service_id) values ($1);`, testSvc.ID)
		if err != nil {
			t.Fatal(err)
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			t.Fatal(err)
		}
		if rowsAffected != 1 {
			t.Fatalf("Expected 1 row to be affected, got %d", rowsAffected)
		}

		jobChan := make(chan *repos.SyncJob)

		h := &fakeRepoSyncHandler{
			jobChan: jobChan,
		}
		worker := repos.NewSyncWorker(ctx, db, h, 1)
		go worker.Start()
		defer worker.Stop()

		var job *repos.SyncJob
		select {
		case job = <-jobChan:
			t.Log("Job received")
		case <-time.After(5 * time.Second):
			t.Fatal("Timeout")
		}

		if job.ExternalServiceID != testSvc.ID {
			t.Fatalf("Expected %d, got %d", testSvc.ID, job.ExternalServiceID)
		}
	}
}

type fakeRepoSyncHandler struct {
	jobChan chan *repos.SyncJob
}

func (h *fakeRepoSyncHandler) Handle(ctx context.Context, tx workerutil.Store, record workerutil.Record) error {
	r, ok := record.(*repos.SyncJob)
	if !ok {
		return fmt.Errorf("expected repos.SyncJob, got %T", record)
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case h.jobChan <- r:
		return nil
	}
}
