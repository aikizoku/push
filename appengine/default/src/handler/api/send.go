package api

import (
	"context"
	"encoding/json"

	"gopkg.in/go-playground/validator.v9"

	"github.com/rabee-inc/push/appengine/default/src/config"
	"github.com/rabee-inc/push/appengine/default/src/lib/cloudtasks"
	"github.com/rabee-inc/push/appengine/default/src/model"
)

// SendAction ...  送信のアクション
type SendAction struct {
	tCli *cloudtasks.Client
}

type sendParams struct {
	AppID   string         `json:"app_id"   validate:"required"`
	UserIDs []string       `json:"user_ids" validate:"required"`
	Message *model.Message `json:"message"  validate:"required"`
}

type sendResponse struct {
	Success bool `json:"success"`
}

// DecodeParams ... 受け取ったJSONパラメータをデコードする
func (h *SendAction) DecodeParams(ctx context.Context, msg *json.RawMessage) (interface{}, error) {
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
func (h *SendAction) Exec(ctx context.Context, method string, params interface{}) (interface{}, error) {
	ps := params.(sendParams)

	src := &model.CloudTasksParamSendUserIDs{
		AppID:   ps.AppID,
		UserIDs: ps.UserIDs,
		Message: ps.Message,
	}

	err := h.tCli.AddTask(ctx, config.QueueSendUser, "/worker/send/users", src)
	if err != nil {
		return nil, err
	}

	return sendResponse{
		Success: true,
	}, nil
}

// NewSendAction ... SendActionを作成する
func NewSendAction(tCli *cloudtasks.Client) *SendAction {
	return &SendAction{
		tCli: tCli,
	}
}