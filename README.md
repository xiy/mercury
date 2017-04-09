# Mercury

> Mercury, the winged messenger

Mercury is a simple email delivery service that pulls pre-rendered email content (both plain and HTML) off an AMQP queue and delivers it via SMTP.

It expects a very simple, but specific, [MessagePack](http://msgpack.org/)-encoded payload from the queue. Here is an example in JSON:

```json
{
    "from": "noreply@domain.com",
    "to": "bob@bobby.com",
    "subject": "Hi :)",
    "html_body": "<h1>Hi!</h1>",
    "text_body": "Hi!"
}
```

## Usage

```
Usage of bin/mercury:
  -amqpHost string
    	The address of the AMQP server to consume from. (default "amqp://guest:guest@localhost:5672")
  -amqpQueue string
    	The queue on the AMQP server to consume email messages from. (default "email")
  -smtpHost string
    	The hostname and port of the SMTP server to send through. (default "localhost:25")
  -smtpPassword string
    	The SMTP password for logging into the SMTP server.
  -smtpUser string
    	The SMTP username for logging into the SMTP server.
```