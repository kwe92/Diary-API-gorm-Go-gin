package utility

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// TODO: replace all ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) with call to SendBadRequestResponse

func SendBadRequestResponse(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

}
