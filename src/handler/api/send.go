package api

import (
	"context"
	"encoding/json"

	"github.com/aikizoku/push/src/config"
	"github.com/aikizoku/push/src/lib/log"
	"github.com/aikizoku/push/src/lib/taskqueue"
	"github.com/aikizoku/push/src/model"
	"github.com/aikizoku/push/src/service"
	"github.com/go-playground/validator"
)

// SendHandler ...  送信のハンドラ
type SendHandler struct {
	Svc service.Sender
}

type sendParams struct {
	UserIDs []string       `json:"user_ids" validate:"required"`
	Message *model.Message `json:"message" validate:"required"`
}

// DecodeParams ... 受け取ったJSONパラメータをデコードする
func (h *SendHandler) DecodeParams(ctx context.Context, msg *json.RawMessage) (interface{}, error) {
	var params sendParams
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
func (h *SendHandler) Exec(ctx context.Context, method string, params interface{}) (interface{}, error) {
	ps := params.(sendParams)

	src := &model.SendUserIDs{
		UserIDs: ps.UserIDs,
		Message: ps.Message,
	}

	err := taskqueue.AddTaskByJSON(ctx, config.QueueSendUser, "/worker/send/users", src)
	if err != nil {
		log.Errorf(ctx, "taskqueue.NewJSONPostTask error: %s", err.Error())
		return nil, err
	}

	return struct {
		Success bool `json:"success"`
	}{
		Success: true,
	}, nil
}

// NewSendHandler ... SendHandlerを作成する
func NewSendHandler(svc service.Sender) *SendHandler {
	return &SendHandler{
		Svc: svc,
	}
}
