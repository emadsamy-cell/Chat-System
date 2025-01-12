class CreateChats < ActiveRecord::Migration[8.0]
  def change
    create_table :chats, id: false do |t|
      t.integer :chat_number, null: false
      t.string :application_token, null: false
      t.string :name, null: false
      t.integer :messages_count, null: false, default: 0, comment: 'Number of messages in that chat'
    end
    add_index :chats, [ :application_token, :chat_number ], unique: true
    add_foreign_key :chats, :applications, column: :application_token, primary_key: :token
  end
end
