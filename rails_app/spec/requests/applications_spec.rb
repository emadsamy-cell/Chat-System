# spec/requests/applications_spec.rb
require 'rails_helper'

RSpec.describe "Applications", type: :request do
  let!(:application) { Application.create(name: "Test App", token: SecureRandom.hex) }

  describe "GET" do
    it "Successful returns all applications" do
      applications = Application.all
      get "/applications"
      expect(response).to have_http_status(:success)
      expect(JSON.parse(response.body).size).to eq(applications.size)
    end
  end

  describe "Create New Application" do
    it "Successfully create new application" do
      expect {
        post "/applications", params: { name: "New App" }
      }.to change { Application.count }.by(1)


      expect(response).to have_http_status(:created)
      json_response = JSON.parse(response.body)
      expect(json_response['name']).to eq("New App")
      expect(json_response['token']).not_to be_nil
    end

    it "Successfully create key in Redis after create new application" do
      expect(RedisKeyService).to receive(:set_new_application).once
      post "/applications", params: { name: "New App" }
    end

    it "Fails to create new application without name" do
      post "/applications"
      expect(response).to have_http_status(:unprocessable_entity)
    end

    it "Fails to create new application with empty name" do
      post "/applications", params: { name: "" }
      expect(response).to have_http_status(:unprocessable_entity)
    end
  end


  describe "Update Application" do
    it "Successfully update " do
      patch "/applications/#{application.token}", params: { name: "Updated App" }
      expect(response).to have_http_status(:success)
      expect(JSON.parse(response.body)['message']).to eq("Application updated successfully")


      updatedApp = Application.find_by(token: application.token)
      expect(updatedApp.name).to eq("Updated App")
    end

    it "Fails to update application without name" do
      patch "/applications/#{application.token}"
      expect(response).to have_http_status(:bad_request)
      expect(JSON.parse(response.body)['error']).to eq("Name is required")
    end

    it "Fails to update application with empty name" do
      patch "/applications/#{application.token}", params: { name: "" }
      expect(response).to have_http_status(:bad_request)
      expect(JSON.parse(response.body)['error']).to eq("Name is required")
    end

    it "Fails to update application with invalid token" do
      patch "/applications/invalid_token", params: { name: "Updated App" }
      expect(response).to have_http_status(:not_found)
      expect(JSON.parse(response.body)['error']).to eq("Application not found")
    end
  end
end