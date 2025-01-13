class ChatsController < ApplicationController
  def index
    @chats = Chat.where(application_token: params[:application_id]).select(:chat_number, :name, :messages_count)
    render json: @chats, status: :ok
  end

  def create
    application_token = params[:application_id]
    name = params[:name]

    # Check if the application exists
    if RedisKeyService.get_application(application_token).nil?
      render json: { error: "Application not found" }, status: :not_found
      return
    end

    if name.nil? || name.empty?
      render json: { error: "Name is required" }, status: :bad_request
      return
    end

    chat_number = RedisKeyService.add_new_chat(application_token)

    RedisKeyService.set_new_chat(application_token, chat_number)

    MessageQueueService.enqueue_chat_creation(application_token, chat_number, name)

    render json: { chat_number: chat_number }, status: :created
  end

  def update
    application_token = params[:application_id]
    chat_number = params[:id]
    name = params[:name]

    if name.nil? || name.empty?
      render json: { error: "Name is required" }, status: :bad_request
      return
    end

    isUpdated = Chat.where(application_token: application_token, chat_number: chat_number).update_all(name: name)

    if isUpdated == 0
      render json: { error: "Chat not found" }, status: :not_found
    else
      render json: { message: "Chat updated successfully" }, status: :ok
    end
  end
end
