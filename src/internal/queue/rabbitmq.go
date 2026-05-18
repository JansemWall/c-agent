package queue

import (
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"ckrof-agent/internal/models"
)

func Publish(ctx context.Context, amqpURL string, payload models.Payload, queueName string) error {
	// 1. Conecta ao Rabbit
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// 2. Garante que a fila existe
	q, err := ch.QueueDeclare(
		queueName,		   // nome da fila
		true,              // durável
		false,             // auto-delete
		false,             // exclusiva
		false,             // no-wait
		nil,
	)

	body, _ := json.Marshal(payload)

	// 3. Publica a mensagem
	return ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
}