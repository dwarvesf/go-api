package portal

import (
	"net/http"

	"github.com/dwarvesf/go-api/pkg/handler/v1/view"
	"github.com/dwarvesf/go-api/pkg/model"
	"github.com/dwarvesf/go-api/pkg/util"
	"github.com/gin-gonic/gin"
)

// Login godoc
// @Summary Login to portal
// @Description Login to portal by email
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param Body body LoginRequest true "Body"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /portal/auth/login [post]
func (h Handler) Login(c *gin.Context) {
	var loginReq view.LoginRequest
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		util.HandleError(c, view.ErrBadRequest(err))
		return
	}

	rs, err := h.authCtrl.Login(model.LoginRequest{
		Email:    loginReq.Email,
		Password: loginReq.Password,
	})
	if err != nil {
		h.log.Error(err)
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, view.LoginResponse{
		Data: view.Auth{
			ID:          rs.ID,
			Email:       rs.Email,
			AccessToken: rs.AccessToken,
		},
	})
}

// Signup godoc
// @Summary Signup
// @Description Signup
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param Body body SignupRequest true "Body"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /portal/auth/signup [post]
func (h Handler) Signup(c *gin.Context) {
	var loginReq view.SignupRequest
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		util.HandleError(c, view.ErrBadRequest(err))
		return
	}

	err := h.authCtrl.Signup(model.SignupRequest{
		Email:    loginReq.Email,
		Password: loginReq.Password,
	})
	if err != nil {
		h.log.Error(err)
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, view.MessageResponse{
		Data: view.Message{
			Message: "success",
		},
	})
}
