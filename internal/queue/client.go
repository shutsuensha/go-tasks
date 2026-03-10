package queue

import (
    "encoding/json"

    "github.com/hibiken/asynq"
	"time"
)

type Client struct {
    client *asynq.Client
}

func NewClient(redisAddr string) *Client {
    return &Client{
        client: asynq.NewClient(asynq.RedisClientOpt{
            Addr: redisAddr,
        }),
    }
}

type TaskCreatedPayload struct {
    TaskID int32
}

func (c *Client) EnqueueTaskCreated(taskID int32) error {

    payload, err := json.Marshal(TaskCreatedPayload{
        TaskID: taskID,
    })
    if err != nil {
        return err
    }

    task := asynq.NewTask(TypeTaskCreated, payload)

    _, err = c.client.Enqueue(task, asynq.Queue(QueueDefault), asynq.MaxRetry(5), asynq.Timeout(30*time.Second))

    return err
}