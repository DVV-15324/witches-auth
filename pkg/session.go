package utils

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type SessionConfig struct {
	SessionTTL  int64
	IdleTimeout int64
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

func (s *SessionService) cacheKeyByToken(userID int64, accessToken string) string {
	h := sha256.Sum256([]byte(accessToken))
	return fmt.Sprintf("session:user:%d:token:%x", userID, h)
}

func (s *SessionService) CreateSession(ctx context.Context, session *SessionCache) error {
	if session == nil {
		return fmt.Errorf("session cannot be nil")
	}
	if session.UserID == 0 {
		return fmt.Errorf("user_id required")
	}
	if session.AccessToken == "" {
		return fmt.Errorf("access_token required")
	}
	// Chỉ mobile mới cần JKT; web thì rỗng là bình thường
	if session.SessionID == "" {
		session.SessionID = uuid.New().String()
	}
	if session.LastActive == 0 {
		session.LastActive = time.Now().Unix()
	}

	key := s.cacheKeyByToken(session.UserID, session.AccessToken)
	data, err := json.Marshal(session)
	if err != nil {
		return err
	}

	ttl := time.Duration(s.SessionTTL) * time.Second
	return s.redis.Set(ctx, key, data, ttl).Err()
}

func (s *SessionService) GetSessionByToken(ctx context.Context, userID int64, accessToken string) (*SessionCache, error) {
	return s.getByKey(ctx, s.cacheKeyByToken(userID, accessToken))
}

func (s *SessionService) getByKey(ctx context.Context, key string) (*SessionCache, error) {
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

func (s *SessionService) DeleteSessionByToken(ctx context.Context, userID int64, accessToken string) error {
	return s.redis.Del(ctx, s.cacheKeyByToken(userID, accessToken)).Err()
}

func (s *SessionService) UpdateSession(ctx context.Context, userID int64, oldAccessToken, newAccessToken string) error {
	var (
		session *SessionCache
		err     error
	)
	session, err = s.GetSessionByToken(ctx, userID, oldAccessToken)
	if err != nil {
		return err
	}
	if session == nil {
		return fmt.Errorf("session not found for user %d", userID)
	}

	if oldAccessToken != newAccessToken {
		_ = s.DeleteSessionByToken(ctx, userID, oldAccessToken)
	}
	session.AccessToken = newAccessToken
	session.LastActive = time.Now().Unix()
	return s.CreateSession(ctx, session)
}

func (s *SessionService) UpdateLastActive(ctx context.Context, userID int64, accessToken string) error {
	var (
		session *SessionCache
		err     error
	)
	session, err = s.GetSessionByToken(ctx, userID, accessToken)
	if err != nil {
		return err
	}
	if session == nil {
		return fmt.Errorf("session not found for user %d", userID)
	}

	session.LastActive = time.Now().Unix()
	return s.CreateSession(ctx, session)
}

func (s *SessionService) ValidateSession(ctx context.Context, userID int64, jkt string, accessToken string) (*SessionCache, error) {
	var (
		session *SessionCache
		err     error
	)
	session, err = s.GetSessionByToken(ctx, userID, accessToken)
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
