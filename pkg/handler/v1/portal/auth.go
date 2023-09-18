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
// @id login
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
	const spanName = "loginHandler"
	ctx, span := h.monitor.Start(c.Request.Context(), spanName)
	defer span.End()

	var req view.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.HandleError(c, view.ErrBadRequest(err))
		return
	}

	rs, err := h.authCtrl.Login(ctx, model.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
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
// @id signup
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
	const spanName = "signupHandler"
	ctx, span := h.monitor.Start(c.Request.Context(), spanName)
	defer span.End()

	var req view.SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.HandleError(c, view.ErrBadRequest(err))
		return
	}

	err := h.authCtrl.Signup(ctx, model.SignupRequest{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.FullName,
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
