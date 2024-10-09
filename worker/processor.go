package worker

import (
	"context"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
	db "github.com/sherifzaher/clone-simplebank/db/sqlc"
	"github.com/sherifzaher/clone-simplebank/mail"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
)

type TaskProcessor interface {
	Start() error
	ProcessSendVerifyEmailTask(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	redis  *asynq.Server
	store  db.Store
	mailer mail.EmailSender
}

func NewRedisTaskProcessor(redisOpts asynq.RedisClientOpt, store db.Store, mailer mail.EmailSender) TaskProcessor {
	server := asynq.NewServer(redisOpts, asynq.Config{
		Queues: map[string]int{
			QueueCritical: 10,
			QueueDefault:  5,
		},
		ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
			log.Error().
				Err(err).
				Str("type", task.Type()).
				Bytes("payload", task.Payload()).
				Msg("process task failed")
		}),
		Logger: NewRedisLogger(),
	})
	return &RedisTaskProcessor{
		redis:  server,
		store:  store,
		mailer: mailer,
	}
}

func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(TaskSendEmailVerify, processor.ProcessSendVerifyEmailTask)

	return processor.redis.Start(mux)
}
