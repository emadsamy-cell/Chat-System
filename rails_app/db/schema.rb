# This file is auto-generated from the current state of the database. Instead
# of editing this file, please use the migrations feature of Active Record to
# incrementally modify your database, and then regenerate this schema definition.
#
# This file is the source Rails uses to define your schema when running `bin/rails
# db:schema:load`. When creating a new database, `bin/rails db:schema:load` tends to
# be faster and is potentially less error prone than running all of your
# migrations from scratch. Old migrations may fail to apply correctly if those
# migrations use external dependencies or application code.
#
# It's strongly recommended that you check this file into your version control system.

ActiveRecord::Schema[8.0].define(version: 2025_01_09_190531) do
  create_table "applications", id: false, charset: "utf8mb4", collation: "utf8mb4_0900_ai_ci", force: :cascade do |t|
    t.string "token", null: false
    t.string "name", null: false
    t.integer "chats_count", default: 0, null: false, comment: "Number of chats of that application"
    t.index ["token"], name: "index_applications_on_token", unique: true
  end

  create_table "chats", id: false, charset: "utf8mb4", collation: "utf8mb4_0900_ai_ci", force: :cascade do |t|
    t.integer "chat_number", null: false
    t.string "application_token", null: false
    t.string "name", null: false
    t.integer "messages_count", default: 0, null: false, comment: "Number of messages in that chat"
    t.index ["application_token", "chat_number"], name: "index_chats_on_application_token_and_chat_number", unique: true
  end

  create_table "messages", id: false, charset: "utf8mb4", collation: "utf8mb4_0900_ai_ci", force: :cascade do |t|
    t.integer "message_number", null: false
    t.string "application_token", null: false
    t.integer "chat_number", null: false
    t.string "body", null: false
    t.index ["application_token", "chat_number", "message_number"], name: "idx_on_application_token_chat_number_message_number_51bfd3c604", unique: true
  end

  add_foreign_key "chats", "applications", column: "application_token", primary_key: "token"
  add_foreign_key "messages", "chats", column: ["application_token", "chat_number"], primary_key: ["application_token", "chat_number"]
end
