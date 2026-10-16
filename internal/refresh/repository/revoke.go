package repository

import (
	"context"
	"time"

	modelRefresh "example/internal/refresh/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Revoke — thu hồi 1 token, clear cache, return cache error nếu fail.
func (r *RefreshTokenRepository) Revoke(ctx context.Context, token, reason string) error {
	err := r.core.DB.WithContext(ctx).
		Model(&modelRefresh.RefreshToken{}).
		Where("token = ?", token).
		Updates(map[string]interface{}{
			"revoked":        true,
			"revoked_at":     time.Now(),
			"revoked_reason": reason,
		}).Error
	if err != nil {
		return err
	}
	return r.deleteCacheToken(ctx, token)
}

// RevokeAllByUser — revoke tất cả token của user, trả []string để blacklist.
func (r *RefreshTokenRepository) RevokeAllByUser(
	ctx context.Context, userID int64, reason string,
) ([]string, error) {
	var revoked []modelRefresh.RefreshToken

	err := r.core.DB.WithContext(ctx).
		Model(&modelRefresh.RefreshToken{}).
		Clauses(clause.Returning{Columns: []clause.Column{{Name: "token"}}}).
		Where("user_id = ? AND revoked = ?", userID, false).
		Updates(map[string]interface{}{
			"revoked":        true,
			"revoked_at":     time.Now(),
			"revoked_reason": reason,
		}).Error
	if err != nil {
		return nil, err
	}
	_ = revoked // GORM Returning không scan vào slice qua Updates, xem note bên dưới

	return r.collectTokensAndClearCache(ctx, userID, "", reason)
}

// RevokeByUserAndJKT — revoke tất cả token còn sống của (userID, jkt).
func (r *RefreshTokenRepository) RevokeByUserAndJKT(
	ctx context.Context, userID int64, jkt, reason string,
) ([]string, error) {
	return r.collectTokensAndClearCache(ctx, userID, jkt, reason)
}

// collectTokensAndClearCache — helper dùng transaction + lock để atomic.
func (r *RefreshTokenRepository) collectTokensAndClearCache(
	ctx context.Context, userID int64, jkt, reason string,
) ([]string, error) {
	var tokens []string

	err := r.core.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Lock + list token còn sống
		var rows []modelRefresh.RefreshToken
		q := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ? AND revoked = ?", userID, false)
		if jkt != "" {
			q = q.Where("jkt = ?", jkt)
		}
		if err := q.Find(&rows).Error; err != nil {
			return err
		}

		if len(rows) == 0 {
			return nil
		}

		// 2. Gom token + id
		ids := make([]int64, 0, len(rows))
		for _, row := range rows {
			ids = append(ids, row.ID)
			tokens = append(tokens, row.Token)
		}

		// 3. Update theo id
		return tx.Model(&modelRefresh.RefreshToken{}).
			Where("id IN ?", ids).
			Updates(map[string]interface{}{
				"revoked":        true,
				"revoked_at":     time.Now(),
				"revoked_reason": reason,
			}).Error
	})
	if err != nil {
		return nil, err
	}

	// 4. Clear cache — best effort, không fail cả hàm
	for _, t := range tokens {
		_ = r.deleteCacheToken(ctx, t)
	}
	return tokens, nil
}
