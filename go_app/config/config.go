package config

const (
	MySQLDSN     = "root:12345678@tcp(db:3306)/chat-system"
	RabbitMQURL  = "amqp://guest:guest@rabbitmq:5672/"
	RedisAddress = "redis:6379"
	ChatQueue    = "chat_creation"
	MessageQueue = "message_creation"
)
