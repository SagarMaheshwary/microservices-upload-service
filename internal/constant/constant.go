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
	HeaderUserId = "x-user-id"
)

const RawVideosDirectory = "raw-videos"
