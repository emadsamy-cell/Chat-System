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
        },
        _source: ["body", "message_number"]
      }
    )
  end

  def as_json(options = {})
    super(options.merge(except: [:id]))
  end
end
