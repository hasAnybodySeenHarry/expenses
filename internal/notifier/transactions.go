package notifier

import (
	"encoding/json"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"harry2an.com/expenses/internal/core"
	"harry2an.com/expenses/internal/data"
)

const (
	transactionCreated = "transaction_created"
)

type transaction struct {
	Metadata metadata        `json:"metadata"`
	Data     transactionData `json:"data"`
}

type transactionData struct {
	ID          int64       `json:"id"`
	Lender      core.Entity `json:"lender"`
	Borrower    core.Entity `json:"borrower"`
	DebtID      int64       `json:"debt_id"`
	Amount      float64     `json:"amount"`
	Description string      `json:"description"`
	CreatedAt   time.Time   `json:"created_at"`
	Version     uuid.UUID   `json:"version"`
}

type TransactionNotifier struct {
	p sarama.SyncProducer
}

func (m *TransactionNotifier) Send(t *data.Transaction) error {
	payload := transaction{}

	payload.Metadata = metadata{
		ID:        uuid.New().String(),
		Timestamp: time.Now().Format(time.RFC3339),
		Type:      transactionCreated,
		Source:    systemName,
		Version:   version,
	}

	payload.Data = transactionData{
		ID:          t.ID,
		DebtID:      t.DebtID,
		Amount:      t.Amount,
		Description: t.Description,
		CreatedAt:   t.CreatedAt,
		Version:     t.Version,
		Lender:      t.Lender,
		Borrower:    t.Borrower,
	}

	json, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: transactionTopic,
		Value: sarama.ByteEncoder(json),
	}

	_, _, err = m.p.SendMessage(msg)
	if err != nil {
		return err
	}

	return nil
}
