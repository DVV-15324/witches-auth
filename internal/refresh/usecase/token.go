package usecase

// import (
// 	"context"
// 	w_resp "github.com/DVV-15324/witches/pkg/core/response"
// 	w_utils "github.com/DVV-15324/witches/pkg/core/utils"
// 	"time"
// )

// // IntrospectToken - Kiểm tra token
// func (uc *RefreshUseCase) IntrospectToken(ctx context.Context, accessToken string) (*w_utils.JwtClaims, *w_resp.AppError) {
// 	// Parse token
// 	claims, err := uc.Core.JWT.ParseToken(ctx, accessToken)
// 	if err != nil {
// 		return nil, w_resp.NewAppError(401, err, time.Now())
// 	}

// 	return claims, nil
// }
