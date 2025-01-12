# rails_app/Dockerfile
FROM ruby:3.3.5

# Install dependencies
RUN apt-get update -qq && apt-get install -y \
    build-essential \
    libpq-dev \
    nodejs \
    default-mysql-client

# Set working directory
WORKDIR /app

# Install gems
COPY Gemfile Gemfile.lock ./
RUN bundle install


# Remove Server file
RUN rm -f /app/tmp/pids/server.pid

# Copy the application
COPY . .

# Start the main process
CMD ["rails", "server", "-b", "0.0.0.0"]