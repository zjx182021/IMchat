package models

import (
	"TM_chat/utils"
	"context"
	"time"
)

func SetUserOnlineInfo(key string, value []byte, timeTTL time.Duration) {
	ctx := context.Background()
	utils.REDIS.Set(ctx, key, value, timeTTL)
}
