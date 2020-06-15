package libs

import (
	"github.com/gin-gonic/gin"
)

// ResponseRESTAPI func
func ResponseRESTAPI(responseData gin.H, c *gin.Context, status int) {
	ResponseType := c.Request.Header.Get("ResponseType")
	if ResponseType == "application/xml" {
		c.XML(status, responseData)
	} else {
		c.JSON(status, responseData)
	}
}
