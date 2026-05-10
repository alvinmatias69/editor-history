package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/alvinmatias69/editor-history/internal/config"
	"github.com/alvinmatias69/editor-history/internal/entity"
	"github.com/redis/go-redis/v9"
)

const keyFormat = "document:%s"

type SessionRepository struct {
	client         *redis.Client
	expiryDuration time.Duration
}

func NewSessionRepository(client *redis.Client, config *config.Redis) *SessionRepository {
	return &SessionRepository{
		client:         client,
		expiryDuration: config.KeyExpiry,
	}
}

func (r *SessionRepository) Get(ctx context.Context, documentID string) (string, error) {
	key := fmt.Sprintf(keyFormat, documentID)
	editorID, err := r.client.Get(ctx, key).Result()
	if errors.Is(redis.Nil, err) {
		return "", nil
	}

	return editorID, err
}

func (r *SessionRepository) Set(ctx context.Context, request entity.SessionRequest) error {
	key := fmt.Sprintf(keyFormat, request.DocumentID)
	return r.client.Set(ctx, key, request.EditorID, r.expiryDuration).Err()
}
