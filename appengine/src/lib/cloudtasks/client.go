package cloudtasks

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"google.golang.org/api/option"
	taskspb "google.golang.org/genproto/googleapis/cloud/tasks/v2"

	"github.com/rabee-inc/push/appengine/src/lib/deploy"
	"github.com/rabee-inc/push/appengine/src/lib/httpclient"
	"github.com/rabee-inc/push/appengine/src/lib/log"
)

// Client ... CloudTasksのクライアント
type Client struct {
	cli        *cloudtasks.Client
	port       int
	deploy     string
	projectID  string
	locationID string
	authToken  string
}

// AddTask ... リクエストをEnqueueする
func (c *Client) AddTask(ctx context.Context, queue string, serviceID string, path string, params interface{}) error {
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": c.authToken,
	}
	body, err := json.Marshal(params)
	if err != nil {
		log.Errorm(ctx, "json.Marshal", err)
		return err
	}
	req := &taskspb.AppEngineHttpRequest{
		AppEngineRouting: &taskspb.AppEngineRouting{
			Service: serviceID,
		},
		HttpMethod:  taskspb.HttpMethod_POST,
		RelativeUri: path,
		Headers:     headers,
		Body:        body,
	}
	return c.addTask(ctx, queue, req)
}

func (c *Client) addTask(ctx context.Context, queue string, aeReq *taskspb.AppEngineHttpRequest) error {
	if deploy.IsLocal() {
		url := fmt.Sprintf("http://localhost:%d%s", c.port, aeReq.RelativeUri)
		status, _, err := httpclient.PostJSON(ctx, url, aeReq.Body, nil)
		if err != nil {
			log.Errorm(ctx, "httpclient.PostJSON", err)
			return err
		}
		if status != http.StatusOK {
			err = log.Errore(ctx, "task http status: %d", status)
			return err
		}
	} else {
		req := &taskspb.CreateTaskRequest{
			Parent: fmt.Sprintf("projects/%s/locations/%s/queues/%s", c.projectID, c.locationID, queue),
			Task: &taskspb.Task{
				MessageType: &taskspb.Task_AppEngineHttpRequest{
					AppEngineHttpRequest: aeReq,
				},
			},
		}
		_, err := c.cli.CreateTask(ctx, req)
		if err != nil {
			log.Errorm(ctx, "c.cli.CreateTask", err)
			return err
		}
	}
	return nil
}

// NewClient ... クライアントを作成する
func NewClient(
	credentialsPath string,
	port int,
	deploy string,
	projectID string,
	locationID string,
	authToken string) *Client {
	ctx := context.Background()
	opt := option.WithCredentialsFile(credentialsPath)
	cli, err := cloudtasks.NewClient(ctx, opt)
	if err != nil {
		panic(err)
	}
	return &Client{
		cli:        cli,
		port:       port,
		deploy:     deploy,
		projectID:  projectID,
		locationID: locationID,
		authToken:  authToken,
	}
}
