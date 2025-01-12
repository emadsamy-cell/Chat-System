class Message < ApplicationRecord
  def self.search_message(application_token, chat_number, query)
    Elasticsearch::Model.client.search(
      index: "messages",
      body: {
        query: {
          bool: {
            must: [
              { match: { application_token: application_token } },
              { match: { chat_number: chat_number } },
              { match: { body: query } }
            ]
          }
        }
      }
    )
  end
end
