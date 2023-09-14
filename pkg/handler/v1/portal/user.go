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
// @ID getMe
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
// @ID updateUser
// @Tags User
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param body body UpdateUserRequest true "Update user"
// @Success 200 {object} UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /portal/users [put]
func (h Handler) UpdateUser(c *gin.Context) {
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
// @ID updatePassword
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

// GetUsersList godoc
// @Summary Get users list
// @Description get users list
// @ID getUsersList
// @Tags User
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param page query int false "Page"
// @Param pageSize query int false "Page size"
// @Success 200 {object} UsersListResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /portal/users [get]
func (h Handler) GetUsersList(c *gin.Context) {
	var req view.GetUsersListRequest
	err := c.ShouldBindQuery(&req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	rs, err := h.userCtrl.UserList(
		c.Request.Context(),
		model.ListQuery{
			Page:     req.Page,
			PageSize: req.PageSize,
		})
	if err != nil {
		util.HandleError(c, err)
		return
	}

	users := make([]view.User, 0, len(rs.Data))
	for _, u := range rs.Data {
		users = append(users, view.User{
			ID:         u.ID,
			Email:      u.Email,
			FullName:   u.FullName,
			Avatar:     u.Avatar,
			Status:     u.Status,
			Title:      u.Title,
			Department: u.Department,
			Role:       u.Role,
		})
	}

	c.JSON(http.StatusOK, view.UsersListResponse{
		Data: users,
		Metadata: view.Metadata{
			Page:         rs.Pagination.Page,
			PageSize:     rs.Pagination.PageSize,
			TotalPages:   rs.Pagination.TotalPages,
			TotalRecords: rs.Pagination.TotalRecords,
			Sort:         rs.Pagination.Sort,
		},
	})
}
