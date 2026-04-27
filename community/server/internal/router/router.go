package router

import (
	_ "net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	return r
}
