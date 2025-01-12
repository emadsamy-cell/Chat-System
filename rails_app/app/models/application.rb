class Application < ApplicationRecord
  validates :name, presence: true, length: { maximum: 255, minimum: 3 }
end
