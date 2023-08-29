package portal

import (
	"net/http"

	"github.com/dwarvesf/go-api/pkg/handler/v1/view"
	mw "github.com/dwarvesf/go-api/pkg/middleware"
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
