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

	// Process in batches of 500
	batchSize := 500
	for i := 0; i < len(chatKeys); i += batchSize {
		end := i + batchSize
		if end > len(chatKeys) {
			end = len(chatKeys)
		}

		currentBatch := chatKeys[i:end]
		updates := make([]string, 0)
		tokens := make([]interface{}, 0)
		inPlaceholders := make([]string, 0)

		for _, key := range currentBatch {
			applicationToken, _ := extractKey(key)
			chatsCount := rdb.Get(ctx, applicationToken).Val()
			updates = append(updates, "WHEN token = ? THEN ?")
			tokens = append(tokens, applicationToken, chatsCount)
		}

		for _, key := range currentBatch {
			applicationToken, _ := extractKey(key)
			inPlaceholders = append(inPlaceholders, "?")
			tokens = append(tokens, applicationToken)
		}

		query := `
				UPDATE applications
				SET chats_count = CASE ` + strings.Join(updates, " ") + `
				END
				WHERE token IN (` + strings.Join(inPlaceholders, ", ") + `);
			`

		_, err = db.Exec(query, tokens...)
		if err != nil {
			log.Printf("Error executing bulk chat count update for batch: %v", err)
			continue
		}

		// Remove Redis keys after successful update
		for _, key := range currentBatch {
			rdb.Del(ctx, key)
		}

		log.Printf("Batch chat count update completed for %d keys.", len(currentBatch))
	}
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

	// Process in batches of 500
	batchSize := 500
	for i := 0; i < len(messageKeys); i += batchSize {
		end := i + batchSize
		if end > len(messageKeys) {
			end = len(messageKeys)
		}

		currentBatch := messageKeys[i:end]
		updates := make([]string, 0)
		tokens := make([]interface{}, 0)
		inPlaceholders := make([]string, 0)

		for _, key := range currentBatch {
			applicationToken, chatNumber := extractKey(key)
			messagesCount := rdb.Get(ctx, applicationToken+":"+chatNumber).Val()

			updates = append(updates, "WHEN application_token = ? AND chat_number = ? THEN ?")
			tokens = append(tokens, applicationToken, chatNumber, messagesCount)
		}

		for _, key := range currentBatch {
			applicationToken, chatNumber := extractKey(key)
			inPlaceholders = append(inPlaceholders, "(?, ?)")
			tokens = append(tokens, applicationToken, chatNumber)
		}

		query := `
			UPDATE chats
			SET messages_count = CASE ` + strings.Join(updates, " ") + `
			END
			WHERE (application_token, chat_number) IN (` + strings.Join(inPlaceholders, ", ") + `);
		`

		_, err = db.Exec(query, tokens...)
		if err != nil {
			log.Printf("Error executing bulk message count update for batch: %v", err)
			continue
		}

		// Remove Redis keys after successful update
		for _, key := range currentBatch {
			rdb.Del(ctx, key)
		}

		log.Printf("Batch message count update completed for %d keys.", len(currentBatch))
	}
}

func extractKey(key string) (string, string) {
	parts := strings.Split(key, ":")
	if len(parts) == 2 {
		return parts[1], ""
	}
	return parts[1], parts[2]
}
