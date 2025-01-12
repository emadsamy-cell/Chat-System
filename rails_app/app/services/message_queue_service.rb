require "bunny"

class MessageQueueService
  def self.enqueue_chat_creation(application_token, chat_number, name)
    connection = Bunny.new
    connection.start

    channel = connection.create_channel
    queue = channel.queue("chat_creation", durable: true)

    message = { application_token: application_token, chat_number: chat_number, name:name }.to_json
    queue.publish(message, persistent: true)

    connection.close
  end

  def self.enqueue_message_creation(application_token, chat_number, body, message_number)
    connection = Bunny.new
    connection.start

    channel = connection.create_channel
    queue = channel.queue("message_creation", durable: true)

    message = { application_token: application_token, chat_number: chat_number, body: body, message_number: message_number }.to_json
    queue.publish(message, persistent: true)

    connection.close
  end
end
