package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"mercury/email"
	"net/smtp"
	"net/url"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/domodwyer/mailyak"
	"github.com/streadway/amqp"
)

var logger = logrus.New()

// TODO: Command-line config
// TODO: Signal traps
// TODO:

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

	log.Printf("Using mailhost=%s", mailhost.Hostname())
	yak := mailyak.New(*smtpHost, smtp.PlainAuth("", *smtpUser, *smtpPassword, mailhost.Hostname()))

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
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range emails {
			e := email.Email{}

			err := json.Unmarshal(d.Body, &e)
			failOnError(err, "Error decoding JSON payload for Email. Check payload!")

			e.Send(yak)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		os.Exit(-1)
	}
}
