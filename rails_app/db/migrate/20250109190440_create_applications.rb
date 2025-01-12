class CreateApplications < ActiveRecord::Migration[8.0]
  def change
    create_table :applications, id: false do |t|
      t.string :token, null: false
      t.string :name, null: false
      t.integer :chats_count, null: false, default: 0, comment: 'Number of chats of that application'
    end
    add_index :applications, :token, unique: true
  end
end
