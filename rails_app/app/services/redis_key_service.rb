class RedisKeyService
  def self.set_new_application(application_token) 
    $redis.set(application_token, 0)
  end

  def self.get_application(application_token)
    return $redis.get(application_token)
  end

  def self.add_new_chat(application_token)
    # Add flag to update applications table in go app every one hour
    $redis.set("new_chat:#{application_token}", true)

    # Increment the chat number
    return $redis.incr(application_token)
  end

  def self.set_new_chat(application_token, chat_number)
    $redis.set("#{application_token}:#{chat_number}", 0)
  end

  def self.get_chat(application_token, chat_number)
    return $redis.get("#{application_token}:#{chat_number}")
  end

  def self.add_new_message(application_token, chat_number)
    # Add flag to update chats table in go app every one hour
    $redis.set("new_message:#{application_token}:#{chat_number}", true)

    # Increment the message number
    return $redis.incr("#{application_token}:#{chat_number}")
  end
end
