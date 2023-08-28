package portal

import (
	"net/http"

	"github.com/dwarvesf/go-api/pkg/handler/v1/view"
	mw "github.com/dwarvesf/go-api/pkg/middleware"
	"github.com/dwarvesf/go-api/pkg/model"
	"github.com/dwarvesf/go-api/pkg/repository/orm"
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
	uID, err := mw.UserIDFromContext(c.Request.Context())
	if err != nil {
		util.HandleError(c, err)
		return
	}
	rs, err := h.userCtrl.Me(uID)
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
// @Success 200 {object} User
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /portal/users [put]
func (h Handler) UpdateUser(c *gin.Context) {
	uID, err := mw.UserIDFromContext(c.Request.Context())
	if err != nil {
		util.HandleError(c, err)
		return
	}

	var req view.UpdateUserRequest
	err = c.ShouldBindJSON(&req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	rs, err := h.userCtrl.UpdateUser(uID, model.UpdateUserRequest{
		FullName: req.FullName,
		Status:   req.Status,
		Avatar:   req.Avatar,
	})
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, orm.User{
		ID:    rs.ID,
		Email: rs.Email,
	})
}

// UpdatePassword godoc
// @Summary Update user's password
// @Description Update user's password
// @Tags User
// @Accept  json
// @Produce  json
// @Success 200 {object} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /portal/users/change-passwords [put]
func (h Handler) UpdatePassword(c *gin.Context) {
	var req view.UpdatePasswordRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	err = h.userCtrl.UpdatePassword(model.UpdatePasswordRequest{
		Email:          req.Email,
		NewPassword:    req.NewPassword,
		RetypePassword: req.RetypePassword,
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
