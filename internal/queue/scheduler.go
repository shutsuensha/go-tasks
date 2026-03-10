package queue

import (
	"log"

	"github.com/hibiken/asynq"
)

func StartScheduler(redisAddr string) {

	scheduler := asynq.NewScheduler(
		asynq.RedisClientOpt{
			Addr: redisAddr,
		},
		&asynq.SchedulerOpts{},
	)

	task := asynq.NewTask(TypeCronEmail, nil)

	entryID, err := scheduler.Register(
		"* * * * *",
		task,
		asynq.Queue(QueueLow),
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("cron registered:", entryID)

	if err := scheduler.Run(); err != nil {
		log.Fatal(err)
	}
}