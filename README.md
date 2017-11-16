# Mercury

> Mercury, the winged messenger

Mercury is a simple email delivery service that pulls pre-rendered email content (both plain and HTML) off an AMQP queue and delivers it via SMTP.

It expects a very simple JSON payload from the queue, where the email body is gzip compressed and base64 encoded. Here is a very simple example:

```json
{
    "from": "noreply@domain.com",
    "to": "bob@bobby.com",
    "subject": "Hi :)",
    "html_body": "H4sIAAAAAAAE/wAPAPD/PGgxPkhlbGxvITwvaDE+AQAA//9WO9LMDwAAAA==",
    "text_body": "H4sIAAAAAAAE/wAGAPn/SGVsbG8hAQAA//9WzCqdBgAAAA=="
}
```

## Usage

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