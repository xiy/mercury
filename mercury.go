package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/smtp"
	"net/url"
	"os"

	"github.com/domodwyer/mailyak"
	"github.com/streadway/amqp"
)

var logger = NewCoLogLogger("mercury")

// TODO: Signal traps

func main() {
	amqpUrL := flag.String(
		"amqpHost",
		"amqp://guest:guest@localhost:5672",
		"The address of the AMQP server to consume from.")
	amqpQueue := flag.String(
		"amqpQueue",
		"email",
		"The queue on the AMQP server to consume email messages from.")
	smtpHost := flag.String(
		"smtpHost",
		"localhost:25",
		"The hostname and port of the SMTP server to send through.")
	smtpUser := flag.String(
		"smtpUser",
		"",
		"The SMTP username for logging into the SMTP server.")
	smtpPassword := flag.String(
		"smtpPassword",
		"",
		"The SMTP password for logging into the SMTP server.")
	flag.Parse()

	// strip port from smtp host for plain PlainAuth
	mailhost, err := url.Parse(fmt.Sprintf("smtp://%s", *smtpHost))
	failOnError(err, "Error parsing mailhost")

	logger.Printf("Using mailhost=%s:%s", mailhost.Hostname(), mailhost.Port())
	logger.Printf("Using amqp=%s", *amqpUrL)

	yak := mailyak.New(
		*smtpHost,
		smtp.PlainAuth("", *smtpUser, *smtpPassword, mailhost.Hostname()),
	)

	// pull message details from rabbit
	conn, err := amqp.Dial(*amqpUrL)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	channel, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer channel.Close()

	q, err := channel.QueueDeclare(
		*amqpQueue,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	emails, err := channel.Consume(
		q.Name,    // queue
		"Mercury", // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)
	go func() {
		for d := range emails {
			go func(d *amqp.Delivery) {
				e := Email{}

				if err := json.Unmarshal(d.Body, &e); err != nil {
					logger.Println("Error decoding JSON payload for Email. Check payload!")
					return
				}

				if err := e.Send(yak); err != nil {
					logger.Printf("Failed sending email, requeing message err=%s", err)
					return
				}

				// Send ACK to remove the message from the queue
				d.Ack(false)
				fmt.Print("*")
			}(&d)
		}
	}()

	logger.Printf("info: Waiting for messages. To exit press CTRL+C")
	<-forever
}

func failOnError(err error, msg string) {
	if err != nil {
		logger.Fatalf("fatal: %s: %s", msg, err)
		os.Exit(-1)
	}
}
