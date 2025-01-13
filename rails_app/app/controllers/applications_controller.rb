class ApplicationsController < ApplicationController
  def index
    @applications = Application.all
    render json: @applications, status: :ok
  end

  def create
    token = SecureRandom.uuid
    name = params[:name]
    @application = Application.new(name: name, token: token)
    if @application.save
      RedisKeyService.set_new_application(@application.token)
      render json: @application, status: :created
    else
      render json: @application.errors, status: :unprocessable_entity
    end
  end

  def update
    token = params[:id]
    name = params[:name]

    if name.nil? || name == ""
      render json: { error: "Name is required" }, status: :bad_request
      return
    end

    isUpdated = Application.where(token: token).update_all(name: name)

    if isUpdated == 0
      render json: { error: "Application not found" }, status: :not_found
    else
      render json: { message: "Application updated successfully" }, status: :ok
    end
  end
end
