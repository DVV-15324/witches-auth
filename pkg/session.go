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

// Mobile (DPoP): lookup by (userID, jkt)
func (s *SessionService) cacheKeyByJKT(userID int64, jkt string) string {
	return fmt.Sprintf("session:user:%d:jkt:%s", userID, jkt)
}

// Web (Bearer): lookup by (userID, sha256(accessToken)) → multi-session OK
func (s *SessionService) cacheKeyByToken(userID int64, accessToken string) string {
	h := sha256.Sum256([]byte(accessToken))
	return fmt.Sprintf("session:user:%d:web:%x", userID, h)
}

// Dispatch dựa trên JKT có rỗng hay không
func (s *SessionService) cacheKey(userID int64, jkt, accessToken string) string {
	if jkt == "" {
		return s.cacheKeyByToken(userID, accessToken)
	}
	return s.cacheKeyByJKT(userID, jkt)
}

// CreateSession - Tạo/update session.
// Mobile: cần JKT + AccessToken
// Web   : cần AccessToken (JKT rỗng)
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

	key := s.cacheKey(session.UserID, session.JKT, session.AccessToken)
	data, err := json.Marshal(session)
	if err != nil {
		return err
	}

	ttl := time.Duration(s.SessionTTL) * time.Second
	return s.redis.Set(ctx, key, data, ttl).Err()
}

// GetSession - Mobile: lookup theo (userID, jkt)
// Trả (nil, nil) nếu không tồn tại.
func (s *SessionService) GetSession(ctx context.Context, userID int64, jkt string) (*SessionCache, error) {
	return s.getByKey(ctx, s.cacheKeyByJKT(userID, jkt))
}

// GetSessionByToken - Web: lookup theo (userID, accessToken)
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

// DeleteSession - Xóa session (mobile: theo jkt; web: theo accessToken)
func (s *SessionService) DeleteSession(ctx context.Context, userID int64, jkt string) error {
	return s.redis.Del(ctx, s.cacheKeyByJKT(userID, jkt)).Err()
}

func (s *SessionService) DeleteSessionByToken(ctx context.Context, userID int64, accessToken string) error {
	return s.redis.Del(ctx, s.cacheKeyByToken(userID, accessToken)).Err()
}

// UpdateSession - Update access token khi refresh.
// Lưu ý: nếu refresh đổi accessToken (web), key cũng đổi → phải xóa key cũ trước.
func (s *SessionService) UpdateSession(ctx context.Context, userID int64, jkt, oldAccessToken, newAccessToken string) error {
	var (
		session *SessionCache
		err     error
	)
	if jkt != "" {
		session, err = s.GetSession(ctx, userID, jkt)
	} else {
		session, err = s.GetSessionByToken(ctx, userID, oldAccessToken)
	}
	if err != nil {
		return err
	}
	if session == nil {
		return fmt.Errorf("session not found for user %d", userID)
	}

	// Web: key dựa trên token → xóa key cũ, tạo key mới
	if jkt == "" && oldAccessToken != newAccessToken {
		_ = s.DeleteSessionByToken(ctx, userID, oldAccessToken)
	}

	session.AccessToken = newAccessToken
	session.LastActive = time.Now().Unix()
	return s.CreateSession(ctx, session)
}

// UpdateLastActive - mobile: dùng jkt; web: dùng accessToken
func (s *SessionService) UpdateLastActive(ctx context.Context, userID int64, jkt, accessToken string) error {
	var (
		session *SessionCache
		err     error
	)
	if jkt != "" {
		session, err = s.GetSession(ctx, userID, jkt)
	} else {
		session, err = s.GetSessionByToken(ctx, userID, accessToken)
	}
	if err != nil {
		return err
	}
	if session == nil {
		return fmt.Errorf("session not found for user %d", userID)
	}

	session.LastActive = time.Now().Unix()
	return s.CreateSession(ctx, session)
}

// ValidateSession - authorize request
func (s *SessionService) ValidateSession(
	ctx context.Context,
	userID int64,
	jkt string,
	accessToken string,
) (*SessionCache, error) {
	var (
		session *SessionCache
		err     error
	)

	if jkt != "" {
		session, err = s.GetSession(ctx, userID, jkt)
	} else {
		session, err = s.GetSessionByToken(ctx, userID, accessToken)
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

	// Chỉ check JKT khi là mobile (jkt có giá trị)
	if jkt != "" && session.JKT != jkt {
		return nil, fmt.Errorf("jkt mismatch")
	}

	now := time.Now().Unix()
	if now-session.LastActive > s.IdleTimeout {
		return nil, fmt.Errorf("session idle timeout")
	}

	return session, nil
}
