# spec/requests/chats_spec.rb
require 'rails_helper'

RSpec.describe "Chats", type: :request do
  let!(:application) { Application.order(:token).first }
  let!(:chat) { Chat.where(application_token: application.token).first }

  describe "GET By Token" do
    it "returns all chats for a specific application" do
      get "/applications/#{application.token}/chats"
      expect(response).to have_http_status(:ok)
    end
  end

  describe "Create By Token" do
    it "creates a new chat for a specific application" do
      post "/applications/#{application.token}/chats", params: { name: "New Chat" }
      
      expect(response).to have_http_status(:created)
      json_response = JSON.parse(response.body)
      expect(json_response['chat_number']).not_to be_nil
    end

    it "fails to create a new chat for a specific application if the application does not exist" do
      post "/applications/invalid_token/chats", params: { name: "New Chat" }
      expect(response).to have_http_status(:not_found)
    end

    it "fails to create a new chat for a specific application without name" do
      post "/applications/#{application.token}/chats"
      expect(response).to have_http_status(:bad_request)
    end

    it "fails to create a new chat for a specific application with empty name" do
      post "/applications/#{application.token}/chats", params: { name: "" }
      expect(response).to have_http_status(:bad_request)
    end
  end

  describe "Update chat by token and chat_number" do
    it "Successfully update " do
      patch "/applications/#{application.token}/chats/#{chat.chat_number}", params: {name: "updated one"}
      expect(response).to have_http_status(:ok)
    end
    
    it "Fails to update chat without name" do
      patch "/applications/#{application.token}/chats/#{chat.chat_number}"
      expect(response).to have_http_status(:bad_request)
      expect(JSON.parse(response.body)['error']).to eq("Name is required")
    end

    it "Fails to update chat with empty name" do
      patch "/applications/#{application.token}/chats/#{chat.chat_number}", params: { name: "" }
      expect(response).to have_http_status(:bad_request)
      expect(JSON.parse(response.body)['error']).to eq("Name is required")
    end

    it "Fails to update chat with invalid token" do
      patch "/applications/invalid_token/chats/#{chat.chat_number}", params: { name: "Updated App" }
      expect(response).to have_http_status(:not_found)
      expect(JSON.parse(response.body)['error']).to eq("Chat not found")
    end

    it "Fails to update chat with invalid chat number" do
      patch "/applications/#{application.token}/chats/invalid_chat_number", params: { name: "Updated App" }
      expect(response).to have_http_status(:not_found)
      expect(JSON.parse(response.body)['error']).to eq("Chat not found")
    end
  end
end