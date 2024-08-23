package notifier

import (
	"encoding/json"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"harry2an.com/expenses/internal/data"
)

const (
	debtCreated = "debt_created"
)

type debt struct {
	Metadata metadata `json:"metadata"`
	Data     debtData `json:"data"`
}

type debtData = data.Debt

type DebtNotifier struct {
	p sarama.SyncProducer
}

func (m *DebtNotifier) Send(d *data.Debt) error {
	payload := debt{}

	payload.Metadata = metadata{
		ID:        uuid.New().String(),
		Timestamp: time.Now().Format(time.RFC3339),
		Type:      debtCreated,
		Source:    systemName,
		Version:   version,
	}

	payload.Data = *d

	json, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: debtTopic,
		Value: sarama.ByteEncoder(json),
	}

	_, _, err = m.p.SendMessage(msg)
	if err != nil {
		return err
	}

	return nil
}
