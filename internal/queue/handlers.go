package queue

import (
    "context"
    "encoding/json"
    "log"

    "github.com/hibiken/asynq"
	"fmt"
)

type Handlers struct{}

func NewHandlers() *Handlers {
    return &Handlers{}
}

func (h *Handlers) HandleTaskCreated(ctx context.Context, t *asynq.Task) error {

    var payload TaskCreatedPayload

    if err := json.Unmarshal(t.Payload(), &payload); err != nil {
        return err
    }

    log.Printf("processing task id=%d", payload.TaskID)

    return fmt.Errorf("test retry")

}


func (h *Handlers) HandleCronEmail(ctx context.Context, t *asynq.Task) error {

	log.Println("cron job: sending emails")

	return nil
}