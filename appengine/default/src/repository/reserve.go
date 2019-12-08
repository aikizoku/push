package repository

import (
	"context"

	"cloud.google.com/go/firestore"

	"github.com/rabee-inc/push/appengine/default/src/config"
	"github.com/rabee-inc/push/appengine/default/src/model"
)

// Reserve ... 予約のリポジトリ
type Reserve interface {
	Get(
		ctx context.Context,
		appID string,
		reserveID string) (*model.Reserve, error)
	ListByCursor(
		ctx context.Context,
		appID string,
		limit int,
		cursor string) ([]*model.Reserve, string, error)
	ListBySend(
		ctx context.Context,
		appID string,
		now int64,
		limit int,
		cursor *firestore.DocumentSnapshot) ([]*model.Reserve, *firestore.DocumentSnapshot, error)
	Create(
		ctx context.Context,
		appID string,
		msg *model.Message,
		reservedAt int64,
		status config.ReserveStatus,
		createdAt int64) (*model.Reserve, error)
	Update(
		ctx context.Context,
		appID string,
		src *model.Reserve,
		updatedAt int64) (*model.Reserve, error)
	BtUpdate(
		ctx context.Context,
		bt *firestore.WriteBatch,
		appID string,
		src *model.Reserve,
		updatedAt int64) *model.Reserve
}
