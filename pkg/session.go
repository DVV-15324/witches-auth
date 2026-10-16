package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type SessionConfig struct {
	SessionTTL  int64
	IdleTimeout int64 // seconds
}

type SessionCache struct {
	SessionID   string `json:"session_id"`
	UserID      int64  `json:"user_id"`
	JKT         string `json:"jkt"`
	AccessToken string `json:"access_token"`
	LastActive  int64  `json:"last_active"`
}

type SessionService struct {
	redis       *redis.Client
	SessionTTL  int64
	IdleTimeout int64
}

func NewSessionService(redis *redis.Client, SessionTTL, IdleTimeout int64) *SessionService {
	return &SessionService{
		redis:       redis,
		SessionTTL:  SessionTTL,
		IdleTimeout: IdleTimeout,
	}
}

func (s *SessionService) cacheKeySession(userID int64, jkt string) string {
	return fmt.Sprintf("session:user:%d:jkt:%s", userID, jkt)
}

// CreateSession - Tạo session mới. Nếu session.SessionID đã có thì giữ nguyên (update).
func (s *SessionService) CreateSession(ctx context.Context, session *SessionCache) error {
	if session == nil {
		return fmt.Errorf("session cannot be nil")
	}
	if session.UserID == 0 {
		return fmt.Errorf("user_id required")
	}
	if session.JKT == "" {
		return fmt.Errorf("jkt required")
	}
	if session.SessionID == "" {
		session.SessionID = uuid.New().String()
	}
	if session.LastActive == 0 {
		session.LastActive = time.Now().Unix()
	}

	key := s.cacheKeySession(session.UserID, session.JKT)
	data, err := json.Marshal(session)
	if err != nil {
		return err
	}

	ttl := time.Duration(s.SessionTTL) * time.Second
	return s.redis.Set(ctx, key, data, ttl).Err()
}

// GetSession - Trả về (nil, nil) nếu không tồn tại, KHÔNG trả error.
func (s *SessionService) GetSession(ctx context.Context, userID int64, jkt string) (*SessionCache, error) {
	key := s.cacheKeySession(userID, jkt)
	data, err := s.redis.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var session SessionCache
	if err := json.Unmarshal([]byte(data), &session); err != nil {
		return nil, err
	}
	return &session, nil
}

// DeleteSession - Xóa session theo userID + jkt
func (s *SessionService) DeleteSession(ctx context.Context, userID int64, jkt string) error {
	key := s.cacheKeySession(userID, jkt)
	return s.redis.Del(ctx, key).Err()
}

// UpdateSession - Update access token khi refresh
func (s *SessionService) UpdateSession(ctx context.Context, userID int64, jkt, accessToken string) error {
	session, err := s.GetSession(ctx, userID, jkt)
	if err != nil {
		return err
	}
	if session == nil {
		return fmt.Errorf("session not found for user %d, jkt %s", userID, jkt)
	}

	session.AccessToken = accessToken
	session.LastActive = time.Now().Unix()
	return s.CreateSession(ctx, session)
}

// UpdateLastActive - Update last active time
func (s *SessionService) UpdateLastActive(ctx context.Context, userID int64, jkt string) error {
	session, err := s.GetSession(ctx, userID, jkt)
	if err != nil {
		return err
	}
	if session == nil {
		return fmt.Errorf("session not found for user %d, jkt %s", userID, jkt)
	}

	session.LastActive = time.Now().Unix()
	return s.CreateSession(ctx, session)
}

// IsSessionIdle - (true, nil) nếu không tồn tại hoặc idle quá lâu
func (s *SessionService) IsSessionIdle(ctx context.Context, userID int64, jkt string) (bool, error) {
	session, err := s.GetSession(ctx, userID, jkt)
	if err != nil {
		return false, err
	}
	if session == nil {
		return true, nil
	}

	now := time.Now().Unix()
	return (now - session.LastActive) > s.IdleTimeout, nil
}

// ValidateSession - Check tồn tại + token match + idle.
// Middleware gọi hàm này để authorize request.
func (s *SessionService) ValidateSession(
	ctx context.Context,
	userID int64,
	jkt string,
	accessToken string,
) (*SessionCache, error) {
	session, err := s.GetSession(ctx, userID, jkt)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, fmt.Errorf("session not found")
	}

	if session.AccessToken != accessToken {
		return nil, fmt.Errorf("token mismatch")
	}
	if session.JKT != jkt {
		return nil, fmt.Errorf("jkt mismatch")
	}

	now := time.Now().Unix()
	if now-session.LastActive > s.IdleTimeout {
		return nil, fmt.Errorf("session idle timeout")
	}

	return session, nil
}
