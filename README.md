# Chat System

This system is designed to manage applications, chats, and messages. It provides functionality to read, create, and update records in these tables.

## Table of Contents

1. [Documentation](#documentation)
   - [Database Schema](#database-schema)
   - [API Controllers](#api-controllers)
   - [Tasks](#tasks)
   - [Workers](#workers)
3. [Running the System](#running-the-system)
4. [Ruby on Rails Application](#ruby-on-rails-application)
   - [Elasticsearch](#elasticsearch)
   - [Redis](#redis)
   - [Concurrency Configuration](#concurrency-configuration)
   - [Testing](#testing)
5. [Go Application](#go-application)
6. [Features](#features)

## Documentation

#### Database Schema 
The schema file can be found in the `db/schema.rb` file, which contains the structure of the database, including tables and their columns.

The database schema for this system includes indexed columns to optimize query performance:
- **Applications Table**: Indexed on `token`.
- **Chats Table**: Indexed on `application_token` and `chat_number`.
- **Messages Table**: Indexed on `application_token`, `chat_number`, and `message_number`.


![Database-schema](https://github.com/user-attachments/assets/55bfd680-cb61-4144-b69f-387f7a447a82)

Here is link to schema: [Database schema](https://drive.google.com/file/d/155w9YRYPV5CAFI96tlRcOvbu-NYCSQGW/view?usp=sharing)




### API Controllers

The API controllers implemented in this system are located in the `rails_app/app/controllers` directory. Below are the specific controllers for each resource:
- **Applications Controller**: `rails_app/app/controllers/applications_controller.rb`
- **Chats Controller**: `rails_app/app/controllers/chats_controller.rb`
- **Messages Controller**: `rails_app/app/controllers/messages_controller.rb`
  
For detailed information on the API endpoints and documentation, please visit the [Postman documentation](https://documenter.getpostman.com/view/40896519/2sAYQXosoD)

### Tasks
This section describes the tasks that are enqueued or performed later in the system:

- **Enqueue Message for Creating Chats and Messages**: `rails_app/app/services/message_queue_service.rb` `In both create chat and message controller`
  - Implemented in the `MessageQueueService` using RabbitMQ to communicate between the Rails app and the Go app.
  
- **Perform Later for Updating Elasticsearch Index**:  `In update message controller`
  - Implemented in the `ElasticsearchUpdateJob` to update the Elasticsearch index asynchronously.

- **Setting Flags for Updates in Go App**: `rails_app/app/services/redis_key_service.rb` `In both create application and chat controller`

  - Flags are set in Redis to indicate new chats and messages, which are later processed in the Go app to update counts.

### Workers

This section describes the workers that process tasks in the background:

- **Elasticsearch Update Job**:  `rails_app/app/jobs/elasticsearch_update_job.rb`
  - Located in `app/jobs/elasticsearch_update_job.rb`, it updates the Elasticsearch index for messages.

- **Go Application Workers**:
  - **Chat Consumer**: Listens for new chat creation tasks and processes them in batches. `go_app/consumer/chat_consumer.go`
  - **Message Consumer**: Listens for new message creation tasks and processes them in batches. `go_app/consumer/message_consumer.go`
  - **Batch Update Counts Job**: Updates `chats_count` and `messages_count` in the database every hour based on Redis data. `go_app/jobs/job.go`



## Running the System
To run the entire system using Docker, ensure you have Docker and Docker Compose installed, then execute the following command:
```bash
docker-compose up
```
This command will build and start all the necessary services defined in the docker-compose.yml file, including the Rails and Go applications, databases, and any other dependencies.


## Ruby on Rails Application

### Elasticsearch
Elasticsearch is used to enable partial matching of message bodies by `chat_number` and `application_token`. This allows for efficient searching and retrieval of messages based on these criteria.

### Redis
- The `RedisKeyService` manages Redis keys for tracking new chats and messages.
- It ensures atomic operations to prevent race conditions, especially when adding new chats.
- It is used to check if an application or chat exists when creating new chats or messages.
- Flags are set in Redis to indicate new chats and messages, which are later processed in the Go app to update counts.

### Concurrency Configuration
- The Ruby on Rails application is configured to work concurrently by adjusting the number of threads and workers.
- This configuration helps improve performance and handle multiple requests simultaneously.
- Ensure that your `config/puma.rb` file is set up to specify the desired number of threads and workers.
- Adjust these values based on your server's capabilities and the expected load.


### Testing
- RSpec is used for testing all endpoints in the Rails application. 
- The test files are located in the `spec` directory.
- To run tests:
  ```
  docker exec -it <rails-containerID> bundle exec rspec spec/requests/applications_spec.rb
  docker exec -it <rails-containerID> bundle exec rspec spec/requests/chats_spec.rb
  docker exec -it <rails-containerID> bundle exec rspec spec/requests/messages_spec.rb
  ```
## Go Application
The Go application includes components for processing tasks and updating data:
- **Consumers**: Listen to message queues for creating chats and messages.
- **Jobs**: Periodically update chat and message counts in the database using a ticker that checks every hour.
- **Controllers (Processors)**: Handle batch processing of chat and message data using goroutines for concurrent execution.


## Features
- **Applications**: Create, read, and update applications.
- **Chats**: Create, read, and update chats associated with applications.
- **Messages**: Create, read, update, and search messages within chats.
