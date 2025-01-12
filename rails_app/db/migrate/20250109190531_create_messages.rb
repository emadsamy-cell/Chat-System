class CreateMessages < ActiveRecord::Migration[8.0]
  def change
    create_table :messages, id: false do |t|
      t.integer :message_number, null: false
      t.string :application_token, null: false
      t.integer :chat_number, null: false
      t.string :body, null: false
    end
    add_index :messages, [ :application_token, :chat_number, :message_number ], unique: true
    add_foreign_key :messages, :chats, column: [ :application_token, :chat_number ], primary_key: [ :application_token, :chat_number ]
  end
end
