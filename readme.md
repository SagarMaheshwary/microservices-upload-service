# MICROSERVICES - UPLOAD SERVICE

Upload Service for the [Microservices](https://github.com/SagarMaheshwary/microservices) project.

### OVERVIEW

- Golang
- ZeroLog
- gRPC – Serves as the main server for service communication
- RabbitMQ - Enables asynchronous communication with the [encode service](https://github.com/SagarMaheshwary/microservices-encode-service)
- Amazon S3 - Handles generating signed urls from s3 for video uploads that are later processed by encode service
- Prometheus Client – Exports default and custom metrics for Prometheus server monitoring
- Jaeger – Distributed request tracing

### SETUP

Follow the instructions in the [README](https://github.com/SagarMaheshwary/microservices?tab=readme-ov-file#setup) of the main microservices repository to run this service along with others using Docker Compose or Kubernetes (KIND).

### APIs (gRPC)

Proto files are located in the **internal/proto** directory.

| SERVICE       | RPC                | BODY                                                                                                                                                                                                               | METADATA | DESCRIPTION                                                       |
| ------------- | ------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | -------- | ----------------------------------------------------------------- |
| UploadService | CreatePresignedUrl | -                                                                                                                                                                                                                  | -        | Generates a presigned url that can be used to upload videos to s3 |
| UploadService | UploadedWebhook    | {"video_id": "string - s3 upload id from presigned-url process", "thumbnail_id": "string - s3 upload id from presigned-url process", "title": "string - video title", "description": "string - video description"} | -        | Send video data to encode service via rabbitmq for video encoding |
| Health        | Check              | -                                                                                                                                                                                                                  | -        | Service health check                                              |

### APIs (REST)

| API      | METHOD | BODY | Headers | Description                 |
| -------- | ------ | ---- | ------- | --------------------------- |
| /metrics | GET    | -    | -       | Prometheus metrics endpoint |

### RABBITMQ MESSAGES

#### Sent Messages (Published to the Queue)

| MESSAGE NAME        | SENT TO                                                                           | DESCRIPTION                                                          |
| ------------------- | --------------------------------------------------------------------------------- | -------------------------------------------------------------------- |
| EncodeUploadedVideo | [Encode Service](https://github.com/SagarMaheshwary/microservices-encode-service) | Notifies the Encode service that a new video is ready for processing |
