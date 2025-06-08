package constant

// Response messages
const (
	MessageOK                  = "Success"
	MessageCreated             = "Created New Resource"
	MessageBadRequest          = "Bad Request"
	MessageUnauthorized        = "Unauthorized"
	MessageForbidden           = "Forbidden"
	MessageNotFound            = "Resource Not Found"
	MessageInternalServerError = "Internal Server Error"
)

const (
	QueueEncodeService = "EncodeService"
)

const (
	MessageTypeEncodeUploadedVideo = "EncodeUploadedVideo"
)

const (
	ContentTypeJSON = "application/json"
)

const (
	ProtocolTCP  = "tcp"
	ProtocolAMQP = "amqp"
)

const (
	HeaderUserID = "x-user-id"
)

const (
	S3RawVideosDirectory  = "raw-videos"
	S3ThumbnailsDirectory = "thumbnails"
)

const ServiceName = "Upload Service"

const TraceTypeRabbitMQPublish = "RabbitMQ Publish"
