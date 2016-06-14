package rabbit

import (
	"github.com/streadway/amqp"
	"testing"
)

func BenchmarkNewTask(b *testing.B) {
	conn, _ := amqp.Dial("amqp://user:password@192.168.99.100:5672/")
	defer conn.Close()

	ch, _ := conn.Channel()
	//ch.Confirm(false)
	defer ch.Close()

	q, _ := ch.QueueDeclare(
		"task_queue", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)

	body := `{"text":"aloxaxaxaxa"}`
	b.SetParallelism(100)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			err := ch.Publish(
				"",     // exchange
				q.Name, // routing key
				false,  // mandatory
				false,
				amqp.Publishing{
					DeliveryMode: amqp.Persistent,
					ContentType:  "application/json",
					Body:         []byte(body),
				},
			)
			if err != nil {
				b.Error(err)
			}
		}
	})
}
