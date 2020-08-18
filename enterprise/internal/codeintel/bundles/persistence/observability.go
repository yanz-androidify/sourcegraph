package persistence

import (
	"context"

	"github.com/sourcegraph/sourcegraph/enterprise/internal/codeintel/bundles/types"
	"github.com/sourcegraph/sourcegraph/internal/metrics"
	"github.com/sourcegraph/sourcegraph/internal/observation"
)

// An ObservedReader wraps another Reader with error logging, Prometheus metrics, and tracing.
type ObservedReader struct {
	reader                     Store
	readMetaOperation          *observation.Operation
	pathsWithPrefixOperation   *observation.Operation
	readDocumentOperation      *observation.Operation
	readResultChunkOperation   *observation.Operation
	readDefinitionsOperation   *observation.Operation
	readReferencesOperation    *observation.Operation
	writeMetaOperation         *observation.Operation
	writeDocumentsOperation    *observation.Operation
	writeResultChunksOperation *observation.Operation
	writeDefinitionsOperation  *observation.Operation
	writeReferencesOperation   *observation.Operation
}

var _ Store = &ObservedReader{}

// singletonMetrics ensures that the operation metrics required by ObservedReader are
// constructed only once as there may be many readers instantiated by a single replica
// of precise-code-intel-bundle-manager.
var singletonMetrics = &metrics.SingletonOperationMetrics{}

// NewObservedReader wraps the given Reader with error logging, Prometheus metrics, and tracing.
func NewObserved(reader Store, observationContext *observation.Context) Store {
	metrics := singletonMetrics.Get(func() *metrics.OperationMetrics {
		return metrics.NewOperationMetrics(
			observationContext.Registerer,
			"bundle_reader",
			metrics.WithLabels("op"),
			metrics.WithCountHelp("Total number of results returned"),
		)
	})

	return &ObservedReader{
		reader: reader,
		readMetaOperation: observationContext.Operation(observation.Op{
			Name:         "Reader.ReadMeta",
			MetricLabels: []string{"read_meta"},
			Metrics:      metrics,
		}),
		pathsWithPrefixOperation: observationContext.Operation(observation.Op{
			Name:         "Reader.PathsWithPrefix",
			MetricLabels: []string{"paths_with_prefix"},
			Metrics:      metrics,
		}),
		readDocumentOperation: observationContext.Operation(observation.Op{
			Name:         "Reader.ReadDocument",
			MetricLabels: []string{"read_document"},
			Metrics:      metrics,
		}),
		readResultChunkOperation: observationContext.Operation(observation.Op{
			Name:         "Reader.ReadResultChunk",
			MetricLabels: []string{"read_result-chunk"},
			Metrics:      metrics,
		}),
		readDefinitionsOperation: observationContext.Operation(observation.Op{
			Name:         "Reader.ReadDefinitions",
			MetricLabels: []string{"read_definitions"},
			Metrics:      metrics,
		}),
		readReferencesOperation: observationContext.Operation(observation.Op{
			Name:         "Reader.ReadReferences",
			MetricLabels: []string{"read_references"},
			Metrics:      metrics,
		}),
		writeMetaOperation: observationContext.Operation(observation.Op{
			Name:         "Reader.WriteMeta",
			MetricLabels: []string{"write_meta"},
			Metrics:      metrics,
		}),
		writeDocumentsOperation: observationContext.Operation(observation.Op{
			Name:         "Reader.WriteDocuments",
			MetricLabels: []string{"write_documents"},
			Metrics:      metrics,
		}),
		writeResultChunksOperation: observationContext.Operation(observation.Op{
			Name:         "Reader.WriteResultChunks",
			MetricLabels: []string{"write_result_chunks"},
			Metrics:      metrics,
		}),
		writeDefinitionsOperation: observationContext.Operation(observation.Op{
			Name:         "Reader.WriteDefinitions",
			MetricLabels: []string{"write_definitions"},
			Metrics:      metrics,
		}),
		writeReferencesOperation: observationContext.Operation(observation.Op{
			Name:         "Reader.WriteReferences",
			MetricLabels: []string{"write_references"},
			Metrics:      metrics,
		}),
	}
}

// ReadMeta calls into the inner Reader and registers the observed results.
func (r *ObservedReader) ReadMeta(ctx context.Context) (_ types.MetaData, err error) {
	ctx, endObservation := r.readMetaOperation.With(ctx, &err, observation.Args{})
	defer endObservation(1, observation.Args{})
	return r.reader.ReadMeta(ctx)
}

