class Application < ApplicationRecord
  has_many :chats, foreign_key: :application_token, primary_key: :token, dependent: :destroy

  validates :name, presence: true, length: { maximum: 255, minimum: 3 }

  def as_json(options = {})
    super(options.merge(except: [:id]))
  end
end
