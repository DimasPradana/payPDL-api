package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)


func Testing(c *gin.Context) {

	c.JSON(http.StatusOK, "berhasil");
}
