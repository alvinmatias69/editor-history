package repository

import (
	"context"
	"errors"

	"github.com/alvinmatias69/editor-history/internal/constants"
	"github.com/alvinmatias69/editor-history/internal/entity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DocumentRepository struct {
	pool *pgxpool.Pool
}

func NewDocumentRepository(pool *pgxpool.Pool) *DocumentRepository {
	return &DocumentRepository{
		pool: pool,
	}
}

func (r *DocumentRepository) Create(ctx context.Context, request entity.Document) error {
	_, err := r.pool.Exec(ctx,
		"INSERT INTO document(id, read_only_id, is_finished) VALUES ($1, $2, $3)",
		request.ID, request.ReadOnlyID, request.IsFinished)

	return err
}

func (r *DocumentRepository) Finish(ctx context.Context, documentID string) error {
	commandTag, err := r.pool.Exec(ctx, "UPDATE document SET is_finished = true WHERE id = $1", documentID)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() != 1 {
		return constants.DocumentNotFound
	}

	return nil
}

func (r *DocumentRepository) GetForRead(ctx context.Context, readOnlyID string) (entity.Document, error) {
	var (
		isFinished bool
		id         string
	)

	err := r.pool.QueryRow(ctx, "SELECT id, is_finished FROM document WHERE read_only_id = $1", readOnlyID).Scan(&id, &isFinished)
	if errors.Is(pgx.ErrNoRows, err) {
		return entity.Document{}, constants.DocumentNotFound
	}

	return entity.Document{
		ID:         id,
		ReadOnlyID: readOnlyID,
		IsFinished: isFinished,
	}, err
}

func (r *DocumentRepository) GetForWrite(ctx context.Context, id string) (entity.Document, error) {
	var (
		isFinished bool
		readOnlyID string
	)

	err := r.pool.QueryRow(ctx, "SELECT read_only_id, is_finished FROM document WHERE id = $1", id).Scan(&readOnlyID, &isFinished)
	if errors.Is(pgx.ErrNoRows, err) {
		return entity.Document{}, constants.DocumentNotFound
	}

	return entity.Document{
		ID:         id,
		ReadOnlyID: readOnlyID,
		IsFinished: isFinished,
	}, err
}
