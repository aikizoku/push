package api

import (
	"context"
	"encoding/json"

	"github.com/rabee-inc/push/src/lib/log"
	"github.com/rabee-inc/push/src/service"
	"gopkg.in/go-playground/validator.v9"
)

// EntryHandler ... エントリーのハンドラ
type EntryHandler struct {
	Svc service.Register
}

type entryParams struct {
	UserID   string `json:"user_id"   validate:"required"`
	Platform string `json:"platform"  validate:"required"`
	DeviceID string `json:"device_id"`
	Token    string `json:"token"     validate:"required"`
}

type entryResponse struct {
	Success bool `json:"success"`
}

// DecodeParams ... 受け取ったJSONパラメータをデコードする
func (h *EntryHandler) DecodeParams(ctx context.Context, msg *json.RawMessage) (interface{}, error) {
	var params entryParams
	err := json.Unmarshal(*msg, &params)
	if err != nil {
		return params, err
	}

	// Validation
	v := validator.New()
	if err := v.Struct(params); err != nil {
		return params, err
	}
	return params, nil
}

// Exec ... 処理をする
func (h *EntryHandler) Exec(ctx context.Context, method string, params interface{}) (interface{}, error) {
	// パラメータを取得
	ps := params.(entryParams)

	err := h.Svc.SetToken(ctx, ps.UserID, ps.Platform, ps.DeviceID, ps.Token)
	if err != nil {
		log.Errorm(ctx, "h.Svc.SetToken", err)
		return nil, err
	}

	return entryResponse{
		Success: true,
	}, nil
}

// NewEntryHandler ... ハンドラを作成する
func NewEntryHandler(svc service.Register) *EntryHandler {
	return &EntryHandler{
		Svc: svc,
	}
}
