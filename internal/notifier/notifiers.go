package notifier

import "github.com/IBM/sarama"

const (
	transactionTopic = "transactions"
	debtTopic        = "debts"
	version          = "1.0"
	systemName       = "expenses"
)

type Notifiers struct {
	Debts        DebtNotifier
	Transactions TransactionNotifier
}

func New(producer sarama.SyncProducer) *Notifiers {
	return &Notifiers{
		Debts:        DebtNotifier{p: producer},
		Transactions: TransactionNotifier{p: producer},
	}
}

type metadata struct {
	ID        string `json:"id"`
	Timestamp string `json:"timestamp"`
	Type      string `json:"type"`
	Source    string `json:"source"`
	Version   string `json:"version"`
}
