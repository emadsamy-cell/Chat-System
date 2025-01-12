package jobs

import (
	"chat_with_go/utils"
	"context"
	"database/sql"
	"log"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

func BatchUpdateCounts(db *sql.DB) {
	ctx := context.Background()
	rdb := utils.GetRedisClient()

	ticker := time.NewTicker(59 * time.Minute) // Every 1 hour
	defer ticker.Stop()

	for range ticker.C {
		// Update chat counts from Redis
		go updateChatCounts(db, rdb, ctx)

		go updateMessageCounts(db, rdb, ctx)
	}
}

func updateChatCounts(db *sql.DB, rdb *redis.Client, ctx context.Context) {
	chatKeys, err := rdb.Keys(ctx, "new_chat:*").Result()
	if err != nil {
		log.Printf("Error fetching chat keys: %v", err)
		return
	}

	if len(chatKeys) == 0 {
		return
	}

	updates := make([]string, 0)
	tokens := make([]interface{}, 0)
	inPlaceholders := make([]string, 0)

	for _, key := range chatKeys {
		applicationToken, _ := extractKey(key)
		chatsCount := rdb.Get(ctx, applicationToken).Val()
		// Add a CASE condition
		updates = append(updates, "WHEN token = ? THEN ?")
		tokens = append(tokens, applicationToken, chatsCount)
	}

	for _, key := range chatKeys {
		applicationToken, _ := extractKey(key)
		// Add the token for the IN clause
		inPlaceholders = append(inPlaceholders, "?")
		tokens = append(tokens, applicationToken)
	}

	// Build the bulk update query
	query := `
		UPDATE applications
		SET chats_count = CASE ` + strings.Join(updates, " ") + `
		END
		WHERE token IN (` + strings.Join(inPlaceholders, ", ") + `);
	`

	// Execute the query
	_, err = db.Exec(query, tokens...)
	if err != nil {
		log.Printf("Error executing bulk chat count update: %v", err)
		return
	}

	// Remove Redis keys after successful update
	for _, key := range chatKeys {
		rdb.Del(ctx, key)
	}

	log.Printf("Bulk chat count update completed for %d keys.", len(chatKeys))
}

func updateMessageCounts(db *sql.DB, rdb *redis.Client, ctx context.Context) {
	messageKeys, err := rdb.Keys(ctx, "new_message:*:*").Result()
	if err != nil {
		log.Printf("Error fetching message keys: %v", err)
		return
	}

	if len(messageKeys) == 0 {
		return
	}

	updates := make([]string, 0)
	tokens := make([]interface{}, 0)
	inPlaceholders := make([]string, 0)

	for _, key := range messageKeys {
		applicationToken, chatNumber := extractKey(key)
		messagesCount := rdb.Get(ctx, applicationToken+":"+chatNumber).Val()

		// Add the `CASE` condition
		updates = append(updates, "WHEN application_token = ? AND chat_number = ? THEN ?")
		tokens = append(tokens, applicationToken, chatNumber, messagesCount)
	}

	for _, key := range messageKeys {
		applicationToken, chatNumber := extractKey(key)
		inPlaceholders = append(inPlaceholders, "(?, ?)")
		tokens = append(tokens, applicationToken, chatNumber)
	}

	// Build the bulk update query
	query := `
		UPDATE chats
		SET messages_count = CASE ` + strings.Join(updates, " ") + `
		END
		WHERE (application_token, chat_number) IN (` + strings.Join(inPlaceholders, ", ") + `);
	`

	// Execute the query
	_, err = db.Exec(query, tokens...)
	if err != nil {
		log.Printf("Error executing bulk message count update: %v", err)
		return
	}

	// Remove Redis keys after successful update
	for _, key := range messageKeys {
		rdb.Del(ctx, key)
	}

	log.Printf("Bulk message count update completed for %d keys.", len(messageKeys))
}

func extractKey(key string) (string, string) {
	parts := strings.Split(key, ":")
	if len(parts) == 2 {
		return parts[1], ""
	}
	return parts[1], parts[2]
}
