package mapping

import (
	modelRefresh "github.com/DVV-15324/witches-auth/internal/refresh/model"
	domainRefresh "github.com/DVV-15324/witches-auth/internal/shared/domain"
)

func FromDomainToModelRefresh(domain *domainRefresh.RefreshToken) *modelRefresh.RefreshToken {
	return &modelRefresh.RefreshToken{
		ID:            domain.ID,
		UserID:        domain.UserID,
		Token:         domain.Token,
		JKT:           domain.JKT,
		IPAddress:     domain.IPAddress,
		UserAgent:     domain.UserAgent,
		ExpiresAt:     domain.ExpiresAt,
		Revoked:       domain.Revoked,
		RevokedAt:     domain.RevokedAt,
		RevokedReason: domain.RevokedReason,
		Locale:        domain.Locale,
		Timezone:      domain.Timezone,
		CreatedAt:     domain.CreatedAt,
		UpdatedAt:     domain.UpdatedAt,
	}
}

func FromModelToDomainRefresh(model *modelRefresh.RefreshToken) *domainRefresh.RefreshToken {
	return &domainRefresh.RefreshToken{
		ID:            model.ID,
		UserID:        model.UserID,
		Token:         model.Token,
		JKT:           model.JKT,
		IPAddress:     model.IPAddress,
		UserAgent:     model.UserAgent,
		ExpiresAt:     model.ExpiresAt,
		Revoked:       model.Revoked,
		RevokedAt:     model.RevokedAt,
		RevokedReason: model.RevokedReason,
		Locale:        model.Locale,
		Timezone:      model.Timezone,
		CreatedAt:     model.CreatedAt,
		UpdatedAt:     model.UpdatedAt,
	}
}
