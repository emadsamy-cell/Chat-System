class Chat < ApplicationRecord
  belongs_to :application, foreign_key: :application_token, primary_key: :token

  def as_json(options = {})
    super(options.merge(except: [:id]))
  end
end
