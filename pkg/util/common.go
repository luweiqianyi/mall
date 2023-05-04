package util

import "github.com/google/uuid"

// GenerateToken token = sha256(appId + appKey + channelId + userId + nonce + timestamp)
func GenerateToken() string {
	return uuid.NewString()
}
