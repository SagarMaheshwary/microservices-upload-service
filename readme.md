# MICROSERVICES - UPLOAD SERVICE

This service is a part of the Microservices project built for handling video uploads.

### TECHNOLOGIES

- Golang (1.22.2)
- gRPC
- RabbitMQ
- Amazon S3

### SETUP

cd into the project directory and copy **.env.example** to **.env** and update the required variables.

Create executable and start the server:

```bash
go build cmd/server/main.go && ./main
```

Or install "[air](https://github.com/cosmtrek/air)" and run it to autoreload when making file changes:

```bash
air -c .air-toml
```

### APIs (RPC)

| SERVICE       | RPC                | METADATA | DESCRIPTION                                                                                                                                       |
| ------------- | ------------------ | -------- | ------------------------------------------------------------------------------------------------------------------------------------------------- |
| UploadService | CreatePresignedUrl | -        | Generates a presigned url that can be used to upload videos to s3.                                                                                |
| UploadService | UploadedWebhook    | -        | Once the video is uploaded via the presigned url, we will inform "encode" service by this RPC to start encoding the uploaded video for streaming. |
