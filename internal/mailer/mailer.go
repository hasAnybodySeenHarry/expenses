package mailer

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/streadway/amqp"
)

type Mailer struct {
	conn  *amqp.Connection
	ch    *amqp.Channel
	queue string
}

func New(conn *amqp.Connection, queue string) (*Mailer, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	_, err = ch.QueueDeclare(
		queue, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		ch.Close()
		return nil, fmt.Errorf("failed to declare a queue: %w", err)
	}

	m := &Mailer{
		conn:  conn,
		ch:    ch,
		queue: queue,
	}

	return m, nil
}

func (m *Mailer) Send(tmpl string, name, recipient string, token string) error {
	msg := struct {
		Email    string `json:"email"`
		Name     string `json:"name"`
		Token    string `json:"token"`
		Template string `json:"template"`
	}{
		Email:    recipient,
		Name:     name,
		Token:    token,
		Template: tmpl,
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	for i := 0; i < 3; i++ {
		err = m.ch.Publish(
			"",            // exchange
			"email_queue", // routing key (queue name)
			false,
			false,
			amqp.Publishing{
				ContentType: "application/json",
				Body:        body,
			},
		)
		if nil == err {
			return nil
		}

		time.Sleep(2 * time.Second)
	}

	return err
}

func (m *Mailer) Reconnect() error {
	newCh, err := m.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a new channel: %w", err)
	}

	m.ch.Close()
	m.ch = newCh

	return nil
}

func (m *Mailer) Close() {
	if m.ch != nil {
		m.ch.Close()
	}
	if m.conn != nil {
		m.conn.Close()
	}
}
