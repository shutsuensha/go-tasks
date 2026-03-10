package queue

import "github.com/hibiken/asynq"

type Worker struct {
    server *asynq.Server
}

func NewWorker(redisAddr string) *Worker {

    server := asynq.NewServer(
        asynq.RedisClientOpt{
            Addr: redisAddr,
        },
        asynq.Config{
            Concurrency: 10,
            Queues: map[string]int{
				QueueCritical: 6,
				QueueDefault:  3,
				QueueLow:      1,
			},
        },
    )

    return &Worker{server: server}
}

func (w *Worker) Run(handlers *Handlers) error {

    mux := asynq.NewServeMux()

    mux.HandleFunc(TypeTaskCreated, handlers.HandleTaskCreated)
	mux.HandleFunc(TypeCronEmail, handlers.HandleCronEmail)

    return w.server.Run(mux)
}