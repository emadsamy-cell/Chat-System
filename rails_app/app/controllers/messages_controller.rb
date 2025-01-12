class MessagesController < ApplicationController
  def index
    @messages = Message.where(application_token: params[:application_id], chat_number: params[:chat_id])
    render json: @messages
  end

  def create
    application_token = params[:application_id]
    chat_number = params[:chat_id]
    body = params[:body]

    if RedisKeyService.get_chat(application_token, chat_number).nil?
      render json: { error: "Chat not found" }, status: :not_found
      return
    end

    if body.nil?
      render json: { error: "Message body cannot be empty" }, status: :bad_request
      return
    end

    chat_number = chat_number.to_i

    message_number = RedisKeyService.add_new_message(application_token, chat_number)
    MessageQueueService.enqueue_message_creation(application_token, chat_number, body, message_number)
    
    render json: { message_number: message_number }, status: :accepted
  end

  def search
    application_token = params[:application_id]
    chat_number = params[:chat_id]
    query = params[:query]

    results = Message.search_message(application_token, chat_number, query)
    render json: results["hits"]["hits"].map { |hit| hit["_source"] }
  end

  def update
    application_token = params[:application_id]
    chat_number = params[:chat_id].to_i
    message_number = params[:id].to_i
    new_body = params[:body]

    if new_body.nil?
      render json: { error: "Message body cannot be empty" }, status: :bad_request
      return
    end
  
    # Update the message directly in the database
    rows_affected = Message.where(
      application_token: application_token,
      chat_number: chat_number,
      message_number: message_number
    ).update_all(body: new_body)
  
    if rows_affected > 0
      
      ElasticsearchUpdateJob.perform_later(application_token, chat_number, message_number, new_body)
      render json: { message: "Message updated successfully" }, status: :ok
    else
      render json: { error: "Message not found" }, status: :not_found
    end
  end

end
