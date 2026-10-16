package usecase

import (
	"context"
	"errors"
	"time"

	pkg "example/pkg"
	w_resp "github.com/DVV-15324/witches/pkg/core/response"

	domainAuth "example/internal/shared/domain"
	domainUser "example/internal/shared/domain"
)

func (a *AuthUseCase) Register(ctx context.Context, auth *domainAuth.Auth, name string) *w_resp.AppError {
	err := a.Core.TxManager.WithinTransaction(ctx, func(ctx context.Context) error {
		authEmail, err := a.AuthReponsitory.GetAuthByEmail(ctx, auth.Email)
		if authEmail != nil {
			return errors.New("email da ton tai")
		}
		if err != nil {
			return err
		}
		randomString, err := pkg.RandomStr(a.Core.HashLen)
		if err != nil {
			return err
		}

		hash, err := a.Core.Hash.GenerateFromPassword(auth.Password, randomString)
		if err != nil {
			return err
		}
		userUid, err_create := a.UserUseCase.CreateUser(ctx, &domainUser.User{
			Email: auth.Email,
			Name:  name,
		})
		if err_create != nil {
			return err_create
		}
		if userUid == 0 {
			return errors.New("user id is 0, cannot create auth")
		}
		err_auth := a.AuthReponsitory.CreateAuth(ctx, &domainAuth.Auth{
			Salt:     randomString,
			Email:    auth.Email,
			Password: hash,
			UserId:   userUid,
		})
		if err_auth != nil {
			return err_auth
		}
		return nil
	})
	if err != nil {
		if err == errors.New("email da ton tai") {
			return w_resp.NewAppError(409, err, time.Now())
		}
		return w_resp.NewAppError(500, err, time.Now())
	}

	return nil
}