// PathsWithPrefix calls into the inner Reader and registers the observed results.
func (r *ObservedReader) PathsWithPrefix(ctx context.Context, prefix string) (_ []string, err error) {
	ctx, endObservation := r.pathsWithPrefixOperation.With(ctx, &err, observation.Args{})
	defer endObservation(1, observation.Args{})
	return r.reader.PathsWithPrefix(ctx, prefix)
}

// ReadDocument calls into the inner Reader and registers the observed results.
func (r *ObservedReader) ReadDocument(ctx context.Context, path string) (_ types.DocumentData, _ bool, err error) {
	ctx, endObservation := r.readDocumentOperation.With(ctx, &err, observation.Args{})
	defer endObservation(1, observation.Args{})
	return r.reader.ReadDocument(ctx, path)
}

// ReadResultChunk calls into the inner Reader and registers the observed results.
func (r *ObservedReader) ReadResultChunk(ctx context.Context, id int) (_ types.ResultChunkData, _ bool, err error) {
	ctx, endObservation := r.readResultChunkOperation.With(ctx, &err, observation.Args{})
	defer endObservation(1, observation.Args{})
	return r.reader.ReadResultChunk(ctx, id)
}

// ReadDefinitions calls into the inner Reader and registers the observed results.
func (r *ObservedReader) ReadDefinitions(ctx context.Context, scheme, identifier string, skip, take int) (locations []types.Location, _ int, err error) {
	ctx, endObservation := r.readDefinitionsOperation.With(ctx, &err, observation.Args{})
	defer func() { endObservation(float64(len(locations)), observation.Args{}) }()
	return r.reader.ReadDefinitions(ctx, scheme, identifier, skip, take)
}

// ReadReferences calls into the inner Reader and registers the observed results.
func (r *ObservedReader) ReadReferences(ctx context.Context, scheme, identifier string, skip, take int) (locations []types.Location, _ int, err error) {
	ctx, endObservation := r.readReferencesOperation.With(ctx, &err, observation.Args{})
	defer func() { endObservation(float64(len(locations)), observation.Args{}) }()
	return r.reader.ReadReferences(ctx, scheme, identifier, skip, take)
}

// WriteMeta calls into the inner Reader and registers the observed results.
func (r *ObservedReader) WriteMeta(ctx context.Context, meta types.MetaData) (err error) {
	ctx, endObservation := r.writeMetaOperation.With(ctx, &err, observation.Args{})
	defer endObservation(1, observation.Args{})
	return r.reader.WriteMeta(ctx, meta)
}

// WriteDocuments calls into the inner Reader and registers the observed results.
func (r *ObservedReader) WriteDocuments(ctx context.Context, documents chan KeyedDocumentData) (err error) {
	ctx, endObservation := r.writeDocumentsOperation.With(ctx, &err, observation.Args{})
	defer endObservation(1, observation.Args{})
	return r.reader.WriteDocuments(ctx, documents)
}

// WriteResultChunks calls into the inner Reader and registers the observed results.
func (r *ObservedReader) WriteResultChunks(ctx context.Context, resultChunks chan IndexedResultChunkData) (err error) {
	ctx, endObservation := r.writeResultChunksOperation.With(ctx, &err, observation.Args{})
	defer endObservation(1, observation.Args{})
	return r.reader.WriteResultChunks(ctx, resultChunks)
}

// WriteDefinitions calls into the inner Reader and registers the observed results.
func (r *ObservedReader) WriteDefinitions(ctx context.Context, monikerLocations chan types.MonikerLocations) (err error) {
	ctx, endObservation := r.writeDefinitionsOperation.With(ctx, &err, observation.Args{})
	defer endObservation(1, observation.Args{})
	return r.reader.WriteDefinitions(ctx, monikerLocations)
}

// WriteReferences calls into the inner Reader and registers the observed results.
func (r *ObservedReader) WriteReferences(ctx context.Context, monikerLocations chan types.MonikerLocations) (err error) {
	ctx, endObservation := r.writeReferencesOperation.With(ctx, &err, observation.Args{})
	defer endObservation(1, observation.Args{})
	return r.reader.WriteReferences(ctx, monikerLocations)
}

func (r *ObservedReader) Close(err error) error {
	return r.reader.Close(err)
}
