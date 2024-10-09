package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
	db "github.com/sherifzaher/clone-simplebank/db/sqlc"
	"github.com/sherifzaher/clone-simplebank/util"
)

const (
	TaskSendEmailVerify = "task:send_email_verify"
)

type PayloadSendEmailTask struct {
	Username string `json:"username"`
}

func (distributor *RedisTaskDistributor) DistributeSendVerifyEmailTask(
	ctx context.Context,
	payload *PayloadSendEmailTask,
	opts ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	task := asynq.NewTask(TaskSendEmailVerify, jsonPayload, opts...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("cannot enqueue task: %w", err)
	}
	log.Info().
		Str("type", info.Type).
		Bytes("payload", info.Payload).
		Str("queue", info.Queue).
		Int("max_retry", info.MaxRetry).
		Msg("enqueued task")
	return nil
}

func (processor *RedisTaskProcessor) ProcessSendVerifyEmailTask(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendEmailTask
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return err
	}

	user, err := processor.store.GetUser(ctx, payload.Username)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// TODO: send email using google smtp

	verifyEmail, err := processor.store.CreateVerifyEmail(ctx, db.CreateVerifyEmailParams{
		Email:      user.Email,
		Username:   user.Username,
		SecretCode: util.RandomString(32),
	})
	if err != nil {
		return fmt.Errorf("failed to create verify_email: %w", err)
	}

	subject := "Welcome to Simple Bank"
	verifyUrl := fmt.Sprintf("http://localhost:8080/v1/verify_email?email_id=%d&secret_code=%s", verifyEmail.ID, verifyEmail.SecretCode)
	to := []string{user.Email}
	content := fmt.Sprintf(`Hello %s, <br />
Thank you for registering with us!<br />
Please <a href="%s">click here </a> to verify your email address. <br /> 
`, user.FullName, verifyUrl)
	err = processor.mailer.SendEmail(subject, content, to, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to send verify email: %w", err)
	}

	log.Info().
		Str("type", task.Type()).
		Bytes("payload", task.Payload()).
		Str("email", user.Email).
		Msg("processed task")
	return nil
}
