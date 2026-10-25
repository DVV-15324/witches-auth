package mapping

import (
	dtoAuth "github.com/DVV-15324/witches-auth/internal/auth/dto/request"
	modelAuth "github.com/DVV-15324/witches-auth/internal/auth/model"
	domainAuth "github.com/DVV-15324/witches-auth/internal/shared/domain"
	pkg_dpop "github.com/DVV-15324/witches-auth/pkg/dpop"
)

func FromRegisterToDomainAuth(dtoAuth *dtoAuth.Register) *domainAuth.Auth {
	return &domainAuth.Auth{
		Email:    dtoAuth.Email,
		Password: dtoAuth.Password,
	}
}

func FromLoginToDomainAuth(dtoAuth *dtoAuth.Login) (*domainAuth.Auth, *pkg_dpop.JWK) {
	return &domainAuth.Auth{
			Email:    dtoAuth.Email,
			Password: dtoAuth.Password,
		},
		&pkg_dpop.JWK{
			Kty: dtoAuth.JWK.Kty,
			Crv: dtoAuth.JWK.Crv,
			X:   dtoAuth.JWK.X,
			Y:   dtoAuth.JWK.Y,
		}
}

func FromDomainToModelAuth(domain *domainAuth.Auth) *modelAuth.Auth {
	return &modelAuth.Auth{
		ID:        domain.ID,
		Email:     domain.Email,
		UserId:    domain.UserId,
		Password:  domain.Password,
		Salt:      domain.Salt,
		Banned:    domain.Banned,
		AuthType:  domain.AuthType,
		CreatedAt: domain.CreatedAt,
		UpdatedAt: domain.UpdatedAt,
	}
}

func FromModelToDomainAuth(model *modelAuth.Auth) *domainAuth.Auth {
	return &domainAuth.Auth{
		ID:        model.ID,
		Email:     model.Email,
		UserId:    model.UserId,
		Password:  model.Password,
		Salt:      model.Salt,
		Banned:    model.Banned,
		AuthType:  model.AuthType,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}
