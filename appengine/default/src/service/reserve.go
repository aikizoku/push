package service

import (
	"context"

	"github.com/rabee-inc/push/appengine/default/src/config"
	"github.com/rabee-inc/push/appengine/default/src/model"
)

// Reserve ... 予約送信のサービス
type Reserve interface {
	Get(
		ctx context.Context,
		appID string,
		reserveID string) (*model.Reserve, error)
	List(
		ctx context.Context,
		appID string,
		limit int,
		cursor string) ([]*model.Reserve, string, error)
	Create(
		ctx context.Context,
		appID string,
		msg *model.Message,
		reservedAt int64) (*model.Reserve, error)
	Update(
		ctx context.Context,
		appID string,
		reserveID string,
		msg *model.Message,
		reservedAt int64,
		status config.ReserveStatus) (*model.Reserve, error)
}
