package task

import (
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"golang.org/x/net/context"
	"log"
)

const (
	TypeEmailDelivery = "email:delivery"
	TypeImageResize   = "image:resize"
)

type EmailDeliveryPayload struct {
	UserID     int    `json:"user_id,omitempty"`
	TemplateID string `json:"template_id,omitempty"`
}

type ImageResizePayload struct {
	SourceURL string `json:"source_url,omitempty"`
}

type ImageProcessor struct {
}

func NewEmailDeliveryTask(userID int, tmplID string) (*asynq.Task, error) {
	payload, err := json.Marshal(EmailDeliveryPayload{UserID: userID, TemplateID: tmplID})
	if err != nil {
		return nil, fmt.Errorf("marshal payload: %v", err)
	}

	return asynq.NewTask(TypeEmailDelivery, payload), nil
}

func HandleEmailDeliveryTask(ctx context.Context, task *asynq.Task) error {
	var payload EmailDeliveryPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("unmarshal %s: %v", string(task.Payload()), err)
	}

	log.Printf("Sending Email to User: user_id=%d, template_id=%s\n", payload.UserID, payload.TemplateID)
	return nil
}

func NewImageResizeTask(src string) (*asynq.Task, error) {
	payload, err := json.Marshal(ImageResizePayload{SourceURL: src})
	if err != nil {
		return nil, fmt.Errorf("marshal payload: %v", err)
	}

	return asynq.NewTask(TypeImageResize, payload), nil
}

func (p *ImageProcessor) ProcessTask(ctx context.Context, task *asynq.Task) error {
	var payload ImageResizePayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("unmarshal %s: %v", string(task.Payload()), err)
	}

	log.Printf("Resizing image: src=%s\n", payload.SourceURL)

	return nil
}

func NewImageProcessor() *ImageProcessor {
	return &ImageProcessor{}
}
