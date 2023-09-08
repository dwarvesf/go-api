package portal

import (
	"net/http"

	"github.com/dwarvesf/go-api/pkg/handler/v1/view"
	"github.com/dwarvesf/go-api/pkg/model"
	"github.com/dwarvesf/go-api/pkg/util"
	"github.com/gin-gonic/gin"
)

// Me godoc
// @Summary Retrieve my information
// @Description Retrieve my information
// @Tags User
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Success 200 {object} MeResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /portal/me [get]
func (h Handler) Me(c *gin.Context) {
	const spanName = "meHandler"
	newCtx, span := h.monitor.NewSpan(c.Request.Context(), spanName)
	defer span.End()

	// Update c ctx to newCtx
	c.Request = c.Request.WithContext(newCtx)

	rs, err := h.userCtrl.Me(c.Request.Context())
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, view.MeResponse{
		Data: view.Me{
			ID:    rs.ID,
			Email: rs.Email,
		},
	})
}

// UpdateUser godoc
// @Summary Update user
// @Description Update user
// @Tags User
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param body body UpdateUserRequest true "Update user"
// @Success 200 {object} User
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /portal/users [put]
func (h Handler) UpdateUser(c *gin.Context) {
	const spanName = "updateUserHandler"
	newCtx, span := h.monitor.NewSpan(c.Request.Context(), spanName)
	defer span.End()

	// Update c ctx to newCtx
	c.Request = c.Request.WithContext(newCtx)

	var req view.UpdateUserRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	rs, err := h.userCtrl.UpdateUser(
		c.Request.Context(),
		model.UpdateUserRequest{
			FullName: req.FullName,
			Avatar:   req.Avatar,
		})
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, view.UserResponse{
		Data: view.User{
			ID:       rs.ID,
			Email:    rs.Email,
			FullName: rs.FullName,
			Avatar:   rs.Avatar,
		},
	})
}

// UpdatePassword godoc
// @Summary Update user's password
// @Description Update user's password
// @Tags User
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param body body UpdatePasswordRequest true "Update user"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /portal/users/password [put]
func (h Handler) UpdatePassword(c *gin.Context) {
	const spanName = "updatePasswordHandler"
	newCtx, span := h.monitor.NewSpan(c.Request.Context(), spanName)
	defer span.End()

	// Update c ctx to newCtx
	c.Request = c.Request.WithContext(newCtx)

	var req view.UpdatePasswordRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	err = h.userCtrl.UpdatePassword(
		c.Request.Context(),
		model.UpdatePasswordRequest{
			NewPassword: req.NewPassword,
			OldPassword: req.OldPassword,
		})
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, view.MessageResponse{
		Data: view.Message{
			Message: "success",
		},
	})
}
