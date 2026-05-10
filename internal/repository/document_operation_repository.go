package repository

import (
	"context"

	"github.com/alvinmatias69/editor-history/internal/entity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DocumentOperationRepository struct {
	pool *pgxpool.Pool
}

func NewDocumentOperationRepository(pool *pgxpool.Pool) *DocumentOperationRepository {
	return &DocumentOperationRepository{
		pool: pool,
	}
}

func (r *DocumentOperationRepository) Get(ctx context.Context, documentID string) ([]entity.DocumentOperation, error) {
	rows, err := r.pool.Query(ctx, "SELECT id, document_id, editor_id, position_start, position_end, val, op, operation_time FROM document_operation WHERE document_id = $1 ORDER BY operation_time ASC", documentID)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[entity.DocumentOperation])
}

func (r *DocumentOperationRepository) Add(ctx context.Context, data []entity.DocumentOperation) (uint64, error) {
	// TODO: change to varchar 10000
	batch := &pgx.Batch{}
	for _, item := range data {
		batch.Queue(
			"INSERT INTO document_operation(id, document_id, editor_id, position_start, position_end, val, op, operation_time) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) ON CONFLICT DO NOTHING",
			item.ID, item.DocumentID, item.EditorID, item.PositionStart, item.PositionEnd, item.Value, item.Operation, item.OperationTimestamp)
	}
	err := r.pool.SendBatch(ctx, batch).Close()
	return 0, err
}
