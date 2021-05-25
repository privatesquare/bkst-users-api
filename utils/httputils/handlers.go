package httputils

// Health health controller
import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Health(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, RestMsg{Message: http.StatusText(http.StatusOK)})
}

// NoRoute no route controller handles request on endpoints that are not configured
func NoRoute(ctx *gin.Context) {
	ctx.JSON(http.StatusNotFound, RestMsg{Message: "Path not found"})
}

// MethodNotAllowed method not allowed controller handles request on known endpoints but on methods that are not configured
func MethodNotAllowed(ctx *gin.Context) {
	ctx.JSON(http.StatusMethodNotAllowed, RestMsg{Message: http.StatusText(http.StatusMethodNotAllowed)})
}
