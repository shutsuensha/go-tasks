package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shutsuensha/go-tasks/internal/db"
	"github.com/shutsuensha/go-tasks/internal/metrics"
	"github.com/shutsuensha/go-tasks/internal/queue"
	"github.com/redis/go-redis/v9"
	
	"fmt"

	"golang.org/x/sync/singleflight"
)

type TaskService struct {
	pool  *pgxpool.Pool
	q     *db.Queries
	rdb   *redis.Client
	sf    singleflight.Group
	queue *queue.Client
}

func NewTaskService(
	pool *pgxpool.Pool,
	q *db.Queries,
	rdb *redis.Client,
	queue *queue.Client,
) *TaskService {

	return &TaskService{
		pool:  pool,
		q:     q,
		rdb:   rdb,
		queue: queue,
	}
}

func (s *TaskService) CreateTask(
	ctx context.Context,
	title string,
	description string,
) (db.Task, error) {

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return db.Task{}, err
	}

	defer tx.Rollback(ctx)

	qtx := s.q.WithTx(tx)

	task, err := qtx.CreateTask(ctx, db.CreateTaskParams{
		Title:       title,
		Description: description,
	})
	if err != nil {
		return db.Task{}, err
	}
	
	_, err = qtx.CreateTaskEvent(ctx, db.CreateTaskEventParams{
		TaskID:    task.ID,
		EventType: "task_created",
	})
	if err != nil {
		return db.Task{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return db.Task{}, err
	}
	
	if s.queue != nil {
		if err := s.queue.EnqueueTaskCreated(task.ID); err != nil {
			return db.Task{}, err
		}
	}

	return task, nil
}

func (s *TaskService) ListTasks(
	ctx context.Context,
	limit int32,
	offset int32,
) ([]db.Task, error) {

	key := fmt.Sprintf("tasks:%d:%d", limit, offset)

	start := time.Now()
	cached, err := s.rdb.Get(ctx, key).Result()
	metrics.ObserveRedisOperation("get", start)

	if err == nil {

		metrics.RedisCacheHits.Inc()

		var tasks []db.Task
		if err := json.Unmarshal([]byte(cached), &tasks); err == nil {
			return tasks, nil
		}
	}

	metrics.RedisCacheMisses.Inc()

	v, err, _ := s.sf.Do(key, func() (interface{}, error) {

		startDB := time.Now()

		tasks, err := s.q.ListTasksPaginated(ctx, db.ListTasksPaginatedParams{
			Limit:  limit,
			Offset: offset,
		})

		metrics.ObserveQuery("list_tasks_paginated", startDB)

		if err != nil {
			return nil, err
		}

		data, err := json.Marshal(tasks)
		if err == nil {
			s.rdb.Set(ctx, key, data, time.Minute)
		}

		return tasks, nil
	})

	if err != nil {
		return nil, err
	}

	return v.([]db.Task), nil
}

func (s *TaskService) GetTask(
	ctx context.Context,
	id int32,
) (db.Task, error) {

	return s.q.GetTask(ctx, id)
}