package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

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
	Role        string `json:"role"`
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

// mobile
func (s *SessionService) cacheKeyByJKT(userID int64, jkt string) string {
	return fmt.Sprintf("session:user:%d:jkt:%s", userID, jkt)
}

// website
func (s *SessionService) cacheKeyBySessionID(userID int64, sessionID string) string {
	return fmt.Sprintf("session:user:%d:sessionID:%s", userID, sessionID)
}

// set cache
func (s *SessionService) cacheKeyFor(session *SessionCache) string {
	if session.JKT != "" {
		return s.cacheKeyByJKT(session.UserID, session.JKT)
	}
	return s.cacheKeyBySessionID(session.UserID, session.SessionID)
}

// get cache
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

func (s *SessionService) CreateSession(ctx context.Context, session *SessionCache) error {
	if session == nil {
		return fmt.Errorf("session cannot be nil")
	}
	if session.SessionID == "" {
		return fmt.Errorf("SessionID required")
	}
	if session.UserID == 0 {
		return fmt.Errorf("user_id required")
	}
	if session.AccessToken == "" {
		return fmt.Errorf("access_token required")
	}
	if session.LastActive == 0 {
		session.LastActive = time.Now().Unix()
	}
	data, err := json.Marshal(session)
	if err != nil {
		return err
	}
	ttl := time.Duration(s.SessionTTL) * time.Second
	return s.redis.Set(ctx, s.cacheKeyFor(session), data, ttl).Err()
}

func (s *SessionService) GetSessionByJKT(ctx context.Context, userID int64, jkt string) (*SessionCache, error) {
	return s.getByKey(ctx, s.cacheKeyByJKT(userID, jkt))
}

func (s *SessionService) GetSessionBySessionID(ctx context.Context, userID int64, sessionID string) (*SessionCache, error) {
	return s.getByKey(ctx, s.cacheKeyBySessionID(userID, sessionID))
}

func (s *SessionService) DeleteSessionBySessionID(ctx context.Context, userID int64, sessionID string) error {
	return s.redis.Del(ctx, s.cacheKeyBySessionID(userID, sessionID)).Err()
}

func (s *SessionService) DeleteSessionByJKT(ctx context.Context, userID int64, jkt string) error {
	return s.redis.Del(ctx, s.cacheKeyByJKT(userID, jkt)).Err()
}

func (s *SessionService) UpdateSessionBySessionID(ctx context.Context, userID int64, sessionID, newAccessToken string) error {
	session, err := s.GetSessionBySessionID(ctx, userID, sessionID)
	if err != nil {
		return err
	}
	if session == nil {
		return fmt.Errorf("session not found for user %d", userID)
	}
	session.AccessToken = newAccessToken
	session.LastActive = time.Now().Unix()
	return s.CreateSession(ctx, session)
}

func (s *SessionService) UpdateSessionByJKT(ctx context.Context, userID int64, JKT, newAccessToken string) error {
	session, err := s.GetSessionByJKT(ctx, userID, JKT)
	if err != nil {
		return err
	}
	if session == nil {
		return fmt.Errorf("session not found for user %d", userID)
	}
	session.AccessToken = newAccessToken
	session.LastActive = time.Now().Unix()
	return s.CreateSession(ctx, session)
}

func (s *SessionService) UpdateLastActiveBySessionID(ctx context.Context, userID int64, sessionID string) error {
	session, err := s.GetSessionBySessionID(ctx, userID, sessionID)
	if err != nil {
		return err
	}
	if session == nil {
		return fmt.Errorf("session not found for user %d", userID)
	}
	session.LastActive = time.Now().Unix()
	return s.CreateSession(ctx, session)
}

func (s *SessionService) UpdateLastActiveByJKT(ctx context.Context, userID int64, JKT string) error {
	session, err := s.GetSessionByJKT(ctx, userID, JKT)
	if err != nil {
		return err
	}
	if session == nil {
		return fmt.Errorf("session not found for user %d", userID)
	}
	session.LastActive = time.Now().Unix()
	return s.CreateSession(ctx, session)
}

func (s *SessionService) ValidateSession(ctx context.Context, userID int64, jkt string, sessionID string, accessToken string) (*SessionCache, error) {
	var session *SessionCache
	var err error

	if jkt != "" {
		session, err = s.GetSessionByJKT(ctx, userID, jkt)
	} else {
		session, err = s.GetSessionBySessionID(ctx, userID, sessionID)
	}
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, fmt.Errorf("session not found")
	}
	if session.AccessToken != accessToken {
		return nil, fmt.Errorf("token mismatch")
	}
	if jkt != "" && session.JKT != jkt {
		return nil, fmt.Errorf("jkt mismatch")
	}
	now := time.Now().Unix()
	if now-session.LastActive > s.IdleTimeout {
		return nil, fmt.Errorf("session idle timeout")
	}
	return session, nil
}
