class ElasticsearchUpdateJob < ApplicationJob
  queue_as :default

  def perform(application_token, chat_number, message_number, new_body)
    # Do something later
    doc_id = "#{application_token}_#{chat_number}_#{message_number}"

    client = Elasticsearch::Model.client
    client.update(
      index: "messages",
      id: doc_id,
      body: {
        doc: {
          body: new_body
        }
      }
    )

    rescue => e
      Rails.logger.error("Elasticsearch update failed for Message ID #{message_id}: #{e.message}")
  end
end
