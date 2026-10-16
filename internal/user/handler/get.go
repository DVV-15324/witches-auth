package handler

import (
	"errors"
	utils "example/internal/shared/utils"
	dtoUser "example/internal/user/dto/response"
	mapping "example/internal/user/mapping"
	"time"

	w_resp "github.com/DVV-15324/witches/pkg/core/response"
	w_utils "github.com/DVV-15324/witches/pkg/core/utils"
	"github.com/gin-gonic/gin"
)

func (h *UserHandle) HandleGetAllUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		req := w_utils.PaginationRequest{}
		// Bind JSON first
		if err := c.ShouldBindJSON(&req); err != nil {
			w_resp.WriteErrorWithLog(c, h.Core.Logger, h.Core.Config, w_resp.NewAppError(400, errors.New("invalid request body"), time.Now()))
			return
		}
		user, total, resp := h.UserUseCase.GetAllUser(c, &req)
		if resp != nil {
			w_resp.WriteErrorWithLog(c, h.Core.Logger, h.Core.Config, resp)
			return
		}
		//Chuyển user từ model -> dto
		dto := mapping.FromDomainToDtoUserList(user)
		paginationResp := req.ToPaginationResponse(total)
		w_resp.WriteSuccessWithPagination(c, dto, &paginationResp)
	}
}

func (h *UserHandle) HandleGetUserById() func(c *gin.Context) {
	return func(c *gin.Context) {
		//uid := shared.DecodeFromBase58(id)
		id := w_utils.GetRequestContext(c.Request.Context(), h.Core.Config).Sub
		uid := utils.DecodeFromBase58(id)
		user, resp := h.UserUseCase.GetUserById(c, uid.LocalID)
		if resp != nil {
			w_resp.WriteErrorWithLog(c, h.Core.Logger, h.Core.Config, resp)
			return
		}
		dto := &dtoUser.User{}
		dto = mapping.FromDomainToDtoUser(user)
		w_resp.WriteSuccess(c, dto)
	}
}
