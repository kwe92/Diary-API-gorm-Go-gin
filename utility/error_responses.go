package utility

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SendBadRequestResponse(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

}
