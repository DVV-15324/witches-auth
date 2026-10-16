package utils

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)

// Định nghĩa các lỗi cụ thể
var (
	ErrTokenExpired     = errors.New("token has expired")
	ErrInvalidSignature = errors.New("invalid token signature")
	ErrMalformedToken   = errors.New("malformed token")
	ErrInvalidAlgorithm = errors.New("invalid algorithm")
	ErrTokenNotYetValid = errors.New("token not yet valid")
	ErrInvalidToken     = errors.New("invalid token")
	ErrMissingCNF       = errors.New("token missing cnf claim")
	ErrMissingJKT       = errors.New("token missing cnf.jkt")
)

type Token struct {
	Token    string `json:"token"`
	ExpireAt int64  `json:"expire_at"` // Unix timestamp (giây)
}

type TokenResponse struct {
	AccessToken  Token `json:"access_token"`
	RefreshToken Token `json:"refresh_token"`
}

// Confirmation — RFC 7800
type Confirmation struct {
	JKT string `json:"jkt,omitempty"`
}

type JwtClaims struct {
	jwt.RegisteredClaims
	CNF *Confirmation `json:"cnf,omitempty"`
}

func (c *JwtClaims) GetJKT() string {
	if c == nil || c.CNF == nil {
		return ""
	}
	return c.CNF.JKT
}

type JwtService struct {
	secretKey          string
	AccessTokenExpiry  int64 // giây
	RefreshTokenExpiry int64 // giây
}

func NewJwtService(secretKey string, accessExpiry, refreshExpiry int64) *JwtService {
	return &JwtService{
		secretKey:          secretKey,
		AccessTokenExpiry:  accessExpiry,
		RefreshTokenExpiry: refreshExpiry,
	}
}

func (j *JwtService) IssueAccessToken(ctx context.Context, sub, tid, jkt string) (*Token, error) {
	now := time.Now()
	exp := now.Add(time.Duration(j.AccessTokenExpiry) * time.Second)

	claims := JwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tid,
			Subject:   sub,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
		},
		CNF: &Confirmation{JKT: jkt},
	}

	signedToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString([]byte(j.secretKey))
	if err != nil {
		return nil, err
	}

	return &Token{Token: signedToken, ExpireAt: exp.Unix()}, nil
}

func (j *JwtService) IssueRefreshToken(ctx context.Context, sub, tid, jkt string) (*Token, error) {
	now := time.Now()
	exp := now.Add(time.Duration(j.RefreshTokenExpiry) * time.Second)

	claims := JwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tid,
			Subject:   sub,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
		},
		CNF: &Confirmation{JKT: jkt},
	}

	signedToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString([]byte(j.secretKey))
	if err != nil {
		return nil, err
	}

	return &Token{Token: signedToken, ExpireAt: exp.Unix()}, nil
}

func (j *JwtService) IssueTokenPair(ctx context.Context, sub, tid, jkt string) (*TokenResponse, error) {
	accessToken, err := j.IssueAccessToken(ctx, sub, tid, jkt)
	if err != nil {
		return nil, err
	}

	refreshToken, err := j.IssueRefreshToken(ctx, sub, tid, jkt)
	if err != nil {
		return nil, err
	}

	return &TokenResponse{
		AccessToken:  *accessToken,
		RefreshToken: *refreshToken,
	}, nil
}

func (j *JwtService) ParseToken(ctx context.Context, tokenStr string) (*JwtClaims, error) {
	var claims JwtClaims

	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, ErrInvalidAlgorithm
		}
		return []byte(j.secretKey), nil
	})

	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			return nil, ErrTokenExpired
		case errors.Is(err, jwt.ErrTokenSignatureInvalid):
			return nil, ErrInvalidSignature
		case errors.Is(err, jwt.ErrTokenMalformed):
			return nil, ErrMalformedToken
		case errors.Is(err, jwt.ErrTokenNotValidYet):
			return nil, ErrTokenNotYetValid
		case errors.Is(err, jwt.ErrTokenUnverifiable):
			return nil, ErrInvalidToken
		case errors.Is(err, ErrInvalidAlgorithm):
			return nil, ErrInvalidAlgorithm
		default:
			if errors.Is(err, jwt.ErrTokenInvalidClaims) {
				return nil, ErrInvalidToken
			}
			return nil, errors.Wrap(err, "parse token failed")
		}
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	now := time.Now()
	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(now) {
		return nil, ErrTokenExpired
	}
	if claims.NotBefore != nil && claims.NotBefore.After(now) {
		return nil, ErrTokenNotYetValid
	}

	return &claims, nil
}
