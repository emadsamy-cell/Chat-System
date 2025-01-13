# spec/requests/messages_spec.rb
require 'rails_helper'

RSpec.describe "Messages", type: :request do
  let!(:application) { Application.order(:token).first }
  let!(:chat) { Chat.where(application_token: application.token).first }
  let!(:message) { Message.where(application_token: application.token, chat_number: chat.chat_number).first }

  describe "Index" do
    it "returns all messages for a specific chat" do
      get "/applications/#{application.token}/chats/#{chat.chat_number}/messages"
      expect(response).to have_http_status(:success)
      json_response = JSON.parse(response.body)
    end
  end

  describe "Create" do
    it "creates a new message for a specific chat" do
      post "/applications/#{application.token}/chats/#{chat.chat_number}/messages", params: { body: "New Message" }

      expect(response).to have_http_status(:created)
      json_response = JSON.parse(response.body)
      expect(json_response['message_number']).not_to be_nil
    end

    it "fails not found invalid application token" do
      post "/applications/Invalid_token/chats/#{chat.chat_number}/messages", params: { body: "New Chat" }
      expect(response).to have_http_status(:not_found)
    end

    it "fails not found invalid chat number" do
      post "/applications/#{application.token}/chats/Invalid_Chat_Number/messages", params: { body: "New Chat" }
      expect(response).to have_http_status(:not_found)
    end

    it "fails create with empty body" do
      post "/applications/#{application.token}/chats/#{chat.chat_number}/messages", params: { body: "" }
      expect(response).to have_http_status(:bad_request)
    end

    it "fails create without body" do
      post "/applications/#{application.token}/chats/#{chat.chat_number}/messages", params: { body: "" }
      expect(response).to have_http_status(:bad_request)
    end
  end

  describe "Update" do
    it "Success " do
      patch "/applications/#{application.token}/chats/#{chat.chat_number}/messages/#{message.message_number}", params: {body: "updated one"}
      expect(response).to have_http_status(:ok)
    end
    
    it "Fails to update chat without body" do
      patch "/applications/#{application.token}/chats/#{chat.chat_number}/messages/#{message.message_number}"
      expect(response).to have_http_status(:bad_request)
      expect(JSON.parse(response.body)['error']).to eq("Message body cannot be empty")
    end

    it "Fails to update chat with empty body" do
      patch "/applications/#{application.token}/chats/#{chat.chat_number}/messages/#{message.message_number}", params: { body: "" }
      expect(response).to have_http_status(:bad_request)
      expect(JSON.parse(response.body)['error']).to eq("Message body cannot be empty")
    end

    it "Fails to update message with invalid token" do
      patch "/applications/invalid_token/chats/#{chat.chat_number}/messages/#{message.message_number}", params: { body: "Updated App" }
      expect(response).to have_http_status(:not_found)
      expect(JSON.parse(response.body)['error']).to eq("Message not found")
    end

    it "Fails to update message with invalid chat number" do
      patch "/applications/#{application.token}/chats/invalid_chat_number/messages/#{message.message_number}", params: { body: "Updated App" }
      expect(response).to have_http_status(:not_found)
      expect(JSON.parse(response.body)['error']).to eq("Message not found")
    end

    it "Fails to update message with invalid message number" do
      patch "/applications/#{application.token}/chats/#{chat.chat_number}/messages/invalid_message_number", params: { body: "Updated App" }
      expect(response).to have_http_status(:not_found)
      expect(JSON.parse(response.body)['error']).to eq("Message not found")
    end
  end
end