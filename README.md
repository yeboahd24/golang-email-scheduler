# golang-email-scheduler

## Usage

```shell
go run .
```

## Sample Request

```shell
curl -X POST \
  http://localhost:8000/schedule \
  -H 'Content-Type: application/json' \
  -d '{
    "to": "yeboahd24@example.com",
    "subject": "Email Scheduler",
    "body": "This is a test email mesika.",
    "send_at": "2024-06-21T21:22:00Z"
}'
```

## Sample Response

```json
{"id":"b011349e-ff05-4f24-9dd3-25406f0fbb35","message":"Email scheduled successfully"}
```
